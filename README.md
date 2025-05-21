# HandballScore API

HandballScore API is a backend application designed to manage handball tournaments, matches, scores, teams, players, and related data. It is built in Go and intended for deployment as an AWS Lambda function, utilizing MongoDB for data storage and AWS services for infrastructure and configuration management.

## Features

The HandballScore API provides a comprehensive suite of features for managing handball competitions:

*   **User Management:** Secure registration, login (JWT-based), and password management.
*   **Tournament Management:** Create and update tournaments, manage tournament categories (e.g., different age groups or divisions within a tournament).
*   **Team Management:** Add, update, and delete teams, including associating them with tournaments.
*   **Player Management:** Manage player profiles, including personal details and team associations. Track player statistics within matches.
*   **Coach Management:** Manage coach profiles and their assignments.
*   **Referee Management:** Manage referee profiles and their assignments to matches.
*   **Match Management:**
    *   Schedule matches with specific dates, times, and locations.
    *   Live score tracking and updates.
    *   Record match events such as goals, yellow cards, red cards, blue cards, exclusions, and timeouts.
    *   Manage player and coach participation in matches.
    *   Start, pause, and end matches.
*   **League Phase Management:** Handle league-specific logic within tournaments.
*   **News & Announcements:** Create, update, and delete news articles or announcements related to the tournaments or associations.
*   **Association Management:** Manage details of handball associations.
*   **Category Management:** Manage general categories (e.g., age groups, skill levels) that can be applied to tournaments.

## Technologies Used

*   **Go (Golang):** The core programming language for the API.
*   **MongoDB:** NoSQL database used for data storage.
*   **AWS Lambda:** Serverless compute service for deploying the API.
*   **AWS Secrets Manager:** Used for securely managing database credentials and other secrets.
*   **AWS S3 (Simple Storage Service):** Utilized for storing files such as user, player, coach, and team avatars.
*   **JWT (JSON Web Tokens):** For securing API endpoints through authentication and authorization.

## Prerequisites for Setup

Before you can build or deploy the HandballScore API, ensure you have the following installed and configured:

*   **Go:** Version 1.18 or higher (or the specific version mentioned in `go.mod` if available).
*   **AWS CLI (Command Line Interface):** Configured with your AWS account credentials and default region. This is necessary for deploying to AWS Lambda and interacting with other AWS services.
*   **Access to an AWS Account:** You will need an AWS account with permissions to create and manage:
    *   AWS Lambda functions
    *   AWS Secrets Manager secrets
    *   AWS S3 buckets

## Environment Variables

The application relies on the following environment variables for its configuration, especially when deployed on AWS Lambda:

*   `SecretName`: (Required) The name or ARN of the secret stored in AWS Secrets Manager. This secret should contain the necessary credentials and configuration parameters.
*   `BucketName`: (Required) The name of the AWS S3 bucket used for storing uploaded files, such as avatars for users, players, coaches, and teams.
*   `UrlPrefix`: (Required) A URL prefix that might be used by API Gateway when routing requests to the Lambda function (e.g., `/api` or a stage name). This prefix is stripped from the request path to correctly route the request within the application.

### Secret Structure in AWS Secrets Manager

The secret identified by `SecretName` is expected to be a JSON object containing the following keys:

*   `Username`: Username for the MongoDB database.
*   `Password`: Password for the MongoDB database.
*   `Host`: The host address for the MongoDB Atlas cluster (e.g., `your-cluster.mongodb.net`).
*   `Database`: The name of the MongoDB database to use.
*   `JWTSign`: The secret key used to sign and verify JWTs.

## Database Setup

The HandballScore API uses **MongoDB** as its primary data store.

*   **Connection:** The application is configured to connect to a MongoDB Atlas cluster.
*   **Credentials:** Database connection details (username, password, host, and database name) are not hardcoded in the application. Instead, they are securely fetched at runtime from **AWS Secrets Manager** using the `SecretName` environment variable, as detailed in the "Environment Variables" section.
*   **Collections:** The application will create and use various collections within the specified MongoDB database to store data for tournaments, teams, players, matches, users, etc., as per the defined models.
*   **Local Development:** For local development, you would typically need to either:
    *   Connect to a MongoDB instance (local or cloud-based) by setting up your local environment to mimic the AWS Lambda environment regarding environment variables and secret access (e.g., by using a local secrets file and tools like `aws-vault` or by directly setting environment variables if your local setup bypasses actual AWS Secrets Manager calls).
    *   Or, if you have a development/staging AWS environment, configure your local setup to point to those AWS resources (Secrets Manager, S3).

## Build

The project includes a `Makefile` to simplify the build process. To build the application:

1.  Ensure Go is installed and your Go environment is correctly set up.
2.  Navigate to the root directory of the project in your terminal.
3.  Run the following command:

    ```bash
    make build
    ```

This command will:
*   Compile the Go application for a Linux AMD64 environment (suitable for AWS Lambda).
*   Create a directory named `bin` if it doesn't already exist.
*   Package the compiled binary (named `main`) into a zip file named `handballscore-app.zip` within the `bin` directory. This zip file is what you will deploy to AWS Lambda.

## Deployment

The project is designed to be deployed as an AWS Lambda function. The `Makefile` also includes a command to facilitate this process.

**Prerequisites for Deployment:**

*   You must have an existing AWS Lambda function configured (e.g., `handball-score` as specified in the `Makefile`).
*   The AWS CLI must be configured with credentials that have permission to update the specified Lambda function.
*   The Lambda function's execution role must have permissions to:
    *   Read from AWS Secrets Manager (for the secret specified by `SecretName`).
    *   Read/Write to the AWS S3 bucket (specified by `BucketName`).
    *   Write logs to AWS CloudWatch Logs.
*   The necessary environment variables (`SecretName`, `BucketName`, `UrlPrefix`) must be configured in the AWS Lambda function's settings.

**To deploy the application:**

1.  First, build the application using `make build` to create the deployment package (`bin/handballscore-app.zip`).
2.  Then, run the following command from the project's root directory:

    ```bash
    make deploy
    ```

This command uses the AWS CLI to update the code of the Lambda function specified by `LAMBDA_FUNCTION_NAME` (defaulting to `handball-score` in the `Makefile`) with the newly built `bin/handballscore-app.zip`.

**Note:** You might need to adjust the `LAMBDA_FUNCTION_NAME` in the `Makefile` if your Lambda function has a different name.
