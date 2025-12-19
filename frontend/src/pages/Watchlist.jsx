import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { useWatchlist } from '../context/WatchlistContext.jsx'
import { apiUtils } from '../services/api.js'
// import jsPDF from 'jspdf' // Removed - not installed

const Watchlist = () => {
  const { items, removeFromWatchlist, markAsWatched, markAsUnwatched } = useWatchlist()
  const [filter, setFilter] = useState('all') // 'all', 'watched', 'unwatched'

  const filteredItems = items.filter(item => {
    switch (filter) {
      case 'watched':
        return item.isWatched
      case 'unwatched':
        return !item.isWatched
      default:
        return true
    }
  })

  const stats = {
    total: items.length,
    watched: items.filter(item => item.isWatched).length,
    unwatched: items.filter(item => !item.isWatched).length
  }

  const handleRemove = (contentId, contentType) => {
    if (window.confirm('Remove this item from your watchlist?')) {
      removeFromWatchlist(contentId, contentType)
    }
  }

  const handleToggleWatched = (contentId, contentType, isWatched) => {
    if (isWatched) {
      markAsUnwatched(contentId, contentType)
    } else {
      markAsWatched(contentId, contentType)
    }
  }

  const exportCSV = () => {
    const csvRows = [
      ['Title', 'Type', 'Watched', 'Added At', 'Watched At'],
      ...items.map(item => [
        item.title || item.name,
        item.contentType,
        item.isWatched ? 'Yes' : 'No',
        item.added_at,
        item.watched_at || ''
      ])
    ]
    const csvContent = csvRows.map(e => e.join(",")).join("\n")
    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'watchlist.csv'
    a.click()
    URL.revokeObjectURL(url)
  }

  const exportPDF = () => {
    alert('PDF export feature requires additional setup. Use CSV export instead.')
  }

  const WatchlistCard = ({ item }) => (
    <div className="card">
      <img 
        src={apiUtils.getPosterURL(item.details?.poster_path || item.posterPath)} 
        alt={item.details?.title || item.details?.name || item.title}
        className="card-image"
        onError={(e) => {
          e.target.src = '/placeholder-poster.jpg'
        }}
      />
      <div className="card-content">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 'var(--spacing-sm)' }}>
          <h3 className="card-title" style={{ margin: 0, flex: 1 }}>
            {apiUtils.truncateText(item.details?.title || item.details?.name || item.title, 30)}
          </h3>
          <span style={{ 
            fontSize: 'var(--font-size-sm)', 
            color: item.isWatched ? 'var(--success-color)' : 'var(--text-muted)',
            fontWeight: 'bold'
          }}>
            {item.isWatched ? '✓ Watched' : '⏳ To Watch'}
          </span>
        </div>
        <p className="card-text">
          {apiUtils.truncateText(item.details?.overview || item.overview, 100)}
        </p>
        <div className="card-meta">
          <span>⭐ {item.details?.vote_average?.toFixed(1) || item.voteAverage?.toFixed(1) || 'N/A'}</span>
          <span>{apiUtils.formatDate(item.details?.release_date || item.releaseDate)}</span>
        </div>
        <div style={{ display: 'flex', gap: 'var(--spacing-sm)', marginTop: 'var(--spacing-sm)' }}>
          <Link 
            to={`/${item.contentType}/${item.contentId}`} 
            className="btn btn-primary btn-sm"
            style={{ flex: 1 }}
          >
            View Details
          </Link>
          <button
            onClick={() => handleToggleWatched(item.contentId, item.contentType, item.isWatched)}
            className={`btn btn-sm ${item.isWatched ? 'btn-outline' : 'btn-secondary'}`}
          >
            {item.isWatched ? 'Mark Unwatched' : 'Mark Watched'}
          </button>
          <button
            onClick={() => handleRemove(item.contentId, item.contentType)}
            className="btn btn-outline btn-sm"
            style={{ color: 'var(--error-color)', borderColor: 'var(--error-color)' }}
          >
            Remove
          </button>
        </div>
      </div>
    </div>
  )

  return (
    <div className="page">
      <div className="page-header">
        <h1 className="page-title">My Watchlist</h1>
        <p className="page-subtitle">
          Manage your personal collection of movies and TV shows
        </p>
      </div>

      <button onClick={exportCSV}>Export as CSV</button>
      <button onClick={exportPDF}>Export as PDF</button>

      {/* Stats */}
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', 
        gap: 'var(--spacing-md)', 
        marginBottom: 'var(--spacing-xl)' 
      }}>
        <div className="card">
          <div className="card-content" style={{ textAlign: 'center' }}>
            <h3 style={{ margin: 0, color: 'var(--text-primary)' }}>{stats.total}</h3>
            <p style={{ margin: 0, color: 'var(--text-secondary)' }}>Total Items</p>
          </div>
        </div>
        <div className="card">
          <div className="card-content" style={{ textAlign: 'center' }}>
            <h3 style={{ margin: 0, color: 'var(--success-color)' }}>{stats.watched}</h3>
            <p style={{ margin: 0, color: 'var(--text-secondary)' }}>Watched</p>
          </div>
        </div>
        <div className="card">
          <div className="card-content" style={{ textAlign: 'center' }}>
            <h3 style={{ margin: 0, color: 'var(--warning-color)' }}>{stats.unwatched}</h3>
            <p style={{ margin: 0, color: 'var(--text-secondary)' }}>To Watch</p>
          </div>
        </div>
      </div>

      {/* Filter Buttons */}
      <div style={{ display: 'flex', gap: 'var(--spacing-md)', marginBottom: 'var(--spacing-lg)' }}>
        <button
          className={`btn ${filter === 'all' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => setFilter('all')}
        >
          All ({stats.total})
        </button>
        <button
          className={`btn ${filter === 'watched' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => setFilter('watched')}
        >
          Watched ({stats.watched})
        </button>
        <button
          className={`btn ${filter === 'unwatched' ? 'btn-primary' : 'btn-outline'}`}
          onClick={() => setFilter('unwatched')}
        >
          To Watch ({stats.unwatched})
        </button>
      </div>

      {/* Watchlist Items */}
      {filteredItems.length > 0 ? (
        <div className="grid grid-3">
          {filteredItems.map((item) => (
            <WatchlistCard key={`${item.contentId}-${item.contentType}`} item={item} />
          ))}
        </div>
      ) : (
        <div className="text-center">
          {items.length === 0 ? (
            <>
              <h3>Your watchlist is empty</h3>
              <p>Start adding movies and TV shows to your watchlist!</p>
              <Link to="/search" className="btn btn-primary btn-lg" style={{ marginTop: 'var(--spacing-md)' }}>
                Start Searching
              </Link>
            </>
          ) : (
            <>
              <h3>No items match your filter</h3>
              <p>Try changing the filter to see more items</p>
            </>
          )}
        </div>
      )}
    </div>
  )
}

export default Watchlist 