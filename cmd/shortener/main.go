package main

import (
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	//"log"
	//"github.com/april1858/shortener-go/internal/pkg/app"
)

var BaseURL string = "http://localhost:8080/"

var data = make(map[string]string)

func createShortened(c *gin.Context) {
	body, _ := c.GetRawData()
	shortened := service()
	data[shortened] = string(body)
	c.Data(http.StatusOK, "", []byte(BaseURL + shortened))

	return
}

func service() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "error in createCode()"
	}

	return hex.EncodeToString(b)
}

func getOriginalURL(c *gin.Context) {
	shortened := c.Param("id")
	answer := data[shortened]
	c.Redirect(http.StatusTemporaryRedirect, answer)
}

func main() {
	fmt.Println("server runing")
	router := gin.Default()
	router.POST("/", createShortened)
	router.GET("/:id", getOriginalURL)
	router.Run(":8080")
}
