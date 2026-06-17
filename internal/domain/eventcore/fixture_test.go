package eventcore

import (
	"testing"
	"time"
)

func TestFixtureStatus_Valid(t *testing.T) {
	testCases := []struct {
		name   string
		status FixtureStatus
		want   bool
	}{
		{"scheduled", FixtureScheduled, true},
		{"live", FixtureLive, true},
		{"finished", FixtureFinished, true},
		{"cancelled", FixtureCancelled, true},
		{"postponed", FixturePostponed, true},
		{"empty", "", false},
		{unknownProvider, unknownProvider, false},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if got := tCase.status.Valid(); got != tCase.want {
				t.Fatalf("Valid() = %v, want %v", got, tCase.want)
			}
		})
	}
}

func TestCompleteByEndTime(t *testing.T) {
	endsAt := time.Date(2026, 6, 18, 23, 45, 0, 0, time.UTC)
	grace := time.Hour
	beforeGrace := endsAt.Add(30 * time.Minute)
	afterGrace := endsAt.Add(90 * time.Minute)

	testCases := []struct {
		name string
		in   FixtureStatus
		now  time.Time
		want FixtureStatus
	}{
		{"live within grace stays live", FixtureLive, beforeGrace, FixtureLive},
		{"live after grace finishes", FixtureLive, afterGrace, FixtureFinished},
		{"scheduled after grace finishes", FixtureScheduled, afterGrace, FixtureFinished},
		{"finished untouched", FixtureFinished, afterGrace, FixtureFinished},
		{"cancelled untouched", FixtureCancelled, afterGrace, FixtureCancelled},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if got := CompleteByEndTime(tCase.in, endsAt, tCase.now, grace); got != tCase.want {
				t.Fatalf("CompleteByEndTime() = %q, want %q", got, tCase.want)
			}
		})
	}
}

func TestFixture_ComputeChecksum_Stable(t *testing.T) {
	base := Fixture{
		Ext:      ExternalID{Provider: "olympics", Key: "EV-001"},
		Name:     "100m Final",
		StartsAt: time.Date(2024, 8, 4, 20, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2024, 8, 4, 21, 0, 0, 0, time.UTC),
		Status:   FixtureScheduled,
	}
	testCases := []struct {
		name string
		a, b Fixture
		same bool
	}{
		{"identical fixtures produce same checksum", base, base, true},
		{
			"different status differs",
			base,
			func() Fixture { f := base; f.Status = FixtureFinished; return f }(),
			false,
		},
		{
			"different name differs",
			base,
			func() Fixture { f := base; f.Name = "200m Final"; return f }(),
			false,
		},
		{"UTC normalisation", base, func() Fixture {
			f := base
			f.StartsAt = base.StartsAt.In(time.FixedZone("UTC+3", 3*3600))
			return f
		}(), true},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			ca, cb := tCase.a.ComputeChecksum(), tCase.b.ComputeChecksum()
			if (ca == cb) != tCase.same {
				t.Fatalf("checksum match=%v, want %v (a=%q b=%q)", ca == cb, tCase.same, ca, cb)
			}
		})
	}
}

func TestFixture_ComputeChecksumWith_Results(t *testing.T) {
	fixture := Fixture{
		Ext:      ExternalID{Provider: "olympics", Key: "EV-001"},
		Name:     "100m Final",
		StartsAt: time.Date(2024, 8, 4, 20, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2024, 8, 4, 21, 0, 0, 0, time.UTC),
		Status:   FixtureFinished,
	}
	pidA := NewID("olympics", "ATH-A")
	pidB := NewID("olympics", "ATH-B")
	gold := []Result{
		{ParticipantID: pidA, Score: "9.81", Outcome: OutcomeMedalGold},
		{ParticipantID: pidB, Score: "9.90", Outcome: OutcomeMedalSilver},
	}
	silverSwapped := []Result{
		{ParticipantID: pidA, Score: "9.81", Outcome: OutcomeMedalSilver},
		{ParticipantID: pidB, Score: "9.90", Outcome: OutcomeMedalGold},
	}
	reordered := []Result{gold[1], gold[0]}

	t.Run("results fold into the checksum", func(t *testing.T) {
		if fixture.ComputeChecksum() == fixture.ComputeChecksumWith(gold) {
			t.Fatal("expected results to change the checksum vs the fixture-only checksum")
		}
	})
	t.Run("a medal change re-flips the checksum", func(t *testing.T) {
		if fixture.ComputeChecksumWith(gold) == fixture.ComputeChecksumWith(silverSwapped) {
			t.Fatal("expected a medal/outcome change to produce a different checksum")
		}
	})
	t.Run("result ordering is stable", func(t *testing.T) {
		if fixture.ComputeChecksumWith(gold) != fixture.ComputeChecksumWith(reordered) {
			t.Fatal("expected checksum to be independent of result ordering")
		}
	})
}

func TestCountry_EmojiFlag_IsThis(t *testing.T) {
	br := Country{ISO2: "BR", ISO3: bra, IOCCode: bra, Name: brazil}
	if got := br.EmojiFlag(); got != flagBR {
		t.Fatalf("EmojiFlag() = %q, want %s", got, flagBR)
	}
	if got := (Country{}).EmojiFlag(); got != "" {
		t.Fatalf("EmojiFlag() on empty country = %q, want empty", got)
	}
	for _, v := range []string{"brazil", bra, "br", "bra"} {
		if !br.IsThis(v) {
			t.Fatalf("IsThis(%q) = false, want true", v)
		}
	}
	for _, v := range []string{"", "argentina", "usa"} {
		if br.IsThis(v) {
			t.Fatalf("IsThis(%q) = true, want false", v)
		}
	}
}
