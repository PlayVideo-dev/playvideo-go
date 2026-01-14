# Contributing to PlayVideo Go SDK

Thank you for your interest in contributing to the PlayVideo Go SDK!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/PlayVideo-dev/playvideo-go.git
   cd playvideo-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test -v ./...
   ```

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -v -race ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run go vet
go vet ./...
```

## Project Structure

```
playvideo-go/
├── client.go           # Main PlayVideo client
├── collections.go      # Collections resource
├── videos.go           # Videos resource
├── webhooks.go         # Webhooks resource
├── embed.go            # Embed resource
├── apikeys.go          # API Keys resource
├── account.go          # Account resource
├── usage.go            # Usage resource
├── errors.go           # Error types
├── types.go            # Type definitions
├── webhook.go          # Webhook signature verification
├── options.go          # Client configuration options
├── *_test.go           # Test files
└── examples/           # Example code
```

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test -v ./...`
6. Ensure vet passes: `go vet ./...`
7. Format code: `go fmt ./...`
8. Commit your changes: `git commit -m "Add my feature"`
9. Push to your fork: `git push origin feature/my-feature`
10. Open a Pull Request

## Code Style

- Follow standard Go conventions
- Use `go fmt` for formatting
- Use `go vet` for static analysis
- Add godoc comments to all exported types and functions
- Use meaningful variable and function names
- Keep functions small and focused

## Pull Request Guidelines

- Include a clear description of the changes
- Reference any related issues
- Add tests for new functionality
- Update documentation if needed
- Keep PRs focused on a single change

## Reporting Issues

When reporting issues, please include:

- SDK version
- Go version (`go version`)
- Operating system
- Minimal code to reproduce the issue
- Expected vs actual behavior

## Questions?

If you have questions, feel free to open an issue or reach out at support@playvideo.dev.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
