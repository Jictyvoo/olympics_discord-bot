package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type VenueRepo struct{ *repoMySQL }

func NewVenueRepo(base *repoMySQL) VenueRepo { return VenueRepo{base} }

// GetVenueByFixture resolves the venue hosting a fixture, or sql.ErrNoRows when
// the fixture has no venue.
func (r VenueRepo) GetVenueByFixture(
	fixtureID eventcore.CanonicalID,
) (eventcore.Venue, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetVenueByFixture(qctx, fixtureID.Bytes())
	if err != nil {
		return eventcore.Venue{}, err
	}
	return eventcore.Venue{
		ID:         mapper.IDFromBytes(row.ID),
		Ext:        eventcore.ExternalID{Provider: row.ProviderID, Key: row.ExternalKey},
		Name:       row.Name,
		City:       row.City.String,
		CountryISO: row.CountryIso.String,
	}, nil
}

func (r VenueRepo) UpsertVenue(v eventcore.Venue) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertVenue(qctx, dbgen.UpsertVenueParams{
		ID:          v.ID.Bytes(),
		ProviderID:  v.Ext.Provider,
		ExternalKey: v.Ext.Key,
		Name:        v.Name,
		City:        mapper.NSStr(v.City),
		CountryIso:  mapper.NSStr(v.CountryISO),
	})
}
