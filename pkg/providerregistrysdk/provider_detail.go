package providerregistrysdk

// Base returns the basic details about a provider
// and converts a ProviderDetail object to a Provider object.
func (p ProviderDetail) Base() Provider {
	return Provider{
		Name:      p.Name,
		Publisher: p.Publisher,
		Version:   p.Version,
	}
}
