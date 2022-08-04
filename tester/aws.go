package tester

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// use this email for testing with SES
const VerifiedEmail = "no-reply@hellodarwin.io"
const localstackImage = "localstack/localstack:0.12.5"

type amazon struct {
	SES *ses.SES
	S3  *s3.S3
}

var awsInstance = &amazon{}
var awsOnce sync.Once

func AWS() *amazon {
	awsOnce.Do(initAws)

	return awsInstance
}

func initAws() {
	endpoint := mustStartLocalStack()

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("foo", "bar", "bar"),
		Region:           aws.String(endpoints.UsEast1RegionID),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})

	awsInstance.SES = ses.New(sess)
	awsInstance.S3 = s3.New(sess)

	// localstack SES requires email verification for the sender before we can send emails
	_, err = awsInstance.SES.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{EmailAddress: aws.String(VerifiedEmail)})

	if err != nil {
		panic(fmt.Errorf("failed to verify identity: %w", err))
	}
}

func mustStartLocalStack() string {
	// find the current directory so we can infer the migrations folder
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filename)
	if err != nil {
		logger.Panic(zap.Error(err))
	}

	ctx := context.Background()

	waiter := wait.ForLog("make_bucket: hellodarwin-dev-assets")

	hostSource := filepath.Join(dir, "../../../docker/docker-entrypoint-initaws.d")
	containerDestination := "/docker-entrypoint-initaws.d"

	req := testcontainers.ContainerRequest{
		Image:        localstackImage,
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   waiter,
		Env: map[string]string{
			"SERVICES": "ses,s3",
		},
		BindMounts: map[string]string{
			hostSource: containerDestination,
		},
		ReaperImage: testcontainersReaper,
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(fmt.Errorf("can't create localstack container: %w", err))
	}

	port, err := c.MappedPort(ctx, "4566/tcp")

	if err != nil {
		panic(fmt.Errorf("failed to get port for localstack: %w", err))
	}

	var host string
	// when running in Gitlab CI, the database host must be "docker" because we start the container with DinD
	_, exists := os.LookupEnv("CI")
	if exists {
		host = "docker"
	} else {
		host = "localhost"
	}

	return "http://" + host + ":" + port.Port()
}
