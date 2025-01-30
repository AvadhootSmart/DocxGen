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

		//Extracts owner and repo from repo_url
		owner, repo, err := extractRepoDetails(req.RepoURL)
		if err != nil {
			log.Println("Error extracting repo details: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error extracting repo details",
			})
		}

		//Gets the entire file tree
		files, err := handlers.GetRepoFiles(owner, repo, "")
		if err != nil {
			log.Println("Error getting repo files: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error getting repo files",
			})
		}

		allFilesContent := handlers.ProcessFiles(files, owner, repo)


		return c.JSON(fiber.Map{
			"message": "success",
			// "files":   files,
			"fileContent": allFilesContent,
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

