# build.ps1 -- the ONE top-level build script: compile the backend (XINPUT1_3.dll) and the WebView2
# (HTML) frontend (snapmap-plus-ui.dll) together: backend first, frontend second, always in lockstep.
# Pure ASCII.
#
# Why lockstep: the backend and the frontend share ONE ABI header (src\common\snapmap_plus_iface.h). Every
# vtable extension the frontend has picked up over time (apply_class_inherit, enum_valid_classes/
# enum_inherits, id_dev_layer_hidden, and now push_to_stack for the backend-hosted SnapStack port) must
# be resolved by a MATCHING backend build -- an old backend paired with a freshly-built frontend that
# reads a newer vtable slot is undefined behavior (the frontend reads past the end of the backend's
# smaller vtable struct). Running src\ui\build.ps1 alone NEVER rebuilds the backend, so it's easy to end
# up mismatched if you only touch the frontend side. This script removes that footgun: one command,
# backend first, frontend second. -BackendOnly skips the frontend build (safe in the other direction:
# a frontend built against an OLDER, smaller vtable never reads past a newer backend's).
#
# Usage:
#   pwsh -File build.ps1                     # backend + frontend (the default; what CI runs)
#   pwsh -File build.ps1 -BackendOnly        # backend only (faster loop when iterating on backend code)
#   pwsh -File build.ps1 -BackendOnly -Diag  # diagnostic backend build (DO NOT DISTRIBUTE)
#
# Any extra args (e.g. -Diag) forward straight through to src\backend\build.ps1.
#
# Needs: Build Tools for Visual Studio 2022 (C++ workload). No Qt required.
param(
    [switch]$BackendOnly,
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$Rest
)
$ErrorActionPreference = "Stop"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$here\src\backend\build.ps1" @Rest
if ($LASTEXITCODE -ne 0) { throw "backend build failed" }
if ($BackendOnly) {
    Write-Host "built: build/XINPUT1_3.dll (backend only)"
} else {
    powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$here\src\ui\build.ps1"
    if ($LASTEXITCODE -ne 0) { throw "frontend build failed" }
    Write-Host "built: build/XINPUT1_3.dll + build/webview/snapmap-plus-ui.dll"
    Write-Host "package with: package.ps1"
}
