package subscriptions

//go:generate go tool -modfile=../../../../tools/go.mod mockgen -source=interfaces.go -destination=interfaces_mock_test.go -package=subscriptions

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Repository persists and queries notification subscriptions. The context is
// bound at resolution time (remy.GetWithContext), so these methods take none.
type Repository interface {
	AddSubscription(sub eventcore.Subscription) error
	RemoveSubscription(sub eventcore.Subscription) error
	ListByGuild(guildID string) ([]eventcore.Subscription, error)
	ListByGuildUser(guildID, userID string) ([]eventcore.Subscription, error)
	ListAll() ([]eventcore.Subscription, error)
}
