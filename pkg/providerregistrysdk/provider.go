package providerregistrysdk

import "fmt"

// String prints the provider as a string,
// in the format 'publisher/name@version'
func (p Provider) String() string {
	return fmt.Sprintf("%s/%s@%s", p.Publisher, p.Name, p.Version)
}

// String prints the provider as a string,
// in the format 'publisher/name@version'
func (p ProviderDetail) String() string {
	return fmt.Sprintf("%s/%s@%s", p.Publisher, p.Name, p.Version)
}
