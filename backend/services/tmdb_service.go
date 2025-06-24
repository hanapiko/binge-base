package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"binge-base/config"
	"binge-base/models"
)

type TMDBService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewTMDBService(cfg *config.Config) *TMDBService {
	return &TMDBService{
		apiKey:  cfg.TMDBAPIKey,
		baseURL: "https://api.themoviedb.org/3",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchMovies searches for movies using TMDB API
func (s *TMDBService) SearchMovies(query string, page int) (*models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/movie", s.baseURL)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("query", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("include_adult", "false")
	params.Add("language", "en-US")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// SearchTVShows searches for TV shows using TMDB API
func (s *TMDBService) SearchTVShows(query string, page int) (*models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/tv", s.baseURL)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("query", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("include_adult", "false")
	params.Add("language", "en-US")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetMovieProviders fetches streaming providers for a movie
func (s *TMDBService) GetMovieProviders(movieID int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/movie/%d/watch/providers", s.baseURL, movieID)
	params := url.Values{}
	params.Add("api_key", s.apiKey)
	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get movie providers: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode providers: %w", err)
	}
	return result, nil
}

// GetMovieDetails gets detailed information about a movie
func (s *TMDBService) GetMovieDetails(movieID int) (*models.Movie, error) {
	endpoint := fmt.Sprintf("%s/movie/%d", s.baseURL, movieID)
	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("language", "en-US")
	params.Add("append_to_response", "credits,videos,images")
	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}
	var movie models.Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	// Fetch providers
	providers, _ := s.GetMovieProviders(movieID)
	if providers != nil {
		movie.Providers = providers["results"]
	}
	// Extract YouTube trailer key
	if videos, ok := movie.Videos["results"].([]interface{}); ok {
		for _, v := range videos {
			if video, ok := v.(map[string]interface{}); ok {
				if video["site"] == "YouTube" && video["type"] == "Trailer" {
					movie.Trailer = video["key"].(string)
					break
				}
			}
		}
	}
	return &movie, nil
}

// GetTVDetails gets detailed information about a TV show
func (s *TMDBService) GetTVDetails(tvID int) (*models.TVShow, error) {
	endpoint := fmt.Sprintf("%s/tv/%d", s.baseURL, tvID)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("language", "en-US")
	params.Add("append_to_response", "credits,videos,images")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tvShow models.TVShow
	if err := json.Unmarshal(body, &tvShow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tvShow, nil
}

// GetTrendingMovies gets trending movies
func (s *TMDBService) GetTrendingMovies(page int) (*models.TrendingResult, error) {
	endpoint := fmt.Sprintf("%s/trending/movie/week", s.baseURL)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("page", strconv.Itoa(page))
	params.Add("language", "en-US")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.TrendingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetTrendingTVShows gets trending TV shows
func (s *TMDBService) GetTrendingTVShows(page int) (*models.TrendingResult, error) {
	endpoint := fmt.Sprintf("%s/trending/tv/week", s.baseURL)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("page", strconv.Itoa(page))
	params.Add("language", "en-US")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get trending TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.TrendingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetGenres gets movie and TV show genres
func (s *TMDBService) GetGenres() ([]models.Genre, error) {
	endpoint := fmt.Sprintf("%s/genre/movie/list", s.baseURL)

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("language", "en-US")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Genres []models.Genre `json:"genres"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Genres, nil
}
