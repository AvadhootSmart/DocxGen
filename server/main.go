package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/process-repo", func(c *fiber.Ctx) error {

		type Request struct {
			RepoURL string `json:"repo_url"`
		}

		req := new(Request)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Request payload",
			})
		}

		//Cloning and processing repo
		tempDir := "./tempRepo"
		defer os.RemoveAll(tempDir) //Cleanup

		if err := cloneRepository(req.RepoURL, tempDir); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		fileData := processRepositoryFiles(tempDir)

		return c.JSON(fiber.Map{
			"message": "Repo processed successfully",
			"data":    fileData,
		})
	})

    

	log.Println("Server started on http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Error starting server, %v", err)
	}
}

func cloneRepository(url, destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

var excludedExtensions = map[string]bool{
	".exe":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".xml":  true,
	".yaml": true,
}

var excludedFileNames = map[string]bool{
	"LICENSE":           true,
	"Dockerfile":        true,
	".gitignore":        true,
	"package-lock.json": true,
}

func processRepositoryFiles(basePath string) map[string][]string {
	fileData := make(map[string][]string)

	filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return nil
		}

		if d.IsDir() {
			if d.Name() == ".git" {
				log.Printf("Skipping directory: %v", path)
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if excludedExtensions[ext] {
			log.Printf("Skipping file: %v with extension: %v", path, ext)
			return nil
		}

		if excludedFileNames[d.Name()] {
			log.Printf("Skipping file: %v with name: %v", path, d.Name())
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %v, error: %v", path, err)
			return nil
		}

		cleanedContent := preprocessContent(string(content))
		chunks := chunkContent(cleanedContent, 500)

		relativePath, _ := filepath.Rel(basePath, path)
		fileData[relativePath] = chunks

		return nil
	})

	return fileData
}

func preprocessContent(content string) string {
	return content
}

func chunkContent(content string, chunkSize int) []string {
	var chunks []string
	runes := []rune(content)

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks

}
