package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/valyala/fastjson"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskDefinitions(t *testing.T) {
	mockService := &MockECSService{
		TaskDefOutput: &ecs.DescribeTaskDefinitionOutput{
			TaskDefinition: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					{
						Name: aws.String("TestContainer"),
						Environment: []*ecs.KeyValuePair{
							{
								Name:  aws.String("TestKey"),
								Value: aws.String("TestValue"),
							},
						},
					},
				},
			},
		},
	}

	provider := AWSProvider{
		EcsService: mockService,
	}

	containerDefs, err := provider.GetTaskDefinitions("TestTaskDef")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(containerDefs))
	assert.Equal(t, "TestContainer", containerDefs[0].Name)
	assert.Equal(t, "TestValue", containerDefs[0].Environment["TestKey"])
}

func TestGetSecretValue(t *testing.T) {
	mockSecretsService := &MockSecretsManagerService{
		SecretOutput: &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`{"key": "value"}`),
		},
	}

	provider := AWSProvider{
		SecretsManagerService: mockSecretsService,
		SecretCache:           make(map[string]*fastjson.Value),
	}

	val, err := provider.GetSecretValue("arn:test:secret", "testSecret")
	assert.NoError(t, err)
	assert.Equal(t, "value", string(val.GetStringBytes("key")))
}
