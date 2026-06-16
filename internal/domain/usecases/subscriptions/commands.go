package subscriptions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// HandleCommand runs a subscription slash command and returns the ephemeral reply.
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
	case "countries":
		return s.handleCountries()
	default:
		return "", fmt.Errorf("unknown action %q", action)
	}
}

func (s Service) handleCountries() (string, error) {
	if s.countries == nil {
		return "The country list is unavailable.", nil
	}
	all, err := s.countries.ListCountries()
	if err != nil {
		return "", fmt.Errorf("list countries: %w", err)
	}
	codes := make([]string, 0, len(all))
	for _, c := range all {
		if code := countryCode(c); code != "" {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)
	return "Available countries (subscribe by name or 3-letter code):\n" +
		strings.Join(codes, ", "), nil
}

func (s Service) handleAdd(guildID, userID, kind, value string) (string, error) {
	// Adding requires a real country so we never store something unmatchable.
	sub, err := s.buildSubscription(guildID, userID, kind, value, true)
	if err != nil {
		return "", err
	}
	if err = s.Add(sub); err != nil {
		return "", fmt.Errorf("add subscription: %w", err)
	}
	return fmt.Sprintf("Subscribed you to %s.", s.describe(sub)), nil
}

func (s Service) handleRemove(guildID, userID, kind, value string) (string, error) {
	// Removing is lenient: it must clear a stored value even if that country is
	// no longer (or was never) in the catalog.
	sub, err := s.buildSubscription(guildID, userID, kind, value, false)
	if err != nil {
		return "", err
	}
	if err = s.Remove(sub); err != nil {
		return "", fmt.Errorf("remove subscription: %w", err)
	}
	return fmt.Sprintf("Unsubscribed you from %s.", s.describe(sub)), nil
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
		fmt.Fprintf(&b, "- %s\n", s.describe(sub))
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// buildSubscription validates and canonicalizes the value (country to code,
// discipline upper-cased). With strictCountry an unrecognized country is
// rejected; otherwise the typed value is kept so a stale row can still be removed.
func (s Service) buildSubscription(
	guildID, userID, kind, value string, strictCountry bool,
) (eventcore.Subscription, error) {
	k := eventcore.SubscriptionKind(kind)
	if !k.Valid() {
		return eventcore.Subscription{}, fmt.Errorf("invalid kind %q", kind)
	}
	value = strings.Join(strings.Fields(value), " ")
	if k != eventcore.SubscribeAllResults && value == "" {
		return eventcore.Subscription{}, fmt.Errorf("kind %q requires a value", kind)
	}

	switch k {
	case eventcore.SubscribeCountry:
		if country, ok := s.resolveCountry(value); ok {
			value = countryCode(country)
		} else if strictCountry {
			return eventcore.Subscription{}, fmt.Errorf(
				"%q is not a recognized country; use its name or 3-letter code (e.g. BRA)", value,
			)
		}
	case eventcore.SubscribeDiscipline:
		value = strings.ToUpper(value)
	}

	return eventcore.Subscription{GuildID: guildID, UserID: userID, Kind: k, Value: value}, nil
}

// describe renders a subscription for display, resolving a country code back to
// its flag and name.
func (s Service) describe(sub eventcore.Subscription) string {
	switch sub.Kind {
	case eventcore.SubscribeAllResults:
		return "all results"
	case eventcore.SubscribeCountry:
		if c, ok := s.resolveCountry(sub.Value); ok {
			if flag := c.EmojiFlag(); flag != "" {
				return "country " + flag + " " + c.Name
			}
			return "country " + c.Name
		}
		return "country " + sub.Value
	case eventcore.SubscribeDiscipline:
		return "discipline " + sub.Value
	}
	return string(sub.Kind)
}

// countryCode returns the code used to match participants: IOC first, then ISO3,
// then ISO2.
func countryCode(c eventcore.Country) string {
	switch {
	case c.IOCCode != "":
		return c.IOCCode
	case c.ISO3 != "":
		return c.ISO3
	default:
		return c.ISO2
	}
}
