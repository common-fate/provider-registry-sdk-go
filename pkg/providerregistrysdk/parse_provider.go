package providerregistrysdk

import (
	"fmt"
	"regexp"
)

// ParseProvider parses a provider in string format.
// The expected format is publisher/name@version
func ParseProvider(input string) (Provider, error) {
	re, err := regexp.Compile(`([\w-]+)/([\w-]+)@(.+)`)
	if err != nil {
		return Provider{}, err
	}

	matches := re.FindStringSubmatch(input)
	if len(matches) != 4 {
		return Provider{}, fmt.Errorf("invalid provider format: %s", input)
	}

	p := Provider{
		Publisher: matches[1],
		Name:      matches[2],
		Version:   matches[3],
	}

	return p, nil
}
