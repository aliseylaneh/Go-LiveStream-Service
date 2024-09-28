package server

import (
	"log"                   // Import the log package for logging.
	"net"                   // Import the net package for network-related operations.
	"safir/libs/appconfigs" // Import the appconfigs package for application configuration.
	"safir/libs/appstates"  // Import the appstates package for managing application states.
	"vpeer_room/internal/database"
	"vpeer_room/internal/handlers"
	"vpeer_room/internal/repository"
	"vpeer_room/internal/services"
	pb "vpeer_room/proto/api"

	"google.golang.org/grpc" // Import the gRPC package.
)

// RunServer starts the gRPC server to handle incoming requests.
func RunServer() {
	var (
		listenAddress = appconfigs.String("listen-address", "Server listen address") // Define a variable for the server's listen address.
		dbHost        = appconfigs.String("db-host", "Database host address")        // Define a variable for the database host address.
		dbPort        = appconfigs.Int("db-port", "Database port")                   // Define a variable for the database port.
		dbName        = appconfigs.String("db-name", "Database name")                // Define a variable for the database name.
		dbUsername    = appconfigs.String("db-username", "Database username")        // Define a variable for the database username.
		dbPassword    = appconfigs.String("db-password", "Database password")        // Define a variable for the database password.
	)

	// Handle configuration errors.
	if err := appconfigs.Parse(); err != nil {
		appstates.PanicMissingEnvParams(err.Error()) // Log an error if there are missing environment parameters.
	}

	// Connect to the PostgreSQL database.
	db, err := database.ConnectToPostgres(*dbHost, *dbPort, *dbName, *dbUsername, *dbPassword)
	if err != nil {
		appstates.PanicDBConnectionFailed(err.Error()) // Log an error if the database connection fails.
	}

	var (
		repository  repository.RoomRepository = repository.NewRoomRepository(db) // Create a repository instance.
		roomService services.RoomService      = services.NewRoomService(repository)
	)

	// Listen on the specified address for incoming connections.
	listener, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("error: %v", err)                    // Log an error if listening on the address fails.
		appstates.PanicServerSocketFailure(err.Error()) // Log an error for server socket failure.
	}

	// Create a new gRPC server instance.
	grpcServer := grpc.NewServer()

	roomHandler := handlers.NewRoomHandler(roomService)
	pb.RegisterRoomServiceServer(grpcServer, roomHandler)
	// Start serving the gRPC server.
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("error: %v", err) // Log an error if serving the gRPC server fails.
	}
}
