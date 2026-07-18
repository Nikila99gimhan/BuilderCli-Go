package imagepusher

import (
	"os"
	"os/exec"
)

// DockerPusher implements ImagePusher using the `docker` CLI.
type DockerPusher struct{}

// NewDockerPusher creates a new DockerPusher.
func NewDockerPusher() *DockerPusher {
	return &DockerPusher{}
}

// LoginToRegistry runs `docker login`. Pass an empty repository to default to DockerHub.
func (p *DockerPusher) LoginToRegistry(username, password, repository string) error {
	args := []string{"login", "--username", username, "--password-stdin"}
	if repository != "" {
		args = append(args, repository)
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdin = newStringReader(password)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ReTagImage creates a registry-qualified tag (e.g. username/imagename:tag).
func (p *DockerPusher) ReTagImage(repositoryName, imageName string) (string, error) {
	newImageName := repositoryName + "/" + imageName
	cmd := exec.Command("docker", "tag", imageName, newImageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return newImageName, nil
}

// PushImage runs `docker push` for the given fully-qualified image name.
func (p *DockerPusher) PushImage(imageName string) error {
	cmd := exec.Command("docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
