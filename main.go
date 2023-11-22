package main

import (
	db "darkness-awakens/db"
	auth_service "darkness-awakens/service/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, []byte("Hello World"))
	}
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
		panic("Failed to load .env file")
	}
	db := db.SetupDB()

	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})
	router.POST("/auth/login", func(ctx *gin.Context) {
		auth_service.Login(ctx, db)
	})
	router.POST("/auth/register", func(ctx *gin.Context) {
		auth_service.Register(ctx, db)
	})

	router.Run("127.0.0.1:3000")

}
