$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Get-SwagCommand {
    $cmd = Get-Command swag -ErrorAction SilentlyContinue
    if ($cmd) { return $cmd.Source }

    $goPath = (& go env GOPATH) 2>$null
    if (-not $goPath) { throw "Go is not available in PATH. Install Go and ensure 'go' is callable." }

    $swagExe = Join-Path $goPath 'bin\swag.exe'
    if (Test-Path $swagExe) { return $swagExe }

    Write-Host "swag not found. Installing github.com/swaggo/swag/cmd/swag@latest ..."
    & go install github.com/swaggo/swag/cmd/swag@latest

    if (Test-Path $swagExe) { return $swagExe }
    throw "swag was installed but '$swagExe' was not found. Ensure GOPATH/bin is correct."
}

$swag = Get-SwagCommand

Write-Host "Using swag: $swag"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$backendDir = Join-Path $repoRoot 'backend'
$swaggerDir = $PSScriptRoot

if (-not (Test-Path $backendDir)) {
    throw "Backend directory not found: $backendDir"
}

Push-Location $backendDir
try {
    Write-Host "Generating backend swagger package (for Swagger UI) into backend/docs ..."
    & $swag init -g main.go -o ./docs

    Write-Host "Generating swagger artifacts (for frontend/codegen) into $swaggerDir ..."
    & $swag init -g main.go -o $swaggerDir --outputTypes json,yaml
}
finally {
    Pop-Location
}

Write-Host "Done. Files updated:"
Write-Host " - backend/docs/* (go + json + yaml)"
Write-Host " - swagger/swagger.json, swagger/swagger.yaml (json + yaml only)"


