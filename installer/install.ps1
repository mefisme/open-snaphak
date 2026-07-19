# install.ps1 -- one-liner bootstrap for the Snapmap+ installer.
#
# Once the origin repo is set, end users run:
#   irm https://github.com/doom-snapmap/snapmap-plus/releases/latest/download/install.ps1 | iex
# It downloads snapmap-plus.exe from the latest release into %LOCALAPPDATA%\snapmap-plus and runs
# `snapmap-plus install`. After that, run snapmap-plus.exe directly for update / uninstall / status.
$ErrorActionPreference = "Stop"

$repo = "doom-snapmap/snapmap-plus"

$dest = Join-Path $env:LOCALAPPDATA "snapmap-plus"
New-Item -ItemType Directory -Force $dest | Out-Null
$exe = Join-Path $dest "snapmap-plus.exe"

Write-Host "Downloading snapmap-plus.exe ..."
Invoke-WebRequest -Uri "https://github.com/$repo/releases/latest/download/snapmap-plus.exe" -OutFile $exe

Write-Host "Installing Snapmap+ ..."
& $exe install

Write-Host ""
Write-Host "snapmap-plus.exe is at $exe"
Write-Host "Run it for: snapmap-plus update | snapmap-plus uninstall | snapmap-plus status"
