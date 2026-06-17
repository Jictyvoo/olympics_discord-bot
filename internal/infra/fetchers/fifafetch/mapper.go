package fifafetch

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Upstream MatchStatus codes.
const (
	statusFinished  = 0
	statusScheduled = 1
	statusLive      = 3
)

// footballMatchDuration is 90 min play plus 15 min halftime; the feed carries no
// kickoff end, so fixtures derive their end from the start.
const footballMatchDuration = 105 * time.Minute

// footballCompletionGrace covers stoppage and extra time before a match the feed
// never finalizes is treated as finished.
const footballCompletionGrace = 60 * time.Minute

// mappedMatches is the relational chain produced from a matches response,
// deduped by external key.
type mappedMatches struct {
	competitions []eventcore.Competition
	seasons      []eventcore.Season
	stages       []eventcore.Stage
	groups       []eventcore.Group
	venues       []eventcore.Venue
	participants []eventcore.Participant
	fixtures     []eventcore.Fixture
	results      []eventcore.Result
	stageKeys    []string // distinct stage keys, in first-seen order
}

type matchBuilder struct {
	out        mappedMatches
	compSeen   map[string]struct{}
	seasonSeen map[string]struct{}
	stageSeen  map[string]struct{}
	groupSeen  map[string]struct{}
	partSeen   map[string]struct{}
	venueSeen  map[string]struct{}
}

func newMatchBuilder() *matchBuilder {
	return &matchBuilder{
		compSeen:   make(map[string]struct{}),
		seasonSeen: make(map[string]struct{}),
		stageSeen:  make(map[string]struct{}),
		groupSeen:  make(map[string]struct{}),
		partSeen:   make(map[string]struct{}),
		venueSeen:  make(map[string]struct{}),
	}
}

// seasonMeta holds the season name and span from the season endpoint; the
// matches feed carries no season dates.
type seasonMeta struct {
	name     string
	startsOn time.Time
	endsOn   time.Time
}

func mapMatches(
	resp apiMatchesResponse,
	lang string,
	season seasonMeta,
	now time.Time,
) mappedMatches {
	b := newMatchBuilder()
	for _, m := range resp.Results {
		// Knockout slots are published without teams; skip until both sides exist.
		if m.Home.IdTeam == "" || m.Away.IdTeam == "" {
			continue
		}

		stageID, groupID := b.emitHierarchy(m, lang, season)
		venueID := b.emitVenue(m.Stadium, lang)
		fixParts := b.emitParticipants(m, lang)

		startsAt := m.Date.UTC()
		endsAt := startsAt.Add(footballMatchDuration)
		status := eventcore.CompleteByEndTime(
			mapStatus(m.MatchStatus), endsAt, now, footballCompletionGrace,
		)
		f := eventcore.Fixture{
			ID:      eventcore.NewID(eventcore.ProviderFIFA, m.IdMatch),
			Ext:     eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: m.IdMatch},
			StageID: stageID,
			GroupID: groupID,
			VenueID: venueID,
			Name: localized(m.Home.TeamName, lang) + " vs " +
				localized(m.Away.TeamName, lang),
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

func (b *matchBuilder) emitHierarchy(
	m apiMatch, lang string, season seasonMeta,
) (stageID eventcore.CanonicalID, groupID *eventcore.CanonicalID) {
	compID := b.emitCompetition(m, lang)
	seasonID := b.emitSeason(m, compID, lang, season)
	stageID = b.emitStage(m, seasonID, lang)
	groupID = b.emitGroup(m, stageID, lang)
	return stageID, groupID
}
