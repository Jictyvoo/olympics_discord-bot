package repolympicfetch

type DisciplineResp struct {
	Code         string `json:"code"`
	Slug         string `json:"slug"`
	IsParalympic bool   `json:"isParalympic"`
	Description  string `json:"description"`
	IsSport      bool   `json:"isSport"`
	Order        string `json:"order"`
}
