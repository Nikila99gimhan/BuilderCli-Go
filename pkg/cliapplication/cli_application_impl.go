package cliapplication

import (
	"cliapp/internal/utils"
	"cliapp/pkg/imagebuilder"
	"cliapp/pkg/imagepusher"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CliApplicationImpl holds the runtime state and dependencies of the CLI.
// Dependencies (imgBuilder, imgPusher) are injected from main.go to allow
// mocking in tests.
type CliApplicationImpl struct {
	language   string
	repoURL    string
	display    *Display
	imgBuilder imagebuilder.ImageBuilder
	imgPusher  imagepusher.ImagePusher

	// Non-interactive flags (set by cobra)
	flagRepo     string
	flagImage    string
	flagTag      string
	flagPush     bool
	flagUsername string
	flagPassword string
}

func NewCliApplicationImpl() *CliApplicationImpl {
	return &CliApplicationImpl{
		display:   NewDisplay(),
		imgPusher: imagepusher.NewDockerPusher(),
	}
}

// Start is the entry point. It builds the cobra command tree and executes it.
// It returns a non-nil error if the workflow fails, so main.go can exit(1).
func (app *CliApplicationImpl) Start() error {
	rootCmd := &cobra.Command{
		Use:   "buildercli",
		Short: "BuilderCLI — Clone, containerize, and push your applications",
		Long: `BuilderCLI simplifies the process of cloning a Git repository,
auto-detecting its language, building a container image using Cloud Native
Buildpacks, and pushing it to a container registry.`,
		RunE: app.run,
	}

	// --- Non-interactive flags ---
	rootCmd.Flags().StringVarP(&app.flagRepo, "repo", "r", "", "Git repository URL to clone (skips prompt)")
	rootCmd.Flags().StringVarP(&app.flagImage, "image", "i", "", "Name for the container image (skips prompt)")
	rootCmd.Flags().StringVarP(&app.flagTag, "tag", "t", "", "Tag for the container image (skips prompt)")
	rootCmd.Flags().BoolVar(&app.flagPush, "push", false, "Automatically push the image after building (skips prompt)")
	rootCmd.Flags().StringVarP(&app.flagUsername, "username", "u", "", "Registry username (skips prompt)")
	rootCmd.Flags().StringVarP(&app.flagPassword, "password", "p", "", "Registry password (use env var REGISTRY_PASSWORD in CI instead)")

	return rootCmd.Execute()
}

// run contains the full build+push workflow. It is invoked by cobra.
func (app *CliApplicationImpl) run(cmd *cobra.Command, args []string) error {
	app.display.PrintLogo()

	// ── Step 1: Resolve repository ─────────────────────────────────────────
	if err := app.resolveRepo(); err != nil {
		return fmt.Errorf("repository step failed: %w", err)
	}

	// ── Step 2: Detect language ────────────────────────────────────────────
	if err := app.resolveLanguage(); err != nil {
		return fmt.Errorf("language detection failed: %w", err)
	}

	// ── Step 3: Build ──────────────────────────────────────────────────────
	app.imgBuilder = imagebuilder.NewPackBuilder(app.language)

	imageName, err := app.resolveImageName()
	if err != nil {
		return fmt.Errorf("image name step failed: %w", err)
	}

	app.display.PrintGreenBold("\nSelecting appropriate builder...")
	if builder := app.imgBuilder.SelectBuilder(); builder == "" {
		return fmt.Errorf("unsupported language: %q — please specify a builder manually", app.language)
	}

	app.display.PrintGreenBold("\nBuilding and containerizing application...")
	if err := app.imgBuilder.Build(imageName); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	imageTag, err := app.resolveImageTag()
	if err != nil {
		return fmt.Errorf("image tag step failed: %w", err)
	}

	if err := app.imgBuilder.TagImage(imageName, imageTag); err != nil {
		return fmt.Errorf("tagging failed: %w", err)
	}

	app.display.PrintGreenBold("\nBuild, containerization, and tagging complete!")

	// ── Step 4: Push (optional) ────────────────────────────────────────────
	shouldPush, err := app.resolvePush()
	if err != nil {
		return err
	}
	if !shouldPush {
		app.display.PrintGreenBold("\nThank you for using BuilderCLI!")
		return nil
	}

	if err := app.pushImage(imageName, imageTag); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	app.display.PrintGreenBold("\nThank you for using BuilderCLI!")
	return nil
}

// ── Resolver helpers ────────────────────────────────────────────────────────

func (app *CliApplicationImpl) resolveRepo() error {
	if app.flagRepo != "" {
		app.repoURL = app.flagRepo
		return utils.Clone(app.repoURL)
	}

	var response string
	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: []string{"Clone a new repository", "Select from current repositories"},
	}
	if err := survey.AskOne(prompt, &response, nil); err != nil {
		return err
	}

	switch response {
	case "Clone a new repository":
		app.display.Print("Enter the repository URL to clone:")
		app.repoURL = utils.GetInput()
		if app.repoURL == "" {
			return errors.New("no repository URL provided")
		}
		return utils.Clone(app.repoURL)

	case "Select from current repositories":
		repos, err := utils.ListDirectories()
		if err != nil {
			return fmt.Errorf("could not list local repositories: %w", err)
		}
		if len(repos) == 0 {
			return errors.New("no local directories found — try cloning a repository first")
		}
		repoSelect := &survey.Select{
			Message: "Select a repository:",
			Options: repos,
		}
		return survey.AskOne(repoSelect, &app.repoURL, nil)
	}
	return nil
}

func (app *CliApplicationImpl) resolveLanguage() error {
	app.display.PrintGreenBold("\nDetecting repository language...")
	repoName := utils.GetRepoNameFromURL(app.repoURL)
	app.language, _ = utils.DetectLanguage(repoName)

	if app.language == "unknown" || app.language == "" {
		app.display.Print("Could not auto-detect the programming language. Please specify it manually:")
		app.language = utils.GetInput()
		if app.language == "" {
			return errors.New("language cannot be empty")
		}
	}
	fmt.Printf("Detected language: %s\n", app.language)
	return nil
}

func (app *CliApplicationImpl) resolveImageName() (string, error) {
	if app.flagImage != "" {
		return app.flagImage, nil
	}
	app.display.Print("\nEnter the name for your image:")
	name := utils.GetInput()
	if name == "" {
		return "", errors.New("image name cannot be empty")
	}
	return name, nil
}

func (app *CliApplicationImpl) resolveImageTag() (string, error) {
	if app.flagTag != "" {
		return app.flagTag, nil
	}
	app.display.PrintGreenBold("\nEnter a tag for your image:")
	tag := utils.GetInput()
	if tag == "" {
		return "", errors.New("image tag cannot be empty")
	}
	return tag, nil
}

func (app *CliApplicationImpl) resolvePush() (bool, error) {
	if app.flagPush {
		return true, nil
	}
	var response string
	prompt := &survey.Select{
		Message: "Do you want to push the image to a registry?",
		Options: []string{"Yes", "No"},
	}
	if err := survey.AskOne(prompt, &response, nil); err != nil {
		return false, err
	}
	return response == "Yes", nil
}

func (app *CliApplicationImpl) pushImage(imageName, imageTag string) error {
	username, password, err := app.resolveRegistryCredentials()
	if err != nil {
		return err
	}

	// If password is provided, login first.
	if password != "" {
		app.display.PrintGreenBold("\nLogging in to registry...")
		if err := app.imgPusher.LoginToRegistry(username, password, ""); err != nil {
			return fmt.Errorf("registry login failed: %w", err)
		}
	}

	newImageName, err := app.imgPusher.ReTagImage(username, imageName+":"+imageTag)
	if err != nil {
		return fmt.Errorf("re-tag failed: %w", err)
	}

	app.display.PrintGreenBold("\nPushing image...")
	if err := app.imgPusher.PushImage(newImageName); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	return nil
}

func (app *CliApplicationImpl) resolveRegistryCredentials() (string, string, error) {
	// Non-interactive: use flags / env var (REGISTRY_PASSWORD is recommended in CI).
	if app.flagUsername != "" {
		return app.flagUsername, app.flagPassword, nil
	}

	// Interactive: ask if already logged in.
	var loggedIn string
	loginPrompt := &survey.Select{
		Message: "Are you already logged in to the registry?",
		Options: []string{"Yes", "No"},
	}
	if err := survey.AskOne(loginPrompt, &loggedIn, nil); err != nil {
		return "", "", err
	}

	app.display.Print("\nEnter the DockerHub username:")
	username := utils.GetInput()

	if loggedIn == "Yes" {
		return username, "", nil
	}

	// Secure password prompt — password will NOT echo to the terminal.
	app.display.Print("\nEnter the DockerHub password:")
	password, err := utils.GetSecretInput()
	if err != nil {
		return "", "", fmt.Errorf("could not read password: %w", err)
	}
	return username, password, nil
}
