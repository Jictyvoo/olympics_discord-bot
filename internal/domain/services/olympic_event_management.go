package services

import (
	"bytes"
	"fmt"
	"log/slog"
	"text/template"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

type (
	NotifierFacade interface {
		InitMessageChannel(channelName string) error
		SendMessage(content string) error
	}
)

type OlympicEventManager struct {
	notifier    NotifierFacade
	msgTemplate *template.Template
}

func discRelativeHour(timestamp time.Time) string {
	return fmt.Sprintf("<t:%d:R>", timestamp.Unix())
}

func NewOlympicEventManager(facade NotifierFacade) (EventObserver, error) {
	const tmpl = `
🏅 **Event:** {{.EventName}}
**Discipline:** {{.DisciplineName}}
**Phase:** {{.Phase}}
**Gender:** {{.Gender}}
**Starts at:** {{discRelativeHour .StartAt}}
**Ends at:** {{discRelativeHour .EndAt}}
**Status:** {{.Status}}
**Competitors:**
{{range .Competitors}}- :{{emojiFlag .Country}}: {{.Name}}, Age: {{.Age}}
{{end}}`

	t, err := template.New("event").Funcs(
		template.FuncMap{
			"emojiFlag":        entities.CountryInfo.EmojiFlag,
			"discRelativeHour": discRelativeHour,
		},
	).Parse(tmpl)
	if err != nil {
		return OlympicEventManager{}, err
	}
	return OlympicEventManager{notifier: facade, msgTemplate: t}, nil
}

func (oen OlympicEventManager) OnEvent(event entities.OlympicEvent) {
	// Create the needed message structure to be sent
	var buf bytes.Buffer
	err := oen.msgTemplate.Execute(&buf, event)
	if err != nil {
		slog.Error("Error executing template for Olympic event", slog.String("error", err.Error()))
		return
	}

	content := buf.String()
	err = oen.notifier.SendMessage(content)
	if err != nil {
		slog.Error(
			"Error sending message using notifier",
			slog.String("message", content),
			slog.String("error", err.Error()),
		)
	}
}
