import React from 'react'
import { Link } from 'react-router-dom'

const Footer = () => {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="footer">
      <div className="footer-container">
        <div className="footer-content">
          <div className="footer-section">
            <h3>🎬 BingeBase</h3>
            <p>
              Your ultimate destination for discovering movies and TV shows. 
              Find your next favorite entertainment with our comprehensive search 
              and personalized recommendations.
            </p>
          </div>
          
          <div className="footer-section">
            <h3>Quick Links</h3>
            <Link to="/">Home</Link>
            <Link to="/search">Search</Link>
            <Link to="/trending">Trending</Link>
            <Link to="/watchlist">Watchlist</Link>
          </div>
          
          <div className="footer-section">
            <h3>Features</h3>
            <p>Movie & TV Search</p>
            <p>Personal Watchlist</p>
            <p>Trending Content</p>
            <p>Genre Filtering</p>
            <p>Multi-Source Ratings</p>
          </div>
          
          <div className="footer-section">
            <h3>Resources</h3>
            <a href="https://www.themoviedb.org/" target="_blank" rel="noopener noreferrer">
              TMDB API
            </a>
            <a href="http://www.omdbapi.com/" target="_blank" rel="noopener noreferrer">
              OMDB API
            </a>
            <a href="https://github.com" target="_blank" rel="noopener noreferrer">
              GitHub Repository
            </a>
            <p>Documentation</p>
          </div>
        </div>
        
        <div className="footer-bottom">
          <p>
            © {currentYear} BingeBase. Built with React, Go, and powered by TMDB & OMDB APIs.
            This project is for educational purposes and follows collaborative development practices.
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer 