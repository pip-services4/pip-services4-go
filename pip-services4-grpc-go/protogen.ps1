#!/usr/bin/env pwsh

##Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Get component data and set necessary variables
$component = Get-Content -Path "component.json" | ConvertFrom-Json

$docsImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-proto"
$container=$component.name

# Remove old generate files
Remove-Item -Path "./protos/*" -Force -Include *.go
Remove-Item -Path "test/protos/*" -Force -Include *.go

# Build docker image
docker build -f docker/Dockerfile.proto -t $docsImage .

# Create and copy compiled files, then destroy
docker create --name $container $docsImage
docker cp "$($container):/app/protos" .
docker cp "$($container):/app/test/protos" ./test/
docker rm $container
