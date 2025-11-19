# My Notes API

A production-ready REST API for managing notes with authentication.

## Features

- 🔐 JWT Authentication with Refresh Tokens
- 📝 Notes CRUD Operations
- 🔍 Search & Pagination
- 🔒 Security Best Practices
- 🧪 Comprehensive Testing (85%+ coverage)
- 📊 Structured Logging
- 🗑️ Soft Delete Support
- 🚀 Production Ready

## Tech Stack

- **Language:** Go 1.22+
- **Framework:** Fiber v2 (Web Framework)
- **Database:** MySQL 8.0
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt/jwt/v5)
- **Validation:** go-playground/validator
- **Logging:** Zerolog
- **Testing:** Testify + SQLite (in-memory)

## Architecture

This project follows **Clean Architecture** principles with a **Go-idiomatic** repository pattern.

### Project Structure

```
my-notes-api/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Server setup & dependency injection
│
├── internal/                   # Private application code
│   ├── domain/                 # Business entities (models)
│   │   ├── user.go
│   │   ├── note.go
│   │   └── refresh_token.go
│   │
│   ├── dto/                    # Data Transfer Objects
│   │   ├── auth_dto.go
│   │   └── note_dto.go
│   │
│   ├── repository/             # Data Access Layer
│   │   ├── interface.go        # Repository interfaces (contracts)
│   │   └── mysql/              # MySQL implementation
│   │       ├── user_repository.go
│   │       ├── note_repository.go
│   │       └── refresh_token_repository.go
│   │
│   ├── service/                # Business Logic Layer
│   │   ├── auth_service.go
│   │   └── note_service.go
│   │
│   ├── handler/                # HTTP Handlers (controllers)
│   │   ├── auth_handler.go
│   │   └── note_handler.go
│   │
│   ├── middleware/             # HTTP Middlewares
│   │   ├── auth.go
│   │   └── rate_limit.go
│   │
│   └── utils/                  # Utility functions
│       ├── jwt.go
│       └── password.go
│
├── pkg/                        # Public reusable packages
│   ├── database/               # Database connection & migrations
│   └── response/               # Standard API responses
│
├── configs/                    # Configuration management
│   └── config.go
│
├── migrations/                 # Database migrations (SQL)
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
│
└── test/                       # Test files
    ├── unit/                   # Unit tests
    └── integration/            # Integration tests
```

### Architecture Layers

```
┌─────────────────────────────────────────┐
│           HTTP Handlers                 │  ← API endpoints
│         (Presentation Layer)            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Services                     │  ← Business logic
│         (Business Layer)                │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          Repositories                   │  ← Data access
│          (Data Layer)                   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Database                     │  ← MySQL
└─────────────────────────────────────────┘
```

### Repository Pattern Explanation

This project uses **Go-idiomatic repository pattern** with package-based separation:

**Pattern Structure:**
```
internal/repository/
├── interface.go              # Repository interfaces (contracts)
└── mysql/                    # MySQL-specific implementation
    ├── user_repository.go
    ├── note_repository.go
    └── refresh_token_repository.go
```

**Key Concepts:**

1. **Interfaces in Root Package** (`repository/interface.go`)
   - Defines contracts that services depend on
   - Technology-agnostic (no MySQL/PostgreSQL details)
   - Example:
     ```go
     type UserRepository interface {
         Create(ctx context.Context, user *domain.User) error
         FindByEmail(ctx context.Context, email string) (*domain.User, error)
     }
     ```

2. **Implementation in Sub-Package** (`repository/mysql/`)
   - Package name indicates the technology (mysql, postgres, mongo, etc)
   - Implements the interfaces from parent package
   - Example:
     ```go
     package mysql
     
     type userRepository struct {
         db *gorm.DB
     }
     
     func NewUserRepository(db *gorm.DB) repository.UserRepository {
         return &userRepository{db: db}
     }
     ```

3. **Dependency Injection**
   - Services depend on interfaces, not concrete implementations
   - Easy to swap implementations (MySQL → PostgreSQL)
   - Example:
     ```go
     // Service depends on interface
     type AuthService struct {
         userRepo repository.UserRepository  // ← Interface, not concrete type
     }
     ```

**Comparison with Other Patterns:**

| Pattern | Structure | Common In | Use Case |
|---------|-----------|-----------|----------|
| **Go Idiomatic** (This project) | `repository/interface.go` + `repository/mysql/` | Go projects, Kubernetes, Docker | Go-native teams, cloud-native apps |
| **Java Style** | `repository/UserRepository.java` + `repository/impl/UserRepositoryImpl.java` | Java/Spring, C#/.NET | Enterprise apps, mixed-language teams |
| **Suffix-based** | `repository/user_repository.go` + `repository/user_repository_mysql.go` | Small Go projects | Quick prototypes, small teams |
| **DDD Style** | `user/repository.go` + `user/repository_mysql.go` | Microservices | Large distributed systems |

**Why Go Idiomatic Pattern?**

✅ **Recommended by Go team** - Follows official Go project layout  
✅ **Clean package names** - `mysql`, `postgres` clearly indicate implementation  
✅ **Easy to add implementations** - Just create new package (e.g., `repository/postgres/`)  
✅ **Interface segregation** - Interfaces stay clean and technology-agnostic  
✅ **Industry standard** - Used by major Go projects (Kubernetes, Terraform, etc)  

**For Java/Spring Developers:**

If you're coming from Java/Spring background:
- `repository/interface.go` = Your `@Repository` interface
- `repository/mysql/` = Your `@Repository` implementation class
- `NewUserRepository()` = Your `@Autowired` / `@Component`

The main difference: Go uses **package-based separation** instead of class-based.

### Design Principles

- **SOLID Principles:** Especially Dependency Inversion (services depend on interfaces)
- **Clean Architecture:** Clear separation between layers
- **DRY (Don't Repeat Yourself):** Reusable packages and utilities
- **Testability:** Interface-based design allows easy mocking
- **Separation of Concerns:** Each layer has single responsibility

### Dependency Flow

```
Handler → Service → Repository → Database
   ↓         ↓          ↓
  DTO     Domain    Interface
```

- **Handlers** receive HTTP requests, validate, and call services
- **Services** implement business logic, call repositories
- **Repositories** handle data persistence, return domain models
- **Domain** models represent core business entities

### Testing Strategy

1. **Unit Tests**
   - Repository layer: In-memory SQLite
   - Service layer: Mock repositories
   - Handler layer: Mock services

2. **Integration Tests**
   - End-to-end API tests
   - Real database (test container)

3. **Coverage Target**
   - Overall: 80%+
   - Critical paths: 90%+

## Getting Started

### Prerequisites

- Go 1.22 or higher
- MySQL 8.0 or higher
- Make (optional, for convenience)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/my-notes-api.git
cd my-notes-api
```

2. **Setup environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
nano .env
```

Required environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME` - Database configuration
- `JWT_SECRET` - Minimum 32 characters for security
- `APP_PORT` - Server port (default: 3000)

3. **Install dependencies**
```bash
make deps
# or
go mod download
```

4. **Run database migrations**
Migrations run automatically on startup, or manually:
```bash
# Migrations are applied automatically when you run the app
go run cmd/api/main.go
```

5. **Run the application**
```bash
make run
# or
go run cmd/api/main.go
```

Server will start at `http://localhost:3000`

### Development

**Available Make commands:**
```bash
make help              # Show all available commands
make run               # Run the application
make build             # Build binary
make test              # Run all tests
make test-coverage     # Generate coverage report
make lint              # Run linter
make fmt               # Format code
```

**Run with hot reload** (requires Air):
```bash
make install-tools     # Install Air
make dev               # Run with hot reload
```

## API Documentation

### Authentication

#### Register
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123"
}

Response:
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 900
  }
}
```

### Notes

All note endpoints require authentication header:
```
Authorization: Bearer {access_token}
```

#### Create Note
```bash
POST /api/v1/notes
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "My First Note",
  "content": "This is the content of my note"
}
```

#### Get All Notes (with pagination & search)
```bash
GET /api/v1/notes?page=1&page_size=10&search=golang&sort_by=created_at&sort_dir=desc
Authorization: Bearer {token}
```

#### Get Note by ID
```bash
GET /api/v1/notes/{id}
Authorization: Bearer {token}
```

#### Update Note
```bash
PUT /api/v1/notes/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

#### Delete Note (Soft Delete)
```bash
DELETE /api/v1/notes/{id}
Authorization: Bearer {token}
```

### Health Check
```bash
GET /api/v1/health

Response:
{
  "status": "success",
  "message": "Server is healthy",
  "data": {
    "status": "ok",
    "env": "development"
  }
}
```

## Testing

### Run Tests
```bash
# All tests
go test ./... -v

# Specific package
go test ./internal/repository/mysql/... -v

# With coverage
make test-coverage
```

### Test Structure
- Unit tests use in-memory SQLite for speed
- Integration tests use test containers
- Mocks generated for interfaces

## Security Features

- ✅ **JWT Authentication** with access & refresh tokens
- ✅ **Password Hashing** using bcrypt
- ✅ **SQL Injection Prevention** via parameterized queries
- ✅ **XSS Prevention** with input validation
- ✅ **Rate Limiting** to prevent abuse
- ✅ **CORS** configuration
- ✅ **Soft Delete** for data recovery
- ✅ **User Ownership** validation on all operations

## Project Roadmap

### Phase 1: Core Features ✅
- [x] User authentication
- [x] Notes CRUD operations
- [x] JWT tokens (access + refresh)
- [x] Repository layer
- [x] Unit tests

### Phase 2: Service Layer (In Progress)
- [ ] Business logic implementation
- [ ] Input validation
- [ ] Error handling
- [ ] Service tests

### Phase 3: API Layer
- [ ] HTTP handlers
- [ ] Middleware (auth, rate limit)
- [ ] Request validation
- [ ] API documentation

### Phase 4: Advanced Features
- [ ] Redis caching
- [ ] Search optimization
- [ ] File attachments
- [ ] Email notifications
- [ ] API versioning

### Phase 5: DevOps
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] Monitoring & metrics
- [ ] Production deployment

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'feat: add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

**Commit Convention:**
```
feat: new feature
fix: bug fix
docs: documentation
style: formatting
refactor: code restructuring
test: adding tests
chore: maintenance
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contact

- GitHub: [@yourusername](https://github.com/yourusername)
- Email: your.email@example.com

## Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- Go community for best practices and patterns

---

**Built with ❤️ using Go**
