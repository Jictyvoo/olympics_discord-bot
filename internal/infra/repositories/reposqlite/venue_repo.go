package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type VenueRepo struct{ *repoSQLite }

func NewVenueRepo(base *repoSQLite) VenueRepo { return VenueRepo{base} }

func (r VenueRepo) UpsertVenue(v eventcore.Venue) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertVenue(qctx, dbgen.UpsertVenueParams{
		ID:          v.ID.Bytes(),
		ProviderID:  v.Ext.Provider,
		ExternalKey: v.Ext.Key,
		Name:        v.Name,
		City:        mapper.OptString(v.City),
		CountryIso:  mapper.OptString(v.CountryISO),
	})
}
