import React, { createContext, useContext, useReducer, useEffect } from 'react'
import { watchlistAPI } from '../services/api.js'

const WatchlistContext = createContext()

const initialState = {
  items: [],
  loading: false,
  error: null
}

const watchlistReducer = (state, action) => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload }
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false }
    case 'SET_WATCHLIST':
      return { ...state, items: action.payload, loading: false, error: null }
    default:
      return state
  }
}

export const WatchlistProvider = ({ children }) => {
  const [state, dispatch] = useReducer(watchlistReducer, initialState)
  const userId = 'default_user'

  // Load watchlist from backend on mount
  useEffect(() => {
    const fetchWatchlist = async () => {
      dispatch({ type: 'SET_LOADING', payload: true })
      try {
        const response = await watchlistAPI.getWatchlist(userId)
        dispatch({ type: 'SET_WATCHLIST', payload: response.data.data || [] })
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to load watchlist' })
      }
    }
    fetchWatchlist()
  }, [])

  const addToWatchlist = async (content) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    try {
      await watchlistAPI.addToWatchlist(userId, content.id, content.media_type || 'movie')
      // Refresh watchlist
      const response = await watchlistAPI.getWatchlist(userId)
      dispatch({ type: 'SET_WATCHLIST', payload: response.data.data || [] })
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to add to watchlist' })
    }
  }

  const removeFromWatchlist = async (contentId, contentType) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    try {
      await watchlistAPI.removeFromWatchlist(userId, contentId, contentType)
      // Refresh watchlist
      const response = await watchlistAPI.getWatchlist(userId)
      dispatch({ type: 'SET_WATCHLIST', payload: response.data.data || [] })
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to remove from watchlist' })
    }
  }

  const markAsWatched = async (contentId, contentType) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    try {
      await watchlistAPI.markAsWatched(userId, contentId, contentType, true)
      // Refresh watchlist
      const response = await watchlistAPI.getWatchlist(userId)
      dispatch({ type: 'SET_WATCHLIST', payload: response.data.data || [] })
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to mark as watched' })
    }
  }

  const markAsUnwatched = async (contentId, contentType) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    try {
      await watchlistAPI.markAsWatched(userId, contentId, contentType, false)
      // Refresh watchlist
      const response = await watchlistAPI.getWatchlist(userId)
      dispatch({ type: 'SET_WATCHLIST', payload: response.data.data || [] })
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to mark as unwatched' })
    }
  }

  const isInWatchlist = (contentId, contentType) => {
    return state.items.some(item => 
      item.contentId === contentId && item.contentType === contentType
    )
  }

  const getWatchlistStats = () => {
    const total = state.items.length
    const watched = state.items.filter(item => item.isWatched).length
    const unwatched = total - watched
    return { total, watched, unwatched }
  }

  const value = {
    ...state,
    addToWatchlist,
    removeFromWatchlist,
    markAsWatched,
    markAsUnwatched,
    isInWatchlist,
    getWatchlistStats
  }

  return (
    <WatchlistContext.Provider value={value}>
      {children}
    </WatchlistContext.Provider>
  )
}

export const useWatchlist = () => {
  const context = useContext(WatchlistContext)
  if (!context) {
    throw new Error('useWatchlist must be used within a WatchlistProvider')
  }
  return context
} 