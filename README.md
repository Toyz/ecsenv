# ECS Environment Generator

The ECS Environment Generator is a utility program that extracts environment variables and secrets from an AWS ECS task definition and generates a `.env` file.

## Features

- Extract environment variables and secrets from ECS task definitions.
- Cache secrets to improve performance and reduce AWS Secret Manager API calls.
- Support for different cloud providers (currently AWS is supported).
- Verbose logging for debugging purposes.

## Prerequisites

- Go (v1.16 or later)
- AWS credentials set up, either as environment variables (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) or via the AWS credentials file.

## Usage

### Build

To build the binary:

```bash
go build -o ecsenv
```

### Run
To generate a `.env` file from an ECS task definition:

```bash
./ecsenv generate-env [taskDefinitionName] --region [region] --provider aws
```

To list all support cloud providers:

```bash
./ecsenv list-providers
```

### Options

- `-v, --verbose`: Enable verbose output.
- `-r, --region`: Specify the AWS region. Defaults to `us-west-2`.
- `-p, --provider`: Specify the cloud provider. Defaults to `aws`.

## Testing

To run the unit tests:

```bash
go test -v ./...
```

## Contributing
Feel free to submit pull requests for new features or bug fixes. Please ensure that any changes include relevant unit tests.

## License
MIT