# Git Commit Message Linter

A command-line tool that enforces consistent Git commit message formats across your projects. This linter helps teams maintain clean and meaningful commit histories by validating commit messages against predefined rules.

## Features

- Validates commit message format (`type(scope): message`)
- Enforces message length limits
- Configurable commit types and scopes
- Detailed error messages with fix instructions
- Git hooks integration support
- Cross-platform compatibility (Windows, macOS, Linux)

## Installation

A standalone installer will be available soon. In the meantime, you can install the by building it from source.

### From Source

```bash
# Clone the repository
git clone https://github.com/randilt/git-commit-linter.git

# Navigate to the project directory
cd git-commit-linter

# Build the binary
go build -o git-commit-linter

# Move to a directory in your PATH so you can run it from anywhere (optional)
mv git-commit-linter /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Check the last commit
git-commit-linter

# Check last N commits
git-commit-linter --check="HEAD~5..HEAD"

# Use custom config file
git-commit-linter --config=path/to/config.yaml
```

### Command Line Flags

- `--check`: Specify commit range to check (default: "HEAD^..HEAD")
- `--config`: Path to custom configuration file
- `--help`: Display help information

## Configuration

Create a `config.yaml` file to customize the linter rules:

```yaml
types:
  - feat
  - fix
  - docs
  - style
  - refactor
  - test
  - chore

scopes:
  - auth
  - api
  - ui
  - db
  - config

rules:
  require_scope: false
  max_message_length: 72
```

### Default Rules

- Valid commit types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Maximum message length: 72 characters
- Scope is optional by default

## Git Hooks Integration

### Pre-commit Hook

1. Rename `commit-msg.sample` file in `.git/hooks` directory to `commit-msg` in your project directory to create a pre-commit hook.

2. Add the following script to the `commit-msg` file:

```bash
#!/bin/sh
commit_msg_file="$1"

# Run the linter with the commit message file
git-commit-linter lint-file "$commit_msg_file" || exit 1
```

3. Make it executable:

```bash
chmod +x .git/hooks/commit-msg
```

Now, the linter will run automatically before each commit making sure your commit messages are properly formatted.

## Valid Commit Message Format

```
type(scope): message

Examples:
feat(auth): add OAuth2 support
fix(api): handle null response from server
docs(readme): update installation instructions
```

### Components

- `type`: The type of change being made (required)
- `scope`: The area of the codebase affected (optional)
- `message`: A concise description of the change (required)

## Error Messages and Fixes

When the linter finds issues, it provides clear error messages and fix instructions:

```
Linting Issues Found:
==================

Commit abc1234: message too long (80 chars, max 72)
Fix Instructions:
- Older commit: Use interactive rebase
  git rebase -i abc1234~1
  Change 'pick' to 'reword' for the target commit

Reference Information:
====================
Valid commit format: type(scope): message (max 72 chars)
Allowed types: feat, fix, docs, style, refactor, test, chore
```

### Common Fixes

1. **Fix Latest Commit**

```bash
git commit --amend -m "feat(scope): your message"
```

2. **Fix Older Commits**

```bash
git rebase -i <commit-hash>~1
# Change 'pick' to 'reword' for the target commit
```

3. If you have already pushed the commit, you will need to force push:

```bash
git push --force-with-lease
```

`--force-with-lease` will prevent you from overwriting changes on the remote branch that you are not aware of.
But if you are sure that you want to overwrite the remote branch, you can use `--force` instead.

```bash
git push --force
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Building for Different Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o git-commit-linter.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o git-commit-linter-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o git-commit-linter-linux
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using the proper format
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by conventional commit messages
- Built with Go and ❤️

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
