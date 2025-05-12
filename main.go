package main

import (
	"context"
	"fmt"
	"net/http"

	pbconnect "buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ConnectService struct {
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
	path, handler := pbconnect.NewServiceHandler(&ConnectService{})
	mux.Handle(path, handler)
	// Runs as long as the server is alive.
	http.ListenAndServe("0.0.0.0:25004", h2c.NewHandler(mux, &http2.Server{}))
}
