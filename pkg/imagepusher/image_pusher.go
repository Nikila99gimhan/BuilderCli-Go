package imagepusher

type ImagePusher interface {
	LoginToRegistry(username, password, repository string) error
	ReTagImage(repositoryName, imageName string) (string, error)
	PushImage(imageName string) error
}
