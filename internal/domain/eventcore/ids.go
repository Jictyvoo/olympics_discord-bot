package eventcore

import "github.com/jictyvoo/olhojogo/pkg/idgen"

// ProviderID identifies the data source (e.g. "olympics", "fifa").
type ProviderID = idgen.ProviderID

// CanonicalID is a 16-byte deterministic identifier derived from (ProviderID, ExternalKey).
type CanonicalID = idgen.CanonicalID

// ExternalID pairs a provider with its opaque key for a record.
type ExternalID = idgen.ExternalID

// Known provider codes.
const (
	ProviderOlympics ProviderID = "olympics"
	ProviderFIFA     ProviderID = "fifa"
	ProviderVNL      ProviderID = "vnl"
)

func NewID(provider ProviderID, key string) CanonicalID {
	return idgen.From(provider, key)
}
