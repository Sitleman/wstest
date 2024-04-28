package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Важно: не используйте этот метод в продакшене без добавления проверки источника запроса
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Ошибка при получении cookie:", err)
		//return
	} else {
		log.Printf("Получено cookie: %s\n", sessionCookie.Value)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке WebSocket соединения:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка при чтении сообщения:", err)
			break
		}
		log.Printf("Получено сообщение: %s\n", message)

		// Эхо-ответ
		responseMessage := "Получено: " + string(message)
		if err := conn.WriteMessage(messageType, []byte(responseMessage)); err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
			break
		}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "session",
		Value:   fmt.Sprintf("session_%d", rand.Intn(10000)),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // Cookie will expire after 24 hours
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("Hello!"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/v1/ws", wsHandler)
	http.HandleFunc("/api/v1/hello", helloHandler)
	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
