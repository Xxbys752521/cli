// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"context"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var DocsCreate = common.Shortcut{
	Service:     "docs",
	Command:     "+create",
	Description: "Create a Lark document",
	Risk:        "write",
	AuthTypes:   []string{"user", "bot"},
	Scopes:      []string{"docx:document:create"},
	Flags: []common.Flag{
		{Name: "title", Desc: "document title"},
		{Name: "markdown", Desc: "Markdown content (Lark-flavored)", Required: true},
		{Name: "folder-token", Desc: "parent folder token"},
		{Name: "wiki-node", Desc: "wiki node token"},
		{Name: "wiki-space", Desc: "wiki space ID (use my_library for personal library)"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		count := 0
		if runtime.Str("folder-token") != "" {
			count++
		}
		if runtime.Str("wiki-node") != "" {
			count++
		}
		if runtime.Str("wiki-space") != "" {
			count++
		}
		if count > 1 {
			return common.FlagErrorf("--folder-token, --wiki-node, and --wiki-space are mutually exclusive")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		body := map[string]interface{}{}
		if v := runtime.Str("title"); v != "" {
			body["title"] = v
		}
		return common.NewDryRunAPI().
			POST("/open-apis/docx/v1/documents").
			Desc("Create docx document via OpenAPI").
			Body(body).
			Set("backend", "openapi")
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if runtime.Str("folder-token") != "" || runtime.Str("wiki-node") != "" || runtime.Str("wiki-space") != "" {
			return output.ErrValidation("OpenAPI adaptation does not yet support --folder-token/--wiki-node/--wiki-space")
		}
		result, err := createDocxViaOpenAPI(runtime, runtime.Str("title"), runtime.Str("markdown"))
		if err != nil {
			return err
		}

		runtime.Out(result, nil)
		return nil
	},
}
