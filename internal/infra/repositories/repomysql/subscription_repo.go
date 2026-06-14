package repomysql

import (
	"database/sql"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type SubscriptionRepo struct{ *repoMySQL }

func NewSubscriptionRepo(base *repoMySQL) SubscriptionRepo { return SubscriptionRepo{base} }

func (r SubscriptionRepo) AddSubscription(sub eventcore.Subscription) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().AddSubscription(qctx, dbgen.AddSubscriptionParams{
		GuildID: sub.GuildID,
		UserID:  sub.UserID,
		Kind:    string(sub.Kind),
		Value:   sql.NullString{String: sub.Value, Valid: true},
	})
}

func (r SubscriptionRepo) RemoveSubscription(
	sub eventcore.Subscription,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().RemoveSubscription(qctx, dbgen.RemoveSubscriptionParams{
		GuildID: sub.GuildID,
		UserID:  sub.UserID,
		Kind:    string(sub.Kind),
		Value:   sql.NullString{String: sub.Value, Valid: true},
	})
}

func (r SubscriptionRepo) ListByGuildUser(
	guildID, userID string,
) ([]eventcore.Subscription, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().
		ListSubscriptionsByGuildUser(qctx, dbgen.ListSubscriptionsByGuildUserParams{
			GuildID: guildID,
			UserID:  userID,
		})
	if err != nil {
		return nil, err
	}
	return rowsToSubscriptions(rows), nil
}

func (r SubscriptionRepo) ListByGuild(
	guildID string,
) ([]eventcore.Subscription, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListSubscriptionsByGuild(qctx, guildID)
	if err != nil {
		return nil, err
	}
	return rowsToSubscriptions(rows), nil
}

func (r SubscriptionRepo) ListAll() ([]eventcore.Subscription, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListSubscriptions(qctx)
	if err != nil {
		return nil, err
	}
	return rowsToSubscriptions(rows), nil
}

func rowToSubscription(row dbgen.Subscription) eventcore.Subscription {
	return eventcore.Subscription{
		ID:      row.ID,
		GuildID: row.GuildID,
		UserID:  row.UserID,
		Kind:    eventcore.SubscriptionKind(row.Kind),
		Value:   row.Value.String,
	}
}

func rowsToSubscriptions(rows []dbgen.Subscription) []eventcore.Subscription {
	out := make([]eventcore.Subscription, len(rows))
	for i, row := range rows {
		out[i] = rowToSubscription(row)
	}
	return out
}
