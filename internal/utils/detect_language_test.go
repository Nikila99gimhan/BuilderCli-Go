package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// createTempDir creates a temporary directory with the given files and returns its path.
func createTempDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write file %s: %v", name, err)
		}
	}
	return dir
}

func TestDetectLanguage_GoManifest(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"go.mod":        "module example.com/app\n\ngo 1.21\n",
		"main.go":       "package main",
		"index.html":    "<html></html>",
		"style.css":     "body {}",
		"config.xml":    "<config/>",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "Go" {
		t.Errorf("expected Go, got %q (manifest-based detection should win over extension counting)", lang)
	}
}

func TestDetectLanguage_JavaManifest(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"pom.xml":              "<project/>",
		"src/main/App.java":    "public class App {}",
		"resources/index.html": "<html></html>",
		"resources/style.css":  "body {}",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "Java" {
		t.Errorf("expected Java, got %q", lang)
	}
}

func TestDetectLanguage_NodeManifest(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"package.json": `{"name":"app","version":"1.0.0"}`,
		"index.js":     "console.log('hello')",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "JavaScript" {
		t.Errorf("expected JavaScript, got %q", lang)
	}
}

func TestDetectLanguage_PythonManifest(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"requirements.txt": "flask==2.3.0\n",
		"app.py":           "print('hello')",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "Python" {
		t.Errorf("expected Python, got %q", lang)
	}
}

func TestDetectLanguage_FallbackExtension(t *testing.T) {
	// No manifest files — falls back to extension counting.
	dir := createTempDir(t, map[string]string{
		"main.rs":   "fn main() {}",
		"lib.rs":    "pub fn lib() {}",
		"helper.rs": "pub fn helper() {}",
		"README.md": "# Rust app",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "Rust" {
		t.Errorf("expected Rust (extension fallback), got %q", lang)
	}
}

func TestDetectLanguage_Unknown(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"README.md":  "# Generic project",
		"Makefile":   "all: build",
	})
	lang, _ := DetectLanguage(dir)
	if lang != "unknown" {
		t.Errorf("expected unknown, got %q", lang)
	}
}
