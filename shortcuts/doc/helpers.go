// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/internal/validate"
	"github.com/larksuite/cli/shortcuts/common"
)

type documentRef struct {
	Kind  string
	Token string
}

func parseDocumentRef(input string) (documentRef, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return documentRef{}, output.ErrValidation("--doc cannot be empty")
	}

	if token, ok := extractDocumentToken(raw, "/wiki/"); ok {
		return documentRef{Kind: "wiki", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/docx/"); ok {
		return documentRef{Kind: "docx", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/doc/"); ok {
		return documentRef{Kind: "doc", Token: token}, nil
	}
	if strings.Contains(raw, "://") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx URL/token or a wiki URL that resolves to docx", raw)
	}
	if strings.ContainsAny(raw, "/?#") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx token or a wiki URL", raw)
	}

	return documentRef{Kind: "docx", Token: raw}, nil
}

func extractDocumentToken(raw, marker string) (string, bool) {
	idx := strings.Index(raw, marker)
	if idx < 0 {
		return "", false
	}
	token := raw[idx+len(marker):]
	if end := strings.IndexAny(token, "/?#"); end >= 0 {
		token = token[:end]
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}

func buildDriveRouteExtra(docID string) (string, error) {
	extra, err := json.Marshal(map[string]string{"drive_route_token": docID})
	if err != nil {
		return "", output.Errorf(output.ExitInternal, "internal_error", "failed to marshal upload extra data: %v", err)
	}
	return string(extra), nil
}

func createDocxViaOpenAPI(runtime *common.RuntimeContext, title, markdown string) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if strings.TrimSpace(title) != "" {
		body["title"] = title
	}
	data, err := runtime.CallAPI("POST", "/open-apis/docx/v1/documents", nil, body)
	if err != nil {
		return nil, err
	}
	document, _ := data["document"].(map[string]interface{})
	docID, _ := document["document_id"].(string)
	if docID == "" {
		return nil, output.Errorf(output.ExitAPI, "api_error", "docx create returned no document_id")
	}
	if strings.TrimSpace(markdown) != "" {
		if err := appendMarkdownAsPlainText(runtime, docID, markdown); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"document_id": docID,
		"title":       document["title"],
		"revision_id": document["revision_id"],
		"backend":     "openapi",
	}, nil
}

func fetchDocxViaOpenAPI(runtime *common.RuntimeContext, docInput string) (map[string]interface{}, error) {
	docID, err := resolveDocxDocumentID(runtime, docInput)
	if err != nil {
		return nil, err
	}
	meta, err := runtime.CallAPI("GET", fmt.Sprintf("/open-apis/docx/v1/documents/%s", validate.EncodePathSegment(docID)), nil, nil)
	if err != nil {
		return nil, err
	}
	document, _ := meta["document"].(map[string]interface{})
	blocks, err := runtime.CallAPI("GET", fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks", validate.EncodePathSegment(docID)), nil, nil)
	if err != nil {
		return nil, err
	}
	items, _ := blocks["items"].([]interface{})
	lines := make([]string, 0, len(items))
	for _, item := range items {
		block, _ := item.(map[string]interface{})
		if block == nil {
			continue
		}
		if text := extractBlockText(block); strings.TrimSpace(text) != "" {
			lines = append(lines, text)
		}
	}
	title, _ := document["title"].(string)
	return map[string]interface{}{
		"document_id": docID,
		"title":       title,
		"markdown":    strings.Join(lines, "\n\n"),
		"backend":     "openapi",
	}, nil
}

func updateDocxViaOpenAPI(runtime *common.RuntimeContext, docInput, mode, markdown string) (map[string]interface{}, error) {
	docID, err := resolveDocxDocumentID(runtime, docInput)
	if err != nil {
		return nil, err
	}
	if mode != "append" && mode != "overwrite" {
		return nil, output.ErrValidation("private deployment fallback currently supports --mode append|overwrite")
	}
	if mode == "overwrite" {
		if err := clearDocxRootChildren(runtime, docID); err != nil {
			return nil, err
		}
	}
	if err := appendMarkdownAsPlainText(runtime, docID, markdown); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"document_id": docID,
		"mode":        mode,
		"backend":     "openapi",
	}, nil
}

func appendMarkdownAsPlainText(runtime *common.RuntimeContext, docID, markdown string) error {
	blocks := markdownBlocks(markdown, openBaseURL(runtime))
	if len(blocks) == 0 {
		return nil
	}
	baseIndex, err := getDocxRootChildrenCount(runtime, docID)
	if err != nil {
		return err
	}
	createData, err := runtime.CallAPI("POST",
		fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
		nil,
		map[string]interface{}{
			"children": blocks,
			"index":    baseIndex,
		})
	if err != nil {
		return err
	}

	requests, err := buildBatchUpdateTextRequests(blocks, createData)
	if err != nil {
		return err
	}
	if len(requests) == 0 {
		return nil
	}

	if _, err := runtime.CallAPI("PATCH",
		fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/batch_update", validate.EncodePathSegment(docID)),
		nil,
		map[string]interface{}{"requests": requests}); err != nil {
		return err
	}
	return nil
}

func buildBatchUpdateTextRequests(sourceBlocks []map[string]interface{}, createData map[string]interface{}) ([]interface{}, error) {
	rawChildren, _ := createData["children"].([]interface{})
	if len(rawChildren) == 0 {
		return nil, output.Errorf(output.ExitAPI, "api_error", "docx create children returned no blocks")
	}

	requests := make([]interface{}, 0, len(sourceBlocks))
	for idx, source := range sourceBlocks {
		if idx >= len(rawChildren) {
			break
		}
		elements := extractElementsForBatchUpdate(source)
		if len(elements) == 0 {
			continue
		}

		child, _ := rawChildren[idx].(map[string]interface{})
		blockID, _ := child["block_id"].(string)
		if blockID == "" {
			return nil, output.Errorf(output.ExitAPI, "api_error", "docx create children returned empty block_id")
		}

		requests = append(requests, map[string]interface{}{
			"block_id": blockID,
			"update_text_elements": map[string]interface{}{
				"elements": elements,
			},
		})
	}
	return requests, nil
}

func extractElementsForBatchUpdate(block map[string]interface{}) []interface{} {
	for _, key := range []string{
		"text",
		"heading1",
		"heading2",
		"heading3",
		"bullet",
		"ordered",
		"quote",
		"code",
		"todo",
	} {
		node, _ := block[key].(map[string]interface{})
		if node == nil {
			continue
		}
		elements, _ := node["elements"].([]interface{})
		if len(elements) == 0 {
			continue
		}
		return elements
	}
	return nil
}

func getDocxRootChildrenCount(runtime *common.RuntimeContext, docID string) (int, error) {
	rootData, err := runtime.CallAPI("GET",
		fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
		nil, nil)
	if err != nil {
		return 0, err
	}
	children, _ := rootData["items"].([]interface{})
	return len(children), nil
}

func clearDocxRootChildren(runtime *common.RuntimeContext, docID string) error {
	for {
		count, err := getDocxRootChildrenCount(runtime, docID)
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		if _, err := runtime.CallAPI("DELETE",
			fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children/batch_delete", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
			nil,
			map[string]interface{}{"start_index": 0, "end_index": count}); err != nil {
			return err
		}
	}
}

var orderedListPattern = regexp.MustCompile(`^\d+\.\s+`)
var todoUncheckedPattern = regexp.MustCompile(`^- \[ \]\s+`)
var todoCheckedPattern = regexp.MustCompile(`^- \[[xX]\]\s+`)
var markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
var mentionUserPattern = regexp.MustCompile(`@(?:user|用户)\[([^\]]+)\]`)
var mentionDocPattern = regexp.MustCompile(`@(?:doc|文档)\[([^\]]+)\]`)

func markdownBlocks(markdown, openBase string) []map[string]interface{} {
	rawLines := strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n")
	blocks := make([]map[string]interface{}, 0, len(rawLines))
	inCode := false
	codeLines := make([]string, 0)

		flushCode := func() {
			if !inCode {
				return
			}
			content := strings.TrimSpace(strings.Join(codeLines, "\n"))
			if content != "" {
				blocks = append(blocks, buildTextualBlock(14, "code", content, openBase))
			}
		inCode = false
		codeLines = codeLines[:0]
	}

	for _, raw := range rawLines {
		line := strings.TrimRight(raw, " \t")
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				flushCode()
			} else {
				inCode = true
				codeLines = codeLines[:0]
			}
			continue
		}
		if inCode {
			codeLines = append(codeLines, line)
			continue
		}
		if trimmed == "" {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "# "):
			blocks = append(blocks, buildTextualBlock(3, "heading1", strings.TrimSpace(strings.TrimPrefix(trimmed, "# ")), openBase))
		case strings.HasPrefix(trimmed, "## "):
			blocks = append(blocks, buildTextualBlock(4, "heading2", strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), openBase))
		case strings.HasPrefix(trimmed, "### "):
			blocks = append(blocks, buildTextualBlock(5, "heading3", strings.TrimSpace(strings.TrimPrefix(trimmed, "### ")), openBase))
		case trimmed == "---":
			blocks = append(blocks, map[string]interface{}{
				"block_type": 22,
				"divider":    map[string]interface{}{},
			})
		case todoUncheckedPattern.MatchString(trimmed):
			content := strings.TrimSpace(todoUncheckedPattern.ReplaceAllString(trimmed, ""))
			blocks = append(blocks, buildTodoBlock(content, false, openBase))
		case todoCheckedPattern.MatchString(trimmed):
			content := strings.TrimSpace(todoCheckedPattern.ReplaceAllString(trimmed, ""))
			blocks = append(blocks, buildTodoBlock(content, true, openBase))
		case strings.HasPrefix(trimmed, "- "), strings.HasPrefix(trimmed, "* "):
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* "))
			blocks = append(blocks, buildTextualBlock(12, "bullet", content, openBase))
		case orderedListPattern.MatchString(trimmed):
			content := orderedListPattern.ReplaceAllString(trimmed, "")
			blocks = append(blocks, buildTextualBlock(13, "ordered", strings.TrimSpace(content), openBase))
		case strings.HasPrefix(trimmed, "> "):
			content := strings.TrimSpace(strings.TrimPrefix(trimmed, "> "))
			blocks = append(blocks, buildTextualBlock(15, "quote", content, openBase))
		default:
			blocks = append(blocks, buildTextualBlock(2, "text", trimmed, openBase))
		}
	}

	flushCode()
	return blocks
}

func buildTextualBlock(blockType int, fieldName, content, openBase string) map[string]interface{} {
	return map[string]interface{}{
		"block_type": blockType,
		fieldName: map[string]interface{}{
			"elements": parseTextElements(content, openBase),
		},
	}
}

func buildTodoBlock(content string, done bool, openBase string) map[string]interface{} {
	return map[string]interface{}{
		"block_type": 17,
		"todo": map[string]interface{}{
			"style": map[string]interface{}{
				"done": done,
			},
			"elements": parseTextElements(content, openBase),
		},
	}
}

func parseTextElements(content, openBase string) []interface{} {
	if content == "" {
		return []interface{}{parseTextElementRun("", "")}
	}

	elements := make([]interface{}, 0, 8)
	cursor := 0
	for cursor < len(content) {
		start, end, kind, groups := findNextInlineElement(content, cursor)
		if start < 0 {
			elements = append(elements, parseTextElementRun(content[cursor:], ""))
			break
		}
		if start > cursor {
			elements = append(elements, parseTextElementRun(content[cursor:start], ""))
		}

		switch kind {
		case "link":
			if len(groups) >= 2 && strings.TrimSpace(groups[0]) != "" && strings.TrimSpace(groups[1]) != "" {
				elements = append(elements, parseTextElementRun(groups[0], groups[1]))
			} else {
				elements = append(elements, parseTextElementRun(content[start:end], ""))
			}
		case "mention_user":
			if len(groups) >= 1 && strings.TrimSpace(groups[0]) != "" {
				elements = append(elements, map[string]interface{}{
					"mention_user": map[string]interface{}{
						"user_id": strings.TrimSpace(groups[0]),
					},
				})
			} else {
				elements = append(elements, parseTextElementRun(content[start:end], ""))
			}
		case "mention_doc":
			if len(groups) >= 1 {
				if elem, ok := parseMentionDocElement(groups[0], openBase); ok {
					elements = append(elements, elem)
				} else {
					elements = append(elements, parseTextElementRun(content[start:end], ""))
				}
			}
		default:
			elements = append(elements, parseTextElementRun(content[start:end], ""))
		}
		cursor = end
	}
	if len(elements) == 0 {
		return []interface{}{parseTextElementRun(content, "")}
	}
	return elements
}

func parseTextElementRun(content, linkTarget string) map[string]interface{} {
	run := map[string]interface{}{
		"content": content,
	}
	if strings.TrimSpace(linkTarget) != "" {
		run["text_element_style"] = map[string]interface{}{
			"link": map[string]interface{}{
				"url": url.QueryEscape(linkTarget),
			},
		}
	}
	return map[string]interface{}{"text_run": run}
}

func findNextInlineElement(content string, cursor int) (start, end int, kind string, groups []string) {
	start, end = -1, -1
	search := content[cursor:]
	candidates := []struct {
		re   *regexp.Regexp
		kind string
	}{
		{re: markdownLinkPattern, kind: "link"},
		{re: mentionUserPattern, kind: "mention_user"},
		{re: mentionDocPattern, kind: "mention_doc"},
	}
	for _, candidate := range candidates {
		loc := candidate.re.FindStringSubmatchIndex(search)
		if len(loc) == 0 {
			continue
		}
		candidateStart := cursor + loc[0]
		candidateEnd := cursor + loc[1]
		if start == -1 || candidateStart < start {
			start, end, kind = candidateStart, candidateEnd, candidate.kind
			submatches := candidate.re.FindStringSubmatch(search)
			groups = groups[:0]
			if len(submatches) > 1 {
				groups = append(groups, submatches[1:]...)
			}
		}
	}
	return start, end, kind, groups
}

func parseMentionDocElement(raw, openBase string) (map[string]interface{}, bool) {
	ref, canonicalURL, objType, ok := resolveMentionDoc(raw, openBase)
	if !ok {
		return nil, false
	}
	return map[string]interface{}{
		"mention_doc": map[string]interface{}{
			"token":    ref.Token,
			"obj_type": objType,
			"url":      url.QueryEscape(canonicalURL),
		},
	}, true
}

func resolveMentionDoc(raw, openBase string) (documentRef, string, int, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return documentRef{}, "", 0, false
	}
	if strings.Contains(raw, "://") {
		ref, err := parseDocumentRef(raw)
		if err != nil {
			return documentRef{}, "", 0, false
		}
		objType, ok := mentionObjType(ref.Kind)
		if !ok {
			return documentRef{}, "", 0, false
		}
		return ref, raw, objType, true
	}

	ref := inferMentionDocRef(raw)
	objType, ok := mentionObjType(ref.Kind)
	if !ok {
		return documentRef{}, "", 0, false
	}
	return ref, mentionDocURL(ref, openBase), objType, true
}

func inferMentionDocRef(raw string) documentRef {
	switch {
	case strings.HasPrefix(raw, "dox"):
		return documentRef{Kind: "docx", Token: raw}
	case strings.HasPrefix(raw, "doc"):
		return documentRef{Kind: "doc", Token: raw}
	case strings.HasPrefix(raw, "shtr"), strings.HasPrefix(raw, "sht"), strings.HasPrefix(raw, "shr"):
		return documentRef{Kind: "sheet", Token: raw}
	case strings.HasPrefix(raw, "bas"):
		return documentRef{Kind: "bitable", Token: raw}
	default:
		return documentRef{Kind: "docx", Token: raw}
	}
}

func mentionObjType(kind string) (int, bool) {
	switch kind {
	case "doc":
		return 1, true
	case "sheet":
		return 3, true
	case "bitable":
		return 8, true
	case "wiki":
		return 16, true
	case "docx":
		return 22, true
	default:
		return 0, false
	}
}

func mentionDocURL(ref documentRef, openBase string) string {
	base := strings.TrimRight(openBase, "/")
	switch ref.Kind {
	case "doc":
		return base + "/doc/" + ref.Token
	case "sheet":
		return base + "/sheets/" + ref.Token
	case "bitable":
		return base + "/base/" + ref.Token
	case "wiki":
		return base + "/wiki/" + ref.Token
	default:
		return base + "/docx/" + ref.Token
	}
}

func openBaseURL(runtime *common.RuntimeContext) string {
	return strings.TrimRight(core.ResolveOpenBaseURL(runtime.Config.Brand), "/")
}

func extractBlockText(block map[string]interface{}) string {
	type prefixRule struct {
		key    string
		prefix string
		suffix string
	}
	for _, rule := range []prefixRule{
		{key: "page"},
		{key: "text"},
		{key: "heading1", prefix: "# "},
		{key: "heading2", prefix: "## "},
		{key: "heading3", prefix: "### "},
		{key: "bullet", prefix: "- "},
		{key: "ordered", prefix: "1. "},
		{key: "quote", prefix: "> "},
		{key: "code", prefix: "```text\n", suffix: "\n```"},
	} {
		if node, ok := block[rule.key].(map[string]interface{}); ok {
			if text := extractElementsText(node["elements"]); text != "" {
				return rule.prefix + text + rule.suffix
			}
		}
	}
	if node, ok := block["todo"].(map[string]interface{}); ok {
		text := extractElementsText(node["elements"])
		if text == "" {
			return ""
		}
		done := false
		if style, ok := node["style"].(map[string]interface{}); ok {
			done, _ = style["done"].(bool)
		}
		if done {
			return "- [x] " + text
		}
		return "- [ ] " + text
	}
	if _, ok := block["divider"].(map[string]interface{}); ok {
		return "---"
	}
	return ""
}

func extractElementsText(raw interface{}) string {
	elements, _ := raw.([]interface{})
	parts := make([]string, 0, len(elements))
	for _, item := range elements {
		elem, _ := item.(map[string]interface{})
		if elem == nil {
			continue
		}
		if run, ok := elem["text_run"].(map[string]interface{}); ok {
			if content, ok := run["content"].(string); ok && content != "" {
				if style, ok := run["text_element_style"].(map[string]interface{}); ok {
					if link, ok := style["link"].(map[string]interface{}); ok {
						if encodedURL, ok := link["url"].(string); ok && encodedURL != "" {
							decodedURL := encodedURL
							if rawURL, err := url.QueryUnescape(encodedURL); err == nil {
								decodedURL = rawURL
							}
							parts = append(parts, fmt.Sprintf("[%s](%s)", content, decodedURL))
							continue
						}
					}
				}
				parts = append(parts, content)
			}
		}
		if mention, ok := elem["mention_user"].(map[string]interface{}); ok {
			if userID, ok := mention["user_id"].(string); ok && userID != "" {
				parts = append(parts, "@user["+userID+"]")
			}
		}
		if mention, ok := elem["mention_doc"].(map[string]interface{}); ok {
			if encodedURL, ok := mention["url"].(string); ok && encodedURL != "" {
				decodedURL := encodedURL
				if rawURL, err := url.QueryUnescape(encodedURL); err == nil {
					decodedURL = rawURL
				}
				parts = append(parts, "@doc["+decodedURL+"]")
				continue
			}
			if token, ok := mention["token"].(string); ok && token != "" {
				parts = append(parts, "@doc["+token+"]")
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}
