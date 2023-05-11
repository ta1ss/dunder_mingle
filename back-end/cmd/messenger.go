package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/back-end/models"
	"social-network/back-end/repository"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

// MessengerJSONResponse defines the response sent back from the websocket
type MessengerJSONResponse struct {
	Action      string         `json:"action"`
	Message     string         `json:"message"`
	UserMessage models.Message `json:"userMessage"`
	OnlineUsers []int          `json:"onlineUsers"`
	Error       string         `json:"error"`
}

type MessengerPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	TargetId int                 `json:"targetId"`
	SenderId int                 `json:"senderId"`
	Conn     WebSocketConnection `json:"-"`
}

var wsChan = make(chan MessengerPayload)
var clients = make(map[WebSocketConnection]int)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (app *application) MessengerEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgradeconnecton err: %v", err)
	}

	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = user.ID

	go func() {
		time.Sleep(1 * time.Second)
		BroadcastToAll(MessengerJSONResponse{Action: "online_users", OnlineUsers: GetOnlineUsers()})
	}()

	go app.ListenForMessengerWs(&conn)
}

func (app *application) ListenForMessengerWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload MessengerPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func (app *application) ListenToMessengerWsChannel() {
	for {
		event := <-wsChan
		switch event.Action {
		case "message":
			message := models.Message{
				SenderId:  clients[event.Conn],
				TargetId:  event.TargetId,
				Body:      event.Message,
				CreatedAt: time.Now(),
			}
			err := repository.CreateUserMessage(message)
			if err != nil {
				BroadcastToCurrentUser(event.Conn, MessengerJSONResponse{Error: fmt.Sprintf("websocket create user message error %v:", err)})
				return
			}

			BroadcastToCurrentUser(event.Conn, MessengerJSONResponse{Action: "new_message", UserMessage: message})
			if clients[event.Conn] != event.TargetId {
				BroadcastToUserID(event.TargetId, MessengerJSONResponse{Action: "new_message", UserMessage: message})
			}
		case "read_messages":
			repository.MessagesRead(event.SenderId, clients[event.Conn])
		case "offline":
			delete(clients, event.Conn)
			go func() {
				time.Sleep(1 * time.Second)
				BroadcastToAll(MessengerJSONResponse{Action: "online_users", OnlineUsers: GetOnlineUsers()})
			}()
		}
	}
}

func BroadcastToCurrentUser(conn WebSocketConnection, response MessengerJSONResponse) {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("broadcasttocurrentuser error")
		_ = conn.Close()
		delete(clients, conn)
	}
}

func BroadcastToUserID(userID int, response MessengerJSONResponse) {
	for client := range clients {
		if clients[client] == userID {
			err := client.WriteJSON(response)
			if err != nil {
				if err != websocket.ErrCloseSent {
					log.Printf("broadcasttouserid err: %v", err)
				}
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}

func BroadcastToAll(response MessengerJSONResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			if err != websocket.ErrCloseSent {
				log.Println("websocket error")
			}
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func GetOnlineUsers() []int {
	var users []int
	for _, userID := range clients {
		if userID != 0 {
			users = append(users, userID)
		}
	}
	return users
}

func (app *application) GetUsers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to get users"), http.StatusInternalServerError)
		return
	}

	users, err = SortUsers(user.ID, users)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to sort users"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, users, nil)
}

func SortUsers(userID int, users []models.User) ([]models.User, error) {
	for i := 0; i < len(users); i++ {
		var err error
		var senderID int
		senderID, users[i].MessageRead, users[i].LastMessage, err = repository.GetLastMessageTime(userID, users[i].ID)
		if err != nil {
			return users, err
		}

		if senderID == userID {
			// If my message was last, dont add notification to chat
			users[i].MessageRead = 1
		}
	}

	sort.Slice(users, func(i, j int) bool {
		if users[i].LastMessage.Equal(users[j].LastMessage) {
			return users[i].FirstName + " " + users[i].LastName < users[j].FirstName + " " + users[j].LastName
		}
		return users[i].LastMessage.After(users[j].LastMessage)
	})

	return users, nil
}

func (app *application) LoadMessages(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "targetID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid target id"), http.StatusBadRequest)
		return
	}

	messages, err := repository.GetMessages(user.ID, targetID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot get messages"), http.StatusBadRequest)
		return
	}

	_ = app.writeJson(w, http.StatusOK, messages, nil)
}
