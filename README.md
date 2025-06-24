# BingeBase - Movie/TV Show Discovery Platform

A comprehensive entertainment discovery platform where users can search for movies and TV shows, view detailed information, manage personal watchlists, and discover trending content.

## 🎬 Features

- **Search & Discovery**: Real-time search for movies and TV shows
- **Detailed Information**: Complete movie/show details with ratings, cast, and plot
- **Watchlist Management**: Add/remove titles and mark as watched
- **Trending Dashboard**: Popular movies and shows
- **Genre Filtering**: Browse by categories
- **Multi-Source Ratings**: IMDB, Rotten Tomatoes, TMDB ratings
- **Recommendations**: Personalized suggestions based on watchlist
- **Responsive Design**: Works on mobile and desktop

## Tech Stack

- **Backend**: Go (Gin framework, SQLite database)
- **Frontend**: React (Vite, vanilla CSS)
- **APIs**: TMDB, OMDB, YouTube (optional)
- **Database**: SQLite
- **State Management**: React Context API

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- API Keys for TMDB and OMDB

### Backend Setup
```bash
cd backend
go mod init binge-base
go mod tidy
cp .env.example .env
# Add your API keys to .env
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## 📁 Project Structure

```
binge-base/
├── backend/          # Go server
│   ├── main.go       # Server entry point
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Database models
│   ├── services/     # Business logic
│   ├── config/       # Configuration
│   └── database/     # Database files
├── frontend/         # React app
│   ├── src/
│   │   ├── components/   # Reusable components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API services
│   │   ├── hooks/        # Custom hooks
│   │   ├── utils/        # Utility functions
│   │   └── styles/       # CSS files
│   └── public/       # Static assets
└── docs/            # Documentation
```

## API Keys Required

1. **TMDB API**: https://www.themoviedb.org/
2. **OMDB API**: http://www.omdbapi.com/
3. **YouTube API** (optional): https://console.developers.google.com/

## Core Features Implementation

### Phase 1: Foundation
- [x] Project structure setup
- [ ] Basic Go server with Gin
- [ ] React app with Vite
- [ ] API integration setup
- [ ] Database models

### Phase 2: Core Features
- [ ] Search functionality
- [ ] Movie/TV show details
- [ ] Watchlist management
- [ ] Trending content

### Phase 3: Advanced Features
- [ ] Genre filtering
- [ ] Recommendations
- [ ] User ratings
- [ ] Responsive design

## Contributing

This project follows a collaborative development process with code reviews.

## 📝 License

MIT License 