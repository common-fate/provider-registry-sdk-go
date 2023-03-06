package handler

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type Local struct {
	Path string
}

func (l Local) Execute(ctx context.Context, request msg.Request) (*msg.Result, error) {
	payload := payload{Type: request.Type(), Data: request}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(".venv/bin/commonfate-provider-py", "run", string(payloadbytes))
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var res msg.Result
	err = json.Unmarshal(out, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
