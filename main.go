package main

import (
	"context"
	"fmt"
	"net/http"

	pbconnect "buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/env"
	"github.com/krelinga/video-in-be/thumbs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ConnectService struct {
	thumbQueue *thumbs.Queue
}

func (s *ConnectService) HelloWorld(ctx context.Context, req *connect.Request[pb.HelloWorldRequest]) (*connect.Response[pb.HelloWorldResponse], error) {
	// Create a new HelloWorldResponse
	response := &pb.HelloWorldResponse{
		Message: "Hello, " + req.Msg.Name,
	}

	// Return the response
	return connect.NewResponse(response), nil
}

func main() {
	fmt.Println("Hello, World!")
	mux := http.NewServeMux()

	// Register the connectRPC service
	service := &ConnectService{
		thumbQueue: thumbs.NewQueue(1000),
	}
	path, handler := pbconnect.NewServiceHandler(service)
	mux.Handle(path, handler)

	// Serve static files with CORS enabled
	staticDir := env.ThumbsDir()
	staticHandler := http.StripPrefix("/thumbs/", http.FileServer(http.Dir(staticDir)))
	mux.Handle("/thumbs/", enableCORS(staticHandler))

	// Start the HTTP server
	http.ListenAndServe("0.0.0.0:25004", h2c.NewHandler(mux, &http2.Server{}))
}

// enableCORS is middleware that adds CORS headers to the response
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
