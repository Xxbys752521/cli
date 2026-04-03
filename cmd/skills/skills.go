// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

const skillsDirEnv = "XFCHAT_CLI_SKILLS_DIR"

// NewCmdSkills exposes bundled skill discovery helpers.
func NewCmdSkills() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Inspect bundled AI agent skills",
	}
	cmd.AddCommand(newCmdSkillsPath())
	cmd.AddCommand(newCmdSkillsList())
	return cmd
}

func newCmdSkillsPath() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Print bundled skills directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := bundledSkillsDir()
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(map[string]interface{}{"path": dir})
			}
			fmt.Fprintln(cmd.OutOrStdout(), dir)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "structured JSON output")
	return cmd
}

func newCmdSkillsList() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bundled skill names",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := bundledSkillsDir()
			if err != nil {
				return err
			}
			names, err := listBundledSkills(dir)
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(map[string]interface{}{
					"path":   dir,
					"skills": names,
				})
			}
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "structured JSON output")
	return cmd
}

func bundledSkillsDir() (string, error) {
	if dir := os.Getenv(skillsDirEnv); dir != "" {
		if resolved, ok := normalizeSkillsDir(dir); ok {
			return resolved, nil
		}
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to locate executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(execPath); err == nil && resolved != "" {
		execPath = resolved
	}
	execDir := filepath.Dir(execPath)
	candidates := []string{
		filepath.Join(execDir, "skills"),
		filepath.Join(execDir, "..", "skills"),
	}
	for _, candidate := range candidates {
		if resolved, ok := normalizeSkillsDir(candidate); ok {
			return resolved, nil
		}
	}
	return "", fmt.Errorf("bundled skills not found; expected a skills directory next to the xfchat_cli package")
}

func listBundledSkills(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read skills dir: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "SKILL.md")); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func normalizeSkillsDir(dir string) (string, bool) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	names, err := listBundledSkills(abs)
	if err != nil || len(names) == 0 {
		return "", false
	}
	return abs, true
}

func writeJSON(value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
