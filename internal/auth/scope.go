// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package auth

import "strings"

// scopeAliases maps scope names that are equivalent across environments.
// Key: scope name used by shortcuts, Value: alternative scope that satisfies the same requirement.
var scopeAliases = map[string]string{
	"contact:user.basic_profile:readonly": "contact:user.base:readonly",
	// calendar:calendar:readonly is a coarse-grained scope in private deployments
	// that covers all calendar read operations.
	"calendar:calendar.event:read":    "calendar:calendar:readonly",
	"calendar:calendar.free_busy:read": "calendar:calendar:readonly",
}

// MissingScopes returns the elements of required that are absent from storedScope.
// storedScope is a space-separated list of granted scope strings (as stored in the token).
// A granted scope is considered to cover a required scope if:
//   - exact match, OR
//   - alias match (e.g. "contact:user.base:readonly" covers "contact:user.basic_profile:readonly"), OR
//   - the granted scope is a prefix parent (e.g. "im:message" covers "im:message.group_msg:get_as_user")
func MissingScopes(storedScope string, required []string) []string {
	grantedList := strings.Fields(storedScope)
	granted := make(map[string]bool, len(grantedList))
	for _, s := range grantedList {
		granted[s] = true
	}
	var missing []string
	for _, req := range required {
		if granted[req] {
			continue
		}
		if alias, ok := scopeAliases[req]; ok && granted[alias] {
			continue
		}
		if coveredByParent(grantedList, req) {
			continue
		}
		missing = append(missing, req)
	}
	return missing
}

// coveredByParent checks if any granted scope is a parent prefix of req.
// e.g. granted "im:message" covers required "im:message.group_msg:get_as_user"
// because "im:message" is a prefix followed by "." or ":".
func coveredByParent(granted []string, req string) bool {
	for _, g := range granted {
		if len(g) < len(req) && strings.HasPrefix(req, g) {
			// Ensure it's a real parent boundary (next char is '.' or ':')
			next := req[len(g)]
			if next == '.' || next == ':' {
				return true
			}
		}
	}
	return false
}
