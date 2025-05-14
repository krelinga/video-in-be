package main

import (
	"context"
	"fmt"
	"time"

	pb "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/tmdb"
)

func convertMovieSearchResult(in *tmdb.MovieSearchResult) *pb.MovieSearchResult {
	formatReleaseDate := func(d time.Time) string {
		if d.IsZero() {
			return ""
		}
		return d.Format(time.DateOnly)
	}
	return &pb.MovieSearchResult{
		Id:            fmt.Sprintf("%d", in.ID),
		OriginalTitle: in.OriginalTitle,
		PosterUrl:     in.PosterUrl,
		Title:         in.Title,
		ReleaseDate:   formatReleaseDate(in.RealaseDate),
		Overview:      in.Overview,
		Genres:        in.Genres,
	}
}

func (*ConnectService) MovieSearch(ctx context.Context, req *connect.Request[pb.MovieSearchRequest]) (*connect.Response[pb.MovieSearchResponse], error) {
	resp := &pb.MovieSearchResponse{}
	movies, err := tmdb.SearchMovies(req.Msg.PartialTitle)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, m := range movies {
		resp.Results = append(resp.Results, convertMovieSearchResult(m))
	}
	return connect.NewResponse(resp), nil
}
