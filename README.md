# Transactions Email Processor

This project is a simple application that processes a file containing debit and credit transactions on an account and generates a summary email with relevant information. The application is built using Go and packaged into a Docker container for easy deployment.

### Prerequisites

- Docker

### How to Run

1. Pull the Docker image from the GitHub Container Registry:

```shell
docker pull --platform linux/x86_64 ghcr.io/zejiran/transactions-email-processor:master
```

2. Create a `.env` file in the same directory as the Dockerfile with the following environment variables:

```
SENDER_MAIL=sender@gmail.com
PASSWORD=secret-password
RECIPIENT_MAIL=recipient@example.com
```

Make sure to replace the email addresses with your own sender and recipient addresses.

3. Run the Docker container with the -v flag to mount a directory for input files:

```shell
docker run --platform linux/x86_64 --env-file ./.env -v /<absolute-host-path>/files:/files ghcr.io/zejiran/transactions-email-processor:master
```

The `-v` flag maps the `files` directory on the host machine to the `/files` directory in the container. This allows you to provide input files from the host machine to the container.

4. Alternatively, you can specify environment variables as command arguments:

```shell
docker run --platform linux/x86_64 --env-file ./.env --env RECIPIENT_MAIL=new.recipient@example.com -v /<absolute-host-path>/files:/files ghcr.io/zejiran/transactions-email-processor:master
```

### File Format

The input file should be in CSV format with the following columns:

```
Id,Date,Transaction
```

Each row represents a transaction, where the `Id` is a unique identifier, `Date` is the date of the transaction in the format `MM/DD`, and `Transaction` is the transaction amount. Credit transactions are indicated with a plus sign (+) and debit transactions with a minus sign (-).

### Customization

The email template used for generating the email body is located in the `email/dist/email_template.html` file. You can modify the template HTML to customize the email content, styling, and layout. Make sure to keep the placeholders for the dynamic data intact.

### Configuration

The application requires the following environment variables to be set:

- `SENDER_MAIL`: The email address of the sender.
- `PASSWORD`: The password or app-specific password for the sender's email account.
- `RECIPIENT_MAIL`: The email address of the recipient.

Make sure to set these environment variables in the `.env` file before running the Docker container.

### Repository Structure

- `main.go`: The main Go application code that reads the input file, processes the transactions, generates the email body, and sends the email.
- `Dockerfile`: The Dockerfile for building the Docker image of the application.
- `email/src/email_template.heml`: The source HEML code used for creating the HTML template.
- `email/dist/email_template.html`: The HTML template used for generating the email body.
- `files/example_transactions.csv`: An example input file with debit and credit transactions.

### Example Received Email

<img width="300" alt="example" src="https://github.com/zejiran/transactions-email-processor/assets/30379522/0eb4f923-94d9-4234-9a2f-a87e2c8b511b">

### Notes

- The email sending functionality is currently configured for Gmail. If you are using a different email provider, you may need to update the SMTP server details in the `sendEmail` function in `main.go`.

Please feel free to reach out if you have any questions or need further assistance!

## License

[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

- **[MIT license](LICENSE)**
- Copyright 2023 © Juan Alegría.
