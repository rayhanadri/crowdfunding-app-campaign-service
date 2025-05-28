package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/rayhanadri/crowdfunding-app-campaign-service/campaign-service/config"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/campaign-service/gen/go/campaign/v1"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/campaign-service/repository"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/campaign-service/service"
)

func main() {
	// Load .env only if running locally
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Warning: .env file not found, using environment variables instead")
		}
	}

	// Initialize database connection
	config.InitDB()
	gorm := config.DB

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Create repository instances
	campaignRepo := repository.NewCampaignRepository(gorm)

	// Inject repositories into services
	campaignService := service.NewCampaignService(campaignRepo)

	// Register server with grpc
	campaign.RegisterCampaignServiceServer(grpcServer, campaignService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5051"
	}
	// Listener for grpcServer without interceptors
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on PORT: %v", err)
	}

	log.Printf("Starting gRPC server on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
