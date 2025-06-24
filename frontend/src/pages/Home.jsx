import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { trendingAPI } from '../services/api.js'
import { apiUtils } from '../services/api.js'

const Home = () => {
  const [trendingMovies, setTrendingMovies] = useState([])
  const [trendingTV, setTrendingTV] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchTrendingContent = async () => {
      try {
        setLoading(true)
        setError(null)
        const [moviesResponse, tvResponse] = await Promise.all([
          trendingAPI.getTrendingMovies(1),
          trendingAPI.getTrendingTV(1)
        ])
        setTrendingMovies(moviesResponse.data.results || [])
        setTrendingTV(tvResponse.data.results || [])
      } catch (err) {
        setError('Failed to load trending content. Please try again later.')
      } finally {
        setLoading(false)
      }
    }
    fetchTrendingContent()
  }, [])

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

  if (loading) {
    return (
      <div className="page">
        <div className="page-header">
          <h1 className="page-title">Welcome to BingeBase</h1>
          <p className="page-subtitle">
            Discover your next favorite movie or TV show
          </p>
        </div>
        <div className="text-center">
          <div className="loading-spinner"></div>
          <p>Loading trending content...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="page">
      {/* Hero Section */}
      <div className="page-header">
        <h1 className="page-title">Welcome to BingeBase</h1>
        <p className="page-subtitle">
          Your ultimate destination for discovering movies and TV shows. 
          Search, explore, and manage your personal watchlist all in one place.
        </p>
        <div style={{ marginTop: 'var(--spacing-xl)' }}>
          <Link to="/search" className="btn btn-primary btn-lg">
            Start Searching
          </Link>
          <Link to="/trending" className="btn btn-outline btn-lg" style={{ marginLeft: 'var(--spacing-md)' }}>
            View Trending
          </Link>
        </div>
      </div>

      {error && (
        <div className="error">
          {error}
        </div>
      )}

      {/* Quick Actions */}
      <section style={{ marginBottom: 'var(--spacing-2xl)' }}>
        <h2 style={{ marginBottom: 'var(--spacing-lg)' }}>Quick Actions</h2>
        <div className="grid grid-2">
          <Link to="/search" className="card" style={{ textDecoration: 'none' }}>
            <div className="card-content">
              <h3>🔍 Search Movies & TV Shows</h3>
              <p>Find exactly what you're looking for with our powerful search feature</p>
            </div>
          </Link>
          <Link to="/watchlist" className="card" style={{ textDecoration: 'none' }}>
            <div className="card-content">
              <h3>📋 Manage Your Watchlist</h3>
              <p>Keep track of movies and shows you want to watch</p>
            </div>
          </Link>
        </div>
      </section>

      {/* Trending Movies */}
      <section style={{ marginBottom: 'var(--spacing-2xl)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--spacing-lg)' }}>
          <h2>🔥 Trending Movies</h2>
          <Link to="/trending" className="btn btn-outline btn-sm">
            View All
          </Link>
        </div>
        <div className="grid grid-3">
          {trendingMovies.map((movie) => (
            <ContentCard key={movie.id} item={movie} type="movie" />
          ))}
        </div>
      </section>

      {/* Trending TV Shows */}
      <section style={{ marginBottom: 'var(--spacing-2xl)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--spacing-lg)' }}>
          <h2>📺 Trending TV Shows</h2>
          <Link to="/trending" className="btn btn-outline btn-sm">
            View All
          </Link>
        </div>
        <div className="grid grid-3">
          {trendingTV.map((show) => (
            <ContentCard key={show.id} item={show} type="tv" />
          ))}
        </div>
      </section>

      {/* Features Section */}
      <section>
        <h2 style={{ marginBottom: 'var(--spacing-lg)', textAlign: 'center' }}>
          Why Choose BingeBase?
        </h2>
        <div className="grid grid-3">
          <div className="card">
            <div className="card-content">
              <h3>🎬 Huge Database</h3>
              <p>Access a vast collection of movies and TV shows from around the world.</p>
            </div>
          </div>
          <div className="card">
            <div className="card-content">
              <h3>⭐ Ratings & Reviews</h3>
              <p>See ratings from TMDB and more to help you choose what to watch.</p>
            </div>
          </div>
          <div className="card">
            <div className="card-content">
              <h3>📅 Stay Up-to-Date</h3>
              <p>Discover trending and upcoming releases every week.</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default Home 