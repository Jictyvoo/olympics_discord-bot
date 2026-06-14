package eventcore

import "testing"

func TestNotificationStatus_Valid(t *testing.T) {
	testCases := []struct {
		name   string
		status NotificationStatus
		want   bool
	}{
		{"pending", NotificationPending, true},
		{"sent", NotificationSent, true},
		{"failed", NotificationFailed, true},
		{"cancelled", NotificationCancelled, true},
		{"skipped", NotificationSkipped, true},
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

func TestParseNotificationStatus_Defaults(t *testing.T) {
	testCases := []struct {
		input string
		want  NotificationStatus
	}{
		{"sent", NotificationSent},
		{"failed", NotificationFailed},
		{"", NotificationPending},
		{"bogus", NotificationPending},
	}
	for _, tCase := range testCases {
		t.Run(tCase.input, func(t *testing.T) {
			if got := ParseNotificationStatus(tCase.input); got != tCase.want {
				t.Fatalf("got %q want %q", got, tCase.want)
			}
		})
	}
}
