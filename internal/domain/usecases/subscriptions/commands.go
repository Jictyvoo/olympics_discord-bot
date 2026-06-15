package subscriptions

import (
	"fmt"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// HandleCommand executes a subscription slash command and returns a
// human-readable, ephemeral reply. It works on primitives so the domain package
// stays free of any Discord/infra dependency.
func (s Service) HandleCommand(
	action, guildID, userID, kind, value string,
) (string, error) {
	switch action {
	case "add":
		return s.handleAdd(guildID, userID, kind, value)
	case "remove":
		return s.handleRemove(guildID, userID, kind, value)
	case "list":
		return s.handleList(guildID, userID)
	default:
		return "", fmt.Errorf("unknown action %q", action)
	}
}

func (s Service) handleAdd(
	guildID, userID, kind, value string,
) (string, error) {
	sub, err := buildSubscription(guildID, userID, kind, value)
	if err != nil {
		return "", err
	}
	if err = s.Add(sub); err != nil {
		return "", fmt.Errorf("add subscription: %w", err)
	}
	return fmt.Sprintf("Subscribed you to %s.", describe(sub)), nil
}

func (s Service) handleRemove(
	guildID, userID, kind, value string,
) (string, error) {
	sub, err := buildSubscription(guildID, userID, kind, value)
	if err != nil {
		return "", err
	}
	if err = s.Remove(sub); err != nil {
		return "", fmt.Errorf("remove subscription: %w", err)
	}
	return fmt.Sprintf("Unsubscribed you from %s.", describe(sub)), nil
}

func (s Service) handleList(guildID, userID string) (string, error) {
	subs, err := s.ListByUser(guildID, userID)
	if err != nil {
		return "", fmt.Errorf("list subscriptions: %w", err)
	}
	if len(subs) == 0 {
		return "You have no subscriptions on this server.", nil
	}
	var b strings.Builder
	b.WriteString("Your subscriptions:\n")
	for _, sub := range subs {
		fmt.Fprintf(&b, "- %s\n", describe(sub))
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

func buildSubscription(
	guildID, userID, kind, value string,
) (eventcore.Subscription, error) {
	k := eventcore.SubscriptionKind(kind)
	if !k.Valid() {
		return eventcore.Subscription{}, fmt.Errorf("invalid kind %q", kind)
	}
	if k != eventcore.SubscribeAllResults && strings.TrimSpace(value) == "" {
		return eventcore.Subscription{}, fmt.Errorf("kind %q requires a value", kind)
	}
	return eventcore.Subscription{
		GuildID: guildID,
		UserID:  userID,
		Kind:    k,
		Value:   value,
	}, nil
}

func describe(sub eventcore.Subscription) string {
	switch sub.Kind {
	case eventcore.SubscribeAllResults:
		return "all results"
	case eventcore.SubscribeCountry:
		return "country " + sub.Value
	case eventcore.SubscribeDiscipline:
		return "discipline " + sub.Value
	}
	return string(sub.Kind)
}
