package main

import (
	"context"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/state"
)

func (*ConnectService) UnclaimedDiscDirList(ctx context.Context, req *connect.Request[pb.UnclaimedDiscDirListRequest]) (*connect.Response[pb.UnclaimedDiscDirListResponse], error) {
	resp := &pb.UnclaimedDiscDirListResponse{}
	state.UnclaimedDiscDirRead(func(dirs []string) {
		resp.Dirs = dirs
	})
	return connect.NewResponse(resp), nil
}

func (*ConnectService) ProjectAssignDiskDirs(ctx context.Context, req *connect.Request[pb.ProjectAssignDiskDirsRequest]) (*connect.Response[pb.ProjectAssignDiskDirsResponse], error) {
	err := state.ProjectAssignDiskDirs(req.Msg.Project, req.Msg.Dirs)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&pb.ProjectAssignDiskDirsResponse{}), nil
}
