// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package core

// LarkBrand represents the Lark platform brand.
// "feishu" targets China-mainland, "lark" targets international.
// Any other string is treated as a custom base URL.
type LarkBrand string

const (
	BrandFeishu LarkBrand = "feishu"
	BrandLark   LarkBrand = "lark"
)

// Endpoints holds resolved endpoint URLs for different Lark services.
type Endpoints struct {
	Open     string // e.g. "https://open.feishu.cn"
	Accounts string // e.g. "https://accounts.feishu.cn"
	MCP      string // e.g. "https://mcp.feishu.cn"
}

// EndpointOverrides allows overriding resolved endpoints for private deployments.
// Configured via the "endpoints" field in config.json.
type EndpointOverrides struct {
	Open     string `json:"open,omitempty"`
	Accounts string `json:"accounts,omitempty"`
	MCP      string `json:"mcp,omitempty"`
}

// endpointOverrides holds the global endpoint overrides set from config.
var endpointOverrides *EndpointOverrides

// SetEndpointOverrides sets global endpoint overrides for private deployments.
func SetEndpointOverrides(o *EndpointOverrides) {
	endpointOverrides = o
}

// ResolveEndpoints resolves endpoint URLs based on brand,
// then applies any global endpoint overrides from config.
func ResolveEndpoints(brand LarkBrand) Endpoints {
	var ep Endpoints
	switch brand {
	case BrandLark:
		ep = Endpoints{
			Open:     "https://open.larksuite.com",
			Accounts: "https://accounts.larksuite.com",
			MCP:      "https://mcp.larksuite.com",
		}
	default:
		ep = Endpoints{
			Open:     "https://open.feishu.cn",
			Accounts: "https://accounts.feishu.cn",
			MCP:      "https://mcp.feishu.cn",
		}
	}
	if endpointOverrides != nil {
		if endpointOverrides.Open != "" {
			ep.Open = endpointOverrides.Open
		}
		if endpointOverrides.Accounts != "" {
			ep.Accounts = endpointOverrides.Accounts
		}
		if endpointOverrides.MCP != "" {
			ep.MCP = endpointOverrides.MCP
		}
	}
	return ep
}

// ResolveOpenBaseURL returns the Open API base URL for the given brand.
func ResolveOpenBaseURL(brand LarkBrand) string {
	return ResolveEndpoints(brand).Open
}
