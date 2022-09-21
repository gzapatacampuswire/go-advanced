package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Message struct {
	Id   string    `json:"id"`
	Time time.Time `json:"time"`
	Body string    `json:"body"`
	Read bool      `json:"read"`
}

type Chat struct {
	Id       string    `json:"id"`
	User     User      `json:"user"`
	Messages []Message `json:"messages"`
}

var users2 = []User{
	{Id: "1", Name: "Gustavo", Username: "tavo"},
	{Id: "2", Name: "Andres", Username: "andy"},
}
var users = map[string]User{
	"tavo": {Id: "1", Name: "Gustavo", Username: "tavo"},
	"andy": {Id: "2", Name: "Andres", Username: "andy"},
}
var messages = []Message{
	{Id: "1", Time: time.Now(), Body: "hey there", Read: true},
	{Id: "2", Time: time.Now(), Body: "how are you?", Read: false},
	{Id: "3", Time: time.Now(), Body: "are you there?", Read: false},
}

var chats []Chat

func main() {
	router := gin.Default() // this creates the server
	router.GET("/chats", getChats)
	router.POST("/chats", sendMessage)
	router.GET("/chat/:id", getChat)
	router.PATCH("/chat/:id", updateChat)
	router.Run("localhost:9090")
}

func updateChat(context *gin.Context) {
	id := context.Param("id")
	chat, err := getChatById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Chat not found"})
	}

	chat.Messages[0].Body = "updating this msg body"
	context.IndentedJSON(http.StatusOK, chat)
}

func getChat(context *gin.Context) {
	id := context.Param("id")
	chat, err := getChatById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Chat not found"})
	}
	context.IndentedJSON(http.StatusOK, chat)
}

func getChatById(id string) (*Chat, error) {
	for _, chat := range chats {
		if chat.Id == id {
			return &chat, nil
		}
	}
	return nil, errors.New("Chat not found")
}

func getChats(context *gin.Context) { // gin.Context holds all the info about the request
	chats = []Chat{
		{Id: "1", User: users["tavo"], Messages: messages},
		{Id: "2", User: users["andy"], Messages: messages},
	}
	context.IndentedJSON(http.StatusOK, chats) // converts the data struct into json
}

func sendMessage(context *gin.Context) {
	var message Message
	if err := context.BindJSON(&message); err != nil { // gets the json we send on the body and binds (maps) it to the struct, if they don't match, error is returned
		fmt.Println("err", err)
		return
	}

	messages = append(messages, message)
	context.IndentedJSON(http.StatusCreated, message)
}
