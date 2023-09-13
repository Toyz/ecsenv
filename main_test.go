package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	aws2 "github.com/toyz/ecsenv/provider/aws"
	"github.com/valyala/fastjson"
	"testing"
)

func TestGenerateEnvContent(t *testing.T) {
	// Arrange
	taskDefName := "test-task-definition"
	secretArn := "arn:test:secret"

	ecsMock := &aws2.MockECSService{
		TaskDefOutput: &ecs.DescribeTaskDefinitionOutput{
			TaskDefinition: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					{
						Name: aws.String("container1"),
						Environment: []*ecs.KeyValuePair{
							{
								Name:  aws.String("ENV_VAR1"),
								Value: aws.String("value1"),
							},
						},
						Secrets: []*ecs.Secret{
							{
								Name:      aws.String("SECRET_VAR1"),
								ValueFrom: aws.String(secretArn),
							},
						},
					},
				},
			},
		},
	}

	secretsMock := &aws2.MockSecretsManagerService{
		SecretOutput: &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`{"SECRET_VAR1": "secret_value1"}`),
		},
	}

	provider := &aws2.AWSProvider{
		EcsService:            ecsMock,
		SecretsManagerService: secretsMock,
		SecretCache:           make(map[string]*fastjson.Value),
	}

	// Act
	envContent, err := generateEnvContent(taskDefName, provider)
	assert.NoError(t, err)

	// Assert
	expected := "ENV_VAR1=value1\nSECRET_VAR1=secret_value1\n"
	assert.Equal(t, expected, string(envContent))
}
