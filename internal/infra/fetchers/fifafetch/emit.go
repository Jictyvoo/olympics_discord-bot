package fifafetch

import "github.com/jictyvoo/olhojogo/internal/domain/eventcore"

func (b *matchBuilder) emitCompetition(m apiMatch, lang string) eventcore.CanonicalID {
	id := eventcore.NewID(eventcore.ProviderFIFA, m.IdCompetition)
	if _, ok := b.compSeen[m.IdCompetition]; !ok {
		b.compSeen[m.IdCompetition] = struct{}{}
		b.out.competitions = append(b.out.competitions, eventcore.Competition{
			ID: id,
			Ext: eventcore.ExternalID{
				Provider: eventcore.ProviderFIFA,
				Key:      m.IdCompetition,
			},
			Code:       m.IdCompetition,
			Name:       localized(m.CompetitionName, lang),
			Discipline: footballDiscipline(lang),
		})
	}
	return id
}

func (b *matchBuilder) emitSeason(
	m apiMatch, compID eventcore.CanonicalID, lang string, season seasonMeta,
) eventcore.CanonicalID {
	id := eventcore.NewID(eventcore.ProviderFIFA, m.IdSeason)
	if _, ok := b.seasonSeen[m.IdSeason]; !ok {
		b.seasonSeen[m.IdSeason] = struct{}{}
		name := season.name
		if name == "" {
			name = localized(m.SeasonName, lang)
		}
		b.out.seasons = append(b.out.seasons, eventcore.Season{
			ID:            id,
			Ext:           eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: m.IdSeason},
			CompetitionID: compID,
			Name:          name,
			StartsOn:      season.startsOn,
			EndsOn:        season.endsOn,
		})
	}
	return id
}

func (b *matchBuilder) emitStage(
	m apiMatch, seasonID eventcore.CanonicalID, lang string,
) eventcore.CanonicalID {
	id := eventcore.NewID(eventcore.ProviderFIFA, m.IdStage)
	if _, ok := b.stageSeen[m.IdStage]; !ok {
		b.stageSeen[m.IdStage] = struct{}{}
		b.out.stageKeys = append(b.out.stageKeys, m.IdStage)
		b.out.stages = append(b.out.stages, eventcore.Stage{
			ID:       id,
			Ext:      eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: m.IdStage},
			SeasonID: seasonID,
			Name:     localized(m.StageName, lang),
		})
	}
	return id
}

func (b *matchBuilder) emitGroup(
	m apiMatch, stageID eventcore.CanonicalID, lang string,
) *eventcore.CanonicalID {
	if m.IdGroup == "" {
		return nil
	}
	id := eventcore.NewID(eventcore.ProviderFIFA, m.IdGroup)
	if _, ok := b.groupSeen[m.IdGroup]; !ok {
		b.groupSeen[m.IdGroup] = struct{}{}
		b.out.groups = append(b.out.groups, eventcore.Group{
			ID:      id,
			Ext:     eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: m.IdGroup},
			StageID: stageID,
			Name:    localized(m.GroupName, lang),
		})
	}
	return &id
}

func (b *matchBuilder) emitVenue(s apiStadium, lang string) *eventcore.CanonicalID {
	if s.IdStadium == "" {
		return nil
	}
	id := eventcore.NewID(eventcore.ProviderFIFA, s.IdStadium)
	if _, ok := b.venueSeen[s.IdStadium]; !ok {
		b.venueSeen[s.IdStadium] = struct{}{}
		b.out.venues = append(b.out.venues, eventcore.Venue{
			ID:         id,
			Ext:        eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: s.IdStadium},
			Name:       localized(s.Name, lang),
			City:       localized(s.CityName, lang),
			CountryISO: s.IdCountry,
		})
	}
	return &id
}

func (b *matchBuilder) emitParticipants(m apiMatch, lang string) []eventcore.FixtureParticipant {
	parts := make([]eventcore.FixtureParticipant, 0, sidesPerMatch)
	for _, side := range []struct {
		team apiTeam
		role string
	}{
		{m.Home, "home"},
		{m.Away, "away"},
	} {
		t := side.team
		if t.IdTeam == "" {
			continue
		}
		pid := eventcore.NewID(eventcore.ProviderFIFA, t.IdTeam)
		if _, ok := b.partSeen[t.IdTeam]; !ok {
			b.partSeen[t.IdTeam] = struct{}{}
			b.out.participants = append(b.out.participants, eventcore.Participant{
				ID:         pid,
				Ext:        eventcore.ExternalID{Provider: eventcore.ProviderFIFA, Key: t.IdTeam},
				Kind:       eventcore.ParticipantTeam,
				Name:       localized(t.TeamName, lang),
				Code:       t.Abbreviation,
				CountryISO: t.IdCountry,
				Gender:     mapGender(t.Gender),
			})
		}
		parts = append(parts, eventcore.FixtureParticipant{ParticipantID: pid, Role: side.role})
	}
	return parts
}
