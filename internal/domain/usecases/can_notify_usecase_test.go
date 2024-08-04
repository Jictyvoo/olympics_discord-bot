package usecases

import (
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
	mockRepoExpectations func(mockRepo *mocks.MockCanNotifyRepository, eventShaID string)
	expectedEventKey     func(event entities.OlympicEvent) string
	expectedError        error
}

func (tt TestCaseShouldNotify) Run(t testing.TB, mockCtrl *gomock.Controller) {
	mockRepo := mocks.NewMockCanNotifyRepository(mockCtrl)
	tt.mockRepoExpectations(mockRepo, tt.event.SHAIdentifier())
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

	// Define test cases
	tests := []TestCaseShouldNotify{
		{
			name: "Send notification without history",
			event: entities.OlympicEvent{
				ID:          eventID,
				Status:      entities.StatusScheduled,
				StartAt:     staticReturnTime.Add(-10 * time.Minute),
				EndAt:       staticReturnTime.Add(30 * time.Minute),
				SessionCode: "215#__6N8S",
			},
			allowedTimeDiff:  allowedTimeDiff,
			staticReturnTime: staticReturnTime,
			mockRepoExpectations: func(mockRepo *mocks.MockCanNotifyRepository, eventShaID string) {
				mockRepo.EXPECT().
					CheckSentNotifications(eventID, gomock.Eq(eventShaID)).
					Return(entities.Notification{}, nil)
				mockRepo.EXPECT().RegisterNotification(gomock.Any()).Do(
					func(notification entities.Notification) {
						assert.Equal(t, eventID, notification.EventID)
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
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.Run(t, ctrl)
			},
		)
	}
}
