package task

import (
	"testing"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNotificationSenderTaskMatchesKey(t *testing.T) {
	filtered := &NotificationSenderTask{keyPrefix: "yearly-stats"}

	assert.True(t, filtered.matchesKey(domain.NotificationKey("yearly-stats")))
	assert.True(t, filtered.matchesKey(domain.NotificationKey("yearly-stats:42:2026")))
	assert.False(t, filtered.matchesKey(domain.NotificationKey("monthly-stats:42:2026-12")))
	assert.False(t, filtered.matchesKey(domain.NotificationKey("yearly-stats-other:42")))

	unfiltered := &NotificationSenderTask{}
	assert.True(t, unfiltered.matchesKey(domain.NotificationKey("monthly-stats:42:2026-12")))
}
