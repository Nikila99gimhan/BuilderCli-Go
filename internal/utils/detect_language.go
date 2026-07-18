package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// manifestFiles maps the presence of a specific file to the language it signals.
// These are checked FIRST and take priority over extension counting.
var manifestFiles = map[string]string{
	"go.mod":            "Go",
	"package.json":      "JavaScript",
	"pom.xml":           "Java",
	"build.gradle":      "Java",
	"build.gradle.kts":  "Java",
	"requirements.txt":  "Python",
	"Pipfile":           "Python",
	"pyproject.toml":    "Python",
	"Cargo.toml":        "Rust",
	"Gemfile":           "Ruby",
	"composer.json":     "PHP",
}

// knownLanguageExtensions is the fallback when no manifest file is found.
var knownLanguageExtensions = map[string]string{
	".py":    "Python",
	".go":    "Go",
	".java":  "Java",
	".js":    "JavaScript",
	".ts":    "JavaScript",
	".cpp":   "C++",
	".h":     "C",
	".c":     "C",
	".cs":    "C#",
	".php":   "PHP",
	".rb":    "Ruby",
	".rs":    "Rust",
	".swift": "Swift",
}

// DetectLanguage first checks for well-known project manifest files (e.g., go.mod,
// pom.xml, package.json). If none are found, it falls back to counting file
// extensions to find the most common language.
// Returns "unknown" if no language can be determined.
func DetectLanguage(sourceCodeDir string) (string, map[string]int) {
	// Step 1: Check for high-confidence manifest files.
	for manifest, language := range manifestFiles {
		if _, err := os.Stat(filepath.Join(sourceCodeDir, manifest)); err == nil {
			return language, map[string]int{language: 1}
		}
	}

	// Step 2: Fallback — count file extensions across the whole directory tree.
	languageCounts := map[string]int{}
	err := filepath.Walk(sourceCodeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			ext := filepath.Ext(path)
			if language, ok := knownLanguageExtensions[ext]; ok {
				languageCounts[language]++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return "unknown", map[string]int{}
	}

	var mostCommonLanguage string
	var maxCount int
	for language, count := range languageCounts {
		if count > maxCount {
			mostCommonLanguage = language
			maxCount = count
		}
	}

	if mostCommonLanguage == "" {
		return "unknown", languageCounts
	}
	return mostCommonLanguage, languageCounts
}
