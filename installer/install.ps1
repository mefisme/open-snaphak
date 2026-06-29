# install.ps1 -- one-liner bootstrap for the SnapHak installer.
#
# Once the origin repo is set, end users run:
#   irm https://github.com/snaphak/open-snaphak/releases/latest/download/install.ps1 | iex
# It downloads snaphak.exe from the latest release into %LOCALAPPDATA%\open-snaphak and runs `snaphak install`.
# After that, run snaphak.exe directly for update / uninstall / status.
$ErrorActionPreference = "Stop"

$repo = "snaphak/open-snaphak"

$dest = Join-Path $env:LOCALAPPDATA "open-snaphak"
New-Item -ItemType Directory -Force $dest | Out-Null
$exe = Join-Path $dest "snaphak.exe"

Write-Host "Downloading snaphak.exe ..."
Invoke-WebRequest -Uri "https://github.com/$repo/releases/latest/download/snaphak.exe" -OutFile $exe

Write-Host "Installing SnapHak ..."
& $exe install

Write-Host ""
Write-Host "snaphak.exe is at $exe"
Write-Host "Run it for: snaphak update | snaphak uninstall | snaphak status"
