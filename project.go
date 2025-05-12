package main

import (
	"context"
	"errors"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/state"
)

var (
	ErrProjectAlreadyExists = connect.NewError(connect.CodeAlreadyExists, errors.New("project already exists"))
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

func (*ConnectService) ProjectNew(ctx context.Context, req *connect.Request[pb.ProjectNewRequest]) (*connect.Response[pb.ProjectNewResponse], error) {
	// Create a new ProjectNewResponse
	response := &pb.ProjectNewResponse{}

	// Create a new project
	var err error
	state.ProjectsModify(func(projects []*state.Project) []*state.Project {
		for _, project := range projects {
			if project.Name == req.Msg.Name {
				// If the project already exists, return an error
				err = ErrProjectAlreadyExists
				return projects
			}
		}

		projects = append(projects, &state.Project{
			Name: req.Msg.Name,
		})
		return projects
	})
	if err != nil {
		return nil, err
	}

	// Return the response
	return connect.NewResponse(response), nil
}
