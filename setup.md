# BingeBase Setup Guide

## 🚀 Quick Start

Follow these steps to get BingeBase running on your local machine.

### Prerequisites
- Go 1.21+ installed
- Node.js 18+ installed
- API keys for TMDB and OMDB

### Step 1: Environment Setup

1. **Copy the environment file:**
   ```bash
   cd backend
   cp env.example .env
   ```

2. **Add your API keys to `.env`:**
   ```bash
   # Edit backend/.env and add your actual API keys
   TMDB_API_KEY=your_actual_tmdb_api_key_here
   OMDB_API_KEY=your_actual_omdb_api_key_here
   ```

### Step 2: Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Initialize Go module and install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the backend server:**
   ```bash
   go run main.go
   ```

   The server should start on `http://localhost:8080` with SQLite database initialized.

### Step 3: Frontend Setup

1. **Open a new terminal and navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

   The frontend should start on `http://localhost:5173`

### Step 4: Verify Installation

1. **Backend Health Check:**
   - Visit `http://localhost:8080/api/v1/health`
   - Should return: `{"status":"ok","message":"BingeBase API is running","database":"connected"}`

2. **Frontend:**
   - Visit `http://localhost:5173`
   - Should see the BingeBase homepage

3. **Database:**
   - Check that `backend/database/bingebase.db` was created
   - This is a SQLite database with all necessary tables

## 🛠 Development Workflow

### Git Workflow (Following Instructions)

1. **Create feature branches:**
   ```bash
   git checkout -b feature/search-and-discovery
   git checkout -b feature/watchlist-management
   ```

2. **Make commits with descriptive messages:**
   ```bash
   git add .
   git commit -m "feat: implement search functionality with debouncing"
   git commit -m "feat: add watchlist management with localStorage"
   git commit -m "feat: create responsive UI components"
   git commit -m "feat: integrate TMDB and OMDB APIs"
   git commit -m "feat: add trending content dashboard"
   ```

3. **Push branches and create PRs:**
   ```bash
   git push origin feature/search-and-discovery
   git push origin feature/watchlist-management
   ```

## 📁 Project Structure

```
binge-base/
├── backend/                 # Go server (standard library)
│   ├── main.go             # Server entry point
│   ├── config/             # Configuration
│   ├── models/             # Database models
│   ├── services/           # API services
│   ├── database/           # Database operations
│   ├── handlers/           # HTTP handlers (to be implemented)
│   └── env.example         # Environment template
├── frontend/               # React app
│   ├── src/
│   │   ├── components/     # Reusable components
│   │   ├── pages/          # Page components
│   │   ├── context/        # React context
│   │   ├── services/       # API services
│   │   └── styles/         # CSS files
│   ├── package.json
│   └── vite.config.js
├── README.md
└── setup.md
```

## 🔧 Available Scripts

### Backend
- `go run main.go` - Start development server
- `go mod tidy` - Install/update dependencies

### Frontend
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## 🎯 Next Steps

### Phase 1: Core Features (Current)
- ✅ Project structure setup
- ✅ Basic routing and navigation (standard library)
- ✅ Database integration (SQLite)
- ✅ Search functionality (frontend)
- ✅ Watchlist management (frontend + backend)
- ✅ Trending content (frontend)
- ⏳ Backend API implementation with TMDB/OMDB
- ⏳ Movie/TV detail pages

### Phase 2: Advanced Features
- Movie/TV show detail pages
- Genre filtering
- Recommendations engine
- User ratings integration
- Advanced search filters

### Phase 3: Polish & Deploy
- Error handling improvements
- Performance optimizations
- Testing
- Deployment setup

## 🐛 Troubleshooting

### Common Issues

1. **Backend won't start:**
   - Check if port 8080 is available
   - Verify API keys are set in `.env`
   - Run `go mod tidy` to install dependencies
   - Check that Go 1.21+ is installed

2. **Database issues:**
   - Ensure the `backend/database/` directory exists
   - Check file permissions for database creation
   - SQLite should be automatically initialized

3. **Frontend won't start:**
   - Check if port 5173 is available
   - Run `npm install` to install dependencies
   - Clear node_modules and reinstall if needed

4. **API calls failing:**
   - Verify backend is running on port 8080
   - Check API keys are valid
   - Check browser console for CORS errors

5. **Images not loading:**
   - Check internet connection
   - Verify TMDB API key is working
   - Check browser console for errors

## 🔍 API Endpoints

### Available Endpoints (Standard Library)
- `GET /api/v1/health` - Health check
- `GET /api/v1/search?query=<search_term>` - Search content
- `GET /api/v1/search/movies?query=<search_term>` - Search movies
- `GET /api/v1/search/tv?query=<search_term>` - Search TV shows
- `GET /api/v1/movie/{id}` - Get movie details
- `GET /api/v1/tv/{id}` - Get TV show details
- `GET /api/v1/trending` - Get trending content
- `GET /api/v1/trending/movies` - Get trending movies
- `GET /api/v1/trending/tv` - Get trending TV shows
- `GET /api/v1/watchlist?user_id=<user_id>` - Get user watchlist
- `POST /api/v1/watchlist` - Add to watchlist
- `GET /api/v1/genres` - Get all genres

## 📞 Support

If you encounter issues:
1. Check the troubleshooting section above
2. Review the browser console for errors
3. Check the backend logs for API errors
4. Verify all environment variables are set correctly
5. Ensure Go 1.21+ and Node.js 18+ are installed

## 🎉 Success!

Once both servers are running, you should have a fully functional movie/TV show discovery platform with:
- **Vanilla Go backend** using standard library
- **SQLite database** for data persistence
- **React frontend** with modern UI
- **Search functionality** with debouncing
- **Watchlist management** with localStorage + database
- **Trending content** dashboard
- **Responsive design** for mobile and desktop
- **CORS handling** for cross-origin requests

Happy coding! 🎬 