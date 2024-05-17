# aws-apps-go

This Git repository contains various AWS applications showcasing different implementation scenarios, such as utilizing AWS Lambda functions, SNS topics, and CloudWatch cron triggers.
You can find explaination of these applications in respective articles here: https://solutiontoolkit.com

## Overview

The repository provides examples of serverless applications built using AWS Lambda functions, demonstrating how to leverage AWS services to create efficient and scalable solutions. Each application focuses on a specific use case and utilizes various AWS components to achieve its objectives.

## Applications

1. **aws-lambda-cron-go**
   - This application showcases how to create a scheduled event using CloudWatch cron expressions. It demonstrates how Lambda functions can be triggered at specific intervals using CloudWatch Events.
More details can be found here: https://solutiontoolkit.com/2023/02/how-to-invoke-aws-lambda-by-a-scheduled-event/ 

2. **aws-lambda-external-sns-topic-go**
   - This application illustrates the cross-account integration of different AWS lambda functions through SNS (Simple Notification Service) topics.
     https://solutiontoolkit.com/2023/07/cross-account-aws-lambda-functions-integration-with-sns/

## Getting Started

To explore each application and its implementation, navigate to the corresponding folder in the repository. Each application contains its source code, CloudFormation templates, and detailed instructions on how to deploy and run the application.

## Prerequisites

Before deploying any of the applications, ensure you have the following prerequisites:

- An AWS account with appropriate permissions to create and manage Lambda functions, SNS topics, and CloudWatch Events.
- AWS CLI installed and configured with the necessary credentials.
- Basic knowledge of AWS services and serverless architectures.

## License

This repository is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Feel free to explore, learn, and experiment with the AWS applications showcased in this repository. If you have any questions or encounter any issues, feel free to open an issue in the repository. Happy coding!
