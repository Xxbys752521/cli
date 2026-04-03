// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package skills

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBundledSkillsDirUsesEnvOverride(t *testing.T) {
	t.Setenv(skillsDirEnv, createSkillsFixture(t))
	dir, err := bundledSkillsDir()
	if err != nil {
		t.Fatalf("bundledSkillsDir err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "lark-calendar", "SKILL.md")); err != nil {
		t.Fatalf("expected bundled skill fixture, err=%v", err)
	}
}

func TestListBundledSkills(t *testing.T) {
	dir := createSkillsFixture(t)
	names, err := listBundledSkills(dir)
	if err != nil {
		t.Fatalf("listBundledSkills err=%v", err)
	}
	want := []string{"lark-calendar", "lark-im"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("names=%v want=%v", names, want)
	}
}

func createSkillsFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"lark-im", "lark-calendar"} {
		skillDir := filepath.Join(dir, name)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", skillDir, err)
		}
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# test\n"), 0o644); err != nil {
			t.Fatalf("write SKILL.md: %v", err)
		}
	}
	return dir
}
