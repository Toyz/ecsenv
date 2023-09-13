package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/toyz/ecsenv/provider"
	"github.com/valyala/fastjson"
	"sync"
)

type ECSService interface {
	DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error)
}

type SecretsManagerService interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type AWSProvider struct {
	EcsService            ECSService
	SecretsManagerService SecretsManagerService
	SecretCache           map[string]*fastjson.Value // Updated cache type
	mu                    sync.RWMutex               // mutex for concurrent map access
}

func (a *AWSProvider) GetTaskDefinitions(taskDefinitionName string) ([]*provider.ContainerDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionName),
	}

	resp, err := a.EcsService.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	var containerDefinitions []*provider.ContainerDefinition
	for _, def := range resp.TaskDefinition.ContainerDefinitions {
		containerDef := &provider.ContainerDefinition{
			Name:        *def.Name,
			Environment: make(map[string]string),
			Secrets:     make(map[string]string),
		}
		for _, env := range def.Environment {
			containerDef.Environment[*env.Name] = *env.Value
		}
		for _, secret := range def.Secrets {
			containerDef.Secrets[*secret.Name] = *secret.ValueFrom
		}
		containerDefinitions = append(containerDefinitions, containerDef)
	}

	return containerDefinitions, nil
}

func (a *AWSProvider) GetSecretValue(secretArn, secretName string) (*fastjson.Value, error) {
	// Check if the secret value is cached
	a.mu.RLock()
	cachedValue, exists := a.SecretCache[secretArn]
	a.mu.RUnlock()

	if exists {
		return cachedValue, nil
	}

	// If not in cache, fetch it from AWS Secrets Manager
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	}

	resp, err := a.SecretsManagerService.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	// Parse the JSON string into a fastjson.Value
	v, err := fastjson.Parse(*resp.SecretString)
	if err != nil {
		return nil, err
	}

	// Cache the result
	a.mu.Lock()
	a.SecretCache[secretArn] = v
	a.mu.Unlock()

	return v, nil
}

func NewAWSProvider(region string) *AWSProvider {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return &AWSProvider{
		EcsService:            ecs.New(sess),
		SecretsManagerService: secretsmanager.New(sess),
		SecretCache:           make(map[string]*fastjson.Value),
	}
}

func RegisterProvider(region string) {
	provider.RegisterProvider("aws", NewAWSProvider(region))
}
