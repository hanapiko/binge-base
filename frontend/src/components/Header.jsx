import React, { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useWatchlist } from '../context/WatchlistContext.jsx'
import { useTheme } from '../App.jsx'

const Header = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const location = useLocation()
  const { getWatchlistStats } = useWatchlist()
  const stats = getWatchlistStats()
  const { theme, toggleTheme } = useTheme()

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen)
  }

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false)
  }

  const isActive = (path) => {
    return location.pathname === path
  }

  return (
    <header className="header">
      <div className="header-container">
        <Link to="/" className="logo" onClick={closeMobileMenu}>
          🎬 BingeBase
        </Link>
        
        <nav className={`nav ${isMobileMenuOpen ? 'mobile-open' : ''}`}>
          <Link 
            to="/" 
            className={`nav-link ${isActive('/') ? 'active' : ''}`}
            onClick={closeMobileMenu}
          >
            Home
          </Link>
          
          <Link 
            to="/search" 
            className={`nav-link ${isActive('/search') ? 'active' : ''}`}
            onClick={closeMobileMenu}
          >
            Search
          </Link>
          
          <Link 
            to="/trending" 
            className={`nav-link ${isActive('/trending') ? 'active' : ''}`}
            onClick={closeMobileMenu}
          >
            Trending
          </Link>
          
          <Link 
            to="/watchlist" 
            className={`nav-link ${isActive('/watchlist') ? 'active' : ''}`}
            onClick={closeMobileMenu}
          >
            Watchlist
            {stats.total > 0 && (
              <span className="watchlist-badge">
                {stats.total}
              </span>
            )}
          </Link>
        </nav>
        
        <button 
          className="mobile-menu-toggle"
          onClick={toggleMobileMenu}
          aria-label="Toggle mobile menu"
        >
          {isMobileMenuOpen ? '✕' : '☰'}
        </button>
        <button onClick={toggleTheme} style={{ fontSize: '1.5rem', padding: '0.5rem 1rem' }}>
          {theme === 'light' ? '🌞 Light' : '🌙 Dark'}
        </button>
      </div>
    </header>
  )
}

export default Header 