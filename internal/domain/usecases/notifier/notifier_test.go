package usecases

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/mocks"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func TestCanNotifyUseCase_timeDiffAllowed(t *testing.T) {
	fixedNow := time.Now()

	tests := []struct {
		name            string
		event           entities.OlympicEvent
		allowedTimeDiff time.Duration
		expected        bool
	}{
		{
			name: "Event within allowed time diff",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-5 * time.Hour),
				EndAt:   fixedNow.Add(5 * time.Hour),
			},
			allowedTimeDiff: 10 * time.Hour,
			expected:        true,
		},
		{
			name: "Event exactly on the boundary",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-5 * time.Hour),
				EndAt:   fixedNow.Add(5 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        true,
		},
		{
			name: "Event outside allowed time diff",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-6 * time.Hour),
				EndAt:   fixedNow.Add(6 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
		{
			name: "Event in the past",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-12 * time.Hour),
				EndAt:   fixedNow.Add(-6 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
		{
			name: "Event in the future",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(6 * time.Hour),
				EndAt:   fixedNow.Add(12 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				uc := CanNotifyUseCase{
					allowedTimeDiff: tt.allowedTimeDiff, timeNow: func() time.Time {
						return fixedNow
					},
				}
				result := uc.timeDiffAllowed(tt.event)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

type TestCaseShouldNotify struct {
	name                 string
	event                entities.OlympicEvent
	allowedTimeDiff      time.Duration
	staticReturnTime     time.Time
	mockRepoExpectations func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent)
	expectedEventKey     func(event entities.OlympicEvent) string
	expectedError        error
}

func (tt TestCaseShouldNotify) Run(t testing.TB, mockCtrl *gomock.Controller) {
	mockRepo := mocks.NewMockCanNotifyRepository(mockCtrl)
	tt.mockRepoExpectations(mockRepo, tt.event)
	if tt.expectedEventKey == nil {
		tt.expectedEventKey = func(event entities.OlympicEvent) string {
			event.Normalize()
			return event.SHAIdentifier()
		}
	}
	uc := CanNotifyUseCase{
		repo:            mockRepo,
		allowedTimeDiff: tt.allowedTimeDiff,
		timeNow: func() time.Time {
			return tt.staticReturnTime
		},
	}

	validatedKey, err := uc.ShouldNotify(tt.event)
	expectedKey := tt.expectedEventKey(tt.event)

	assert.Equal(t, expectedKey, validatedKey)
	assert.Equal(t, tt.expectedError, err)
}

func TestCanNotifyUseCase_ShouldNotify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	const allowedTimeDiff = time.Hour
	const eventID entities.Identifier = 5185923
	staticReturnTime := time.Date(2024, time.August, 4, 10, 0, 0, 0, time.UTC)

	mountNotifyPendingCases := func(name string, notifyStatus string) TestCaseShouldNotify {
		return TestCaseShouldNotify{
			name: name,
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusScheduled,
				StartAt:     staticReturnTime.Add(-10 * time.Minute),
				EndAt:       staticReturnTime.Add(3 * time.Hour >> 1),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Blue-Eyes White Dragon"}, {Code: "Red-Eyes Black Dragon"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				eventShaID := event.SHAIdentifier()
				returnNotification := entities.Notification{}
				if notifyStatus != "" {
					returnNotification = entities.Notification{
						ID:            4259,
						EventID:       eventID,
						Status:        entities.NotificationStatus(notifyStatus),
						EventChecksum: eventShaID,
						NotifiedAt:    time.Time{},
					}
				}

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(returnNotification, nil)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Do(
					func(notification entities.Notification) {
						assert.Equal(t, eventID, notification.EventID)
						assert.Equal(t, eventShaID, notification.EventChecksum)
						assert.Equal(t, entities.NotificationStatusPending, notification.Status)
						if !notification.NotifiedAt.IsZero() {
							t.Errorf("NotifiedAt should be zero")
						}
					},
				).Return(nil)
			},
		}
	}
	mountPreventNotifyOnStatus := func(notifyStatus string) TestCaseShouldNotify {
		return TestCaseShouldNotify{
			name: fmt.Sprintf("Prevent notification with %s status on history", notifyStatus),
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusScheduled,
				StartAt:     staticReturnTime.Add(-10 * time.Minute),
				EndAt:       staticReturnTime.Add(3 * time.Hour >> 1),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Blue-Eyes White Dragon"}, {Code: "Red-Eyes Black Dragon"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				eventShaID := event.SHAIdentifier()

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(
						entities.Notification{
							ID:            9521,
							EventID:       eventID,
							Status:        entities.NotificationStatus(notifyStatus),
							EventChecksum: eventShaID,
							NotifiedAt:    time.Time{},
						}, nil,
					)
			},
			expectedEventKey: func(_ entities.OlympicEvent) string {
				return ""
			},
		}
	}

	// Define test cases
	tests := [...]TestCaseShouldNotify{
		mountNotifyPendingCases("Send notification without history", ""),
		mountNotifyPendingCases(
			"Send notification with history as pending", string(entities.NotificationStatusPending),
		),
		mountNotifyPendingCases(
			"Send notification without history", string(entities.NotificationStatusFailed),
		),
		mountPreventNotifyOnStatus(string(entities.NotificationStatusSent)),
		mountPreventNotifyOnStatus(string(entities.NotificationStatusSkipped)),
		mountPreventNotifyOnStatus(string(entities.NotificationStatusCancelled)),
		{
			name: "Ongoing (with partial results) notification with history without results",
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusOngoing,
				StartAt:     staticReturnTime.Add(-10 * time.Minute),
				EndAt:       staticReturnTime.Add(3 * time.Hour >> 1),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Kaiba"}, {Code: "Yugi"}, {Code: "Pegasus"},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"Kaiba": {Mark: "1560LP"},
					"Yugi":  {Mark: "2200LP"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				event.ResultPerCompetitor = map[string]entities.Results{}
				event.Competitors = slices.Clone(event.Competitors)
				event.Normalize()

				eventShaID := event.SHAIdentifier()

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(
						entities.Notification{
							ID:            4259,
							EventID:       eventID,
							Status:        entities.NotificationStatusSent,
							EventChecksum: eventShaID,
							NotifiedAt:    staticReturnTime.Add(time.Hour >> 1),
						}, nil,
					)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Times(0)
			},
			expectedEventKey: func(event entities.OlympicEvent) string {
				return ""
			},
		},
		{
			name: "Event finished with partial results and a notification with history without results",
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusFinished,
				StartAt:     staticReturnTime.Add(-10 * time.Minute),
				EndAt:       staticReturnTime.Add(3 * time.Hour >> 1),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Kaiba"}, {Code: "Yugi"}, {Code: "Pegasus"},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"Kaiba": {Mark: "1560LP"},
					"Yugi":  {Mark: "2200LP"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				var eventShaID struct{ old, new string }
				event.Normalize()
				eventShaID.new = event.SHAIdentifier()

				event.ResultPerCompetitor = map[string]entities.Results{}
				event.Competitors = slices.Clone(event.Competitors)
				event.Normalize()

				eventShaID.old = event.SHAIdentifier()

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.All(gomock.Eq(eventShaID.new), gomock.Not(gomock.Eq(eventShaID.old)))).
					Return(
						entities.Notification{}, nil,
					)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Do(
					func(notification entities.Notification) {
						assert.Equal(t, eventID, notification.EventID)
						assert.Equal(t, eventShaID.new, notification.EventChecksum)
						assert.Equal(t, entities.NotificationStatusPending, notification.Status)
						if !notification.NotifiedAt.IsZero() {
							t.Errorf("NotifiedAt should be zero")
						}
					},
				).Return(nil)
			},
			expectedEventKey: func(event entities.OlympicEvent) string {
				event.Normalize()
				return event.SHAIdentifier()
			},
		},
		{
			name: "Event 2x allowedTimeDiff after now is marked as pending",
			event: entities.OlympicEvent{
				ID:      eventID,
				Status:  entities.StatusScheduled,
				StartAt: staticReturnTime.Add(allowedTimeDiff + allowedTimeDiff>>1),
				EndAt: staticReturnTime.Add(
					allowedTimeDiff + (allowedTimeDiff >> 1) + (3 * time.Hour >> 1),
				),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Kaiba"}, {Code: "Yugi"}, {Code: "Pegasus"},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"Kaiba": {Mark: "1560LP"},
					"Yugi":  {Mark: "2200LP"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				event.ResultPerCompetitor = map[string]entities.Results{}
				event.Competitors = slices.Clone(event.Competitors)
				event.Normalize()

				eventShaID := event.SHAIdentifier()

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(
						entities.Notification{
							EventID:       eventID,
							EventChecksum: eventShaID,
							Status:        entities.NotificationStatusPending,
						}, nil,
					)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Do(
					func(notification entities.Notification) {
						assert.Equal(t, eventID, notification.EventID)
						assert.Equal(t, eventShaID, notification.EventChecksum)
						assert.Equal(t, entities.NotificationStatusPending, notification.Status)
						if !notification.NotifiedAt.IsZero() {
							t.Errorf("NotifiedAt should be zero")
						}
					},
				).Return(nil)
			},
			expectedEventKey: func(event entities.OlympicEvent) string {
				return ""
			},
		},
		{
			name: "Event ends before now with a diff of value allowedTimeDiff. Now gets cancelled",
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusFinished,
				StartAt:     staticReturnTime.Add(-allowedTimeDiff << 1),
				EndAt:       staticReturnTime.Add(-allowedTimeDiff),
				SessionCode: "215#__6N8S",
				Competitors: []entities.OlympicCompetitors{
					{Code: "Kaiba"}, {Code: "Yugi"}, {Code: "Pegasus"},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"Kaiba": {Mark: "1560LP"},
					"Yugi":  {Mark: "2200LP"},
				},
			},
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, event entities.OlympicEvent) {
				event.Competitors = slices.Clone(event.Competitors)
				event.Normalize()
				eventShaID := event.SHAIdentifier()

				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(
						entities.Notification{
							EventID:       eventID,
							EventChecksum: eventShaID,
							Status:        entities.NotificationStatusPending,
						}, nil,
					)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Do(
					func(notification entities.Notification) {
						assert.Equal(t, eventID, notification.EventID)
						assert.Equal(t, eventShaID, notification.EventChecksum)
						assert.Equal(t, entities.NotificationStatusCancelled, notification.Status)
						if !notification.NotifiedAt.IsZero() {
							t.Errorf("NotifiedAt should be zero")
						}
					},
				).Return(nil)
			},
			expectedEventKey: func(event entities.OlympicEvent) string {
				return ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.allowedTimeDiff = allowedTimeDiff
				tt.staticReturnTime = staticReturnTime

				tt.Run(t, ctrl)
			},
		)
	}
}
