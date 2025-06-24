import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { contentAPI, apiUtils } from '../services/api.js'

const TVDetails = () => {
  const { id } = useParams()
  const [tv, setTV] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchTV = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await contentAPI.getTVDetails(id)
        setTV(response.data.data)
      } catch (err) {
        setError('Failed to load TV show details.')
      } finally {
        setLoading(false)
      }
    }
    fetchTV()
  }, [id])

  if (loading) {
    return (
      <div className="page">
        <div className="text-center">
          <div className="loading-spinner"></div>
          <p>Loading TV show details...</p>
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
  if (!tv) return null

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">{tv.name}</h1>
        <p className="page-subtitle">{tv.tagline}</p>
      </div>
      <div className="details-grid">
        <img
          src={apiUtils.getPosterURL(tv.poster_path)}
          alt={tv.name}
          className="details-poster"
          onError={e => { e.target.src = '/placeholder-poster.jpg' }}
        />
        <div className="details-content">
          <h2>Overview</h2>
          <p>{tv.overview}</p>
          <div className="details-meta">
            <div><strong>First Air Date:</strong> {apiUtils.formatDate(tv.first_air_date)}</div>
            <div><strong>Last Air Date:</strong> {apiUtils.formatDate(tv.last_air_date)}</div>
            <div><strong>Rating:</strong> ⭐ {tv.vote_average?.toFixed(1) || 'N/A'}</div>
            <div><strong>Genres:</strong> {tv.genres?.map(g => g.name).join(', ')}</div>
            <div><strong>Status:</strong> {tv.status}</div>
            <div><strong>Seasons:</strong> {tv.number_of_seasons}</div>
            <div><strong>Episodes:</strong> {tv.number_of_episodes}</div>
          </div>
          {/* Add more details as needed, e.g., cast, ratings from OMDB, etc. */}
        </div>
      </div>
    </div>
  )
}

export default TVDetails 