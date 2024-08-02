package services

import (
	"bytes"
	"log/slog"
	"text/template"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

type (
	NotifierFacade interface {
		InitMessageChannel(channelName string) error
		SendMessage(content string) error
	}
)

type OlympicEventManager struct {
	notifier       NotifierFacade
	msgTemplate    *template.Template
	watchCountries []string
}

func NewOlympicEventManager(
	watchCountries []string,
	facade NotifierFacade,
) (OlympicEventManager, error) {
	const tmpl = `
# {{.Discipline}}
**Event:** {{.EventName}} - {{.Status}}
**Phase:** {{.Phase}}
**Gender:** {{.Gender}}
**Start:** {{discRelativeHour .StartAt}}
**End:** {{discRelativeHour .EndAt}}
**Competitors:**
{{range .Competitors}}- :{{emojiFlag .Country}}: {{.Name}}{{with $result := index $.ResultPerCompetitor .Code}}{{if $result.MedalType}} ({{$result.MedalType}}){{end}}{{end}}
{{end}}`

	t, err := template.New("event").Funcs(
		template.FuncMap{
			"emojiFlag":        entities.CountryInfo.EmojiFlag,
			"discRelativeHour": utils.DiscordTimestamp,
		},
	).Parse(tmpl)
	if err != nil {
		return OlympicEventManager{}, err
	}
	return OlympicEventManager{
		notifier:       facade,
		watchCountries: watchCountries,
		msgTemplate:    t,
	}, nil
}

func (oen OlympicEventManager) NormalizeEvent4Notification(event *entities.OlympicEvent) bool {
	if len(oen.watchCountries) <= 0 {
		return true
	}
	newCompetitorsList := event.Competitors
	if len(event.Competitors) > 0 {
		newCompetitorsList = make([]entities.OlympicCompetitors, 0, len(event.Competitors))
	}

	var foundCountry bool
competitorLoop:
	for _, competitor := range event.Competitors {
		for _, watch := range oen.watchCountries {
			if competitor.Country.IsThis(watch) {
				foundCountry = true
				newCompetitorsList = append(newCompetitorsList, competitor)
				continue competitorLoop
			}
		}
		if _, hasResult := event.ResultPerCompetitor[competitor.Code]; hasResult {
			newCompetitorsList = append(newCompetitorsList, competitor)
		}
	}

	if len(newCompetitorsList) < len(event.Competitors) {
		event.Competitors = newCompetitorsList
	}

	return foundCountry
}

func (oen OlympicEventManager) genContent(event entities.OlympicEvent) string {
	// Create the needed message structure to be sent
	var buf bytes.Buffer
	err := oen.msgTemplate.Execute(&buf, event)
	if err != nil {
		slog.Error("Error executing template for Olympic event", slog.String("error", err.Error()))
		return ""
	}

	return buf.String()
}

func (oen OlympicEventManager) OnEvent(event entities.OlympicEvent) {
	// Check if it should notify the event
	if !oen.NormalizeEvent4Notification(&event) {
		return
	}

	// Create the needed message structure to be sent
	content := oen.genContent(event)
	if content == "" {
		return
	}
	if err := oen.notifier.SendMessage(content); err != nil {
		slog.Error(
			"Error sending message using notifier",
			slog.String("message", content),
			slog.String("error", err.Error()),
		)
		return
	}

	slog.Info("Event successfully sent to notifier", slog.Any("event", event))
}
