package olympicsfetch

import (
	"strconv"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Upstream medal-type codes carried in competitor results.
const (
	medalCodeGold   = "ME_GOLD"
	medalCodeSilver = "ME_SILVER"
)

// mappedSchedule is the relational chain produced from a schedule response,
// deduped by external key.
type mappedSchedule struct {
	competitions []eventcore.Competition
	seasons      []eventcore.Season
	stages       []eventcore.Stage
	groups       []eventcore.Group
	participants []eventcore.Participant
	fixtures     []eventcore.Fixture
	results      []eventcore.Result
}

type scheduleBuilder struct {
	out        mappedSchedule
	compSeen   map[string]struct{}
	seasonSeen map[string]struct{}
	stageSeen  map[string]struct{}
	groupSeen  map[string]struct{}
	partSeen   map[string]struct{}
}

func newScheduleBuilder() *scheduleBuilder {
	return &scheduleBuilder{
		compSeen:   make(map[string]struct{}),
		seasonSeen: make(map[string]struct{}),
		stageSeen:  make(map[string]struct{}),
		groupSeen:  make(map[string]struct{}),
		partSeen:   make(map[string]struct{}),
	}
}

func (b *scheduleBuilder) emitHierarchy(
	unit apiUnit, group apiGroup,
) (stageID eventcore.CanonicalID, groupID *eventcore.CanonicalID) {
	compID := b.emitCompetition(unit)
	seasonID := b.emitSeason(unit, compID)
	stageID = b.emitStage(unit, seasonID)
	groupID = b.emitGroup(unit, group, stageID)
	return stageID, groupID
}

func (b *scheduleBuilder) emitCompetition(unit apiUnit) eventcore.CanonicalID {
	compKey := "comp_" + unit.DisciplineCode
	compID := eventcore.NewID(eventcore.ProviderOlympics, compKey)
	if _, ok := b.compSeen[compKey]; !ok {
		b.compSeen[compKey] = struct{}{}
		b.out.competitions = append(b.out.competitions, eventcore.Competition{
			ID: compID,
			Ext: eventcore.ExternalID{
				Provider: eventcore.ProviderOlympics,
				Key:      compKey,
			},
			Code:       unit.DisciplineCode,
			Name:       unit.DisciplineName,
			Discipline: unit.DisciplineName,
		})
	}
	return compID
}

func (b *scheduleBuilder) emitSeason(
	unit apiUnit, compID eventcore.CanonicalID,
) eventcore.CanonicalID {
	seasonKey := "season_" + unit.DisciplineCode
	seasonID := eventcore.NewID(eventcore.ProviderOlympics, seasonKey)
	if _, ok := b.seasonSeen[seasonKey]; !ok {
		b.seasonSeen[seasonKey] = struct{}{}
		b.out.seasons = append(b.out.seasons, eventcore.Season{
			ID: seasonID,
			Ext: eventcore.ExternalID{
				Provider: eventcore.ProviderOlympics,
				Key:      seasonKey,
			},
			CompetitionID: compID,
			Name:          unit.DisciplineName,
		})
	}
	return seasonID
}

func (b *scheduleBuilder) emitStage(
	unit apiUnit, seasonID eventcore.CanonicalID,
) eventcore.CanonicalID {
	stageKey := unit.PhaseId
	if stageKey == "" {
		stageKey = unit.PhaseCode
	}
	if stageKey == "" {
		stageKey = "stage_" + unit.DisciplineCode
	}
	stageID := eventcore.NewID(eventcore.ProviderOlympics, stageKey)
	if _, ok := b.stageSeen[stageKey]; !ok {
		b.stageSeen[stageKey] = struct{}{}
		b.out.stages = append(b.out.stages, eventcore.Stage{
			ID:       stageID,
			Ext:      eventcore.ExternalID{Provider: eventcore.ProviderOlympics, Key: stageKey},
			SeasonID: seasonID,
			Name:     unit.PhaseName,
			Ord:      unit.EventOrder,
		})
	}
	return stageID
}

// emitGroup returns nil when the unit has no group.
func (b *scheduleBuilder) emitGroup(
	unit apiUnit, group apiGroup, stageID eventcore.CanonicalID,
) *eventcore.CanonicalID {
	if unit.GroupId == "" {
		return nil
	}
	gid := eventcore.NewID(eventcore.ProviderOlympics, unit.GroupId)
	if _, ok := b.groupSeen[unit.GroupId]; !ok {
		b.groupSeen[unit.GroupId] = struct{}{}
		b.out.groups = append(b.out.groups, eventcore.Group{
			ID: gid,
			Ext: eventcore.ExternalID{
				Provider: eventcore.ProviderOlympics,
				Key:      unit.GroupId,
			},
			StageID: stageID,
			Name:    group.Title,
		})
	}
	return &gid
}

func (b *scheduleBuilder) emitParticipants(unit apiUnit) []eventcore.FixtureParticipant {
	gender := mapGender(unit.GenderCode)
	fixParts := make([]eventcore.FixtureParticipant, 0, len(unit.Competitors))
	for _, c := range unit.Competitors {
		pid := eventcore.NewID(eventcore.ProviderOlympics, c.Code)
		if _, ok := b.partSeen[c.Code]; !ok {
			b.partSeen[c.Code] = struct{}{}
			b.out.participants = append(b.out.participants, eventcore.Participant{
				ID: pid,
				Ext: eventcore.ExternalID{
					Provider: eventcore.ProviderOlympics,
					Key:      c.Code,
				},
				Kind:       eventcore.ParticipantAthlete,
				Name:       c.Name,
				Code:       c.Code,
				CountryISO: c.Noc,
				Gender:     gender,
			})
		}
		fixParts = append(fixParts, eventcore.FixtureParticipant{
			ParticipantID: pid,
			Role:          "athlete",
		})
	}
	return fixParts
}

func mapSchedule(resp apiScheduleResponse) mappedSchedule {
	groupMap := make(map[string]apiGroup, len(resp.Groups))
	for _, g := range resp.Groups {
		groupMap[g.Id] = g
	}

	b := newScheduleBuilder()

	for _, unit := range resp.Units {
		group := groupMap[unit.GroupId]

		stageID, groupID := b.emitHierarchy(unit, group)

		fixParts := b.emitParticipants(unit)

		extKey := unit.Id
		if extKey == "" {
			extKey = unit.UnitNum + "_" + unit.SessionCode
		}
		f := eventcore.Fixture{
			ID:           eventcore.NewID(eventcore.ProviderOlympics, extKey),
			Ext:          eventcore.ExternalID{Provider: eventcore.ProviderOlympics, Key: extKey},
			StageID:      stageID,
			GroupID:      groupID,
			Name:         unit.EventUnitName,
			StartsAt:     unit.StartDate.UTC(),
			EndsAt:       unit.EndDate.UTC(),
			Status:       mapStatus(unit.Status, group.IsLive),
			Participants: fixParts,
		}

		// Results are built before the checksum so medal changes re-notify.
		results := mapResults(unit.Competitors, f.ID)
		f.Checksum = f.ComputeChecksumWith(results)

		b.out.fixtures = append(b.out.fixtures, f)
		b.out.results = append(b.out.results, results...)
	}

	return b.out
}

func mapStatus(raw string, isLive bool) eventcore.FixtureStatus {
	switch strings.ToLower(raw) {
	case "scheduled":
		return eventcore.FixtureScheduled
	case "finished":
		return eventcore.FixtureFinished
	case "cancelled":
		return eventcore.FixtureCancelled
	}
	if isLive {
		return eventcore.FixtureLive
	}
	return eventcore.FixtureScheduled
}

// mapGender maps upstream "W" -> "F" and "M" -> "M", else "".
func mapGender(genderCode string) string {
	switch strings.ToUpper(genderCode) {
	case "W":
		return "F"
	case "M":
		return "M"
	}
	return ""
}

// mapResults keeps a competitor only when it has a medal OR a non-empty mark.
func mapResults(competitors []apiCompetitor, fixtureID eventcore.CanonicalID) []eventcore.Result {
	out := make([]eventcore.Result, 0, len(competitors))
	for _, c := range competitors {
		outcome := mapOutcome(c.Results.MedalType, c.Results.WinnerLoserTie)
		hasMedal := outcome == eventcore.OutcomeMedalGold ||
			outcome == eventcore.OutcomeMedalSilver ||
			outcome == eventcore.OutcomeMedalBronze
		if !hasMedal && c.Results.Mark == "" {
			continue
		}

		pid := eventcore.NewID(eventcore.ProviderOlympics, c.Code)
		res := eventcore.Result{
			FixtureID:     fixtureID,
			ParticipantID: pid,
			Score:         c.Results.Mark,
			RawMark:       c.Results.Irm,
			Outcome:       outcome,
		}
		if n, err := strconv.Atoi(c.Results.Position); err == nil {
			res.Position = &n
		}
		out = append(out, res)
	}
	return out
}

func mapOutcome(medalType, winnerLoserTie string) eventcore.Outcome {
	switch strings.ToUpper(medalType) {
	case "GM", medalCodeGold:
		return eventcore.OutcomeMedalGold
	case "SM", medalCodeSilver:
		return eventcore.OutcomeMedalSilver
	case "BM", "ME_BRONZE":
		return eventcore.OutcomeMedalBronze
	}
	switch strings.ToUpper(winnerLoserTie) {
	case "W":
		return eventcore.OutcomeWin
	case "L":
		return eventcore.OutcomeLoss
	case "T":
		return eventcore.OutcomeDraw
	}
	return eventcore.OutcomeNone
}
