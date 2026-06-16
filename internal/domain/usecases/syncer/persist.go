package syncer

import "github.com/jictyvoo/olhojogo/internal/domain/eventcore"

// persist writes the delta in one transaction, ordered so parents exist before
// their children (FK-safe).
func (s *Syncer) persist(delta eventcore.SyncDelta) error {
	tx, err := s.repo.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err = persistEach(delta.Competitions, tx.UpsertCompetition); err != nil {
		return err
	}
	if err = persistEach(delta.Seasons, tx.UpsertSeason); err != nil {
		return err
	}
	if err = persistEach(delta.Stages, tx.UpsertStage); err != nil {
		return err
	}
	if err = persistEach(delta.Groups, tx.UpsertGroup); err != nil {
		return err
	}
	if err = persistEach(delta.Venues, tx.UpsertVenue); err != nil {
		return err
	}
	if err = persistEach(delta.Participants, tx.UpsertParticipant); err != nil {
		return err
	}
	if err = persistFixtures(tx, delta.Fixtures); err != nil {
		return err
	}
	if err = persistEach(delta.Results, tx.UpsertResult); err != nil {
		return err
	}
	if err = persistEach(delta.Standings, tx.UpsertStanding); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	// Emit only durably-committed fixtures so observers react to persisted state.
	if s.events != nil {
		for _, f := range delta.Fixtures {
			s.events.Emit(f)
		}
	}
	return nil
}

// persistFixtures upserts each fixture together with its participant links.
func persistFixtures(tx Tx, fixtures []eventcore.Fixture) error {
	for _, f := range fixtures {
		if err := tx.UpsertFixture(f); err != nil {
			return err
		}
		if err := tx.UpsertFixtureParticipants(f.ID, f.Participants); err != nil {
			return err
		}
	}
	return nil
}

func persistEach[T any](items []T, upsert func(T) error) error {
	for _, item := range items {
		if err := upsert(item); err != nil {
			return err
		}
	}
	return nil
}
