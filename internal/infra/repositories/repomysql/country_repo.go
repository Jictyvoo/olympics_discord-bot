package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type CountryRepo struct{ *repoMySQL }

func NewCountryRepo(base *repoMySQL) CountryRepo { return CountryRepo{base} }

func (r CountryRepo) SeedCountry(c eventcore.Country) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertCountry(qctx, dbgen.UpsertCountryParams{
		Iso2:       c.ISO2,
		Iso3:       c.ISO3,
		IocCode:    mapper.NSStr(c.IOCCode),
		Name:       c.Name,
		CodeNum:    mapper.NSInt(int64(c.CodeNum)),
		Population: mapper.NSInt(c.Population),
		AreaKm2:    mapper.NSFloat(c.AreaKm2),
		GdpUsd:     mapper.NSFloat(c.GDPUSD),
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
	row, err := r.Queries().GetCountryByIOC(qctx, mapper.NSStr(ioc))
	if err != nil {
		return eventcore.Country{}, err
	}
	return rowToCountry(row), nil
}

func rowToCountry(row dbgen.Country) eventcore.Country {
	return eventcore.Country{
		ISO2:       row.Iso2,
		ISO3:       row.Iso3,
		IOCCode:    row.IocCode.String,
		Name:       row.Name,
		CodeNum:    int(row.CodeNum.Int64),
		Population: row.Population.Int64,
		AreaKm2:    row.AreaKm2.Float64,
		GDPUSD:     row.GdpUsd.Float64,
	}
}
