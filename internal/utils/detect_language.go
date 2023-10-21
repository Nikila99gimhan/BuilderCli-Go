package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectLanguage function from the provided code
func DetectLanguage(sourceCodeDir string) (string, map[string]int) {
	languageCounts := map[string]int{}

	err := filepath.Walk(sourceCodeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			ext := filepath.Ext(path)
			language, ok := knownLanguageExtensions[ext]
			if ok {
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

	return mostCommonLanguage, languageCounts
}

var knownLanguageExtensions = map[string]string{
	".py":    "Python",
	".go":    "Go",
	".java":  "Java",
	".js":    "JavaScript",
	".cpp":   "C++",
	".h":     "C",
	".c":     "C",
	".cs":    "C#",
	".php":   "PHP",
	".rb":    "Ruby",
	".rs":    "Rust",
	".swift": "Swift",
	".json":  "JSON",
	".xml":   "XML",
	".html":  "HTML",
	".css":   "CSS",
	".scss":  "SCSS",
}
