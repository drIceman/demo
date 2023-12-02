package main

import (
	"github.com/drIceman/demo/internal/book"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/books", book.GetBooks)
	router.GET("/books/:id", book.GetBookByID)
	router.POST("/books", book.CreateBook)

	router.Run(":8080")
}
