package configure

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type ConfigValue struct {
	// Secret is true if the config value is a secret
	Secret bool

	// Description of the config variable
	Description string

	// Value is the actual value of the config
	Value string

	// Ref is the path to the actual secret (only used if the config value is secret)
	Ref string
}

type Config struct {
	Values map[string]ConfigValue
}

// ConfigFromSchema initialises a Config map with keys from the Provider schema.
func ConfigFromSchema(schema *map[string]providerregistrysdk.Config) Config {
	cfg := Config{
		Values: map[string]ConfigValue{},
	}

	if schema == nil {
		return cfg
	}

	for k, configSchema := range *schema {
		cv := ConfigValue{}

		if configSchema.Description != nil {
			cv.Description = *configSchema.Description
		}

		if configSchema.Secret != nil {
			cv.Secret = *configSchema.Secret
		}

		cfg.Values[k] = cv
	}
	return cfg
}

type EnvVarResolver struct {
	Prefix string
}

func (r EnvVarResolver) Resolve(ctx context.Context, key string, c ConfigValue) (string, error) {
	envVar := r.Prefix + strings.ToUpper(key)
	return os.Getenv(envVar), nil
}

type MapResolver struct {
	kv map[string]string
}

func (r MapResolver) Resolve(ctx context.Context, key string, c ConfigValue) (string, error) {
	val, ok := r.kv[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

type PromptResolver struct {
	Stdin  terminal.FileReader
	Stdout terminal.FileWriter
	Stderr io.Writer
}

func (r PromptResolver) Resolve(ctx context.Context, key string, c ConfigValue) (string, error) {
	var val string
	err := survey.AskOne(&survey.Input{Message: key + ":", Help: c.Description}, &val, survey.WithStdio(r.Stdin, r.Stdout, r.Stderr))
	if err != nil {
		return "", err
	}

	return val, nil
}

func pascalCase(s string) string {
	arg := strings.Split(s, "_")
	var formattedStr []string

	for _, v := range arg {
		formattedStr = append(formattedStr, strings.ToUpper(v[0:1])+v[1:])
	}

	return strings.Join(formattedStr, "")
}
