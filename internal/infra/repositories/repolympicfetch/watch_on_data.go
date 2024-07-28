package repolympicfetch

type WatchOnResp struct {
	Title       string `json:"title"`
	CountryCode string `json:"countryCode"`
	MrhItems    []struct {
		Url                string `json:"url"`
		Label              string `json:"label"`
		Title              string `json:"title"`
		Image              string `json:"image"`
		Target             string `json:"target"`
		Slug               string `json:"slug"`
		EmbeddedWidgetHtml []any  `json:"embeddedWidgetHtml"`
		Tags               []struct {
			Slug string `json:"slug"`
		} `json:"tags"`
	} `json:"mrhItems"`
	Errors []any `json:"errors"`
}
