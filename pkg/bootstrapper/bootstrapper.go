package bootstrapper

import (
	"context"
	_ "embed"
	"net/url"
	"os"
	"path"
	"time"

	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/common-fate/clio"
	"github.com/common-fate/cloudform/deployer"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"golang.org/x/term"
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

// DetectOpts allows the detection of the bootstrap bucket to be customised.
type DetectOpts struct {
	// StackName is the name of the stack to query
	// By default it is set to 'common-fate-bootstrap'
	StackName string

	// Retry to use when querying for the stack if it doesn't exist.
	Retry retry.Backoff
}

// WithRetry uses a fibonacci backoff capped at 20 seconds to retry for the stack output.
var WithRetry = func(do *DetectOpts) {
	do.Retry = retry.WithMaxDuration(time.Second*20, retry.NewFibonacci(time.Second))
}

// Detect an existing bootstrap stack.
func (b *Bootstrapper) Detect(ctx context.Context, opts ...func(*DetectOpts)) (*BootstrapStackOutput, error) {
	o := DetectOpts{
		StackName: BootstrapStackName,
		Retry:     retry.WithMaxRetries(1, retry.NewConstant(time.Second)),
	}

	for _, opt := range opts {
		opt(&o)
	}

	var stack *types.Stack
	err := retry.Do(ctx, o.Retry, func(ctx context.Context) (err error) {
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

// Deployment returns the deployment options with the template filled in
// and the stack name set to the 'common-fate-bootstrap' stack name, ready
// to use with the deployer package to create the bootstrap stack.
//
// Usage:
//
//	d := deployer.NewFromConfig(cfg)
//	deployment := bootstrapper.Deployment()
//	d.Deploy(ctx, deployment)
func Deployment(opts ...deployer.DeployOptFunc) deployer.DeployOpts {
	d := deployer.DeployOpts{
		Template:  BootstrapTemplate,
		StackName: BootstrapStackName,
	}
	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// GetOrDeployBootstrap loads the output if the stack already exists, else it deploys the bootstrap stack first
func (b *Bootstrapper) GetOrDeployBootstrapBucket(ctx context.Context, opts ...deployer.DeployOptFunc) (*BootstrapStackOutput, error) {
	bootstrapStackOutput, err := b.Detect(ctx)
	if err == nil {
		return bootstrapStackOutput, nil
	}

	if err != ErrNotDeployed {
		// some other error which we can't handle
		return nil, err
	}

	deployment := Deployment(opts...)

	// if we get here, we need to deploy the bootstrap stack into the particular AWS account and region.

	// get the current AWS account and region to display it in the info message
	stsClient := sts.NewFromConfig(b.cfg)
	ci, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, errors.Wrap(err, "getting caller identity")
	}

	clio.Debug("the bootstrap stack was not detected")
	clio.Warnf("To get started deploying providers, you need to bootstrap this AWS account and region (%s:%s)", *ci.Account, b.cfg.Region)
	clio.Infof("Bootstrapping will deploy a CloudFormation stack called '%s' which creates an S3 Bucket.\nProvider assets will be copied from the Common Fate Provider Registry into this bucket.\nThese assets can then be deployed into your account.", deployment.StackName)

	if !deployment.Confirm {
		// if the terminal is non-interactive (e.g. in CI/CD systems)
		// return with an error so that we don't cause a deployment to hang forever.
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return nil, errors.New("bootstrapping needs a confirmation but the terminal is non-interactive (you can try including a '--confirm' flag to resolve this)")
		}

		err = survey.AskOne(&survey.Confirm{Message: "Deploy bootstrap stack", Default: true}, &deployment.Confirm)
		if err != nil {
			return nil, err
		}

		if !deployment.Confirm {
			return nil, errors.New("cancelling deployment")
		}
	}

	_, err = b.deployer.Deploy(ctx, deployment)
	if err != nil {
		return nil, err
	}

	bootstrapStackOutput, err = b.Detect(ctx, WithRetry)
	if err != nil {
		return nil, err
	}

	return bootstrapStackOutput, nil
}

type ProviderFiles struct {
	CloudformationTemplateURL string
}

func AssetsExist(ctx context.Context, client *s3.Client, bucket string, key string) (bool, error) {
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {

		return false, err
	}
	return true, nil
}

// CopyProviderFiles will clone the handler and cfn template from the registry bucket to the bootstrap bucket of the current account
func (b *Bootstrapper) CopyProviderFiles(ctx context.Context, provider providerregistrysdk.ProviderDetail) (*ProviderFiles, error) {
	// detect the bootstrap bucket
	out, err := b.Detect(ctx)
	if err != nil {
		return nil, err
	}

	lambdaAssetPath := path.Join("registry.commonfate.io", "v1alpha1", "providers", provider.Publisher, provider.Name, provider.Version)

	//check if asset already exists and if force flag was not passed
	exists, err := AssetsExist(ctx, b.s3Client, out.AssetsBucket, path.Join(lambdaAssetPath, "handler.zip"))
	if err != nil {
		return nil, err
	}
	if !exists {
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
	} else {
		clio.Debugf("already exists, skipped copy of handler asset")

	}

	exists, err = AssetsExist(ctx, b.s3Client, out.AssetsBucket, path.Join(lambdaAssetPath, "cloudformation.json"))
	if err != nil {
		return nil, err
	}

	if !exists {
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

	} else {
		clio.Debugf("already exists, skipped copy of cloudformation asset")

	}

	pf := &ProviderFiles{
		CloudformationTemplateURL: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", out.AssetsBucket, b.cfg.Region, path.Join(lambdaAssetPath, "cloudformation.json")),
	}

	return pf, nil
}
