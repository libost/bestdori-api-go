# Contributing to the Bestdori API Go Project

We welcome contributions from the community! If you would like to contribute to the Bestdori API Go project, please follow these guidelines:

- Create a new issue to discuss your proposed changes before submitting a pull request. This will help us understand the context of your contribution and ensure that it aligns with the project's goals.
- Fork the repository and create a new branch for your changes. Use a descriptive name for your branch that reflects the purpose of your changes.
- Write clear and concise commit messages that explain the purpose of your changes. Follow the conventional commit format if possible.
- Ensure that your code adheres to the project's coding standards and passes all tests before submitting a pull request. We use Go's built-in testing framework, so please write tests for any new functionality you add.

## Security Considerations

When contributing to the project, please keep security in mind. Avoid introducing vulnerabilities or exposing sensitive information in your code. If you discover a security vulnerability, please report it to us using GitHub's security advisories feature instead of creating a public issue or pull request.

## Some Tips for Testing

- A typical test for the whole project is to run `go test ./...` in the root directory. This will run all tests in the project and report any failures. Keep in mind that it may take 40-80 seconds to run all tests, depending on your machine's performance and Internet connection speed. If you want to run tests for a specific package, you can navigate to that package's directory and run `go test` there.
- Tests that require a Bestdori account to run will be skipped if the `BESTDORI_API_USER` and `BESTDORI_API_PSK` environment variables are not set. You can set these variables in your shell or in your IDE's run configuration.