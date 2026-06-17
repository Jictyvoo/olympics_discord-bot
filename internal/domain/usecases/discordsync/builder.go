package discordsync

import (
	"fmt"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/discordfacade"
)

func buildEventInput(f eventcore.Fixture, location string) discordfacade.ScheduledEventInput {
	desc := fmt.Sprintf("Provider: %s | ID: %s", f.Ext.Provider, f.Ext.Key)
	return discordfacade.ScheduledEventInput{
		Name:        f.Name,
		Description: desc,
		StartsAt:    f.StartsAt,
		EndsAt:      f.EndsAt,
		Location:    location,
	}
}
