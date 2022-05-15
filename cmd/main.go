package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homeWorkAdvancedPart2_REST_service/internal/genserve"
	"homeWorkAdvancedPart2_REST_service/internal/grpcclient"
	"homeWorkAdvancedPart2_REST_service/internal/server"
	"homeWorkAdvancedPart2_REST_service/internal/store/inmemorystore"
	"homeWorkAdvancedPart2_REST_service/internal/swaggerdoc"
	"log"
	"net/http"
)

var (
	port            = flag.String("port", "8081", "The server port")
	grpcServiceAddr = flag.String("grpc-addr", "localhost:8080", "The grpc delivery service addr. Format: domain:port")
)

func main() {
	flag.Parse()

	router := chi.NewRouter()
	store := inmemorystore.NewStore()

	conn, err := grpc.Dial(*grpcServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect grpc: %v", err)
	}
	defer conn.Close()
	grpcClient := grpcclient.NewClient(conn)

	s := server.NewServer(store, grpcClient)

	router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	router.Use(middleware.Logger)

	// Да это ужасно, но хотелось и доку попробовать отдавать)
	swag.Register("doc.json", &swaggerdoc.Sw{File: "doc.json"})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.InstanceName("doc.json"),
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	router.Handle("/*", genserve.HandlerWithOptions(s, genserve.ChiServerOptions{
		BaseURL: "/api",
	}))

	httpServer := http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	fmt.Println("HTTP Server Listen on", *port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}
