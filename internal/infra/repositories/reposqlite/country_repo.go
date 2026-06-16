package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type CountryRepo struct{ *repoSQLite }

func NewCountryRepo(base *repoSQLite) CountryRepo { return CountryRepo{base} }

func (r CountryRepo) SeedCountry(c eventcore.Country) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	var codeNum, pop any
	if c.CodeNum != 0 {
		codeNum = int64(c.CodeNum)
	}
	if c.Population != 0 {
		pop = c.Population
	}
	var area, gdp any
	if c.AreaKm2 != 0 {
		area = c.AreaKm2
	}
	if c.GDPUSD != 0 {
		gdp = c.GDPUSD
	}
	return r.Queries().UpsertCountry(qctx, dbgen.UpsertCountryParams{
		Iso2:       c.ISO2,
		Iso3:       c.ISO3,
		IocCode:    mapper.OptString(c.IOCCode),
		Name:       c.Name,
		CodeNum:    codeNum,
		Population: pop,
		AreaKm2:    area,
		GdpUsd:     gdp,
	})
}

func (r CountryRepo) ByISO2(iso2 string) (eventcore.Country, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetCountryByISO2(qctx, iso2)
	if err != nil {
		return eventcore.Country{}, err
	}
	return rowToCountry(row), nil
}

func (r CountryRepo) ByISO3(iso3 string) (eventcore.Country, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetCountryByISO3(qctx, iso3)
	if err != nil {
		return eventcore.Country{}, err
	}
	return rowToCountry(row), nil
}

func (r CountryRepo) ByIOC(ioc string) (eventcore.Country, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetCountryByIOC(qctx, mapper.OptString(ioc))
	if err != nil {
		return eventcore.Country{}, err
	}
	return rowToCountry(row), nil
}

func (r CountryRepo) ListCountries() ([]eventcore.Country, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListCountries(qctx)
	if err != nil {
		return nil, err
	}
	out := make([]eventcore.Country, 0, len(rows))
	for _, row := range rows {
		out = append(out, rowToCountry(row))
	}
	return out, nil
}

func rowToCountry(row dbgen.Country) eventcore.Country {
	codeNum := 0
	if p := mapper.NullInt(row.CodeNum); p != nil {
		codeNum = *p
	}
	pop := int64(0)
	if p := mapper.NullInt(row.Population); p != nil {
		pop = int64(*p)
	}
	return eventcore.Country{
		ISO2:       row.Iso2,
		ISO3:       row.Iso3,
		IOCCode:    mapper.NullStr(row.IocCode),
		Name:       row.Name,
		CodeNum:    codeNum,
		Population: pop,
		AreaKm2:    mapper.FloatOrZero(row.AreaKm2),
		GDPUSD:     mapper.FloatOrZero(row.GdpUsd),
	}
}
