package tenant

import "testing"

func TestNormalizeSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{" TokoBunga ", "tokobunga"},
		{"TOKOBUNGA", "tokobunga"},
		{"toko-bunga", "toko-bunga"},
	}

	for _, tt := range tests {
		if got := normalizeSlug(tt.input); got != tt.want {
			t.Fatalf("normalizeSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlugPattern(t *testing.T) {
	valid := []string{"toko", "toko-bunga", "a1"}
	invalid := []string{"", "-toko", "toko-", "toko_bunga"}

	for _, s := range valid {
		if !slugPattern.MatchString(s) {
			t.Fatalf("expected valid slug: %s", s)
		}
	}
	for _, s := range invalid {
		if slugPattern.MatchString(s) {
			t.Fatalf("expected invalid slug: %s", s)
		}
	}
}

func TestBuildSubdomain(t *testing.T) {
	svc := NewService(nil, "localhost")
	if got := svc.BuildSubdomain("tokobunga"); got != "tokobunga.localhost" {
		t.Fatalf("BuildSubdomain() = %q", got)
	}
}
