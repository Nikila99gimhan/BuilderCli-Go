package imagebuilder

import (
	"os"
	"os/exec"
)

type ImageBuilderImpl struct {
	language string
}

func NewImageBuilderImpl(language string) *ImageBuilderImpl {
	return &ImageBuilderImpl{language: language}
}

func (ib *ImageBuilderImpl) SelectBuilder() string {
	switch ib.language {
	case "Java":
		return "paketobuildpacks/builder-jammy-base"
	case "Go":
		return "paketobuildpacks/builder-jammy-base"
	case "JavaScript":
		return "paketobuildpacks/builder-jammy-base"
	case "Python":

		return "gcr.io/buildpacks/builder:v1"
	default:
		return ""
	}
}

func (ib *ImageBuilderImpl) Build(imageName string) error {
	builder := ib.SelectBuilder()
	cmd := exec.Command("pack", "build", imageName, "--builder", builder, "--pull-policy", "if-not-present")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (ib *ImageBuilderImpl) TagImage(imageName, imageTag string) error {
	cmd := exec.Command("docker", "tag", imageName, imageName+":"+imageTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
