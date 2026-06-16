package subscriptions

import (
	"sort"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const (
	codeBRA = "BRA"
	codeSWM = "SWM"
)

func sub(user string, kind eventcore.SubscriptionKind, value string) eventcore.Subscription {
	return eventcore.Subscription{
		GuildID: "g1",
		UserID:  user,
		Kind:    kind,
		Value:   value,
	}
}

type mentionsForCase struct {
	name       string
	subs       []eventcore.Subscription
	countries  []string
	discipline string
	want       []string
}

var mentionsForCases = []mentionsForCase{
	{
		name:       "all_results matches everything",
		subs:       []eventcore.Subscription{sub("u-all", eventcore.SubscribeAllResults, "")},
		countries:  []string{codeBRA},
		discipline: codeSWM,
		want:       []string{"u-all"},
	},
	{
		name:      "country case-insensitive",
		subs:      []eventcore.Subscription{sub("u-bra", eventcore.SubscribeCountry, "bra")},
		countries: []string{codeBRA},
		want:      []string{"u-bra"},
	},
	{
		name:      "country no match",
		subs:      []eventcore.Subscription{sub("u-bra", eventcore.SubscribeCountry, codeBRA)},
		countries: []string{"USA"},
		want:      []string{},
	},
	{
		name:       "discipline case-insensitive",
		subs:       []eventcore.Subscription{sub("u-swm", eventcore.SubscribeDiscipline, "swm")},
		discipline: codeSWM,
		want:       []string{"u-swm"},
	},
	{
		name: "dedupe across multiple subs on same user",
		subs: []eventcore.Subscription{
			sub("u-x", eventcore.SubscribeAllResults, ""),
			sub("u-x", eventcore.SubscribeCountry, codeBRA),
			sub("u-y", eventcore.SubscribeDiscipline, codeSWM),
		},
		countries:  []string{codeBRA},
		discipline: codeSWM,
		want:       []string{"u-x", "u-y"},
	},
	{
		name: "unknown kind ignored",
		subs: []eventcore.Subscription{
			sub("u-bad", eventcore.SubscriptionKind("athlete"), "X"),
			sub("u-good", eventcore.SubscribeCountry, codeBRA),
		},
		countries: []string{codeBRA},
		want:      []string{"u-good"},
	},
	{
		name:      "empty result",
		subs:      []eventcore.Subscription{sub("u-bra", eventcore.SubscribeCountry, codeBRA)},
		countries: []string{"USA"},
		want:      []string{},
	},
}

func TestService_MentionsFor(t *testing.T) {
	for _, tt := range mentionsForCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			repo.EXPECT().ListByGuild("g1").Return(tt.subs, nil)

			got, err := New(repo, nil).MentionsFor("g1", tt.countries, tt.discipline)
			if err != nil {
				t.Fatalf("MentionsFor: %v", err)
			}
			assertSameUsers(t, got, tt.want)
		})
	}
}

func assertSameUsers(t *testing.T, got, want []string) {
	t.Helper()
	sort.Strings(got)
	want = append([]string(nil), want...)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("users = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("users = %v, want %v", got, want)
		}
	}
}
