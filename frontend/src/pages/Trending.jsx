import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { trendingAPI, apiUtils } from '../services/api.js'

const Trending = () => {
  const [trendingMovies, setTrendingMovies] = useState([])
  const [trendingTV, setTrendingTV] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('all') // 'all', 'movies', 'tv'
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(0)

  useEffect(() => {
    fetchTrendingContent()
  }, [activeTab, page])

  const fetchTrendingContent = async () => {
    try {
      setLoading(true)
      setError(null)

      let response
      switch (activeTab) {
        case 'movies':
          response = await trendingAPI.getTrendingMovies(page)
          setTrendingMovies(response.data.results || [])
          setTrendingTV([])
          break
        case 'tv':
          response = await trendingAPI.getTrendingTV(page)
          setTrendingTV(response.data.results || [])
          setTrendingMovies([])
          break
        default:
          const [moviesResponse, tvResponse] = await Promise.all([
            trendingAPI.getTrendingMovies(page),
            trendingAPI.getTrendingTV(page)
          ])
          setTrendingMovies(moviesResponse.data.results || [])
          setTrendingTV(tvResponse.data.results || [])
          response = moviesResponse // Use for pagination
      }

      setTotalPages(response.data.total_pages || 0)
    } catch (err) {
      console.error('Error fetching trending content:', err)
      setError('Failed to load trending content. Please try again later.')
    } finally {
      setLoading(false)
    }
  }

  const handleTabChange = (tab) => {
    setActiveTab(tab)
    setPage(1)
  }

  const handleLoadMore = () => {
    if (page < totalPages && !loading) {
      setPage(page + 1)
    }
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

  const allContent = [...trendingMovies, ...trendingTV]

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">Trending Now</h1>
        <p className="page-subtitle">
          Discover what's popular in movies and TV shows this week
        </p>
      </div>

      {/* Tab Navigation */}
      <div style={{ display: 'flex', gap: 'var(--spacing-md)', marginBottom: 'var(--spacing-xl)' }}>
        <button
          className={`btn ${activeTab === 'all' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleTabChange('all')}
        >
          All
        </button>
        <button
          className={`btn ${activeTab === 'movies' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleTabChange('movies')}
        >
          Movies
        </button>
        <button
          className={`btn ${activeTab === 'tv' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => handleTabChange('tv')}
        >
          TV Shows
        </button>
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
          <p>Loading trending content...</p>
        </div>
      )}

      {/* Content Grid */}
      {allContent.length > 0 && (
        <div>
          <div className="grid grid-3">
            {allContent.map((item) => (
              <ContentCard 
                key={`${item.id}-${item.media_type || activeTab}`} 
                item={item} 
                type={item.media_type || (activeTab === 'movies' ? 'movie' : 'tv')} 
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

      {/* No Content */}
      {!loading && allContent.length === 0 && !error && (
        <div className="text-center">
          <h3>No trending content available</h3>
          <p>Check back later for the latest trending movies and TV shows</p>
        </div>
      )}
    </div>
  )
}

export default Trending 