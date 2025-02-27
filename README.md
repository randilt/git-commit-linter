# Git Commit Message Linter

A command-line tool that enforces consistent Git commit message formats across your projects. This linter helps teams maintain clean and meaningful commit histories by validating commit messages against predefined rules.

![Demo gif of the tool](https://github.com/user-attachments/assets/e6e8ba0a-dc20-46ee-90ca-b2c7484ee675)

## Features

- Validates commit message format (`type(scope): message`)
- Enforces message length limits
- Configurable commit types and scopes
- Detailed error messages with fix instructions
- Git hooks integration support
- Cross-platform compatibility (Windows, macOS, Linux)

## Installation

You can install the Git Commit Linter using the automated installation scripts or manually from the releases page.

### Quick Install

#### Unix-like Systems (macOS and Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/randilt/git-commit-linter/main/scripts/install.sh | bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/randilt/git-commit-linter/main/scripts/install.sh | bash
```

#### Windows

Open PowerShell as Administrator and run:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/randilt/git-commit-linter/main/scripts/install.ps1'))
```

### Manual Installation

If you prefer to install manually, you can download the appropriate binary for your system from our [releases page](https://github.com/randilt/git-commit-linter/releases).

#### Unix-like Systems (macOS and Linux)

1. Download the appropriate tar.gz file for your system
2. Extract the archive: `tar xzf git-commit-linter_*.tar.gz`
3. Move the binary to your PATH:
   ```bash
   sudo mv git-commit-linter_[your os and arch]/git-commit-linter /usr/local/bin/
   sudo chmod +x /usr/local/bin/git-commit-linter
   ```

#### Windows

1. Download the git-commit-linter_Windows_x86_64.zip file from the releases page
2. Extract the archive to a permanent location (e.g., `C:\Program Files\GitKit`)
3. Add the installation directory to your PATH:
   - Open System Properties (Win + Pause)
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Under "User variables", select "Path" and click "Edit"
   - Click "New" and add the installation directory path
   - Click "OK" to save

### Verifying Installation

After installation, verify that it's working correctly:

```bash
git-commit-linter version
```

### From Source

```bash
# Clone the repository
git clone https://github.com/randilt/git-commit-linter.git

# Navigate to the project directory
cd git-commit-linter

# Build the binary
go build -o git-commit-linter

# Move to a directory in your PATH so you can run it from anywhere (optional)
sudo mv git-commit-linter /usr/local/bin/
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
Linting Issues Found
────────────────────
✗ Commit 91170f3e: message too long (77 chars, max 72)
┌───────────────────────────────────────────────────┐
│ Fix Instructions:                                 │
│ - Older commit: Use interactive rebase            │
│   git rebase -i 91170f3e~1                        │
│   Change 'pick' to 'reword' for the target commit │
└───────────────────────────────────────────────────┘

Reference Information
─────────────────────
ℹ Valid commit format: type(scope): message (max 72 chars)
ℹ Allowed types: feat, fix, docs, style, refactor, test, chore
✗ Some commits failed linting - please fix the issues above
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
