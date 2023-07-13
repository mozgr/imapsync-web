package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type CommandData struct {
	Host1             string
	User1             string
	Password1         string
	Host2             string
	User2             string
	Password2         string
	Automap           bool
	Delete2Duplicates bool
	Output            []string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	maxLines, _ := strconv.Atoi(os.Getenv("MAXLINES"))

	r := gin.Default()

	// Basic Auth middleware
	r.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
			if err == nil && string(decoded) == fmt.Sprintf("%s:%s", username, password) {
				c.Next()
				return
			}
		}

		c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	// Home route to display the form
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	// Command execution route
	r.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins to upgrade
			},
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		var data CommandData

		// Read the JSON payload sent by the client through the WebSocket connection
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read WebSocket message: %v", err)
			return
		}

		// Unmarshal the JSON payload into the CommandData struct
		if err := json.Unmarshal(payload, &data); err != nil {
			log.Printf("Failed to unmarshal JSON payload: %v", err)
			return
		}

		// Construct the command string
		command := fmt.Sprintf("imapsync --host1 %s --user1 %s --password1 %s --host2 %s --user2 %s --password2 %s --dry", data.Host1, data.User1, data.Password1, data.Host2, data.User2, data.Password2)

		// Append optional parameters if checked
		if data.Automap {
			command += " --automap"
		}
		if data.Delete2Duplicates {
			command += " --delete2duplicates"
		}

		// Create command execution
		cmd := exec.Command("bash", "-c", command)

		// Create a pipe to capture the command output
		outputPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Failed to create output pipe: %v", err)
			return
		}

		// Start the command
		err = cmd.Start()
		if err != nil {
			log.Printf("Failed to start command: %v", err)
			return
		}

		// Read and display output line by line
		scanner := bufio.NewScanner(outputPipe)
		lineCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineCount++

			if lineCount > maxLines {
				data.Output = data.Output[1:]
			}

			data.Output = append(data.Output, line)

			outputString := strings.Join(data.Output, "\n") + "\n"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(outputString)); err != nil {
				log.Printf("Failed to send WebSocket message: %v", err)
				break
			}
		}
	})

	// Serve the template files
	r.LoadHTMLGlob("templates/*")

	// Run the server
	log.Fatal(r.Run(":3000"))
}
