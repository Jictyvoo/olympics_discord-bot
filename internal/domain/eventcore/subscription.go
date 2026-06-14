package eventcore

type SubscriptionKind string

const (
	SubscribeAllResults SubscriptionKind = "all_results"
	SubscribeCountry    SubscriptionKind = "country"
	SubscribeDiscipline SubscriptionKind = "discipline"
)

func (k SubscriptionKind) Valid() bool {
	switch k {
	case SubscribeAllResults, SubscribeCountry, SubscribeDiscipline:
		return true
	}
	return false
}

// Subscription is a Discord user's request to be @mentioned about a country, a
// discipline, or all results, scoped to a guild.
type Subscription struct {
	ID      int64
	GuildID string
	UserID  string // discord user id of the subscriber
	Kind    SubscriptionKind
	Value   string // country code (IOC/ISO) or discipline code; empty for all_results
}
