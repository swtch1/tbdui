package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/swtch1/tbdui/logger"
)

// DB is a DynamoDB instance.
type DB struct {
	dynDB       *dynamo.DB
	Environment string
	logger      *logger.UILogger
}

// NewDB instantiates a new Dynamo DB.
func NewDB(awsCredentials *credentials.Credentials, awsRegion, environment string) (*DB, error) {
	sess, err := session.NewSession(&aws.Config{Credentials: awsCredentials})
	if err != nil {
		return &DB{}, fmt.Errorf("error creating new AWS session: %w", err)
	}

	db := dynamo.New(
		sess,
		&aws.Config{Region: aws.String(awsRegion)},
	)
	return &DB{
		dynDB:       db,
		Environment: environment,
		logger:      logger.NewUILogger(), // start with an empty logger so it can be enables selectively
	}, nil
}
