package imagepusher

import (
	"os"
	"os/exec"
)

type ImagePusherImpl struct{}

func NewImagePusherImpl() *ImagePusherImpl {
	return &ImagePusherImpl{}
}

func (p *ImagePusherImpl) LoginToRegistry(username, password, repository string) error {
	cmd := exec.Command("docker", "login", "--username", username, "--password", password, repository)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *ImagePusherImpl) ReTagImage(repositoryName, imageName string) (string, error) {
	newImageName := repositoryName + "/" + imageName

	cmd := exec.Command("docker", "tag", imageName, newImageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return newImageName, nil
}

func (p *ImagePusherImpl) PushImage(imageName string) error {
	cmd := exec.Command("docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
