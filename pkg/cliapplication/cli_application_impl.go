package cliapplication

import (
	"cliapp/internal/utils"
	"cliapp/pkg/imagebuilder"
	"cliapp/pkg/imagepusher"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type CliApplicationImpl struct {
	language   string
	repoURL    string
	display    *Display
	imgBuilder imagebuilder.ImageBuilder
}

func NewCliApplicationImpl() *CliApplicationImpl {
	return &CliApplicationImpl{
		display: NewDisplay(),
	}
}

func (app *CliApplicationImpl) Start() {
	app.display.PrintLogo()
	var response string

	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: []string{"Clone a new repository", "Select from current repositories"},
	}
	survey.AskOne(prompt, &response, nil)

	switch response {
	case "Clone a new repository":
		app.display.Print("Enter the repository URL to clone:")
		app.repoURL = utils.GetInput()
		if err := utils.Clone(app.repoURL); err != nil {
			app.display.Print("Error cloning repository: " + err.Error())
			return
		}
	case "Select from current repositories":
		repos, err := utils.ListDirectories()
		if err != nil {
			app.display.Print("Error fetching repositories: " + err.Error())
			return
		}
		prompt := &survey.Select{
			Message: "Select a repository:",
			Options: repos,
		}
		survey.AskOne(prompt, &app.repoURL, nil)
	}

	if app.repoURL == "" {
		fmt.Println("No repository URL provided.")
		return
	}

	app.display.PrintGreenBold("\nDetecting repository...\n")
	repoName := utils.GetRepoNameFromURL(app.repoURL)
	app.language, _ = utils.DetectLanguage(repoName)
	if app.language == "unknown" {
		app.display.Print("Failed to auto-detect the programming language. Please specify it manually: ")
		app.language = utils.GetInput()
	}

	fmt.Printf("Using language: %s\n", app.language)

	app.imgBuilder = imagebuilder.NewImageBuilderImpl(app.language)

	app.display.Print("\nEnter the name for your image:")
	imageName := utils.GetInput()

	app.display.PrintGreenBold("\nSelecting appropriate builder...\n")
	if builder := app.imgBuilder.SelectBuilder(); builder == "" {
		fmt.Println("Unsupported language:", app.language)
		return
	}

	app.display.PrintGreenBold("\nBuilding and containerizing application...\n")
	if err := app.imgBuilder.Build(imageName); err != nil {
		fmt.Println("Error during build process:", err)
		return
	}

	app.display.PrintGreenBold("\nEnter a tag for your image:")
	imageTag := utils.GetInput()

	if err := app.imgBuilder.TagImage(imageName, imageTag); err != nil {
		fmt.Println("Error during tagging process:", err)
		return
	}

	app.display.PrintGreenBold("\nBuild, containerization, and tagging complete!\n")

	prompt = &survey.Select{
		Message: "Do you want to push the image to a registry?",
		Options: []string{"Yes", "No"},
	}

	survey.AskOne(prompt, &response, nil)

	if response == "No" {
		app.display.PrintGreenBold("\nThank you for using the application!\n")
		return
	}

	prompt = &survey.Select{
		Message: "Are you already logged in to the registry?",
		Options: []string{"Yes", "No"},
	}

	survey.AskOne(prompt, &response, nil)

	var imagePusher imagepusher.ImagePusher = imagepusher.NewImagePusherImpl()
	var username string

	if response == "Yes" {
		app.display.Print("\nEnter the DockerHub username:")
		username = utils.GetInput()
	}

	if response == "No" {
		app.display.Print("\nEnter the DockerHub username:")
		username = utils.GetInput()
		app.display.Print("\nEnter the DockerHub password:")
		password := utils.GetInput()

		// Login to DockerHub
		err := imagePusher.LoginToRegistry(username, password, "")
		if err != nil {
			app.display.Print("Error logging in to registry: " + err.Error())
			return
		}
	}

	// ReTag the image
	newImageName, err := imagePusher.ReTagImage(username, imageName+":"+imageTag)
	if err != nil {
		return
	}

	// Push the image
	err = imagePusher.PushImage(newImageName)
	if err != nil {
		return
	}

	app.display.PrintGreenBold("\nThank you for using the application!\n")
}
