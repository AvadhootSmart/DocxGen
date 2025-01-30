package main

import (
	"fmt"
	// "io/fs"
	"log"
	// "net/http"
	"os"
	// "path/filepath"
	"strings"

	"docxGen/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	// "github.com/gofiber/fiber/v2/middleware/adaptor"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/joho/godotenv"
	// "gopkg.in/src-d/go-git.v4"
)

type Config struct {
	GEMINI_API_KEY string
	PORT           string
}

func LoadConfig() *Config {
	config := &Config{
		GEMINI_API_KEY: os.Getenv("GEMINI_API_KEY"),
		PORT:           os.Getenv("PORT"),
	}

	if config.PORT == "" {
		config.PORT = "3000"
	}

	if config.GEMINI_API_KEY == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	return config

}

// func Handler(w http.ResponseWriter, r *http.Request) {
// 	r.RequestURI = r.URL.String()

// 	handler().ServeHTTP(w, r)
// }

// func main() http.HandlerFunc {

func main() {

	// config := LoadConfig()
	// if config == nil {
	// 	log.Fatal("Error accessing env variables")

	// }

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var GEMINI_API_KEY string = os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	app := fiber.New()

	app.Use(logger.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "https://docxgen.vercel.app",
	// 	AllowMethods: "GET,POST,PUT,DELETE",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Running successfully..")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// app.Post("/process-repo", func(c *fiber.Ctx) error {

	// 	type Request struct {
	// 		RepoURL string `json:"repo_url"`
	// 	}

	// 	req := new(Request)
	// 	if err := c.BodyParser(req); err != nil {
	// 		log.Println("Error parsing request: %v", err)
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "Invalid Request payload",
	// 		})
	// 	}

	// 	//Cloning and processing repo
	// 	tempDir := "/temp/repo"
	// 	defer os.RemoveAll(tempDir) //Cleanup

	// 	if err := cloneRepository(req.RepoURL, tempDir); err != nil {
	// 		log.Println("Error cloning repo: %v", err)
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error":   err.Error(),
	// 			"message": "Error cloning the repository",
	// 		})
	// 	}

	// 	fileData := processRepositoryFiles(tempDir)

	// 	fileJsonData, err := json.Marshal(fileData)
	// 	if err != nil {
	// 		log.Printf("Couldnt convert to json: %v", err)
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error":   err.Error(),
	// 			"message": "Error converting to json",
	// 		})

	// 	}

	// 	inputString := string(fileJsonData)

	// 	GEMINI_API_KEY := config.GEMINI_API_KEY

	// 	docx, err := handlers.GenerateDocx(inputString, GEMINI_API_KEY)
	// 	if err != nil {
	// 		log.Println("Error generating docx: %v", err)
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error":   err.Error(),
	// 			"message": "Error generating docx",
	// 		})
	// 	}

	// 	log.Println("Response: %v", docx)
	// 	// log.Println("Response: %v", inputString)

	// 	return c.JSON(fiber.Map{
	// 		"message": "Success",
	// 		"data":    docx,
	// 	})

	// 	// return c.JSON(fiber.Map{
	// 	// 	"message": "Success",
	// 	// 	"data":    fileData,
	// 	// })
	// })

	// log.Printf("Server started on http://localhost:%s", PORT)
	// if err := app.Listen(":" + PORT); err != nil {
	// 	log.Fatal("Error starting server, %v", err)
	// }

	app.Post("/Repo", func(c *fiber.Ctx) error {

		type Request struct {
			RepoURL string `json:"repo_url"`
		}

		req := new(Request)
		if err := c.BodyParser(req); err != nil {
			log.Println("Error parsing request: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Request payload",
			})
		}

		owner, repo, err := extractRepoDetails(req.RepoURL)
		if err != nil {
			log.Println("Error extracting repo details: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error extracting repo details",
			})
		}

		files, err := handlers.GetRepoFiles(owner, repo)
		if err != nil {
			log.Println("Error getting repo files: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error getting repo files",
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"files":   files,
		})
	})

	log.Printf("Server started on http://localhost:%s", PORT)
	app.Listen(":" + PORT)

	// return adaptor.FiberApp(app)
}

func extractRepoDetails(url string) (string, string, error) {

	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("Invalid repo url")
	}

	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	return owner, repo, nil
}

// func cloneRepository(url, destination string) error {
// 	_, err := git.PlainClone(destination, false, &git.CloneOptions{
// 		URL:      url,
// 		Progress: os.Stdout,
// 	})
// 	return err
// }

// var excludedExtensions = map[string]bool{
// 	".exe":  true,
// 	".png":  true,
// 	".jpg":  true,
// 	".jpeg": true,
// 	".gif":  true,
// 	".svg":  true,
// 	".xml":  true,
// 	".yaml": true,
// 	".html": true,
// 	".md":   true,
// 	".mov":  true,
// 	".ico":  true,
//     ".ttf":  true,
//     ".woff": true,
//     ".woff2": true,
//     ".eot": true,
//     ".zip": true,
//     ".mp4": true,
//     ".mp3": true,
//     ".tar": true,
//     ".gz": true,
//     ".rar": true,
//     ".7z": true,
// }

//	var excludedFileNames = map[string]bool{
//		"LICENSE":           true,
//		"Dockerfile":        true,
//		".gitignore":        true,
//		"package-lock.json": true,
//		"index.html":        true,
//	}

// var excludedExtensions = map[string]bool{
// 	// Binary and Executable Files
// 	".exe": true,
// 	".dll": true,
// 	".so":  true,
// 	".bin": true,
// 	".o":   true,
// 	".out": true,

// 	// Media Files
// 	".png":  true,
// 	".jpg":  true,
// 	".jpeg": true,
// 	".gif":  true,
// 	".svg":  true,
// 	".ico":  true,
// 	".webp": true,
// 	".bmp":  true,

// 	// Archive and Compressed Files
// 	".zip": true,
// 	".tar": true,
// 	".gz":  true,
// 	".bz2": true,
// 	".7z":  true,
// 	".rar": true,

// 	// Documentation and Markup
// 	".html": true,
// 	".xml":  true,
// 	".yaml": true,
// 	".yml":  true,
// 	".md":   true,
// 	".rst":  true,
// 	".txt":  true,

// 	// System and Metadata
// 	".DS_Store": true, // macOS
// 	".log":      true,
// 	".ini":      true,
// 	".cfg":      true,
// 	".conf":     true,

// 	// Fonts and Misc
// 	".ttf":  true,
// 	".woff": true,
// 	".eot":  true,
// 	".otf":  true,
// 	".mov":  true,
// 	".mp4":  true,
// 	".mp3":  true,
// }

// var excludedFileNames = map[string]bool{
// 	"LICENSE":         true,
// 	"README.md":       true,
// 	"README.txt":      true,
// 	"CONTRIBUTING.md": true,
// 	".gitignore":      true,
// 	".gitattributes":  true,

// 	// CI/CD and Build Files
// 	".gitlab-ci.yml": true,
// 	"Jenkinsfile":    true,
// 	"Makefile":       true,

// 	// Lock and Dependency Files
// 	"package-lock.json": true,
// 	"yarn.lock":         true,
// 	"pnpm-lock.yaml":    true,
// 	"composer.lock":     true,
// 	"requirements.txt":  true,
// 	"Pipfile.lock":      true,

// 	// Environment and Secrets
// 	".env": true,

// 	// Miscellaneous
// 	"Dockerfile": true,
// 	"Thumbs.db":  true, // Windows
// 	".DS_Store":  true, // macOS
// }

// func processRepositoryFiles(basePath string) map[string][]string {
// 	fileData := make(map[string][]string)

// 	filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			log.Printf("Error reading file: %v", err)
// 			return nil
// 		}

// 		if d.IsDir() {
// 			if d.Name() == ".git" {
// 				log.Printf("Skipping directory: %v", path)
// 				return filepath.SkipDir
// 			}
// 			return nil
// 		}

// 		ext := filepath.Ext(path)
// 		if excludedExtensions[ext] {
// 			log.Printf("Skipping file: %v with extension: %v", path, ext)
// 			return nil
// 		}

// 		if excludedFileNames[d.Name()] {
// 			log.Printf("Skipping file: %v with name: %v", path, d.Name())
// 			return nil
// 		}

// 		content, err := os.ReadFile(path)
// 		if err != nil {
// 			log.Printf("Failed to read file %v, error: %v", path, err)
// 			return nil
// 		}

// 		cleanedContent := preprocessContent(string(content))
// 		chunks := chunkContent(cleanedContent, 500)

// 		relativePath, _ := filepath.Rel(basePath, path)
// 		fileData[relativePath] = chunks

// 		return nil
// 	})

// 	log.Printf("fileData processed successfully")

// 	return fileData
// }

// func preprocessContent(content string) string {
// 	return content
// }

// func chunkContent(content string, chunkSize int) []string {
// 	var chunks []string
// 	runes := []rune(content)

// 	for i := 0; i < len(runes); i += chunkSize {
// 		end := i + chunkSize
// 		if end > len(runes) {
// 			end = len(runes)
// 		}
// 		chunks = append(chunks, string(runes[i:end]))
// 	}

// 	return chunks

// }
