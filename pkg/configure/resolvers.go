package configure

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// Resolvers resolve config values.
// Resolvers should only return an error if something fatal has occurred
// and the process should exit.
type Resolver interface {
	Resolve(ctx context.Context, key string, c ConfigValue) (string, error)
}

type FillOpts struct {
	ConfigResolvers []Resolver
	SecretResolvers []Resolver
}

func Dev() FillOpts {
	return FillOpts{
		ConfigResolvers: []Resolver{EnvVarResolver{Prefix: "PROVIDER_CONFIG_"}},
		SecretResolvers: []Resolver{EnvVarResolver{Prefix: "PROVIDER_SECRET_"}},
	}
}

func (c *Config) Fill(ctx context.Context, opts FillOpts) error {
	for k, v := range c.Values {
		if v.Secret {
			for _, resolver := range opts.SecretResolvers {
				value, err := resolver.Resolve(ctx, k, v)
				if err != nil {
					return err
				}
				if value != "" {
					v.Ref = value
					break
				}
			}
		} else {
			for _, resolver := range opts.ConfigResolvers {
				value, err := resolver.Resolve(ctx, k, v)
				if err != nil {
					return err
				}
				if value != "" {
					v.Value = value
					break
				}
			}
		}
		c.Values[k] = v
	}
	return nil
}

func (cv Config) CfnParams() []types.Parameter {
	var params []types.Parameter
	for k, v := range cv.Values {
		paramName := pascalCase(k)
		val := v.Value
		if v.Secret {
			paramName += "Secret"
			val = v.Ref
		}

		params = append(params, types.Parameter{
			ParameterKey:   &paramName,
			ParameterValue: &val,
		})
	}
	return params
}

func (cv Config) ToCLIFlag(flag string) []string {
	var flags []string
	for k, v := range cv.Values {
		flag := fmt.Sprintf("%s %s=%s", flag, k, v.Value)
		flags = append(flags, flag)
	}
	return flags
}
