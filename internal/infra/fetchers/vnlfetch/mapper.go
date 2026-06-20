package vnlfetch

import (
	"strconv"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Upstream matchStatus codes.
const (
	statusScheduled = 0
	statusLive      = 1
	statusFinished  = 2
)

const (
	sidesPerMatch = 2
	langKeyLen    = 2
)

// Best-of-five estimate; the feed carries no end time.
const volleyballMatchDuration = 120 * time.Minute

// Grace past the expected end before a stale live match is forced finished.
const volleyballCompletionGrace = 60 * time.Minute

type mappedSchedule struct {
	competitions []eventcore.Competition
	seasons      []eventcore.Season
	stages       []eventcore.Stage
	groups       []eventcore.Group
	venues       []eventcore.Venue
	participants []eventcore.Participant
	fixtures     []eventcore.Fixture
	results      []eventcore.Result
}

type scheduleBuilder struct {
	out         mappedSchedule
	tournaments map[int]apiTournament
	teams       map[int]apiTeam
	lang        string
	day         time.Time
	compSeen    map[string]struct{}
	seasonSeen  map[int]struct{}
	stageSeen   map[string]struct{}
	groupSeen   map[int]struct{}
	venueSeen   map[string]struct{}
	partSeen    map[int]struct{}
}

func newScheduleBuilder(
	resp apiScheduleResponse, lang string, day time.Time,
) *scheduleBuilder {
	tournaments := make(map[int]apiTournament, len(resp.AllTournaments))
	for _, t := range resp.AllTournaments {
		tournaments[t.No] = t
	}
	teams := make(map[int]apiTeam, len(resp.AllTeams))
	for _, t := range resp.AllTeams {
		teams[t.No] = t
	}
	return &scheduleBuilder{
		tournaments: tournaments,
		teams:       teams,
		lang:        lang,
		day:         day,
		compSeen:    make(map[string]struct{}),
		seasonSeen:  make(map[int]struct{}),
		stageSeen:   make(map[string]struct{}),
		groupSeen:   make(map[int]struct{}),
		venueSeen:   make(map[string]struct{}),
		partSeen:    make(map[int]struct{}),
	}
}

func mapSchedule(
	resp apiScheduleResponse, lang string, day, now time.Time,
) mappedSchedule {
	b := newScheduleBuilder(resp, lang, day)
	for _, m := range resp.Matches {
		// Knockout slots are published without teams; skip until both sides exist.
		if m.IsMatchTBD || m.TeamANo == 0 || m.TeamBNo == 0 {
			continue
		}

		stageID := b.emitHierarchy(m)
		groupID := b.emitGroup(m, stageID)
		venueID := b.emitVenue(m)
		fixParts := b.emitParticipants(m)

		startsAt := m.MatchDateUtc.UTC()
		endsAt := startsAt.Add(volleyballMatchDuration)
		status := eventcore.CompleteByEndTime(
			mapStatus(m.MatchStatus), endsAt, now, volleyballCompletionGrace,
		)
		f := eventcore.Fixture{
			ID: eventcore.NewID(eventcore.ProviderVNL, matchKey(m.MatchNo)),
			Ext: eventcore.ExternalID{
				Provider: eventcore.ProviderVNL,
				Key:      matchKey(m.MatchNo),
			},
			StageID:      stageID,
			GroupID:      groupID,
			VenueID:      venueID,
			Name:         b.teamName(m.TeamANo) + " vs " + b.teamName(m.TeamBNo),
			StartsAt:     startsAt,
			EndsAt:       endsAt,
			Status:       status,
			Participants: fixParts,
		}

		// Results feed the checksum so score changes re-notify.
		results := mapResults(m, f.ID, f.Status)
		f.Checksum = f.ComputeChecksumWith(results)

		b.out.fixtures = append(b.out.fixtures, f)
		b.out.results = append(b.out.results, results...)
	}
	return b.out
}

func mapStatus(raw int) eventcore.FixtureStatus {
	switch raw {
	case statusFinished:
		return eventcore.FixtureFinished
	case statusLive:
		return eventcore.FixtureLive
	default:
		return eventcore.FixtureScheduled
	}
}

// Namespaced external keys: the feed's small ids overlap across entity kinds.
func compKey(slug string) string   { return "comp_" + slug }
func seasonKey(no int) string      { return "season_" + strconv.Itoa(no) }
func stageKey(tNo, rNo int) string { return "stage_" + strconv.Itoa(tNo) + "_" + strconv.Itoa(rNo) }
func poolKey(no int) string        { return "pool_" + strconv.Itoa(no) }
func teamKey(no int) string        { return "team_" + strconv.Itoa(no) }
func matchKey(no int) string       { return "match_" + strconv.Itoa(no) }

func venueKey(countryCode, city string) string { return "venue_" + countryCode + "_" + city }
