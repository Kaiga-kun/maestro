# Maestro installer for Windows
# Usage: irm https://raw.githubusercontent.com/uprockcom/maestro/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "uprockcom/maestro"
$InstallDir = "$env:ProgramFiles\Maestro"
$ConfigDir = "$env:APPDATA\maestro"
$DockerImage = "ghcr.io/uprockcom/maestro"

# Colors
function Write-Info { Write-Host "ℹ $args" -ForegroundColor Blue }
function Write-Success { Write-Host "✓ $args" -ForegroundColor Green }
function Write-Warning { Write-Host "⚠ $args" -ForegroundColor Yellow }
function Write-ErrorMsg { Write-Host "✗ $args" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "╔════════════════════════════════════════╗"
Write-Host "║   Maestro Installer                   ║"
Write-Host "║   Multi-Container Claude Manager      ║"
Write-Host "╚════════════════════════════════════════╝"
Write-Host ""

# Check if running as Administrator
$IsAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $IsAdmin) {
    Write-Warning "Not running as Administrator"
    Write-Info "Trying to elevate privileges..."
    Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -Command `"irm https://raw.githubusercontent.com/uprockcom/maestro/main/install.ps1 | iex`"" -Verb RunAs
    exit
}

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ((Get-WmiObject Win32_Processor).Architecture -eq 12) { "arm64" } else { "x86_64" }
} else {
    Write-ErrorMsg "32-bit Windows is not supported"
}

Write-Info "Detected platform: Windows ($Arch)"

# Check for Docker
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-ErrorMsg "Docker is required but not installed. Please install Docker Desktop from https://www.docker.com/products/docker-desktop"
}

# Check if Docker is running
try {
    docker ps | Out-Null
    Write-Success "Docker is available"
} catch {
    Write-ErrorMsg "Docker daemon is not running. Please start Docker Desktop."
}

# Get latest version
Write-Info "Fetching latest version..."
$ReleasesUrl = "https://api.github.com/repos/$Repo/releases/latest"
try {
    $Release = Invoke-RestMethod -Uri $ReleasesUrl
    $Version = $Release.tag_name
    Write-Success "Latest version: $Version"
} catch {
    Write-ErrorMsg "Failed to fetch latest version from GitHub"
}

# Download binary
$FileName = "maestro_${Version}_Windows_${Arch}.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$FileName"

Write-Info "Downloading $FileName..."
$TempDir = Join-Path $env:TEMP "maestro-install"
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null
$ZipPath = Join-Path $TempDir $FileName

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath
    Write-Success "Downloaded successfully"
} catch {
    Write-ErrorMsg "Failed to download from $DownloadUrl"
}

# Extract
Write-Info "Extracting binary..."
Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force

# Install
Write-Info "Installing to $InstallDir..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Copy-Item -Path "$TempDir\maestro.exe" -Destination "$InstallDir\maestro.exe" -Force

# Add to PATH if not already there
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine)
if ($CurrentPath -notlike "*$InstallDir*") {
    Write-Info "Adding to system PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$InstallDir", [EnvironmentVariableTarget]::Machine)
    $env:Path = "$env:Path;$InstallDir"
    Write-Success "Added to PATH"
}

Write-Success "Installed to $InstallDir\maestro.exe"

# Verify installation
try {
    $InstalledVersion = & "$InstallDir\maestro.exe" version 2>$null | Select-Object -First 1
    Write-Success "Verified installation: $InstalledVersion"
} catch {
    Write-Warning "Installation complete but version check failed"
}

# Pull Docker image
Write-Info "Pulling Docker image..."
$VersionTag = $Version -replace '^v', ''
$DockerImageTag = "${DockerImage}:${VersionTag}"

try {
    docker pull $DockerImageTag 2>&1 | Out-Null
    Write-Success "Pulled $DockerImageTag"
} catch {
    Write-Warning "Failed to pull Docker image (you may need to authenticate or pull manually)"
}

# Create config directory
if (-not (Test-Path $ConfigDir)) {
    Write-Info "Creating config directory: $ConfigDir"
    New-Item -ItemType Directory -Force -Path $ConfigDir | Out-Null
    New-Item -ItemType Directory -Force -Path "$ConfigDir\.claude" | Out-Null
    Write-Success "Created $ConfigDir"
}

# Create initial config if doesn't exist
$ConfigFile = Join-Path $ConfigDir "config.yml"
if (-not (Test-Path $ConfigFile)) {
    Write-Info "Creating initial configuration..."

    @"
# Maestro Configuration
# See: https://github.com/uprockcom/maestro

claude:
  auth_path: $env:APPDATA\maestro\.claude

containers:
  prefix: mcl-
  image: ${DockerImage}:latest
  resources:
    memory: "4g"
    cpus: "2.0"

firewall:
  enabled: true
  allowed_domains:
    - "github.com"
    - "githubusercontent.com"
    - "npmjs.org"
    - "registry.npmjs.org"
    - "pypi.org"
    - "files.pythonhosted.org"
    - "anthropic.com"
    - "claude.ai"

sync:
  additional_folders: []

daemon:
  show_nag: true
  token_check_interval: 3600
  token_warning_threshold: 86400
"@ | Out-File -FilePath $ConfigFile -Encoding UTF8

    Write-Success "Config created at $ConfigFile"
}

# Cleanup
Remove-Item -Path $TempDir -Recurse -Force

Write-Host ""
Write-Success "Installation complete!"
Write-Host ""
Write-Info "Next steps:"
Write-Host "  1. Open a new PowerShell window (to refresh PATH)"
Write-Host "  2. Verify installation: maestro version"
Write-Host "  3. Authenticate with Claude: maestro auth"
Write-Host "  4. Create your first container: maestro new `"your task description`""
Write-Host ""
Write-Info "Documentation: https://github.com/uprockcom/maestro"
Write-Host ""
