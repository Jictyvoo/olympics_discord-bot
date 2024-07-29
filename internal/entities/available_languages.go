package entities

import "strings"

type Language struct {
	ID   string
	Name string
	Code string
}

var languageList = [...]Language{
	{ID: "en", Name: "English", Code: "ENG"},
	{ID: "fr", Name: "Français", Code: "FRA"},
	{ID: "de", Name: "Deutsch", Code: "DEU"},
	{ID: "it", Name: "Italiano", Code: "ITA"},
	{ID: "pt", Name: "Português", Code: "POR"},
	{ID: "es", Name: "Español", Code: "SPA"},
	{ID: "ja", Name: "日本語", Code: "JPN"},
	{ID: "ar", Name: "العربية", Code: "ENG"},
	{ID: "zh", Name: "中文", Code: "CHI"},
	{ID: "hi", Name: "हिन्दी", Code: "ENG"},
	{ID: "ko", Name: "한국어", Code: "KOR"},
	{ID: "ru", Name: "Русский", Code: "RUS"},
}

func GetLanguage(idOrCode string) Language {
	for _, lang := range languageList {
		switch strings.ToLower(idOrCode) {
		case lang.ID, lang.Code:
			return lang
		}
	}

	return languageList[0]
}
