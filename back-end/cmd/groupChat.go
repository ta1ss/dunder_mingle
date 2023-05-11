package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/back-end/models"
	"social-network/back-end/repository"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type WebSocketConnectionGC struct {
	*websocket.Conn
}

type GroupChatJSONResponse struct {
	Action             string         `json:"action"`
	Message            string         `json:"message"`
	UserMessage        models.Message `json:"userMessage"`
	OnlineGroupMembers []int          `json:"onlineGroupMemebers"`
	Error              string         `json:"error"`
}

type GroupChatPayload struct {
	Action   string                `json:"action"`
	Message  string                `json:"message"`
	GroupId  int                   `json:"groupId"`
	SenderId int                   `json:"senderId"`
	Conn     WebSocketConnectionGC `json:"-"`
}

var wsChanGC = make(chan GroupChatPayload)
var groupClients = make(map[WebSocketConnectionGC]int)

var upgradeConnectionGC = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (app *application) GroupChatEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnectionGC.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgradeconnectonGC err: %v", err)
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

	connGC := WebSocketConnectionGC{Conn: ws}
	groupClients[connGC] = user.ID

	go func() {
		time.Sleep(1 * time.Second)
		BroadcastToAllGroupMembers(GroupChatJSONResponse{Action: "online_group_members", OnlineGroupMembers: GetOnlineGroupMembers()})
	}()

	go app.ListenForGroupChatWs(&connGC)
}

func (app *application) ListenForGroupChatWs(connGC *WebSocketConnectionGC) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload GroupChatPayload
	for {
		err := connGC.ReadJSON(&payload)
		if err != nil {
			//do nothing
		} else {
			payload.Conn = *connGC
			wsChanGC <- payload
		}
	}
}

func (app *application) ListenToGroupChatWsChannel() {
	for {
		event := <-wsChanGC
		switch event.Action {
		case "message":
			message := models.Message{
				SenderId:  groupClients[event.Conn],
				GroupId:   event.GroupId,
				Body:      event.Message,
				CreatedAt: time.Now(),
			}
			err := repository.CreateGroupMessage(message)
			if err != nil {
				BroadcastToCurrentGroupMember(event.Conn, GroupChatJSONResponse{Error: fmt.Sprintf("websocket create group message error %v:", err)})
				return
			}
			BroadcastToAllGroupMembers(GroupChatJSONResponse{Action: "new_message", UserMessage: message})
		case "offline":
			delete(groupClients, event.Conn)
			go func() {
				time.Sleep(1 * time.Second)
				BroadcastToAllGroupMembers(GroupChatJSONResponse{Action: "online_group_members", OnlineGroupMembers: GetOnlineGroupMembers()})
			}()
		}
	}
}

func BroadcastToCurrentGroupMember(conn WebSocketConnectionGC, response GroupChatJSONResponse) {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("broadcasttocurrentuser error")
		_ = conn.Close()
		delete(groupClients, conn)
	}
}

func BroadcastToAllGroupMembers(response GroupChatJSONResponse) {
	for groupClient := range groupClients {
		err := groupClient.WriteJSON(response)
		if err != nil {
			if err != websocket.ErrCloseSent {
				log.Println("websocket error")
			}
			_ = groupClient.Close()
			delete(groupClients, groupClient)
		}
	}
}

func GetOnlineGroupMembers() []int {
	var users []int
	for _, userID := range clients {
		if userID != 0 {
			users = append(users, userID)
		}
	}
	return users
}

func (app *application) LoadGroupChatMessages(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid group id"), http.StatusBadRequest)
		return
	}

	messages, err := repository.GetGroupChatMessages(groupID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot get group messages"), http.StatusBadRequest)
		return
	}
	_ = app.writeJson(w, http.StatusOK, messages, nil)
}
