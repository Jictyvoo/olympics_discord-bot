package idgen

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

const canonicalIDByteLen = 16

// CanonicalID is a 16-byte identifier, either deterministic (via [From]) or a
// UUIDv7 (via [NewV7]). Both share one raw-bytes column type for a uniform schema.
type CanonicalID [16]byte

type ProviderID = string

var Zero CanonicalID

// From derives a CanonicalID from (provider, externalKey). Same inputs always
// yield the same output, so upserts are safe.
func From(provider ProviderID, externalKey string) CanonicalID {
	h := sha256.New()
	h.Write([]byte(provider))
	h.Write([]byte{0})
	h.Write([]byte(externalKey))
	sum := h.Sum(nil)
	var id CanonicalID
	copy(id[:], sum[:16])
	return id
}

// NewV7 mints a UUIDv7 CanonicalID for rows with no deterministic source
// (notifications, alerts, discord_events). Panics only if the system RNG fails.
func NewV7() CanonicalID {
	u, err := uuid.NewV7()
	if err != nil {
		panic("idgen: NewV7: " + err.Error())
	}
	return CanonicalID(u)
}

// FromBytes copies the first 16 bytes of b, returning Zero if b is shorter.
func FromBytes(b []byte) CanonicalID {
	if len(b) < canonicalIDByteLen {
		return Zero
	}
	var id CanonicalID
	copy(id[:], b[:canonicalIDByteLen])
	return id
}

// Bytes returns a slice aliasing the underlying array; callers must not mutate it.
func (id CanonicalID) Bytes() []byte {
	return id[:]
}

// String returns the hex-encoded ID for logs and debug output, never DB storage.
func (id CanonicalID) String() string {
	return hex.EncodeToString(id[:])
}

func (id CanonicalID) IsZero() bool {
	return id == Zero
}

type ExternalID struct {
	Provider ProviderID
	Key      string
}

func (e ExternalID) Canonical() CanonicalID {
	return From(e.Provider, e.Key)
}
