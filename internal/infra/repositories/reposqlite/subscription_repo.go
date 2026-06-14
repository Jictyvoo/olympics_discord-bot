package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type SubscriptionRepo struct{ *repoSQLite }

func NewSubscriptionRepo(base *repoSQLite) SubscriptionRepo { return SubscriptionRepo{base} }

func (r SubscriptionRepo) AddSubscription(sub eventcore.Subscription) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().AddSubscription(qctx, dbgen.AddSubscriptionParams{
		GuildID: sub.GuildID,
		UserID:  sub.UserID,
		Kind:    string(sub.Kind),
		Value:   sub.Value,
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
		Value:   sub.Value,
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
	value, _ := row.Value.(string)
	return eventcore.Subscription{
		ID:      row.ID,
		GuildID: row.GuildID,
		UserID:  row.UserID,
		Kind:    eventcore.SubscriptionKind(row.Kind),
		Value:   value,
	}
}

func rowsToSubscriptions(rows []dbgen.Subscription) []eventcore.Subscription {
	out := make([]eventcore.Subscription, len(rows))
	for i, row := range rows {
		out[i] = rowToSubscription(row)
	}
	return out
}
