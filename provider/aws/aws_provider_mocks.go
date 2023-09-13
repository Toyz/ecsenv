package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type MockECSService struct {
	TaskDefOutput *ecs.DescribeTaskDefinitionOutput
	Error         error
}

func (m *MockECSService) DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return m.TaskDefOutput, m.Error
}

type MockSecretsManagerService struct {
	SecretOutput *secretsmanager.GetSecretValueOutput
	Error        error
}

func (m *MockSecretsManagerService) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return m.SecretOutput, m.Error
}
