package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"binge-base/config"
	"binge-base/database"
	"binge-base/services"

	"github.com/joho/godotenv"
)

type Server struct {
	config      *config.Config
	db          *database.Database
	tmdbService *services.TMDBService
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize TMDB service
	tmdbService := services.NewTMDBService(cfg)

	// Create server instance
	server := &Server{
		config:      cfg,
		db:          db,
		tmdbService: tmdbService,
	}

	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/v1/health", server.healthHandler)
	mux.HandleFunc("/api/v1/search", server.searchHandler)
	mux.HandleFunc("/api/v1/search/movies", server.searchMoviesHandler)
	mux.HandleFunc("/api/v1/search/tv", server.searchTVHandler)
	mux.HandleFunc("/api/v1/movie/", server.movieDetailsHandler)
	mux.HandleFunc("/api/v1/tv/", server.tvDetailsHandler)
	mux.HandleFunc("/api/v1/trending", server.trendingHandler)
	mux.HandleFunc("/api/v1/trending/movies", server.trendingMoviesHandler)
	mux.HandleFunc("/api/v1/trending/tv", server.trendingTVHandler)
	mux.HandleFunc("/api/v1/watchlist", server.watchlistHandler)
	mux.HandleFunc("/api/v1/genres", server.genresHandler)
	mux.HandleFunc("/api/v1/genres/", server.genresContentHandler)

	// Apply CORS middleware
	handler := corsMiddleware(mux)

	// Get port from environment or use default
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 BingeBase API server starting on port %s", port)
	log.Printf("📊 Database: %s", cfg.DBPath)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set content type for JSON responses
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// Helper function to send JSON responses
func (s *Server) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Helper function to send error responses
func (s *Server) sendError(w http.ResponseWriter, statusCode int, message string) {
	s.sendJSON(w, statusCode, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// Health check handler
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"status":   "ok",
		"message":  "BingeBase API is running",
		"database": "connected",
	})
}

// Search handlers
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		s.sendError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}
	// We'll fetch the first 3 pages for more results
	maxPages := 3
	var allMovieResults, allTVResults []interface{}
	totalMovieResults := 0
	totalTVResults := 0
	totalMoviePages := 0
	totalTVPages := 0
	for page := 1; page <= maxPages; page++ {
		movieResults, err := s.tmdbService.SearchMovies(query, page)
		if err == nil && movieResults != nil {
			for _, m := range movieResults.Results {
				if movie, ok := m.(map[string]interface{}); ok {
					movie["media_type"] = "movie"
				}
				allMovieResults = append(allMovieResults, m)
			}
			totalMovieResults = movieResults.TotalResults
			if movieResults.TotalPages > totalMoviePages {
				totalMoviePages = movieResults.TotalPages
			}
		}
		tvResults, err := s.tmdbService.SearchTVShows(query, page)
		if err == nil && tvResults != nil {
			for _, t := range tvResults.Results {
				if tv, ok := t.(map[string]interface{}); ok {
					tv["media_type"] = "tv"
				}
				allTVResults = append(allTVResults, t)
			}
			totalTVResults = tvResults.TotalResults
			if tvResults.TotalPages > totalTVPages {
				totalTVPages = tvResults.TotalPages
			}
		}
	}
	results := append(allMovieResults, allTVResults...)
	totalResults := totalMovieResults + totalTVResults
	totalPages := totalMoviePages
	if totalTVPages > totalPages {
		totalPages = totalTVPages
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          1,
		"results":       results,
		"total_pages":   totalPages,
		"total_results": totalResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

func (s *Server) searchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		s.sendError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}
	maxPages := 3
	var allMovieResults []interface{}
	totalMovieResults := 0
	totalMoviePages := 0
	for page := 1; page <= maxPages; page++ {
		movieResults, err := s.tmdbService.SearchMovies(query, page)
		if err == nil && movieResults != nil {
			for _, m := range movieResults.Results {
				if movie, ok := m.(map[string]interface{}); ok {
					movie["media_type"] = "movie"
				}
				allMovieResults = append(allMovieResults, m)
			}
			totalMovieResults = movieResults.TotalResults
			if movieResults.TotalPages > totalMoviePages {
				totalMoviePages = movieResults.TotalPages
			}
		}
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          1,
		"results":       allMovieResults,
		"total_pages":   totalMoviePages,
		"total_results": totalMovieResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

func (s *Server) searchTVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		s.sendError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}
	maxPages := 3
	var allTVResults []interface{}
	totalTVResults := 0
	totalTVPages := 0
	for page := 1; page <= maxPages; page++ {
		tvResults, err := s.tmdbService.SearchTVShows(query, page)
		if err == nil && tvResults != nil {
			for _, t := range tvResults.Results {
				if tv, ok := t.(map[string]interface{}); ok {
					tv["media_type"] = "tv"
				}
				allTVResults = append(allTVResults, t)
			}
			totalTVResults = tvResults.TotalResults
			if tvResults.TotalPages > totalTVPages {
				totalTVPages = tvResults.TotalPages
			}
		}
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          1,
		"results":       allTVResults,
		"total_pages":   totalTVPages,
		"total_results": totalTVResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

// Movie details handler
func (s *Server) movieDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract movie ID from URL path (last part)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 2 {
		s.sendError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}
	movieIDStr := pathParts[len(pathParts)-1]
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movie, err := s.tmdbService.GetMovieDetails(movieID)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch movie details: "+err.Error())
		return
	}

	// Optionally: fetch OMDB ratings if IMDB id is available
	// (You can add OMDBService to Server struct if you want to fetch more ratings)

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    movie,
	})
}

// TV details handler
func (s *Server) tvDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract TV show ID from URL path (last part)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 2 {
		s.sendError(w, http.StatusBadRequest, "Invalid TV show ID")
		return
	}
	tvIDStr := pathParts[len(pathParts)-1]
	tvID, err := strconv.Atoi(tvIDStr)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid TV show ID")
		return
	}
	tvShow, err := s.tmdbService.GetTVDetails(tvID)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch TV show details: "+err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tvShow,
	})
}

// Trending handlers
func (s *Server) trendingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	movies, err := s.tmdbService.GetTrendingMovies(page)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch trending movies: "+err.Error())
		return
	}
	tv, err := s.tmdbService.GetTrendingTVShows(page)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch trending TV shows: "+err.Error())
		return
	}
	results := append(movies.Results, tv.Results...)
	totalResults := movies.TotalResults + tv.TotalResults
	totalPages := movies.TotalPages
	if tv.TotalPages > totalPages {
		totalPages = tv.TotalPages
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          page,
		"results":       results,
		"total_pages":   totalPages,
		"total_results": totalResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

func (s *Server) trendingMoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	movies, err := s.tmdbService.GetTrendingMovies(page)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch trending movies: "+err.Error())
		return
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          page,
		"results":       movies.Results,
		"total_pages":   movies.TotalPages,
		"total_results": movies.TotalResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

func (s *Server) trendingTVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	tv, err := s.tmdbService.GetTrendingTVShows(page)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch trending TV shows: "+err.Error())
		return
	}
	resp := map[string]interface{}{
		"success":       true,
		"page":          page,
		"results":       tv.Results,
		"total_pages":   tv.TotalPages,
		"total_results": tv.TotalResults,
	}
	s.sendJSON(w, http.StatusOK, resp)
}

// Watchlist handler
func (s *Server) watchlistHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = "default_user"
		}
		items, err := s.db.GetWatchlist(userID)
		if err != nil {
			s.sendError(w, http.StatusInternalServerError, "Failed to get watchlist")
			return
		}
		// Fetch real details for each item
		var detailedItems []map[string]interface{}
		for _, item := range items {
			wi, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			contentID, _ := wi["content_id"].(int)
			contentType, _ := wi["content_type"].(string)
			isWatched, _ := wi["is_watched"].(bool)
			var details interface{}
			if contentType == "movie" {
				details, _ = s.tmdbService.GetMovieDetails(contentID)
			} else if contentType == "tv" {
				details, _ = s.tmdbService.GetTVDetails(contentID)
			}
			entry := map[string]interface{}{
				"contentId":   contentID,
				"contentType": contentType,
				"isWatched":   isWatched,
				"details":     details,
			}
			detailedItems = append(detailedItems, entry)
		}
		s.sendJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    detailedItems,
		})

	case http.MethodPost:
		var request struct {
			UserID      string `json:"user_id"`
			ContentID   int    `json:"content_id"`
			ContentType string `json:"content_type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.sendError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		if request.UserID == "" {
			request.UserID = "default_user"
		}
		if err := s.db.AddToWatchlist(request.UserID, request.ContentID, request.ContentType); err != nil {
			s.sendError(w, http.StatusInternalServerError, "Failed to add to watchlist")
			return
		}
		s.sendJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Added to watchlist",
		})

	case http.MethodDelete:
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = "default_user"
		}
		contentIDStr := r.URL.Query().Get("content_id")
		contentType := r.URL.Query().Get("content_type")
		contentID, err := strconv.Atoi(contentIDStr)
		if err != nil {
			s.sendError(w, http.StatusBadRequest, "Invalid content ID")
			return
		}
		if err := s.db.RemoveFromWatchlist(userID, contentID, contentType); err != nil {
			s.sendError(w, http.StatusInternalServerError, "Failed to remove from watchlist")
			return
		}
		s.sendJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Removed from watchlist",
		})

	case http.MethodPut:
		var request struct {
			UserID      string `json:"user_id"`
			ContentID   int    `json:"content_id"`
			ContentType string `json:"content_type"`
			IsWatched   bool   `json:"is_watched"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.sendError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		if request.UserID == "" {
			request.UserID = "default_user"
		}
		if err := s.db.MarkAsWatched(request.UserID, request.ContentID, request.ContentType, request.IsWatched); err != nil {
			s.sendError(w, http.StatusInternalServerError, "Failed to update watched status")
			return
		}
		s.sendJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Watch status updated",
		})

	default:
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Genres handlers
func (s *Server) genresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Get genres endpoint - to be implemented",
	})
}

func (s *Server) genresContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract genre ID and content type from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		s.sendError(w, http.StatusBadRequest, "Invalid genre path")
		return
	}
	genreID := pathParts[3]
	contentType := pathParts[4]

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Get content by genre endpoint - to be implemented",
		"genre_id":     genreID,
		"content_type": contentType,
	})
}
