package server

import (
	"context"
	"errors"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/state"
)

func (*ConnectService) ProjectCategorizeFiles(ctx context.Context, req *connect.Request[pb.ProjectCategorizeFilesRequest]) (*connect.Response[pb.ProjectCategorizeFilesResponse], error) {
	resp := &pb.ProjectCategorizeFilesResponse{}

	for _, file := range req.Msg.Files {
		if file.File == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("file name is empty"))
		}

	}

	var err error
	found := state.ProjectModify(req.Msg.Project, func(p *state.Project) {
		for _, fileProto := range req.Msg.Files {
			disc := p.FindDiscByName(fileProto.Disc)
			if disc == nil {
				err = connect.NewError(connect.CodeNotFound, errors.New("disc not found"))
				return
			}
			file := disc.FindFileByName(fileProto.File)
			if file == nil {
				err = connect.NewError(connect.CodeNotFound, errors.New("file not found"))
				return
			}

			isNone := fileProto.Category == string(state.FileCatNone)
			isMainTitle := fileProto.Category == string(state.FileCatMainTitle)
			isExtra := fileProto.Category == string(state.FileCatExtra)
			isTrash := fileProto.Category == string(state.FileCatTrash)
			if !isMainTitle && !isExtra && !isTrash && !isNone {
				err = connect.NewError(connect.CodeInvalidArgument, errors.New("invalid file category"))
				return
			}
			file.Category = state.FileCat(fileProto.Category)
		}
	})
	if !found {
		err = connect.NewError(connect.CodeNotFound, state.ErrUnknownProject)
	}
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
