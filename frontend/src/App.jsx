import React, { useState, useEffect, createContext, useContext } from 'react'
import { Routes, Route } from 'react-router-dom'
import Header from './components/Header.jsx'
import Footer from './components/Footer.jsx'
import Home from './pages/Home.jsx'
import Search from './pages/Search.jsx'
import MovieDetails from './pages/MovieDetails.jsx'
import TVDetails from './pages/TVDetails.jsx'
import Watchlist from './pages/Watchlist.jsx'
import Trending from './pages/Trending.jsx'
import './styles/App.css'

const ThemeContext = createContext()

export const useTheme = () => useContext(ThemeContext)

const App = () => {
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'light')

  useEffect(() => {
    document.body.className = theme
    localStorage.setItem('theme', theme)
  }, [theme])

  const toggleTheme = () => setTheme(theme === 'light' ? 'dark' : 'light')

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      <div className="app">
        <Header />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/search" element={<Search />} />
            <Route path="/movie/:id" element={<MovieDetails />} />
            <Route path="/tv/:id" element={<TVDetails />} />
            <Route path="/watchlist" element={<Watchlist />} />
            <Route path="/trending" element={<Trending />} />
          </Routes>
        </main>
        <Footer />
      </div>
    </ThemeContext.Provider>
  )
}

export default App 