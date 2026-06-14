package eventcore

type Venue struct {
	ID         CanonicalID
	Ext        ExternalID
	Name       string
	City       string
	CountryISO string
}
