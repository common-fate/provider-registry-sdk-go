package handlerclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type Lambda struct {

	/*
		The name of the Lambda function, version, or alias. Name formats

		* Function name – my-function (name-only), my-function:v1 (with alias).

		* Function ARN – arn:aws:lambda:us-west-2:123456789012:function:my-function.

		* Partial ARN – 123456789012:function:my-function.

		You can append a version number or alias to any of the formats. The length constraint applies only to the full ARN. If you specify only the function name, it is limited to 64 characters in length.
	*/
	FunctionName string // Help text is copied from the AWS Lambda Invoke help
	lambdaClient *lambda.Client
}

/*
functionName: The name of the Lambda function, version, or alias. Name formats

* Function name – my-function (name-only), my-function:v1 (with alias).

* Function ARN – arn:aws:lambda:us-west-2:123456789012:function:my-function.

* Partial ARN – 123456789012:function:my-function.

You can append a version number or alias to any of the formats. The length constraint applies only to the full ARN. If you specify only the function name, it is limited to 64 characters in length.
*/
func NewLambdaRuntime(ctx context.Context, functionName string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	l := Lambda{FunctionName: functionName, lambdaClient: lambdaClient}
	return &Client{Executor: l}, nil
}

// NewLambdaRuntimeFromConfig creates a new handler client from a
// provided AWS config.
func NewLambdaRuntimeFromConfig(cfg aws.Config, functionName string) *Client {
	lambdaClient := lambda.NewFromConfig(cfg)

	l := Lambda{FunctionName: functionName, lambdaClient: lambdaClient}
	return &Client{Executor: l}
}

// payload is the actual request JSON sent to the Lambda function.
type payload struct {
	Type msg.RequestType `json:"type"`
	Data any             `json:"data"`
}

func (l Lambda) Execute(ctx context.Context, request msg.Request) (*msg.Result, error) {
	payload := payload{Type: request.Type(), Data: request}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	res, err := l.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(l.FunctionName),
		InvocationType: lambdatypes.InvocationTypeRequestResponse,
		Payload:        payloadbytes,
		LogType:        lambdatypes.LogTypeTail,
	})
	if err != nil {
		return nil, err
	}

	if res.FunctionError != nil {
		var logs string
		if res.LogResult != nil {
			logbyte, err := base64.URLEncoding.DecodeString(*res.LogResult)
			if err != nil {
				logger.Get(ctx).Errorw("error decoding lambda log", "error", err)
			}
			logs = string(logbyte)
		}
		return nil, fmt.Errorf("lambda execution error: %s: %s", *res.FunctionError, logs)
	}

	var result msg.Result
	err = json.Unmarshal(res.Payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
