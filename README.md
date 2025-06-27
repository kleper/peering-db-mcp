# PeeringDB Location Server

A Go-based MCP server that retrieves Internet Exchange (IX) locations for
Autonomous System Numbers (ASNs) using the PeeringDB API.

## Features

- Retrieves IX locations for any given ASN
- Shows operational status of each IX location
- Provides city information for each peering point
- Formatted output with emoji indicators for better readability

## Prerequisites

- Go 1.x or higher
- Access to the PeeringDB API

## Installation

```bash
go mod download
```

## Usage

### Stdio MCP Server

Run the original MCP server which communicates via stdin/stdout:

```bash
go run main.go
```

The tool `get_peering_locations` accepts an ASN and returns formatted
information about all IX locations where that AS is present.

### HTTP API

An HTTP server exposing the same functionality is available in
`server.go`. Create a `.env` file based on `.env.example` and set the
`PORT` variable. Then start the API server:

```bash
go run server.go
```

Endpoints:

- `GET /locations/{asn}` – returns peering locations for the given ASN
  in a formatted text response.
- `GET /openapi.json` – serves the OpenAPI specification describing the
  API.

## API Integration

The service queries the following PeeringDB endpoints:

1. `/api/net` – used to resolve an ASN to the PeeringDB network ID.
2. `/api/netixlan` – fetches IX location information.

## Error Handling

The server handles several error cases:
- Invalid ASN format
- ASN not found in PeeringDB
- API connection errors
- Invalid response data

## Contributing

Contributions are welcome! Please ensure your code follows the existing
structure and includes appropriate error handling.

## Note

This tool requires access to the PeeringDB API. Please ensure you comply
with PeeringDB's terms of service and API usage guidelines.
