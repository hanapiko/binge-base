import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { contentAPI, apiUtils } from '../services/api.js'
import { useWatchlist } from '../context/WatchlistContext.jsx'

const MovieDetails = () => {
  const { id } = useParams()
  const [movie, setMovie] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const watchlist = useWatchlist()

  useEffect(() => {
    const fetchMovie = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await contentAPI.getMovieDetails(id)
        setMovie(response.data.data)
      } catch (err) {
        setError('Failed to load movie details.')
      } finally {
        setLoading(false)
      }
    }
    fetchMovie()
  }, [id])

  const shareUrl = window.location.href
  const handleCopy = () => {
    navigator.clipboard.writeText(shareUrl)
    alert('Link copied!')
  }
  const handleShareTwitter = () => {
    window.open(`https://twitter.com/intent/tweet?url=${encodeURIComponent(shareUrl)}&text=Check%20out%20this%20movie!`, '_blank')
  }

  if (loading) {
    return (
      <div className="page">
        <div className="text-center">
          <div className="loading-spinner"></div>
          <p>Loading movie details...</p>
        </div>
      </div>
    )
  }
  if (error) {
    return (
      <div className="page">
        <div className="error">{error}</div>
      </div>
    )
  }
  if (!movie) return null

  const isInWatchlist = watchlist.isInWatchlist(Number(id), 'movie')
  const handleAdd = () => watchlist.addToWatchlist({ id: Number(id), media_type: 'movie' })
  const handleRemove = () => watchlist.removeFromWatchlist(Number(id), 'movie')

  console.log('Movie details:', movie)

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">{movie.title}</h1>
        <p className="page-subtitle">{movie.tagline}</p>
      </div>
      <div className="details-grid">
        <img
          src={apiUtils.getPosterURL(movie.poster_path)}
          alt={movie.title}
          className="details-poster"
          onError={e => { e.target.src = '/placeholder-poster.jpg' }}
        />
        <div className="details-content">
          <h2>Overview</h2>
          <p>{movie.overview}</p>
          <div className="details-meta">
            <div><strong>Release Date:</strong> {apiUtils.formatDate(movie.release_date)}</div>
            <div><strong>Runtime:</strong> {apiUtils.formatRuntime(movie.runtime)}</div>
            <div><strong>Rating:</strong> ⭐ {movie.vote_average?.toFixed(1) || 'N/A'}</div>
            <div><strong>Genres:</strong> {movie.genres?.map(g => g.name).join(', ')}</div>
            <div><strong>Status:</strong> {movie.status}</div>
            <div><strong>Budget:</strong> ${movie.budget?.toLocaleString()}</div>
            <div><strong>Revenue:</strong> ${movie.revenue?.toLocaleString()}</div>
          </div>
          {/* Add more details as needed, e.g., cast, ratings from OMDB, etc. */}
          <div className="watchlist-action">
            {watchlist.loading ? (
              <button disabled>Updating...</button>
            ) : isInWatchlist ? (
              <button onClick={handleRemove} className="btn btn-secondary">Remove from Watchlist</button>
            ) : (
              <button onClick={handleAdd} className="btn btn-primary">Add to Watchlist</button>
            )}
          </div>
        </div>
      </div>
      <div className="share-buttons">
        <button onClick={handleCopy}>Copy Link</button>
        <button onClick={handleShareTwitter}>Share on Twitter</button>
      </div>
      {movie.providers && Array.isArray(movie.providers.US?.flatrate) ? (
        <div className="providers">
          <h4>Watch on:</h4>
          <ul>
            {movie.providers.US.flatrate.map(p => <li key={p.provider_id}>{p.provider_name}</li>)}
          </ul>
        </div>
      ) : (
        <div className="providers"><em>No streaming providers found.</em></div>
      )}
      {movie.trailer ? (
        <div className="trailer">
          <h4>Trailer:</h4>
          <iframe width="560" height="315" src={`https://www.youtube.com/embed/${movie.trailer}`} frameBorder="0" allowFullScreen></iframe>
        </div>
      ) : (
        <div className="trailer"><em>No trailer available.</em></div>
      )}
    </div>
  )
}

export default MovieDetails 