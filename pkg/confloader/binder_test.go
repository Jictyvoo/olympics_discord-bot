package confloader

import (
	"testing"
	"time"
)

const defaultLiteral = "default"

func TestBindField_String(t *testing.T) {
	testCases := []struct {
		name    string
		envVal  string
		initial string
		want    string
	}{
		{"env set", "hello", defaultLiteral, "hello"},
		{"env empty keeps default", "", defaultLiteral, defaultLiteral},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Setenv("TEST_STR", tCase.envVal)
			got := tCase.initial
			if err := BindField(&got, "TEST_STR", ParseString).Bind(); err != nil {
				t.Fatal(err)
			}
			if got != tCase.want {
				t.Fatalf("got %q want %q", got, tCase.want)
			}
		})
	}
}

func TestBindField_Bool(t *testing.T) {
	testCases := []struct {
		name    string
		envVal  string
		initial bool
		want    bool
		wantErr bool
	}{
		{"true", "true", false, true, false},
		{"false", "false", true, false, false},
		{"empty keeps initial", "", true, true, false},
		{"invalid errors", "notabool", false, false, true},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Setenv("TEST_BOOL", tCase.envVal)
			got := tCase.initial
			err := BindField(&got, "TEST_BOOL", ParseBool).Bind()
			if (err != nil) != tCase.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tCase.wantErr)
			}
			if !tCase.wantErr && got != tCase.want {
				t.Fatalf("got %v want %v", got, tCase.want)
			}
		})
	}
}

func TestBindField_Duration(t *testing.T) {
	testCases := []struct {
		name   string
		envVal string
		want   time.Duration
	}{
		{"30s", "30s", 30 * time.Second},
		{"5m", "5m", 5 * time.Minute},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Setenv("TEST_DUR", tCase.envVal)
			var got time.Duration
			if err := BindField(&got, "TEST_DUR", ParseDuration).Bind(); err != nil {
				t.Fatal(err)
			}
			if got != tCase.want {
				t.Fatalf("got %v want %v", got, tCase.want)
			}
		})
	}
}

func TestBindEnv_AccumulatesErrors(t *testing.T) {
	t.Setenv("TEST_BOOL_A", "bad")
	t.Setenv("TEST_BOOL_B", "alsobad")
	var a, b bool
	err := BindEnv(
		BindField(&a, "TEST_BOOL_A", ParseBool),
		BindField(&b, "TEST_BOOL_B", ParseBool),
	)
	if err == nil {
		t.Fatal("expected error from two bad bools")
	}
}
