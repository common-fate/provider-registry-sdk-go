package handlerruntime

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/mitchellh/mapstructure"
)

type Local struct {
	Path string
}

func (l Local) FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error) {
	b, err := json.Marshal(contx)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	cmd := exec.Command("pdk", "test", "load-resources", "--name="+name, "--ctx="+string(b))
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return LoadResourceResponse{}, err
	}

	var lr LambdaResponse
	err = json.Unmarshal(out, &lr)
	if err != nil {
		return LoadResourceResponse{}, err
	}

	byt, err := json.Marshal(lr.Body)
	if err != nil {
		return LoadResourceResponse{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(byt, &data)
	if err != nil {
		return LoadResourceResponse{}, err
	}

	err = mapstructure.Decode(data, &resources)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	return
}

func (l Local) Describe(ctx context.Context) (*providerregistrysdk.DescribeResponse, error) {
	cmd := exec.Command("pdk", "test", "describe")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var lr LambdaResponse
	err = json.Unmarshal(out, &lr)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(lr.Body)
	if err != nil {
		return nil, err
	}

	var describe providerregistrysdk.DescribeResponse
	err = json.Unmarshal(b, &describe)
	if err != nil {
		return nil, err
	}

	return &describe, nil
}
func (l Local) Grant(ctx context.Context, subject string, target Target) (err error) {
	// @TODO this is untested/ not implemented in the local CLI
	cmd := exec.Command("pdk", "test", "grant")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	_, err = cmd.Output()
	return err

}
func (l Local) Revoke(ctx context.Context, subject string, target Target) (err error) {
	// @TODO this is untested/ not implemented in the local CLI
	cmd := exec.Command("pdk", "test", "revoke")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	_, err = cmd.Output()
	return err

}
