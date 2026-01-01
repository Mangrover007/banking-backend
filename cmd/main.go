package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Mangrover007/banking-backend/database"
)

type User struct {
	FullName    string `json:"fullName"`
	PhoneNumber uint32 `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type Account struct {
	AccountNumber uint32
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumsById(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"message": fmt.Sprintf("no album with id: %s found", id),
	})
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
			"error":   err.Error(),
		})
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusOK, newAlbum)
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong motherfucker",
		})
	})

	r.GET("/albums", getAlbums)
	r.GET("/albums/:id", getAlbumsById)
	r.POST("/albums", postAlbums)

	database.Shit()

	r.Run()
}
