package vnlfetch

import (
	"strconv"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func (b *scheduleBuilder) emitHierarchy(m apiMatch) eventcore.CanonicalID {
	compID := b.emitCompetition(m)
	seasonID := b.emitSeason(m, compID)
	return b.emitStage(m, seasonID)
}

func (b *scheduleBuilder) emitCompetition(m apiMatch) eventcore.CanonicalID {
	key := compKey(m.CompetitionSlug)
	id := eventcore.NewID(eventcore.ProviderVNL, key)
	if _, ok := b.compSeen[key]; !ok {
		b.compSeen[key] = struct{}{}
		b.out.competitions = append(b.out.competitions, eventcore.Competition{
			ID:         id,
			Ext:        eventcore.ExternalID{Provider: eventcore.ProviderVNL, Key: key},
			Code:       m.CompetitionSlug,
			Name:       m.CompetitionFullName,
			Discipline: volleyballDiscipline(b.lang),
		})
	}
	return id
}

func (b *scheduleBuilder) emitSeason(
	m apiMatch, compID eventcore.CanonicalID,
) eventcore.CanonicalID {
	key := seasonKey(m.TournamentNo)
	id := eventcore.NewID(eventcore.ProviderVNL, key)
	if _, ok := b.seasonSeen[m.TournamentNo]; !ok {
		b.seasonSeen[m.TournamentNo] = struct{}{}
		t := b.tournaments[m.TournamentNo]
		name := t.Name
		if name == "" {
			name = m.CompetitionFullName
		}
		// The not-null timestamp columns need a valid span; fall back to the
		// requested day when the tournament metadata is absent.
		startsOn := t.StartDate.UTC()
		endsOn := t.EndDate.UTC()
		if startsOn.IsZero() {
			startsOn = b.day.UTC()
		}
		if endsOn.IsZero() {
			endsOn = startsOn
		}
		b.out.seasons = append(b.out.seasons, eventcore.Season{
			ID:            id,
			Ext:           eventcore.ExternalID{Provider: eventcore.ProviderVNL, Key: key},
			CompetitionID: compID,
			Name:          name,
			StartsOn:      startsOn,
			EndsOn:        endsOn,
		})
	}
	return id
}

func (b *scheduleBuilder) emitStage(
	m apiMatch, seasonID eventcore.CanonicalID,
) eventcore.CanonicalID {
	key := stageKey(m.TournamentNo, m.RoundNo)
	id := eventcore.NewID(eventcore.ProviderVNL, key)
	if _, ok := b.stageSeen[key]; !ok {
		b.stageSeen[key] = struct{}{}
		ord, _ := strconv.Atoi(m.RoundCode)
		b.out.stages = append(b.out.stages, eventcore.Stage{
			ID:       id,
			Ext:      eventcore.ExternalID{Provider: eventcore.ProviderVNL, Key: key},
			SeasonID: seasonID,
			Name:     m.RoundName,
			Ord:      ord,
		})
	}
	return id
}

func (b *scheduleBuilder) emitGroup(
	m apiMatch, stageID eventcore.CanonicalID,
) *eventcore.CanonicalID {
	if m.Pool.No == 0 {
		return nil
	}
	key := poolKey(m.Pool.No)
	id := eventcore.NewID(eventcore.ProviderVNL, key)
	if _, ok := b.groupSeen[m.Pool.No]; !ok {
		b.groupSeen[m.Pool.No] = struct{}{}
		b.out.groups = append(b.out.groups, eventcore.Group{
			ID:      id,
			Ext:     eventcore.ExternalID{Provider: eventcore.ProviderVNL, Key: key},
			StageID: stageID,
			Name:    m.Pool.Name,
		})
	}
	return &id
}

// No arena in the feed: the host city is the venue, else the Discord event
// location shows "TBD". nil when the match carries no city.
func (b *scheduleBuilder) emitVenue(m apiMatch) *eventcore.CanonicalID {
	if m.City == "" {
		return nil
	}
	key := venueKey(m.CountryCode, m.City)
	id := eventcore.NewID(eventcore.ProviderVNL, key)
	if _, ok := b.venueSeen[key]; !ok {
		b.venueSeen[key] = struct{}{}
		b.out.venues = append(b.out.venues, eventcore.Venue{
			ID:         id,
			Ext:        eventcore.ExternalID{Provider: eventcore.ProviderVNL, Key: key},
			City:       m.City,
			CountryISO: m.CountryCode,
		})
	}
	return &id
}

func (b *scheduleBuilder) emitParticipants(m apiMatch) []eventcore.FixtureParticipant {
	gender := mapGender(m.Gender)
	parts := make([]eventcore.FixtureParticipant, 0, sidesPerMatch)
	for _, side := range []struct {
		no   int
		role string
	}{
		{m.TeamANo, "home"},
		{m.TeamBNo, "away"},
	} {
		pid := eventcore.NewID(eventcore.ProviderVNL, teamKey(side.no))
		if _, ok := b.partSeen[side.no]; !ok {
			b.partSeen[side.no] = struct{}{}
			t := b.teams[side.no]
			b.out.participants = append(b.out.participants, eventcore.Participant{
				ID: pid,
				Ext: eventcore.ExternalID{
					Provider: eventcore.ProviderVNL,
					Key:      teamKey(side.no),
				},
				Kind:       eventcore.ParticipantTeam,
				Name:       b.teamName(side.no),
				Code:       t.Code,
				CountryISO: t.Code,
				Gender:     gender,
			})
		}
		parts = append(parts, eventcore.FixtureParticipant{ParticipantID: pid, Role: side.role})
	}
	return parts
}

func (b *scheduleBuilder) teamName(no int) string {
	t, ok := b.teams[no]
	if !ok {
		return "Team " + strconv.Itoa(no)
	}
	if t.Name != "" {
		return t.Name
	}
	if t.TranslatedName != "" {
		return t.TranslatedName
	}
	if t.Code != "" {
		return t.Code
	}
	return "Team " + strconv.Itoa(no)
}
