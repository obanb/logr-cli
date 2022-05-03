package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerStatus int

const (
	Running ServerStatus = iota
	Idle
)

func (s ServerStatus) String() string {
	switch s {
	case Running:
		return "Running"
	case Idle:
		return "Idle"
	default:
		return "Unknown"
	}
}

type HttpServer struct {
	port            string
	healthCheckPath string
	quit            chan os.Signal
}

type Broadcast struct {
	id      string                 `json:"id"`
	body    map[string]interface{} `json:"body"`
	headers http.Header            `json:"headers"`
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPS")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *HttpServer) Start() error {
	fmt.Printf("Starting server %s:", s.port)
	s.quit = make(chan os.Signal)

	newHandler := func() *gin.Engine {
		h := gin.New()
		h.Use(gin.Recovery())
		h.Use(CORS())
		return h
	}

	handler := newHandler()

	handler.POST("/imitate", func(c *gin.Context) {

		var body map[string]interface{}
		headers := c.Request.Header

		err := c.BindJSON(&body)
		if err != nil {
			log.Fatal(err)
		}

		file, _ := json.MarshalIndent(body, "", " ")

		err = ioutil.WriteFile("storage.json", file, 0644)

		if err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadFile("storage.json")

		var request map[string]interface{}

		err = json.Unmarshal(content, &request)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("pes")
		fmt.Println(headers)
		fmt.Println(request)
		fmt.Println("pes")

	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM)
	<-s.quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("server shutting down")
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	server := HttpServer{
		port:            "8080",
		healthCheckPath: "/health",
		quit:            make(chan os.Signal, 1),
	}
	fmt.Println(server)
	var rootCmd = &cobra.Command{
		Use:   "logr-cli",
		Short: "Logr cli application written in Go",
		Long:  `Logr cli application written in Go.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := server.Start(); err != nil {
				return err
			}
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}
