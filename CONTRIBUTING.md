# Contributing to golang-mock-oauth2-server

Thank you for considering contributing to this project! Here are some guidelines to help you get started.

## Code of Conduct

By participating in this project, you are expected to uphold our [Code of Conduct](./CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check our [issue list](https://github.com/chrisw-dev/golang-mock-oauth2-server/issues) to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** (config values, command-line arguments, etc.)
- **Describe the behavior you observed and what you expected to see**
- **Include logs and/or screenshots** if applicable
- **Specify the version** of the software you're using

### Suggesting Enhancements

Enhancement suggestions are tracked as [GitHub issues](https://github.com/chrisw-dev/golang-mock-oauth2-server/issues).

- **Use a clear and descriptive title**
- **Provide a step-by-step description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **Specify which version you're using**

### Pull Requests

- **Fork the repo** and create your branch from `main`
- **Follow our coding style** and ensure tests pass
- **Write or update tests** for new functionality
- **Ensure your code lints** (we use golangci-lint)
- **Document new code** with comments
- **Update documentation** if needed
- **Submit a pull request** with a clear description of the changes

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/chrisw-dev/golang-mock-oauth2-server.git
   cd golang-mock-oauth2-server
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the tests:
   ```bash
   go test ./...
   ```

4. Build the application:
   ```bash
   go build -o mock-oauth2-server ./cmd/server
   ```

## Coding Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use meaningful variable and function names
- Write descriptive comments
- Add unit tests for new functionality
- Format your code using `go fmt`
- Verify your code using `go vet` and `golangci-lint`

## Release Process

This project uses an automated release process. When changes are merged to the `main` branch:

1. **Continuous Integration (CI)** runs all tests and checks
2. **Automatic Tagging**: After tests pass, a new version tag is automatically created
   - The patch version is incremented (e.g., `v1.0.0` → `v1.0.1`)
   - If no tags exist, the first version will be `v0.1.0`
3. **Docker Image Build**: The Docker image is built and pushed to GitHub Container Registry
4. **GitHub Release**: A new GitHub release is created with auto-generated release notes

### Version Numbering

The project follows [Semantic Versioning](https://semver.org/):
- **Major version** (vX.0.0): Breaking changes (manual tag required)
- **Minor version** (v1.X.0): New features, backwards compatible (manual tag required)
- **Patch version** (v1.0.X): Bug fixes, backwards compatible (automated)

By default, each merge to `main` increments the patch version automatically.

### Manual Version Bumps

To create a major or minor version release, manually create and push a tag:

```bash
# For a minor version bump (new features)
git tag -a v1.1.0 -m "Release v1.1.0: Add new features"
git push origin v1.1.0

# For a major version bump (breaking changes)
git tag -a v2.0.0 -m "Release v2.0.0: Breaking changes"
git push origin v2.0.0
```

Manual tags will skip the auto-tagging process and directly trigger the release workflow.

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](./LICENSE).