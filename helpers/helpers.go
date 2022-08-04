package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"go.uber.org/zap"
	"net/url"
	"path/filepath"
)

// MustLoadConfig load configuration files into struct c or panics if it fails
func MustLoadConfig(c interface{}, filenames ...string) {
	loader := aconfig.LoaderFor(&c, aconfig.Config{
		Files: filenames,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
}

func MustLoadSecrets(c interface{}, rawurl string) {
	var err error

	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
		panic(fmt.Errorf("failed to parse secrets url: %w", err))
	}

	switch u.Scheme {
	case "file":
		MustLoadConfig(c, filepath.Join(u.Host, u.Path))
	case "secretsmanager":
		MustLoadSecretsManager(c, u.Host)
	default:
		panic(fmt.Errorf("unknown url scheme to retrieve secrets: %s", u.Scheme))
	}
}

// MustLoadSecretsManager loads the secrets stored in secrets manager into v
func MustLoadSecretsManager(v interface{}, id string) {
	var err error

	svc := secretsmanager.New(session.Must(session.NewSession()))

	var val *secretsmanager.GetSecretValueOutput
	if val, err = svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &id,
	}); err != nil {
		panic(fmt.Errorf("failed to retrieve secrets (%s): %w", id, err))
	}

	if err = json.Unmarshal([]byte(*val.SecretString), &v); err != nil {
		panic(fmt.Errorf("failed to marshal secrets (%s): %w", id, err))
	}
}

// MustAWSSession starts an aws session
func MustAWSSession(endpoint string, forcePathStyle bool, region string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(forcePathStyle),
		Region:           aws.String(region),
	}))
}

// NewLogger creates a new logger
func NewLogger(debug bool) *zap.SugaredLogger {
	var logger *zap.Logger

	// human readable logs outputs line by line instead of JSON
	if debug {
		logger, _ = zap.NewDevelopment(zap.AddStacktrace(zap.FatalLevel))
	} else {
		logger, _ = zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	}

	return logger.Sugar()
}
