package bootstrapper

import (
	"context"
	_ "embed"
	"net/url"
	"path"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/common-fate/clio"
	"github.com/common-fate/cloudform/deployer"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
)

//go:embed cloudformation/bootstrap.json
var BootstrapTemplate string

const BootstrapStackName = "common-fate-bootstrap"

type BootstrapStackOutput struct {
	AssetsBucket string `json:"AssetsBucket"`
}

type Bootstrapper struct {
	cfnClient *cloudformation.Client
	s3Client  *s3.Client
	deployer  *deployer.Deployer
	cfg       aws.Config
}

func NewFromConfig(cfg aws.Config) *Bootstrapper {
	deploy := deployer.NewFromConfig(cfg)

	return &Bootstrapper{
		cfnClient: cloudformation.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
		deployer:  deploy,
		cfg:       cfg,
	}
}

var ErrNotDeployed error = errors.New("bootstrap stack has not yet been deployed in this account and region")

func (b *Bootstrapper) Detect(ctx context.Context, retryOnStackNotExist bool) (*BootstrapStackOutput, error) {
	r := retry.NewFibonacci(time.Second)

	// don't retry unless specified
	if retryOnStackNotExist {
		r = retry.WithMaxDuration(time.Second*20, r)
	} else {
		r = retry.WithMaxRetries(1, r)
	}

	var stack *types.Stack
	err := retry.Do(ctx, r, func(ctx context.Context) (err error) {
		stacks, err := b.cfnClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: aws.String(BootstrapStackName),
		})
		var genericError *smithy.GenericAPIError
		if ok := errors.As(err, &genericError); ok && genericError.Code == "ValidationError" {
			return retry.RetryableError(ErrNotDeployed)
		} else if err != nil {
			return err
		} else if len(stacks.Stacks) != 1 {
			return fmt.Errorf("expected 1 stack but got %d", len(stacks.Stacks))
		}
		stack = &stacks.Stacks[0]
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := ProcessOutputs(*stack)
	return &out, nil
}

func ProcessOutputs(stack types.Stack) BootstrapStackOutput {
	// decode the output variables into the Go struct.
	var out BootstrapStackOutput

	for _, o := range stack.Outputs {
		if *o.OutputKey == "AssetsBucket" {
			out.AssetsBucket = *o.OutputValue
		}
	}

	return out
}

type DeployOpts struct {
	// Tags to associate with the stack
	Tags map[string]string

	// RoleARN is an optional deployment role to use
	RoleARN string

	// Confirm will skip interactive confirmations
	// if set to tru
	Confirm bool
}

// Deployment returns the deployment options
// to use with the deployer package to create the bootstrap stack.
//
// Usage:
//
//	d := deployer.NewFromConfig(cfg)
//	deployment := bootstrapper.Deployment
//	d.Deploy(ctx, deployment)
var Deployment = deployer.DeployOpts{
	Template:  BootstrapTemplate,
	StackName: BootstrapStackName,
}

// GetOrDeployBootstrap loads the output if the stack already exists, else it deploys the bootstrap stack first
func (b *Bootstrapper) GetOrDeployBootstrapBucket(ctx context.Context, confirm bool) (string, error) {
	out, err := b.Detect(ctx, false)
	if err == ErrNotDeployed {
		deployment := Deployment
		deployment.Confirm = confirm
		_, err := b.deployer.Deploy(ctx, deployment)
		if err != nil {
			return "", err
		}
		out, err := b.Detect(ctx, true)
		if err != nil {
			return "", err
		}
		return out.AssetsBucket, nil
	}
	if err != nil {
		return "", err
	}
	return out.AssetsBucket, nil
}

type ProviderFiles struct {
	CloudformationTemplateURL string
}

// CopyProviderFiles will clone the handler and cfn template from the registry bucket to the bootstrap bucket of the current account
func (b *Bootstrapper) CopyProviderFiles(ctx context.Context, provider providerregistrysdk.ProviderDetail) (*ProviderFiles, error) {
	// detect the bootstrap bucket
	out, err := b.Detect(ctx, false)
	if err != nil {
		return nil, err
	}

	lambdaAssetPath := path.Join(provider.Publisher, provider.Name, provider.Version)
	clio.Debugf("Copying the handler.zip into %s", path.Join(out.AssetsBucket, lambdaAssetPath, "handler.zip"))
	_, err = b.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(out.AssetsBucket),
		Key:        aws.String(path.Join(lambdaAssetPath, "handler.zip")),
		CopySource: aws.String(url.QueryEscape(provider.LambdaAssetS3Arn)),
	})
	if err != nil {
		return nil, err
	}
	clio.Debugf("Successfully copied the handler.zip into %s", path.Join(out.AssetsBucket, lambdaAssetPath, "handler.zip"))

	clio.Debugf("Copying the CloudFormation template into %s", path.Join(out.AssetsBucket, lambdaAssetPath, "cloudformation.json"))
	_, err = b.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(out.AssetsBucket),
		Key:        aws.String(path.Join(lambdaAssetPath, "cloudformation.json")),
		CopySource: aws.String(url.QueryEscape(provider.CfnTemplateS3Arn)),
	})
	if err != nil {
		return nil, err
	}
	clio.Debugf("Successfully copied the CloudFormation template into %s", path.Join(out.AssetsBucket, lambdaAssetPath, "cloudformation.json"))

	return &ProviderFiles{
		CloudformationTemplateURL: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", out.AssetsBucket, b.cfg.Region, path.Join(lambdaAssetPath, "cloudformation.json")),
	}, nil
}
