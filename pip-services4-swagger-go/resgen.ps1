#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate image and container names using the data in the "component.json" file
$component = Get-Content -Path "component.json" | ConvertFrom-Json

$resImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-res"
$container=$component.name

# Remove build files
if (Test-Path "./resources") {
    Remove-Item -Recurse -Force -Path "./resources/*"
} else {
    $null = New-Item -ItemType Directory -Force -Path "./resources"
}
if (Test-Path "./example/resources") {
    Remove-Item -Recurse -Force -Path "./example/resources/*"
} else {
    $null = New-Item -ItemType Directory -Force -Path "./example/resources"
}

# Build docker image
docker build -f docker/Dockerfile.resgen -t $resImage .

# Run resgen container
docker run -d --name $container $resImage

# Copy resources from container
docker cp "$($container):/app/src/example/resources" ./example
docker cp "$($container):/app/src/resources" .
# Remove docgen container
docker rm $container --force

# Verify resources
if (!(Test-Path "./resources/*.go")) {
    Write-Host "resources folder doesn't exist in root dir. Watch logs above."
    exit 1
}
