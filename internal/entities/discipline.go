package entities

type Discipline struct {
	ID           Identifier
	Code         string
	Name         string
	Description  string
	IsSport      bool
	IsParalympic bool
}
