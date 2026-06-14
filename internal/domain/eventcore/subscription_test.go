package eventcore

import "testing"

func TestSubscriptionKind_Valid(t *testing.T) {
	tests := []struct {
		name string
		kind SubscriptionKind
		want bool
	}{
		{"all_results", SubscribeAllResults, true},
		{"country", SubscribeCountry, true},
		{"discipline", SubscribeDiscipline, true},
		{"blank", SubscriptionKind(""), false},
		{"unknown", SubscriptionKind("athlete"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
