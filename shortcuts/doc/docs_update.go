// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"context"
	"strings"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var validModes = map[string]bool{
	"append":        true,
	"overwrite":     true,
	"replace_range": true,
	"replace_all":   true,
	"insert_before": true,
	"insert_after":  true,
	"delete_range":  true,
}

var needsSelection = map[string]bool{
	"replace_range": true,
	"replace_all":   true,
	"insert_before": true,
	"insert_after":  true,
	"delete_range":  true,
}

var DocsUpdate = common.Shortcut{
	Service:     "docs",
	Command:     "+update",
	Description: "Update a Lark document",
	Risk:        "write",
	Scopes:      []string{"docx:document", "docx:document:readonly"},
	AuthTypes:   []string{"user", "bot"},
	Flags: []common.Flag{
		{Name: "doc", Desc: "document URL or token", Required: true},
		{Name: "mode", Desc: "update mode: append | overwrite | replace_range | replace_all | insert_before | insert_after | delete_range", Required: true},
		{Name: "markdown", Desc: "new content (Lark-flavored Markdown; create blank whiteboards with <whiteboard type=\"blank\"></whiteboard>, repeat to create multiple boards)"},
		{Name: "selection-with-ellipsis", Desc: "content locator (e.g. 'start...end')"},
		{Name: "selection-by-title", Desc: "title locator (e.g. '## Section')"},
		{Name: "new-title", Desc: "also update document title"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		mode := runtime.Str("mode")
		if !validModes[mode] {
			return common.FlagErrorf("invalid --mode %q, valid: append | overwrite | replace_range | replace_all | insert_before | insert_after | delete_range", mode)
		}

		if mode != "delete_range" && runtime.Str("markdown") == "" {
			return common.FlagErrorf("--%s mode requires --markdown", mode)
		}

		selEllipsis := runtime.Str("selection-with-ellipsis")
		selTitle := runtime.Str("selection-by-title")
		if selEllipsis != "" && selTitle != "" {
			return common.FlagErrorf("--selection-with-ellipsis and --selection-by-title are mutually exclusive")
		}

		if needsSelection[mode] && selEllipsis == "" && selTitle == "" {
			return common.FlagErrorf("--%s mode requires --selection-with-ellipsis or --selection-by-title", mode)
		}

		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		d := common.NewDryRunAPI().Set("backend", "openapi")
		mode := runtime.Str("mode")
		if mode == "overwrite" {
			d.GET("/open-apis/docx/v1/documents/:document_id/blocks/:document_id/children").
				Desc("Query root children before overwrite").
				DELETE("/open-apis/docx/v1/documents/:document_id/blocks/:document_id/children/batch_delete").
				Desc("Delete existing root children")
		}
		return d.POST("/open-apis/docx/v1/documents/:document_id/blocks/:document_id/children").
			Desc("Create all target blocks in a single append request").
			PATCH("/open-apis/docx/v1/documents/:document_id/blocks/batch_update").
			Desc("Batch update block text elements")
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if isWhiteboardCreateMarkdown(runtime.Str("markdown")) {
			return output.ErrValidation("OpenAPI adaptation does not support whiteboard markdown via docs +update")
		}
		if runtime.Str("new-title") != "" {
			return output.ErrValidation("OpenAPI adaptation does not yet support --new-title")
		}
		result, err := updateDocxViaOpenAPI(runtime, runtime.Str("doc"), runtime.Str("mode"), runtime.Str("markdown"))
		if err != nil {
			return err
		}

		normalizeDocsUpdateResult(result, runtime.Str("markdown"))
		runtime.Out(result, nil)
		return nil
	},
}

func normalizeDocsUpdateResult(result map[string]interface{}, markdown string) {
	if !isWhiteboardCreateMarkdown(markdown) {
		return
	}
	result["board_tokens"] = normalizeBoardTokens(result["board_tokens"])
}

func isWhiteboardCreateMarkdown(markdown string) bool {
	lower := strings.ToLower(markdown)
	if strings.Contains(lower, "```mermaid") || strings.Contains(lower, "```plantuml") {
		return true
	}
	return strings.Contains(lower, "<whiteboard") &&
		(strings.Contains(lower, `type="blank"`) || strings.Contains(lower, `type='blank'`))
}

func normalizeBoardTokens(raw interface{}) []string {
	switch v := raw.(type) {
	case nil:
		return []string{}
	case []string:
		return v
	case []interface{}:
		tokens := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				tokens = append(tokens, s)
			}
		}
		return tokens
	case string:
		if v == "" {
			return []string{}
		}
		return []string{v}
	default:
		return []string{}
	}
}
