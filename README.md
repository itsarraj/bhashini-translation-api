# Translation Service

A production-ready Go Fiber-based translation service that integrates with [Bhashini API](https://bhashini.gov.in/) for translating content across Indian languages with intelligent caching.

## 🚀 Features

- ✅ **Bhashini API Integration**: Seamless integration with Bhashini translation API
- ✅ **Intelligent Caching**: Prevents redundant API calls for the same translation
- ✅ **Configurable TTL**: Customizable cache expiration time (default: 24 hours)
- ✅ **PostgreSQL Backend**: Efficient cache storage using PostgreSQL
- ✅ **RESTful API**: Clean and simple API endpoints
- ✅ **Production Ready**: Error handling, logging, and health checks

## 📋 Prerequisites

- **Go**: 1.23 or higher
- **PostgreSQL**: 12+ (for cache storage)
- **Bhashini Account**: API credentials from [Bhashini Dashboard](https://bhashini.gov.in/ulca/dashboard)

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd user-service
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Set Up Environment Variables

Copy `env.example` to `.env` and update with your credentials:

```bash
cp env.example .env
```

Edit `.env` file:

```env
# Database Configuration
DATABASE_URL=postgres://user:password@host:5432/database?sslmode=disable

# Server Configuration
PORT=3001

# Bhashini API Configuration
BHASHINI_BASE_URL=https://meity-auth.ulcacontrib.org
BHASHINI_USER_ID=your_user_id_here
BHASHINI_API_KEY=your_ulcaApiKey_here
BHASHINI_PIPELINE_ID=64392f96daac500b55c543cd  # Optional

# Translation Cache Configuration
TRANSLATION_CACHE_TTL=24h  # Cache TTL (default: 24h)
```

### 4. Get Bhashini API Credentials

1. **Login to Bhashini Dashboard**: https://bhashini.gov.in/ulca/dashboard
2. **Navigate to "My Profile"** section
3. **Copy your credentials**:
   - `userID` → Set as `BHASHINI_USER_ID`
   - `ulcaApiKey` → Set as `BHASHINI_API_KEY` (⚠️ Must be the `ulcaApiKey`, not other API keys)

### 5. Run Database Migration

```bash
psql $DATABASE_URL -f migrations/003_create_translation_cache.sql
```

### 6. Start the Service

```bash
# Development
go run cmd/main.go

# Production
go build -o translation-service cmd/main.go
./translation-service
```

The service will start on `http://localhost:3001` (or your configured PORT).

## 📖 API Documentation

### Base URL

```
http://localhost:3001/api/v1/translation
```

### Translate Text

Translate text from one language to another.

**Endpoint:** `POST /translate`

**Request:**
```bash
curl -X POST http://localhost:3001/api/v1/translation/translate \
  -H "Content-Type: application/json" \
  -d '{
    "source_text": "Hello, how are you?",
    "source_lang": "en",
    "target_lang": "hi"
  }'
```

**Request Body:**
```json
{
  "source_text": "Hello, how are you?",
  "source_lang": "en",
  "target_lang": "hi"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "data": {
    "source_text": "Hello, how are you?",
    "source_lang": "en",
    "target_lang": "hi",
    "translated_text": "नमस्कार, आप कैसे हैं?"
  }
}
```

**Error Response (400/500):**
```json
{
  "status": "error",
  "error": "error message"
}
```

### Clean Cache

Remove expired cache entries manually.

**Endpoint:** `POST /cache/clean`

**Request:**
```bash
curl -X POST http://localhost:3001/api/v1/translation/cache/clean
```

**Response:**
```json
{
  "status": "success",
  "message": "Expired cache entries cleaned"
}
```

### Health Check

Check service health and status.

**Endpoint:** `GET /health`

**Request:**
```bash
curl http://localhost:3001/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "translation-service"
}
```

## 🌐 Supported Languages

Bhashini supports translation between multiple Indian languages. Common language codes (ISO-639):

| Code | Language | Code | Language |
|------|----------|------|----------|
| `en` | English | `hi` | Hindi |
| `mr` | Marathi | `ta` | Tamil |
| `te` | Telugu | `kn` | Kannada |
| `gu` | Gujarati | `bn` | Bengali |
| `pa` | Punjabi | `or` | Odia |
| `ml` | Malayalam | `as` | Assamese |
| `ur` | Urdu | `ne` | Nepali |

For a complete list, refer to [ISO-639 series](https://www.loc.gov/standards/iso639-2/php/code_list.php).

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `PORT` | Server port | No | `3001` |
| `BHASHINI_BASE_URL` | Bhashini API base URL | No | `https://meity-auth.ulcacontrib.org` |
| `BHASHINI_USER_ID` | Bhashini user ID from dashboard | Yes | - |
| `BHASHINI_API_KEY` | Bhashini ulcaApiKey from dashboard | Yes | - |
| `BHASHINI_PIPELINE_ID` | Pipeline ID for translation | No | `64392f96daac500b55c543cd` |
| `TRANSLATION_CACHE_TTL` | Cache TTL duration | No | `24h` |

### Cache Configuration

The service implements intelligent caching:

- **Cache Key**: Based on `source_text`, `source_lang`, and `target_lang`
- **TTL**: Configurable via `TRANSLATION_CACHE_TTL` (supports Go duration format: `24h`, `1h30m`, etc.)
- **Storage**: PostgreSQL table `translation_cache`
- **Cleanup**: Expired entries can be cleaned manually via `/cache/clean` endpoint

## 🗄️ Database Schema

### translation_cache

```sql
CREATE TABLE translation_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_text TEXT NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    translated_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    UNIQUE(source_text, source_lang, target_lang)
);

-- Indexes for faster lookups
CREATE INDEX idx_translation_cache_lookup ON translation_cache(source_text, source_lang, target_lang, expires_at);
CREATE INDEX idx_translation_cache_expires ON translation_cache(expires_at);
```

## 📁 Project Structure

```
user-service/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── db/
│   │   └── db.go                  # Database connection
│   ├── handlers/
│   │   └── translation_handler.go # HTTP handlers
│   ├── models/
│   │   └── bhashini_models.go     # API models
│   ├── repository/
│   │   └── translation_repo.go   # Cache repository
│   ├── router/
│   │   └── router.go              # Route setup
│   └── services/
│       ├── bhashini_client.go     # Bhashini API client
│       └── translation_service.go # Translation business logic
├── migrations/
│   └── 003_create_translation_cache.sql
├── .env                           # Environment variables (not in git)
├── env.example                    # Environment template
├── go.mod                         # Go dependencies
├── go.sum                         # Go checksums
├── Dockerfile                     # Docker configuration
├── test_api_key.sh                # API key testing script
└── README.md                      # This file
```

## 🐳 Docker

### Build Docker Image

```bash
docker build -t translation-service .
```

### Run with Docker

```bash
docker run -p 3001:3001 \
  -e DATABASE_URL=postgres://user:password@host:5432/database \
  -e BHASHINI_USER_ID=your_user_id \
  -e BHASHINI_API_KEY=your_api_key \
  translation-service
```

## 🧪 Testing

### Test Translation API

```bash
curl -X POST http://localhost:3001/api/v1/translation/translate \
  -H "Content-Type: application/json" \
  -d '{
    "source_text": "Hello",
    "source_lang": "en",
    "target_lang": "hi"
  }'
```

## 🔧 Troubleshooting

### Database Connection Issues

**Error**: `DB connection failed`

**Solutions**:
1. Verify `DATABASE_URL` is correct
2. Check database is running: `pg_isready`
3. Verify network connectivity
4. Check SSL mode settings in connection string

### Bhashini API Issues

**Error**: `ulcaApiKey does not exist`

**Solutions**:
1. ✅ **Most Common**: Ensure you're using the `ulcaApiKey` from "My Profile" section, not other API keys
2. Verify API key is active/enabled in dashboard
3. Check for extra spaces or quotes in `.env` file
4. Regenerate API key if expired

**Error**: `API returned status 400/500`

**Solutions**:
1. Verify API credentials are correct
2. Check if pipeline ID is valid (use default: `64392f96daac500b55c543cd`)
3. Ensure source and target language codes are valid ISO-639 codes
4. Check Bhashini API status

### Cache Issues

**Error**: `relation "translation_cache" does not exist`

**Solution**: Run the migration:
```bash
psql $DATABASE_URL -f migrations/003_create_translation_cache.sql
```

**Error**: Cache lookup/storage errors

**Solution**: These are logged but don't fail the request. Check database connectivity and ensure the table exists.

### Port Already in Use

**Error**: `bind: address already in use`

**Solution**: 
1. Change `PORT` in `.env` file, or
2. Kill the process using the port:
```bash
lsof -ti:3001 | xargs kill -9
```

## 📚 Additional Resources

- **Bhashini API Documentation**: https://dibd-bhashini.gitbook.io/bhashini-apis
- **Bhashini Dashboard**: https://bhashini.gov.in/ulca/dashboard
- **ISO-639 Language Codes**: https://www.loc.gov/standards/iso639-2/php/code_list.php

## 🏗️ Development

### Running in Development Mode

```bash
go run cmd/main.go
```

### Building for Production

```bash
go build -o translation-service cmd/main.go
```

### Code Structure

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and API integration
- **Repository**: Database operations
- **Models**: Data structures and API models

## 📝 License

[Your License Here]

## 🤝 Support

For issues and questions:
- Check the [Troubleshooting](#-troubleshooting) section
- Review Bhashini API documentation
- Contact the development team

---

**Built with ❤️ using Go Fiber and Bhashini API**
