package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeJoin(t *testing.T) {
	base := "/tmp/out"

	rejected := []struct {
		name string
		path string
	}{
		{"parent escape", "../evil"},
		{"deep parent escape", "../../../../etc/passwd"},
		{"mid-path parent", "a/../../evil"},
		{"absolute unix", "/etc/passwd"},
		{"windows backslash escape", "..\\..\\evil"},
		{"empty", ""},
		{"dot", "."},
		{"dotdot", ".."},
	}
	for _, tc := range rejected {
		t.Run("reject/"+tc.name, func(t *testing.T) {
			if _, err := safeJoin(base, tc.path); err == nil {
				t.Fatalf("safeJoin(%q, %q) = nil error, want rejection", base, tc.path)
			}
		})
	}

	accepted := []struct {
		name string
		path string
	}{
		{"plain file", "shot.png"},
		{"nested", "docs/arch.png"},
		{"deep nested", "a/b/c/d.png"},
		{"leading dot file", ".keep"},
	}
	for _, tc := range accepted {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			got, err := safeJoin(base, tc.path)
			if err != nil {
				t.Fatalf("safeJoin(%q, %q) unexpected error: %v", base, tc.path, err)
			}
			want := filepath.Join(base, tc.path)
			if got != want {
				t.Fatalf("safeJoin(%q, %q) = %q, want %q", base, tc.path, got, want)
			}
			// Belt-and-suspenders: result must be under base.
			if !strings.HasPrefix(got, base) {
				t.Fatalf("result %q not under base %q", got, base)
			}
		})
	}
}
