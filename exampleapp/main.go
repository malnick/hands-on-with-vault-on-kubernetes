package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/hashicorp/vault/api"
)

var port string

type Config struct {
	Secrets []Secret `json:"secrets"`
	Message string   `json:"message"`
}

type Secret struct {
	Username string
	Password string
}

func main() {
	port = getEnvOrDefault("PORT", "80")
	portWithColon := fmt.Sprintf(":%s", port)

	fmt.Printf("Starting server on http://0.0.0.0:%s\n", port)
	fmt.Println("(Pass as PORT environment variable)")

	router := mux.NewRouter()
	router.PathPrefix("/socket.io/").Handler(startWebsocket())
	router.HandleFunc("/health", HealthHandler)
	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("assets").HTTPBox()))

	log.Fatal(http.ListenAndServe(portWithColon, router))
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HealthHandler returns a succesful status and a message.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

func startWebsocket() *gosocketio.Server {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	fmt.Println("Starting websocket server...")
	server.On(gosocketio.OnConnection, handleConnection)
	server.On("send", handleSend)

	return server
}

func handleConnection(c *gosocketio.Channel) {
	fmt.Println("New client connected")
	c.Join("visits")
	handleSend(c, Config{})
}

func handleSend(c *gosocketio.Channel, msg Config) string {
	config, err := getSecretDataFromVault()
	fmt.Println(config)

	if err != nil {
		config = Config{Secrets: nil, Message: err.Error()}
	}
	fmt.Println("secrets:", config.Secrets)
	c.Ack("message", config, time.Second*10)
	return "OK"
}

func getSecretDataFromVault() (Config, error) {
	c := Config{}
	token, ok := os.LookupEnv("VAULT_TOKEN")
	if !ok {
		return c, fmt.Errorf("unable to get token from env")
	}

	vc := api.DefaultConfig()
	vault, err := api.NewClient(vc)
	if err != nil {
		return c, err
	}

	vault.SetToken(token)
	logical := vault.Logical()
	s, err := logical.Read("secret/data/exampleapp/config")
	if err != nil {
		return c, err
	}

	fmt.Printf("got secrets %+v", s.Data)
	secret := Secret{}
	username, ok := s.Data["username"]
	if !ok {
		return c, fmt.Errorf("username empty")
	}

	secret.Username = username.(string)

	return Config{Secrets: []Secret{secret}, Message: "ok"}, nil
}
