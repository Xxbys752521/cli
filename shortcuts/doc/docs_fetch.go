// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"context"
	"fmt"
	"io"

	"github.com/larksuite/cli/shortcuts/common"
)

var DocsFetch = common.Shortcut{
	Service:     "docs",
	Command:     "+fetch",
	Description: "Fetch Lark document content",
	Risk:        "read",
	Scopes:      []string{"docx:document:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "doc", Desc: "document URL or token", Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/open-apis/docx/v1/documents/:document_id").
			Desc("Fetch docx metadata via OpenAPI").
			GET("/open-apis/docx/v1/documents/:document_id/blocks").
			Desc("Fetch docx block tree via OpenAPI").
			Set("backend", "openapi")
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		result, err := fetchDocxViaOpenAPI(runtime, runtime.Str("doc"))
		if err != nil {
			return err
		}

		runtime.OutFormat(result, nil, func(w io.Writer) {
			if title, ok := result["title"].(string); ok && title != "" {
				fmt.Fprintf(w, "# %s\n\n", title)
			}
			if md, ok := result["markdown"].(string); ok {
				fmt.Fprintln(w, md)
			}
			if hasMore, ok := result["has_more"].(bool); ok && hasMore {
				fmt.Fprintln(w, "\n--- more content available, use --offset and --limit to paginate ---")
			}
		})
		return nil
	},
}
