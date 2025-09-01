package manual

import (
	"flag"
)

var (
	movieDirFlag    = flag.String("movie_dir", "", "Directory containing movie files")
	tmdbMovieIdFlag = flag.Int("tmdb_movie_id", 0, "TMDB movie ID")
	dryRunFlag      = flag.Bool("dry_run", false, "Enable dry run mode")
)
