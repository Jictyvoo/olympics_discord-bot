package render

import (
	"testing"
	"time"
)

func TestDisciplineIcon(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{name: "known discipline", code: athleticsCode, want: athleticsIcon},
		{name: "gymnastics artistic", code: "GAR", want: gymnasticsIcon},
		{name: "gymnastics trampoline", code: "GTR", want: gymnasticsIcon},
		{name: "gymnastics rhythmic", code: "GRY", want: gymnasticsIcon},
		{name: "unknown code", code: "ZZZ", want: ""},
		{name: "empty code", code: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisciplineIcon(tt.code); got != tt.want {
				t.Fatalf("DisciplineIcon(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestDisciplineIconByName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "english", in: "football", want: soccerIcon},
		{name: "portuguese", in: futebol, want: soccerIcon},
		{name: "spanish", in: "Fútbol", want: soccerIcon},
		{name: "german", in: "Fußball", want: soccerIcon},
		{name: "unknown", in: "Cricket", want: ""},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisciplineIconByName(tt.in); got != tt.want {
				t.Fatalf("DisciplineIconByName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDiscordTimestamp(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want string
	}{
		{
			name: "epoch",
			in:   time.Unix(0, 0),
			want: "<t:0:R>",
		},
		{
			name: "fixed time",
			in:   time.Date(2024, 7, 26, 20, 0, 0, 0, time.UTC),
			want: "<t:1722024000:R>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscordTimestamp(tt.in); got != tt.want {
				t.Fatalf("DiscordTimestamp(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
