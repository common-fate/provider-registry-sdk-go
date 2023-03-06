package handlerclient

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type Local struct {
	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Sterr stream to write to.
	// If unset, os.Stderr will be used.
	Stderr io.Writer

	// Env vars to provide to the local process.
	// If Env is nil, the new process uses the current process's environment.
	Env []string
}

func (l Local) Execute(ctx context.Context, request msg.Request) (*msg.Result, error) {
	stderr := l.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	payload := payload{Type: request.Type(), Data: request}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(".venv/bin/commonfate-provider-py", "run", string(payloadbytes))
	cmd.Env = l.Env
	cmd.Dir = l.Dir
	cmd.Stderr = stderr
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
