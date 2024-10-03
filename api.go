package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	subs map[*websocket.Conn]Subscriber
	db   *sql.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer() *Server {
	return &Server{
		subs: make(map[*websocket.Conn]Subscriber),
		db:   NewDBInstance(),
	}
}

func (s *Server) StartServer(port string) {
	http.HandleFunc("/", handleFrontend)
	http.HandleFunc("/ping", s.Ping)
	http.HandleFunc("/subscriber", s.HandleSubscribe)
	http.HandleFunc("/notification", s.HandlerNotification)
	http.HandleFunc("/broadcast", s.HandleBroadcast)

	fmt.Println("Server start on port : 1234")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("error is : ", err)
		panic("unable to start the server")
	}
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}

func handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintln(w, "Get only :)")
	}
	content, err := os.ReadFile("frontend/index.html")
	if err != nil {
		fmt.Fprintln(w, "Error reading index.html to send to browser")
	}

	fmt.Println("content is ", string(content))

	fmt.Fprintln(w, string(content))
}

func (s *Server) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("erros while parsing : ", err)
			return
		}

		name := r.Form.Get("name")
		if name == "" {
			fmt.Fprintln(w, "no name attr found")
			return
		}

		fmt.Println("name is ", name)

		sub, err := s.InsertSub(name)
		if err != nil {
			log.Println("error while insert the sub : ", err)
			return
		}

		fmt.Fprintf(w, "Subscribed %s : password : %s ", name, sub.password)
	} else {
		fmt.Fprintf(w, "Only post request :)")
	}
}

func (s *Server) HandlerNotification(w http.ResponseWriter, r *http.Request) {
	// Authentication
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		http.Error(w, "No auth header found", http.StatusUnauthorized)
		return
	}

	parsed := strings.Split(auth_header, ";")
	if len(parsed) != 2 {
		http.Error(w, "Invalid auth header format", http.StatusBadRequest)
		return
	}

	id := parsed[0]
	password := parsed[1]

	if id == "" || password == "" {
		http.Error(w, "No id or password found", http.StatusBadRequest)
		return
	}

	sub, err := s.GetSub(id)
	if err != nil {
		fmt.Println("error while getting the sub : ", err)
		http.Error(w, "No subscriber found", http.StatusUnauthorized)
		return
	}

	if sub.password != password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	s.subs[ws] = sub

	ws.WriteMessage(websocket.TextMessage, []byte("Connected to notification server"))
}

func (s *Server) HandleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("error while parsing the form : ", err)
			return
		}
		msg := r.Form.Get("msg")
		if msg == "" {
			fmt.Fprintf(w, "msg not found")
			return
		}

		for conn, sub := range s.subs {
			if sub.valid {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s! %s", sub.name, msg))); err != nil {
					log.Println("error while sending the message : ", err)
					return
				}
			}
		}
	} else {
		fmt.Fprintf(w, "Only post request :)")
	}

	fmt.Fprintf(w, "Broadcasted Message")
}
