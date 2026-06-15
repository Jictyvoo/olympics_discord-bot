package subscriptions

import (
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func TestService_HandleCommand_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)

	var got eventcore.Subscription
	repo.EXPECT().
		AddSubscription(gomock.Any()).
		Do(func(s eventcore.Subscription) { got = s }).
		Return(nil)

	reply, err := New(repo).HandleCommand("add", "g1", "u1", "country", "BRA")
	if err != nil {
		t.Fatalf("HandleCommand: %v", err)
	}
	if got.GuildID != "g1" || got.UserID != "u1" ||
		got.Kind != eventcore.SubscribeCountry || got.Value != "BRA" {
		t.Fatalf("built subscription mismatch: %+v", got)
	}
	if !strings.Contains(reply, "country BRA") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestService_HandleCommand_AddAllResultsNoValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	repo.EXPECT().AddSubscription(gomock.Any()).Return(nil)

	reply, err := New(repo).HandleCommand("add", "g1", "u1", "all_results", "")
	if err != nil {
		t.Fatalf("HandleCommand: %v", err)
	}
	if !strings.Contains(reply, "all results") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestService_HandleCommand_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)

	var got eventcore.Subscription
	repo.EXPECT().
		RemoveSubscription(gomock.Any()).
		Do(func(s eventcore.Subscription) { got = s }).
		Return(nil)

	reply, err := New(repo).HandleCommand("remove", "g1", "u1", "discipline", "SWM")
	if err != nil {
		t.Fatalf("HandleCommand: %v", err)
	}
	if got.GuildID != "g1" || got.UserID != "u1" ||
		got.Kind != eventcore.SubscribeDiscipline || got.Value != "SWM" {
		t.Fatalf("built subscription mismatch: %+v", got)
	}
	if !strings.Contains(reply, "discipline SWM") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestService_HandleCommand_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	repo.EXPECT().
		ListByGuildUser("g1", "u1").
		Return([]eventcore.Subscription{
			{UserID: "u1", Kind: eventcore.SubscribeAllResults},
			{UserID: "u1", Kind: eventcore.SubscribeCountry, Value: "BRA"},
		}, nil)

	reply, err := New(repo).HandleCommand("list", "g1", "u1", "", "")
	if err != nil {
		t.Fatalf("HandleCommand: %v", err)
	}
	if !strings.Contains(reply, "all results") || !strings.Contains(reply, "country BRA") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestService_HandleCommand_ListEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	repo.EXPECT().ListByGuildUser("g1", "u1").Return(nil, nil)

	reply, err := New(repo).HandleCommand("list", "g1", "u1", "", "")
	if err != nil {
		t.Fatalf("HandleCommand: %v", err)
	}
	if !strings.Contains(reply, "no subscriptions") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestService_HandleCommand_InvalidKind(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)

	_, err := New(repo).HandleCommand("add", "g1", "u1", "athlete", "X")
	if err == nil {
		t.Fatal("expected error for invalid kind")
	}
}

func TestService_HandleCommand_MissingValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)

	_, err := New(repo).HandleCommand("add", "g1", "u1", "country", "")
	if err == nil {
		t.Fatal("expected error for missing value")
	}
}

func TestService_HandleCommand_UnknownAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)

	_, err := New(repo).HandleCommand("frobnicate", "g1", "u1", "", "")
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestService_HandleCommand_AddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	repo.EXPECT().AddSubscription(gomock.Any()).Return(errors.New("db down"))

	_, err := New(repo).HandleCommand("add", "g1", "u1", "all_results", "")
	if err == nil {
		t.Fatal("expected error propagated from repo")
	}
}
