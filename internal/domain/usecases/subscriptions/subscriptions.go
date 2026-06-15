package subscriptions

import (
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Service manages notification subscriptions and resolves which users should be
// @mentioned for a given fixture.
type Service struct {
	repo Repository
}

func New(repo Repository) Service { return Service{repo: repo} }

func (s Service) Add(sub eventcore.Subscription) error {
	return s.repo.AddSubscription(sub)
}

func (s Service) Remove(sub eventcore.Subscription) error {
	return s.repo.RemoveSubscription(sub)
}

func (s Service) List(
	guildID string,
) ([]eventcore.Subscription, error) {
	return s.repo.ListByGuild(guildID)
}

func (s Service) ListByUser(
	guildID, userID string,
) ([]eventcore.Subscription, error) {
	return s.repo.ListByGuildUser(guildID, userID)
}

// fixtureFacts describes the matchable attributes of a fixture.
type fixtureFacts struct {
	countryCodes   []string
	disciplineCode string
}

// MentionsFor returns the deduped user IDs whose subscriptions in the given
// guild match the given country codes and discipline. Unknown kinds are ignored.
func (s Service) MentionsFor(
	guildID string,
	countryCodes []string,
	disciplineCode string,
) ([]string, error) {
	subs, err := s.repo.ListByGuild(guildID)
	if err != nil {
		return nil, err
	}

	facts := fixtureFacts{countryCodes: countryCodes, disciplineCode: disciplineCode}
	seen := make(map[string]struct{})
	users := make([]string, 0, len(subs))
	for _, sub := range subs {
		if !sub.Kind.Valid() {
			continue
		}
		if !matches(sub, facts) {
			continue
		}
		if _, ok := seen[sub.UserID]; ok {
			continue
		}
		seen[sub.UserID] = struct{}{}
		users = append(users, sub.UserID)
	}
	return users, nil
}

func matches(sub eventcore.Subscription, facts fixtureFacts) bool {
	switch sub.Kind {
	case eventcore.SubscribeAllResults:
		return true
	case eventcore.SubscribeCountry:
		for _, code := range facts.countryCodes {
			if strings.EqualFold(code, sub.Value) {
				return true
			}
		}
		return false
	case eventcore.SubscribeDiscipline:
		return strings.EqualFold(facts.disciplineCode, sub.Value)
	}
	return false
}
