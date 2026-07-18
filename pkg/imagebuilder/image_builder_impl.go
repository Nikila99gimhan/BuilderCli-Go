package imagebuilder

import (
	"os"
	"os/exec"
)

// PackBuilder implements ImageBuilder using the `pack` CLI (Cloud Native Buildpacks).
type PackBuilder struct {
	language string
}

// NewPackBuilder creates a new PackBuilder for the given detected language.
func NewPackBuilder(language string) *PackBuilder {
	return &PackBuilder{language: language}
}

// SelectBuilder returns the appropriate buildpack builder image for the detected language.
func (pb *PackBuilder) SelectBuilder() string {
	switch pb.language {
	case "Java":
		return "paketobuildpacks/builder-jammy-base"
	case "Go":
		return "paketobuildpacks/builder-jammy-base"
	case "JavaScript":
		return "paketobuildpacks/builder-jammy-base"
	case "Python":
		return "gcr.io/buildpacks/builder:v1"
	case "Rust":
		return "paketobuildpacks/builder-jammy-base"
	default:
		return ""
	}
}

// Build runs `pack build` using the selected builder image.
func (pb *PackBuilder) Build(imageName string) error {
	builder := pb.SelectBuilder()
	cmd := exec.Command("pack", "build", imageName, "--builder", builder, "--pull-policy", "if-not-present")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TagImage tags an existing local image with `docker tag`.
func (pb *PackBuilder) TagImage(imageName, imageTag string) error {
	cmd := exec.Command("docker", "tag", imageName, imageName+":"+imageTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
