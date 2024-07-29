package repolympicfetch

type DisciplineResp struct {
	Slug         string `json:"slug"`
	IsParalympic bool   `json:"isParalympic"`
	Code         string `json:"id"`
	Description  string `json:"description"`
	IsSport      bool   `json:"isSport"`
	Order        string `json:"order"`
}
