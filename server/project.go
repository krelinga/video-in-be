package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/publish"
	"github.com/krelinga/video-in-be/state"
	"github.com/krelinga/video-in-be/thumbs"
	"github.com/krelinga/video-in-be/tmdb"
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

func (*ConnectService) ProjectGet(ctx context.Context, req *connect.Request[pb.ProjectGetRequest]) (*connect.Response[pb.ProjectGetResponse], error) {
	// Create a new ProjectGetResponse
	response := &pb.ProjectGetResponse{}

	// Get the project
	var err error
	found := state.ProjectRead(req.Msg.Project, func(project *state.Project) {
		response.Project = project.Name
		if project.TmdbId != 0 {
			movieDetails, err := tmdb.GetMovieDetails(project.TmdbId)
			if err != nil {
				err = connect.NewError(connect.CodeInternal, errors.New("failed to get movie details"))
				return
			}
			response.SearchResult = convertMovieSearchResult(&movieDetails.MovieSearchResult)
		}
		for _, disc := range project.Discs {
			dProto := &pb.ProjectDisc{
				Disc:       disc.Name,
				ThumbState: string(disc.ThumbState),
			}
			response.Discs = append(response.Discs, dProto)
			if disc.ThumbState != state.ThumbStateDone {
				continue
			}
			for _, file := range disc.Files {
				fProto := &pb.DiscFile{
					File:     file.Name,
					Category: string(file.Category),
					Thumb:    file.Thumbnail,
				}
				dProto.DiscFiles = append(dProto.DiscFiles, fProto)
			}
		}
	})
	if !found {
		err = connect.NewError(connect.CodeNotFound, errors.New("project not found"))
	}
	if err != nil {
		return nil, err
	}

	// Return the response
	return connect.NewResponse(response), nil
}

func (*ConnectService) ProjectSetMetadata(ctx context.Context, req *connect.Request[pb.ProjectSetMetadataRequest]) (*connect.Response[pb.ProjectSetMetadataResponse], error) {
	resp := &pb.ProjectSetMetadataResponse{}
	id, err := strconv.Atoi(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid project ID"))
	}
	found := state.ProjectModify(req.Msg.Project, func(project *state.Project) {
		project.TmdbId = id
	})
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("project not found"))
	}
	return connect.NewResponse(resp), nil
}

func (*ConnectService) ProjectAbandon(ctx context.Context, req *connect.Request[pb.ProjectAbandonRequest]) (*connect.Response[pb.ProjectAbandonResponse], error) {
	resp := &pb.ProjectAbandonResponse{}
	found := state.ProjectReadAndRemove(req.Msg.Project, func(project *state.Project) error {
		if err := os.RemoveAll(state.ProjectDir(project.Name)); err != nil {
			return fmt.Errorf("could not remove project directory %s: %w", state.ProjectDir(project.Name), err)
		}

		// Remove thumbs
		thumbsDir := thumbs.ProjectThumbsDir(project.Name)
		if err := os.RemoveAll(thumbsDir); err != nil {
			return fmt.Errorf("could not remove thumbs directory %s: %w", thumbsDir, err)
		}

		return nil
	})
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("project not found"))
	}
	return connect.NewResponse(resp), nil
}

func (*ConnectService) ProjectFinish(ctx context.Context, req *connect.Request[pb.ProjectFinishRequest]) (*connect.Response[pb.ProjectFinishResponse], error) {
	resp := &pb.ProjectFinishResponse{}
	found := state.ProjectReadAndRemove(req.Msg.Project, func(project *state.Project) error {
		return publish.Do(project)
	})
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("project not found"))
	}

	return connect.NewResponse(resp), nil
}
