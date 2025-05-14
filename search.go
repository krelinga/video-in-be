package main

import (
	"context"
	"fmt"
	"time"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/tmdb"
)

func (*ConnectService) MovieSearch(ctx context.Context, req *connect.Request[pb.MovieSearchRequest]) (*connect.Response[pb.MovieSearchResponse], error) {
	resp := &pb.MovieSearchResponse{}
	movies, err := tmdb.SearchMovies(req.Msg.PartialTitle)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	formatReleaseDate := func(d time.Time) string {
		if d.IsZero() {
			return ""
		}
		return d.Format(time.DateOnly)
	}
	for _, m := range movies {
		resp.Results = append(resp.Results, &pb.MovieSearchResult{
			Id:            fmt.Sprintf("%d", m.ID),
			OriginalTitle: m.OriginalTitle,
			PosterUrl:     m.PosterUrl,
			Title:         m.Title,
			ReleaseDate:   formatReleaseDate(m.RealaseDate),
			Overview:      m.Overview,
			Genres:        m.Genres,
		})
	}
	return connect.NewResponse(resp), nil
}
