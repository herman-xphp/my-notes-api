# Manual Testing Documentation

## 1. Application Startup

### Server Startup
| Test | Steps | Expected Result | Status |
|------|--------|------------------|---------|
| Start the server | Run `go run cmd/api/main.go` | Server starts without errors | ✅ Passed |
| Load configuration | Auto-loaded from `.env` | All required environment variables load successfully | ✅ Passed |

---

## 2. Database Connection

### MySQL Connection
| Test | Steps | Expected Result | Status |
|------|--------|------------------|---------|
| Connect to database | Start server | Successful connection to MySQL | ✅ Passed |
| Auto migration | Start server | Tables created/updated according to GORM models | ✅ Passed |

---

## 3. API Endpoints

### GET /health
| Steps | Expected Result | Status |
|--------|------------------|---------|
| Send GET request to `/health` | Returns `{ "status": "ok" }` with status `200` | ✅ Passed |

### GET /
| Steps | Expected Result | Status |
|--------|------------------|---------|
| Send GET request to `/` | Returns welcome message | ✅ Passed |

---

## 4. Middlewares

### CORS Middleware
| Steps | Expected Result | Status |
|--------|------------------|---------|
| Make cross-origin request | Response includes valid CORS headers | ✅ Passed |

---

## 5. Graceful Shutdown

| Steps | Expected Result | Status |
|--------|------------------|---------|
| Press `Ctrl+C` while server is running | Server shuts down gracefully and closes DB connections | ✅ Passed |

---

## Summary

All basic functionality has been manually tested and works as expected:

- ✅ Application starts successfully  
- ✅ Database connection established  
- ✅ Migrations run automatically  
- ✅ Health check endpoint working  
- ✅ Welcome endpoint working  
- ✅ CORS middleware active  
- ✅ Graceful shutdown working  

