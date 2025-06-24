package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type OMDBService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type OMDBResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	IMDBRating string `json:"imdbRating"`
	IMDBVotes  string `json:"imdbVotes"`
	IMDBID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Response   string `json:"Response"`
	Error      string `json:"Error,omitempty"`
}

func NewOMDBService(apiKey string) *OMDBService {
	return &OMDBService{
		apiKey:  apiKey,
		baseURL: "http://www.omdbapi.com/",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetMovieDetails gets detailed information about a movie from OMDB
func (s *OMDBService) GetMovieDetails(title string, year string) (*OMDBResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("t", title)
	if year != "" {
		params.Add("y", year)
	}
	params.Add("plot", "full")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", s.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details from OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response OMDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", response.Error)
	}

	return &response, nil
}

// GetTVShowDetails gets detailed information about a TV show from OMDB
func (s *OMDBService) GetTVShowDetails(title string, year string) (*OMDBResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("t", title)
	if year != "" {
		params.Add("y", year)
	}
	params.Add("type", "series")
	params.Add("plot", "full")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", s.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details from OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response OMDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", response.Error)
	}

	return &response, nil
}

// GetRatingsByIMDBID gets ratings using IMDB ID
func (s *OMDBService) GetRatingsByIMDBID(imdbID string) (*OMDBResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("i", imdbID)
	params.Add("plot", "short")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s?%s", s.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get ratings by IMDB ID: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response OMDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", response.Error)
	}

	return &response, nil
}

// ExtractRottenTomatoesRating extracts Rotten Tomatoes rating from OMDB response
func (s *OMDBService) ExtractRottenTomatoesRating(response *OMDBResponse) string {
	for _, rating := range response.Ratings {
		if rating.Source == "Rotten Tomatoes" {
			return rating.Value
		}
	}
	return ""
}

// ExtractMetascore extracts Metascore from OMDB response
func (s *OMDBService) ExtractMetascore(response *OMDBResponse) string {
	return response.Metascore
}
