package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func Handler (w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close();

	go HandleConnection(conn);
}

func HandleConnection(conn *websocket.Conn){
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		log.Printf("Received message type: %d, content: %s", messageType, string(p))

		processMessage(conn, messageType, p)
	}
}

func processMessage(conn *websocket.Conn, messageType int, data []byte) {
	// Обработка разных типов сообщений
	switch messageType {
	case websocket.TextMessage:
		log.Printf("Text message: %s", string(data))
		// Отправляем ответ
		conn.WriteMessage(websocket.TextMessage, []byte("Server received: "+string(data)))
		
	case websocket.BinaryMessage:
		log.Printf("Binary message, length: %d bytes", len(data))
		conn.WriteMessage(websocket.BinaryMessage, data)
		
	case websocket.CloseMessage:
		log.Println("Client closed connection")
		return
		
	case websocket.PingMessage:
		log.Println("Received ping")
		conn.WriteMessage(websocket.PongMessage, nil)
		
	case websocket.PongMessage:
		log.Println("Received pong")
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request){
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break;
		}
		log.Printf("Received: %s", message)
	
		if err := ws.WriteMessage(messageType, message); err != nil{
			log.Println(err);
			break;
		}
	}
}
