# PKITool

PKITool is a web-based tool to manage Public Key Infrastructures (PKI) based on a REST API. 

## Features

todo

## ROADMAP

See [ROADMAP.md](ROADMAP.md)

## Development

The repo is organized in several folders and Go packages. The main places of interest are: 

 * **api/**: REST API specification file, following the [OpenAPI v3](https://swagger.io/specification/) standard. 
 * **internal/**: crypto operations and internal (storage-agnostic) representations of crypto objects.
 * **models/**: API model structures, generated from the API specification. 
 * **server/**: code of the HTTP server. The router is generated from the API specification.
 * **storage/**: different storage engines and their integration with the server.
 * **templates/**: Mustache templates used to generate Go models and packages from the specification. 
 * **util/**: boilerplate and conversion code. 

Several scripts are also maintained to help with the development:

 * **build.sh**: dockerized compilation.
 * **generate.sh**: regenerate router and models from the specification using [swagger codegen](https://github.com/swagger-api/swagger-codegen). 
 * **test.sh**: use curl against the REST API to run a full PKI manipulation scenario. 