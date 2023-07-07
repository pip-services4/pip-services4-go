# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Google Cloud Platform specific components for Golang Changelog

## <a name="1.1.0"></a> 1.1.0 (2023-03-01)

### Breaking changes
* Renamed descriptors for services:
    - "\*:service:gcp-function\*:1.0" -> "\*:service:cloudfunc\*:1.0"
    - "\*:service:commandable-gcp-function\*:1.0" -> "\*:service:commandable-cloudfunc\*:1.0"

### Features
- Updated dependencies

## <a name="1.0.0"></a> 1.0.0 (2022-07-10) 

Initial public release

### Features

- **clients** - client components for working with Google Cloud Platform
- **connect** - components of installation and connection settings
- **container** - components for creating containers for Google server-side functions
- **services** - contains interfaces and classes used to create Google services

