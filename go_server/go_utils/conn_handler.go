package handlers

import (
		"fmt"
    "context"
    "log"
    "os"
    "strconv"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "path/to/generated/package" // Replace with your actual path
)

// GRPCClientConfig holds configuration for the gRPC client
type GRPCClientConfig struct {
    Host           string
    Port           int
    TimeoutSeconds int
    MaxWorkers     int // For future concurrency control
}

// GRPCClient wraps the connection and client for reuse
type GRPCClient struct {
    conn   *grpc.ClientConn
    client pb.GreeterClient
    config GRPCClientConfig
}

// NewGRPCClient creates a new gRPC client with configuration
func NewGRPCClient(config GRPCClientConfig) (*GRPCClient, error) {
    // Build the server address
    serverAddr := config.Host + ":" + strconv.Itoa(config.Port)
    
    // Establish connection to gRPC server
    conn, err := grpc.Dial(
        serverAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(), // Wait for connection to be established
    )
    if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to connect to gRPC server: %v\n", err)
        return nil, err
    }

    // Create the client stub
    client := pb.NewGreeterClient(conn)

    return &GRPCClient{
        conn:   conn,
        client: client,
        config: config,
    }, nil
}

// Close closes the gRPC connection
func (gc *GRPCClient) Close() error {
    return gc.conn.Close()
}

// GetPyDict makes a GetPyDict request to the gRPC server
func (gc *GRPCClient) GetPyDict(name string) (*pb.HelloReply, error) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(
        context.Background(), 
        time.Duration(gc.config.TimeoutSeconds)*time.Second,
    )
    defer cancel() // Always call cancel to free resources

    // Make the gRPC call
    response, err := gc.client.GetPyDict(ctx, &pb.HelloRequest{
        Name: name,
    })
    if err != nil {
        return nil, err
    }

    return response, nil
}

// Gateway is your main function that can be called from main.go
func Gateway() (*pb.HelloReply, error) {
    // Get configuration from environment variables with defaults
    config := GRPCClientConfig{
        Host:           getEnv("GRPC_HOST", "localhost"),
        Port:           getEnvAsInt("GRPC_PORT", 50051),
        TimeoutSeconds: getEnvAsInt("GRPC_TIMEOUT", 5),
        MaxWorkers:     getEnvAsInt("GRPC_MAX_WORKERS", 1),
    }

    // Create gRPC client
    grpcClient, err := NewGRPCClient(config)
    if err != nil {
        log.Printf("Failed to create gRPC client: %v", err)
        return nil, err
    }
    defer grpcClient.Close() // Ensure connection is closed

    // Make the request
    response, err := grpcClient.GetPyDict("Shayan")
    if err != nil {
        log.Printf("Failed to call GetPyDict: %v", err)
        return nil, err
    }

    log.Printf("Greeting: %s", response.Message)
    return response, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
