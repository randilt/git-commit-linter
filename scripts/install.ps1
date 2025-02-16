# Installation script for Git Commit Linter on Windows

# Configuration
$GitHubUser = "randilt"
$RepoName = "git-commit-linter"
$BinaryName = "git-commit-linter.exe"
$InstallDir = "$env:ProgramFiles\GitKit"
$TempDir = "$env:TEMP\$RepoName-install"

# Error handling
$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "Error: $Message" -ForegroundColor Red
    exit 1
}

# Cleanup function
function Cleanup {
    if (Test-Path $TempDir) {
        Remove-Item -Recurse -Force $TempDir
    }
}

# Ensure running with administrator privileges
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Error "This script requires administrator privileges. Please run as administrator."
}

# Create temporary directory
Write-Step "Creating temporary directory..."
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null

try {
    # Get latest release info
    Write-Step "Fetching latest release information..."
    $releaseUrl = "https://api.github.com/repos/$GitHubUser/$RepoName/releases/latest"
    $release = Invoke-RestMethod -Uri $releaseUrl

    # Find Windows x64 asset
    $asset = $release.assets | Where-Object { $_.name -like "*Windows_x86_64.zip" }
    if (-not $asset) {
        Write-Error "Could not find Windows x64 release asset"
    }

    # Download release
    Write-Step "Downloading latest release..."
    $downloadPath = "$TempDir\$RepoName.zip"
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $downloadPath

    # Extract archive
    Write-Step "Extracting archive..."
    Expand-Archive -Path $downloadPath -DestinationPath $TempDir

    # Get the actual extracted directory name
    $extractedDir = Get-ChildItem -Path $TempDir -Directory | Where-Object { $_.Name -like "*Windows_x86_64" } | Select-Object -First 1
    if (-not $extractedDir) {
        Write-Error "Could not find extracted directory"
    }

    # Create installation directory
    Write-Step "Creating installation directory..."
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    }

    # Check if binary exists in the extracted directory
    $binaryPath = Join-Path $extractedDir.FullName $BinaryName
    if (-not (Test-Path $binaryPath)) {
        Write-Error "Binary not found in the extracted directory"
    }

    # Check if binary already exists in install directory
    $installBinaryPath = Join-Path $InstallDir $BinaryName
    if (Test-Path $installBinaryPath) {
        $confirmation = Read-Host "Binary already exists in $InstallDir. Overwrite? (y/N)"
        if ($confirmation -ne 'y') {
            Write-Error "Installation cancelled by user"
        }
    }

    # Copy files to installation directory
    Write-Step "Installing files..."
    Copy-Item -Path $binaryPath -Destination $InstallDir -Force

    # Add to PATH
    Write-Step "Adding to PATH..."
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$userPath;$InstallDir",
            "User"
        )
    }

    # Refresh PATH in current session
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    # Verify installation
    Write-Step "Verifying installation..."
    try {
        $version = & "$InstallDir\$BinaryName" version
        Write-Success "Installation completed successfully!"
        Write-Success "Version information: $version"
    }
    catch {
        Write-Error "Installation verification failed. Please ensure $InstallDir is in your PATH"
    }
}
catch {
    Write-Error $_.Exception.Message
}
finally {
    Cleanup
}