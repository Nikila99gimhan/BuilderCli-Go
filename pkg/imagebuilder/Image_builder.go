package imagebuilder

type ImageBuilder interface {
	SelectBuilder() string
	Build(imageName string) error
	TagImage(imageName, imageTag string) error
}
