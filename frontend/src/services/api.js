import axios from 'axios'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add loading state or auth token here if needed
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Handle common errors
    if (error.response) {
      // Server responded with error status
      console.error('API Error:', error.response.data)
      
      // Handle specific error codes
      switch (error.response.status) {
        case 401:
          // Unauthorized - could redirect to login
          break
        case 403:
          // Forbidden
          break
        case 404:
          // Not found
          break
        case 429:
          // Rate limited
          console.error('Rate limit exceeded')
          break
        case 500:
          // Server error
          console.error('Server error')
          break
        default:
          console.error(`HTTP ${error.response.status}: ${error.response.data.message || 'Unknown error'}`)
      }
    } else if (error.request) {
      // Request was made but no response received
      console.error('Network error: No response received')
    } else {
      // Something else happened
      console.error('Request error:', error.message)
    }
    
    return Promise.reject(error)
  }
)

// Search API
export const searchAPI = {
  // Search for movies and TV shows
  search: (query, page = 1) => 
    api.get('/search', { params: { query, page } }),
  
  // Search for movies only
  searchMovies: (query, page = 1) => 
    api.get('/search/movies', { params: { query, page } }),
  
  // Search for TV shows only
  searchTV: (query, page = 1) => 
    api.get('/search/tv', { params: { query, page } }),
}

// Content details API
export const contentAPI = {
  // Get movie details
  getMovieDetails: (id) => 
    api.get(`/movie/${id}`),
  
  // Get TV show details
  getTVDetails: (id) => 
    api.get(`/tv/${id}`),
}

// Trending API
export const trendingAPI = {
  // Get trending content (movies and TV shows)
  getTrending: (page = 1) => 
    api.get('/trending', { params: { page } }),
  
  // Get trending movies
  getTrendingMovies: (page = 1) => 
    api.get('/trending/movies', { params: { page } }),
  
  // Get trending TV shows
  getTrendingTV: (page = 1) => 
    api.get('/trending/tv', { params: { page } }),
}

// Watchlist API
export const watchlistAPI = {
  // Get user's watchlist
  getWatchlist: (userId = 'default_user') => 
    api.get('/watchlist', { params: { user_id: userId } }),
  
  // Add item to watchlist
  addToWatchlist: (userId, contentId, contentType) => 
    api.post('/watchlist', { user_id: userId, content_id: contentId, content_type: contentType }),
  
  // Remove item from watchlist
  removeFromWatchlist: (userId, contentId, contentType) => 
    api.delete('/watchlist', { params: { user_id: userId, content_id: contentId, content_type: contentType } }),
  
  // Mark item as watched/unwatched
  markAsWatched: (userId, contentId, contentType, isWatched) => 
    api.put('/watchlist', { user_id: userId, content_id: contentId, content_type: contentType, is_watched: isWatched }),
}

// Genres API
export const genresAPI = {
  // Get all genres
  getGenres: () => 
    api.get('/genres'),
  
  // Get movies by genre
  getMoviesByGenre: (genreId, page = 1) => 
    api.get(`/genres/${genreId}/movies`, { params: { page } }),
  
  // Get TV shows by genre
  getTVByGenre: (genreId, page = 1) => 
    api.get(`/genres/${genreId}/tv`, { params: { page } }),
}

// Health check API
export const healthAPI = {
  // Check API health
  checkHealth: () => 
    api.get('/health'),
}

// Utility functions
export const apiUtils = {
  // Debounce function for search
  debounce: (func, wait) => {
    let timeout
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout)
        func(...args)
      }
      clearTimeout(timeout)
      timeout = setTimeout(later, wait)
    }
  },
  
  // Format date
  formatDate: (dateString) => {
    if (!dateString) return 'Unknown'
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  },
  
  // Format runtime
  formatRuntime: (minutes) => {
    if (!minutes) return 'Unknown'
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`
  },
  
  // Get poster URL
  getPosterURL: (posterPath, size = 'w500') => {
    if (!posterPath) return '/placeholder-poster.jpg'
    return `https://image.tmdb.org/t/p/${size}${posterPath}`
  },
  
  // Get backdrop URL
  getBackdropURL: (backdropPath, size = 'w1280') => {
    if (!backdropPath) return '/placeholder-backdrop.jpg'
    return `https://image.tmdb.org/t/p/${size}${backdropPath}`
  },
  
  // Truncate text
  truncateText: (text, maxLength = 150) => {
    if (!text) return ''
    if (text.length <= maxLength) return text
    return text.substring(0, maxLength) + '...'
  }
}

export default api 