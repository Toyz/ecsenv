package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/toyz/ecsenv/provider"
	"github.com/toyz/ecsenv/provider/aws"
	"os"
	"strings"
)

var verbose bool
var region string
var providerName string
var outputFileName string

var rootCmd = &cobra.Command{
	Use:   "your-program-name",
	Short: "A program to generate .env from cloud providers",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		registerSelectedProvider()
	},
}

var generateEnvCmd = &cobra.Command{
	Use:   "generate-env [taskDefinitionName]",
	Short: "Generate .env from ECS task definition",
	Args:  cobra.ExactArgs(1),
	Run:   runECSEnv,
}

var listProvidersCmd = &cobra.Command{
	Use:   "list-providers",
	Short: "List all registered cloud providers",
	Run: func(cmd *cobra.Command, args []string) {
		providers := provider.Providers()
		if len(providers) == 0 {
			log.Info("No providers registered.")
			return
		}

		log.Info("Registered providers:")
		for _, p := range providers {
			log.Infof(" - %s", p)
		}
	},
}

func generateEnvContent(taskDefinitionName string, cloudProvider provider.Provider) ([]byte, error) {
	containerDefinitions, err := cloudProvider.GetTaskDefinitions(taskDefinitionName)
	if err != nil {
		return nil, fmt.Errorf("error fetching task definitions (%s): %v", taskDefinitionName, err)
	}

	var buffer strings.Builder

	// Iterate over each container definition
	for _, containerDefinition := range containerDefinitions {
		// Append environment variables
		for key, value := range containerDefinition.Environment {
			buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		}

		// Fetch and append secret values
		for secretName, secretArn := range containerDefinition.Secrets {
			secretValue, err := cloudProvider.GetSecretValue(secretArn, secretName)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch secret for ARN %s: %v", secretArn, err)
			}

			secretValueByName := secretValue.Get(secretName)
			if secretValueByName == nil {
				return nil, fmt.Errorf("no value found for secret: %s", secretName)
			}

			buffer.WriteString(fmt.Sprintf("%s=%s\n", secretName, string(secretValueByName.GetStringBytes())))
		}
	}
	return []byte(buffer.String()), nil
}

func determineOutputFileName(output, taskDefinition string) string {
	if output != "" {
		return output
	}
	return fmt.Sprintf("%s.env", taskDefinition)
}

func runECSEnv(cmd *cobra.Command, args []string) {
	taskDefinitionName := args[0]

	// Adjust the logging level based on the verbose flag
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	cloudProvider, exists := provider.GetProvider(providerName)
	if !exists {
		log.Fatalf("Provider '%s' not registered", providerName)
		return
	}

	content, err := generateEnvContent(taskDefinitionName, cloudProvider)
	if err != nil {
		log.Fatal(err)
		return
	}

	envFileName := determineOutputFileName(outputFileName, taskDefinitionName)
	if err := os.WriteFile(envFileName, content, 0644); err != nil {
		log.Fatalf("Failed to write to %s: %v", envFileName, err)
		return
	}

	log.Infof(".env file generated: %s", envFileName)
}

func registerSelectedProvider() {
	switch providerName {
	case "aws":
		aws.RegisterProvider(region)
		// Add more providers as cases when you expand your program
	default:
		log.Fatalf("Unsupported provider: %s", providerName)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-west-2", "Provider region (e.g. us-west-2)")
	rootCmd.PersistentFlags().StringVarP(&providerName, "provider", "p", "aws", "Cloud provider (e.g. aws, gcp)")

	rootCmd.AddCommand(listProvidersCmd)
	rootCmd.AddCommand(generateEnvCmd)

	generateEnvCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "Optional output file name (e.g. output.env)")
}

func main() {
	aws.RegisterProvider(region)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
