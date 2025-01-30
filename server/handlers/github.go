package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GithubFile struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Type          string `json:"type"`
	FileExtension string `json:"file_extension"`
}

var excludedExtensions = map[string]bool{
	// Binary and Executable Files
	".exe": true,
	".dll": true,
	".so":  true,
	".bin": true,
	".o":   true,
	".out": true,

	// Media Files
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".ico":  true,
	".webp": true,
	".bmp":  true,

	// Archive and Compressed Files
	".zip": true,
	".tar": true,
	".gz":  true,
	".bz2": true,
	".7z":  true,
	".rar": true,

	// Documentation and Markup
	".html": true,
	".xml":  true,
	".yaml": true,
	".yml":  true,
	".md":   true,
	".rst":  true,
	".txt":  true,

	// System and Metadata
	".DS_Store": true, // macOS
	".log":      true,
	".ini":      true,
	".cfg":      true,
	".conf":     true,

	// Fonts and Misc
	".ttf":  true,
	".woff": true,
	".eot":  true,
	".otf":  true,
	".mov":  true,
	".mp4":  true,
	".mp3":  true,
}

var excludedFileNames = map[string]bool{
	"LICENSE":         true,
	"README.md":       true,
	"README.txt":      true,
	"CONTRIBUTING.md": true,
	".gitignore":      true,
	".gitattributes":  true,

	// CI/CD and Build Files
	".gitlab-ci.yml": true,
	"Jenkinsfile":    true,
	"Makefile":       true,

	// Lock and Dependency Files
	"package-lock.json": true,
	"yarn.lock":         true,
	"pnpm-lock.yaml":    true,
	"composer.lock":     true,
	"requirements.txt":  true,
	"Pipfile.lock":      true,

	// Environment and Secrets
	".env": true,

	// Miscellaneous
	"Dockerfile": true,
	"Thumbs.db":  true, // Windows
	".DS_Store":  true, // macOS
}

func GetRepoFiles(owner, repo, path string) ([]GithubFile, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var files []GithubFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	for i, file := range files {

		if file.Type == "file" {
			extension, err := extractFileExtension(file.Name)
			if err != nil {
				return nil, fmt.Errorf("error extracting file extension: %w", err)
			}
			files[i].FileExtension = extension
		}

		if file.Type == "dir" {
			subDirFiles, err := GetRepoFiles(owner, repo, file.Path)
			if err != nil {
				return nil, fmt.Errorf("error fetching subdirectory files: %w", err)
			}
			files = append(files[:i+1], append(subDirFiles, files[i+1:]...)...)
		}
	}

	return files, nil
}

func extractFileExtension(fileName string) (string, error) {

	if strings.Contains(fileName, ".") {
		parts := strings.Split(fileName, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Invalid file name")
		}

		extension := fmt.Sprintf(".%v", parts[len(parts)-1])
		return extension, nil
	}
	return fileName, nil
}

func GetFileContent(owner, repo, path string) (string, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var fileContent struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(body, &fileContent); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 content: %w", err)
	}

	return string(decodedContent), nil
}

func ProcessFiles(files []GithubFile, owner, repo string) map[string]string {
	fileData := make(map[string]string)

	for _, file := range files {
		if file.Type == "file" {

            if excludedExtensions[file.FileExtension] {
                continue
            }

            if excludedFileNames[file.Name] {
                continue
            }
            
			content, err := GetFileContent(owner, repo, file.Path)
			if err != nil {
				log.Printf("Error getting file content: %v", err)
				continue
			}
			fileData[file.Path] = content
		}
	}
	return fileData
}
