package main

import (
	"context"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/state"
)

func (*ConnectService) ProjectList(ctx context.Context, req *connect.Request[pb.ProjectListRequest]) (*connect.Response[pb.ProjectListResponse], error) {
	// Create a new ProjectListResponse
	response := &pb.ProjectListResponse{}
	state.ProjectsRead(func(projects []*state.Project) {
		response.Projects = make([]string, len(projects))
		for i := range projects {
			response.Projects[i] = projects[i].Name
		}
	})

	// Return the response
	return connect.NewResponse(response), nil
}
