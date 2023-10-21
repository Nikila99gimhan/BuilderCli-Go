# MyCLIApp - Containerized App Builder

A CLI tool to simplify the process of cloning, building, and containerizing applications using Cloud Native Buildpacks.

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Technical Details](#technical-details)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

- Go (version 1.15 or newer)
- Git
- Docker or a compatible container runtime

### Installation

1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```
2. Navigate to the project directory:
   ```bash
   cd path/to/project
   ```
3. Build the application:
   ```bash
   go build -o mycliapp
   ```

## Usage

```bash
./mycliapp [repository_url]
```

**Example**:

```bash
./mycliapp https://github.com/user/sample-app.git
```

## Technical Details

- **Modular Architecture**: The application is designed using Object-Oriented Principles in Go, with distinct modules for displaying content (`Display`), managing repositories (`RepoManager`), and building & tagging container images (`ImageBuilder`).

- **Repository Management (`RepoManager`)**: Handles git operations, including cloning the specified repository.

- **Image Building (`ImageBuilder`)**: Manages the process of building container images using buildpacks and optionally tagging them. Uses the `pack` CLI under the hood.

- **Display (`Display`)**: A module dedicated to handling all user interface interactions, including printing messages and logos.

- **Cache Optimization**: Uses caching to optimize build times when building container images.

1. **`App` (located in `app.go`):**
   This is the main application module that controls the flow of the program.

   - `Run()`: Executes the main application flow, including repository cloning, getting user input, and triggering the image build process.

2. **`Display` (located in `display.go`):**
   This module handles all display-related functionalities, such as printing to the console.

   - `PrintLogo()`: Outputs the application logo.
   - `PrintGreenBold(message string)`: Outputs a message in green and bold.

3. **`RepoManager` (located in `repo_manager.go`):**
   This module manages Git repositories.

   - `Clone(url string) error`: Clones the given Git repository URL.

4. **`ImageBuilder` (located in `image_builder.go`):**
   Handles building the container images using Cloud Native Buildpacks.

   - `Build(imageName string) error`: Builds a container image with the specified image name.
   - `TagImage(imageName string, tag string) error`: Tags the built image with the specified tag.

### Workflow:

1. **Initialization**:
   The user initializes the tool with the desired language and Git repository URL as command-line arguments.

2. **Repository Cloning**:
   The application attempts to clone the specified Git repository.

3. **Image Naming**:
   The user is prompted to enter a desired name for the resulting container image.

4. **Builder Selection**:
   Based on the language specified (or auto-detected), the application selects an appropriate Cloud Native Buildpack builder.

5. **Image Building**:
   Using the selected builder and the cloned repository, the application triggers the build process to create the container image.

6. **Image Tagging**:
   The user is prompted to enter a desired tag for the container image, and the image is tagged accordingly.

### Features:

- **Caching**:
  To optimize build times, the application utilizes a caching mechanism to store build layers. This ensures that only modified layers are rebuilt in subsequent runs, saving both time and resources.

- **Language Auto-detection** (optional feature):
  Instead of requiring the user to specify the language, the application can automatically detect the primary language of the cloned repository based on file extensions.

### Contributing

1. Fork the project.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

---
