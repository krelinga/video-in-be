package fanart

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/krelinga/go-fanart"
	"github.com/krelinga/video-in-be/env"
)

var client fanart.Client = func() fanart.Client {
	key := env.FanartKey()
	return fanart.NewClient(fanart.ClientOptions{
		APIKey: key,
	})
}()

type ArtURLs map[string]string

func GetArtURLs(ctx context.Context, tmdbId int) (ArtURLs, error) {
	if env.FanartKey() == "test-key" {
		// In test mode, don't call the fanart api.
		return ArtURLs{}, nil
	}

	result, err := fanart.GetMovie(ctx, client, fmt.Sprintf("%d", tmdbId))
	if err != nil {
		return nil, err
	}

	urls := make(ArtURLs)

	if banner, ok := getFirstImage(result, fanart.Movie.MovieBanner); ok {
		urls["banner"] = banner
	}

	if clearart, ok := getFirstImage(result, fanart.Movie.HdMovieClearArt); ok {
		urls["clearart"] = clearart
	}

	if clearlogo, ok := getFirstImage(result, fanart.Movie.HdMovieLogo); ok {
		urls["clearlogo"] = clearlogo
	} else if url, ok := getFirstImage(result, fanart.Movie.MovieLogo); ok {
		urls["clearlogo"] = url
	}

	if discart, ok := getFirstImage(result, fanart.Movie.MovieDisc); ok {
		urls["discart"] = discart
	}

	if landscape, ok := getFirstImage(result, fanart.Movie.MovieThumb); ok {
		urls["landscape"] = landscape
	}

	if keyart, ok := getFirstImage(result, fanart.Movie.MoviePoster); ok {
		urls["keyart"] = keyart
	}

	if logo, ok := getFirstImage(result, fanart.Movie.HdMovieLogo); ok {
		urls["logo"] = logo
	} else if url, ok := getFirstImage(result, fanart.Movie.MovieLogo); ok {
		urls["logo"] = url
	}

	if poster, ok := getFirstImage(result, fanart.Movie.MoviePoster); ok {
		urls["poster"] = poster
	}

	return urls, nil
}

func getFirstImage(movie fanart.Movie, getter func(fanart.Movie) ([]fanart.Image, error)) (string, bool) {
	images, err := getter(movie)
	if err != nil || len(images) == 0 {
		return "", false
	}

	type candidate struct {
		url      string
		language string
		likes    int
	}
	candidates := make([]candidate, 0, len(images))
	for _, img := range images {
		c := candidate{}
		if url, err := img.URL(); err != nil {
			continue
		} else {
			c.url = url
		}

		if likesStr, err := img.Likes(); err != nil {
			continue
		} else if likes, err := strconv.Atoi(likesStr); err != nil {
			continue
		} else {
			c.likes = likes
		}

		if lang, err := img.Lang(); err != nil {
			continue
		} else {
			c.language = lang
		}

		candidates = append(candidates, c)
	}

	getLangSortKey := func(c candidate) int {
		switch c.language {
		case "en":
			return 0
		case "":
			return 1
		default:
			return 2
		}
	}

	slices.SortFunc(candidates, func(a, b candidate) int {
		aLangKey := getLangSortKey(a)
		bLangKey := getLangSortKey(b)
		if result := cmp.Compare(aLangKey, bLangKey); result != 0 {
			return result
		}

		return cmp.Compare(b.likes, a.likes)
	})

	if len(candidates) > 0 {
		return candidates[0].url, true
	}

	return "", false
}
