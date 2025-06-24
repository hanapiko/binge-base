package models

import (
	"time"
)

// Movie represents a movie from TMDB
type Movie struct {
	ID                   int                    `json:"id" db:"id"`
	TMDBID               int                    `json:"tmdb_id" db:"tmdb_id"`
	Title                string                 `json:"title" db:"title"`
	Overview             string                 `json:"overview" db:"overview"`
	PosterPath           string                 `json:"poster_path" db:"poster_path"`
	BackdropPath         string                 `json:"backdrop_path" db:"backdrop_path"`
	ReleaseDate          string                 `json:"release_date" db:"release_date"`
	VoteAverage          float64                `json:"vote_average" db:"vote_average"`
	VoteCount            int                    `json:"vote_count" db:"vote_count"`
	Popularity           float64                `json:"popularity" db:"popularity"`
	GenreIDs             []int                  `json:"genre_ids" db:"genre_ids"`
	Runtime              int                    `json:"runtime" db:"runtime"`
	Status               string                 `json:"status" db:"status"`
	Tagline              string                 `json:"tagline" db:"tagline"`
	Budget               int64                  `json:"budget" db:"budget"`
	Revenue              int64                  `json:"revenue" db:"revenue"`
	IMDBRating           string                 `json:"imdb_rating" db:"imdb_rating"`
	RottenTomatoesRating string                 `json:"rotten_tomatoes_rating" db:"rotten_tomatoes_rating"`
	CreatedAt            time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at" db:"updated_at"`
	Providers            interface{}            `json:"providers,omitempty"`
	Videos               map[string]interface{} `json:"videos,omitempty"`
	Trailer              string                 `json:"trailer,omitempty"`
}

// TVShow represents a TV show from TMDB
type TVShow struct {
	ID                   int       `json:"id" db:"id"`
	TMDBID               int       `json:"tmdb_id" db:"tmdb_id"`
	Name                 string    `json:"name" db:"name"`
	Overview             string    `json:"overview" db:"overview"`
	PosterPath           string    `json:"poster_path" db:"poster_path"`
	BackdropPath         string    `json:"backdrop_path" db:"backdrop_path"`
	FirstAirDate         string    `json:"first_air_date" db:"first_air_date"`
	LastAirDate          string    `json:"last_air_date" db:"last_air_date"`
	VoteAverage          float64   `json:"vote_average" db:"vote_average"`
	VoteCount            int       `json:"vote_count" db:"vote_count"`
	Popularity           float64   `json:"popularity" db:"popularity"`
	GenreIDs             []int     `json:"genre_ids" db:"genre_ids"`
	NumberOfSeasons      int       `json:"number_of_seasons" db:"number_of_seasons"`
	NumberOfEpisodes     int       `json:"number_of_episodes" db:"number_of_episodes"`
	Status               string    `json:"status" db:"status"`
	Type                 string    `json:"type" db:"type"`
	IMDBRating           string    `json:"imdb_rating" db:"imdb_rating"`
	RottenTomatoesRating string    `json:"rotten_tomatoes_rating" db:"rotten_tomatoes_rating"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// WatchlistItem represents an item in user's watchlist
type WatchlistItem struct {
	ID          int        `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	ContentID   int        `json:"content_id" db:"content_id"`
	ContentType string     `json:"content_type" db:"content_type"` // "movie" or "tv"
	IsWatched   bool       `json:"is_watched" db:"is_watched"`
	AddedAt     time.Time  `json:"added_at" db:"added_at"`
	WatchedAt   *time.Time `json:"watched_at" db:"watched_at"`
}

// Genre represents a movie/TV show genre
type Genre struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// SearchResult represents a search result
type SearchResult struct {
	Page         int           `json:"page"`
	Results      []interface{} `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

// TrendingResult represents trending content
type TrendingResult struct {
	Page         int           `json:"page"`
	Results      []interface{} `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
	PerPage      int `json:"per_page"`
}
