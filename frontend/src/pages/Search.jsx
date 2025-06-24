import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { searchAPI, apiUtils } from '../services/api.js'

const Search = () => {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [searchType, setSearchType] = useState('all') // 'all', 'movies', 'tv'
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(0)
  const [yearMin, setYearMin] = useState('')
  const [yearMax, setYearMax] = useState('')
  const [ratingMin, setRatingMin] = useState('')
  const [ratingMax, setRatingMax] = useState('')
  const [runtimeMin, setRuntimeMin] = useState('')
  const [runtimeMax, setRuntimeMax] = useState('')

  // Debounced search function
  const debouncedSearch = apiUtils.debounce(async (searchQuery, type, pageNum) => {
    if (!searchQuery.trim()) {
      setResults([])
      setTotalPages(0)
      return
    }

    try {
      setLoading(true)
      setError(null)

      let response
      switch (type) {
        case 'movies':
          response = await searchAPI.searchMovies(searchQuery, pageNum)
          break
        case 'tv':
          response = await searchAPI.searchTV(searchQuery, pageNum)
          break
        default:
          response = await searchAPI.search(searchQuery, pageNum)
      }

      if (pageNum === 1) {
        setResults(response.data.results || [])
      } else {
        setResults(prev => [...prev, ...(response.data.results || [])])
      }
      
      setTotalPages(response.data.total_pages || 0)
    } catch (err) {
      console.error('Search error:', err)
      setError('Failed to search. Please try again.')
    } finally {
      setLoading(false)
    }
  }, 500)

  const fetchResults = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await searchAPI.search(query)
      setResults(response.data.results || [])
    } catch (err) {
      setError('Failed to fetch results.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (query.length > 1) {
      fetchResults()
    } else {
      setResults([])
    }
    // eslint-disable-next-line
  }, [query, searchType, yearMin, yearMax, ratingMin, ratingMax, runtimeMin, runtimeMax])

  const handleLoadMore = () => {
    if (page < totalPages && !loading) {
      const nextPage = page + 1
      setPage(nextPage)
      debouncedSearch(query, searchType, nextPage)
    }
  }

  const handleSearchTypeChange = (type) => {
    setSearchType(type)
    setPage(1)
  }

  const ContentCard = ({ item, type }) => (
    <div className="card">
      <img 
        src={apiUtils.getPosterURL(item.poster_path)} 
        alt={item.title || item.name}
        className="card-image"
        onError={(e) => {
          e.target.src = '/placeholder-poster.jpg'
        }}
      />
      <div className="card-content">
        <h3 className="card-title">
          {apiUtils.truncateText(item.title || item.name, 30)}
        </h3>
        <p className="card-text">
          {apiUtils.truncateText(item.overview, 100)}
        </p>
        <div className="card-meta">
          <span>⭐ {item.vote_average?.toFixed(1) || 'N/A'}</span>
          <span>{apiUtils.formatDate(item.release_date || item.first_air_date)}</span>
        </div>
        <Link 
          to={`/${type}/${item.id}`} 
          className="btn btn-primary btn-sm"
          style={{ marginTop: 'var(--spacing-sm)', width: '100%' }}
        >
          View Details
        </Link>
      </div>
    </div>
  )

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">Search Movies & TV Shows</h1>
        <p className="page-subtitle">
          Find your next favorite entertainment
        </p>
      </div>

      {/* Search Form */}
      <div style={{ marginBottom: 'var(--spacing-xl)' }}>
        <div className="form-group">
          <input
            type="text"
            placeholder="Search for movies, TV shows, actors..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="form-input"
            style={{ fontSize: 'var(--font-size-lg)', padding: 'var(--spacing-lg)' }}
          />
        </div>

        {/* Search Type Filter */}
        <div style={{ display: 'flex', gap: 'var(--spacing-md)', marginBottom: 'var(--spacing-lg)' }}>
          <button
            className={`btn ${searchType === 'all' ? 'btn-primary' : 'btn-outline'}`}
            onClick={() => handleSearchTypeChange('all')}
          >
            All
          </button>
          <button
            className={`btn ${searchType === 'movies' ? 'btn-primary' : 'btn-outline'}`}
            onClick={() => handleSearchTypeChange('movies')}
          >
            Movies
          </button>
          <button
            className={`btn ${searchType === 'tv' ? 'btn-primary' : 'btn-outline'}`}
            onClick={() => handleSearchTypeChange('tv')}
          >
            TV Shows
          </button>
        </div>
      </div>

      {/* Advanced Filters */}
      <div className="advanced-filters">
        <label>Year Min: <input type="number" value={yearMin} onChange={e => setYearMin(e.target.value)} /></label>
        <label>Year Max: <input type="number" value={yearMax} onChange={e => setYearMax(e.target.value)} /></label>
        <label>Rating Min: <input type="number" step="0.1" value={ratingMin} onChange={e => setRatingMin(e.target.value)} /></label>
        <label>Rating Max: <input type="number" step="0.1" value={ratingMax} onChange={e => setRatingMax(e.target.value)} /></label>
        <label>Runtime Min: <input type="number" value={runtimeMin} onChange={e => setRuntimeMin(e.target.value)} /></label>
        <label>Runtime Max: <input type="number" value={runtimeMax} onChange={e => setRuntimeMax(e.target.value)} /></label>
      </div>

      {/* Error Message */}
      {error && (
        <div className="error">
          {error}
        </div>
      )}

      {/* Loading State */}
      {loading && page === 1 && (
        <div className="text-center">
          <div className="loading-spinner"></div>
          <p>Searching...</p>
        </div>
      )}

      {/* Search Results */}
      {results.length > 0 && (
        <div>
          <h2 style={{ marginBottom: 'var(--spacing-lg)' }}>
            Search Results ({results.length} found)
          </h2>
          <div className="grid grid-3">
            {results.map((item) => (
              <ContentCard 
                key={`${item.id}-${item.media_type || searchType}`} 
                item={item} 
                type={item.media_type || (searchType === 'movies' ? 'movie' : 'tv')} 
              />
            ))}
          </div>

          {/* Load More Button */}
          {page < totalPages && (
            <div className="text-center" style={{ marginTop: 'var(--spacing-xl)' }}>
              <button
                onClick={handleLoadMore}
                disabled={loading}
                className="btn btn-outline btn-lg"
              >
                {loading ? (
                  <>
                    <div className="loading-spinner"></div>
                    Loading...
                  </>
                ) : (
                  'Load More'
                )}
              </button>
            </div>
          )}
        </div>
      )}

      {/* No Results */}
      {!loading && query && results.length === 0 && !error && (
        <div className="text-center">
          <h3>No results found</h3>
          <p>Try adjusting your search terms or filters</p>
        </div>
      )}

      {/* Initial State */}
      {!query && !loading && (
        <div className="text-center">
          <h3>Start searching to discover content</h3>
          <p>Enter a movie title, TV show name, or actor to get started</p>
        </div>
      )}
    </div>
  )
}

export default Search 