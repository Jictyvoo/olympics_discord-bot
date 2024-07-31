package entities

import "github.com/jictyvoo/olympics_data_fetcher/internal/utils"

type Discipline struct {
	ID           Identifier
	Code         string
	Name         string
	Description  string
	IsSport      bool
	IsParalympic bool
}

func (disc Discipline) String() string {
	return utils.DisciplineIcon(disc.Code) + " " + disc.Name
}
