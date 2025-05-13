package main

import (
	"context"
	"errors"
	"path/filepath"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/state"
)

func (*ConnectService) DiscCategorizeFiles(ctx context.Context, req *connect.Request[pb.DiscCategorizeFilesRequest]) (*connect.Response[pb.DiscCategorizeFilesResponse], error) {
	resp := &pb.DiscCategorizeFilesResponse{}

	for _, file := range req.Msg.Files {
		if file.File == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("file name is empty"))
		}

		isMainTitle := file.Category == string(state.FileCatMainTitle)
		isExtra := file.Category == string(state.FileCatExtra)
		isTrash := file.Category == string(state.FileCatTrash)
		if !isMainTitle && !isExtra && !isTrash {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid file category"))
		}
	}

	var err error
	found := state.ProjectModify(req.Msg.Project, func(p *state.Project) {
		if _, ok := p.Thumbs[req.Msg.Disc]; !ok {
			err = connect.NewError(connect.CodeNotFound, errors.New("unknown disc"))
			return
		}
		for _, file := range req.Msg.Files {
			fileKey := filepath.Join(req.Msg.Disc, file.File)
			p.Files[fileKey] = state.FileCat(file.Category)
		}
	})
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, state.ErrUnknownProject)
	}
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}