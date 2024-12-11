package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PeeringDB API response structures
type NetResponse struct {
	Data []Network `json:"data"`
}

type Network struct {
	ID   int    `json:"id"`
	ASN  int    `json:"asn"`
	Name string `json:"name"`
}

type NetixlanResponse struct {
	Data []Netixlan `json:"data"`
}

type Netixlan struct {
	IXName string `json:"name"`
	City   string `json:"city"`
	Status bool   `json:"operational"`
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"PeeringDB Location Server",
		"1.0.0",
		server.WithPromptCapabilities(true),
	)

	// Add get locations tool
	locationsTool := mcp.NewTool("get_peering_locations",
		mcp.WithDescription("Get IX locations for an ASN"),
		mcp.WithNumber("asn",
			mcp.Required(),
			mcp.Description("Autonomous System Number"),
		),
	)

	// Register tool
	s.AddTool(locationsTool, getPeeringLocationsHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func getPeeringLocationsHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	asn, ok := arguments["asn"].(float64) // JSON numbers come as float64
	if !ok {
		return mcp.NewToolResultError("ASN must be a number"), nil
	}

	// Get network ID from ASN
	netID, err := getNetworkID(int(asn))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting network ID: %v", err)), nil
	}

	// Get peering locations
	locations, err := getPeeringLocations(netID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting peering locations: %v", err)), nil
	}

	// Format the response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("🌐 IX Locations for AS%d:\n\n", int(asn)))

	if len(locations) == 0 {
		response.WriteString("No IX locations found for this ASN.")
	} else {
		for i, loc := range locations {
			status := "🟢 Operational"
			if !loc.Status {
				status = "🔴 Not Operational"
			}
			response.WriteString(fmt.Sprintf("%d. %s\n   📍 %s\n   Status: %s\n\n",
				i+1,
				loc.IXName,
				loc.City,
				status))
		}
	}

	return mcp.NewToolResultText(response.String()), nil
}

func getNetworkID(asn int) (int, error) {
	url := fmt.Sprintf("https://api.peeringdb.com/api/net?asn=%d", asn)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var netResp NetResponse
	if err := json.Unmarshal(body, &netResp); err != nil {
		return 0, err
	}

	if len(netResp.Data) == 0 {
		return 0, fmt.Errorf("ASN %d not found in PeeringDB", asn)
	}

	return netResp.Data[0].ID, nil
}

func getPeeringLocations(netID int) ([]Netixlan, error) {
	url := fmt.Sprintf("https://api.peeringdb.com/api/netixlan?net_id=%d", netID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var nixResp NetixlanResponse
	if err := json.Unmarshal(body, &nixResp); err != nil {
		return nil, err
	}

	return nixResp.Data, nil
}
