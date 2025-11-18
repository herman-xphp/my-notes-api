# My Notes API

A production-ready REST API for managing notes with authentication.

## Features

- 🔐 JWT Authentication
- 📝 Notes CRUD Operations
- 🔒 Security Best Practices
- 🧪 Comprehensive Testing
- 📊 Structured Logging
- 🚀 Production Ready

## Tech Stack

- Go 1.22+
- Fiber v2 (Web Framework)
- MySQL 8.0 (Database)
- GORM (ORM)
- JWT (Authentication)

## Getting Started

### Prerequisites

- Go 1.22 or higher
- MySQL 8.0 or higher
- Make (optional)

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/my-notes-api.git
cd my-notes-api
```

2. Setup environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies
```bash
make deps
```

4. Run the application
```bash
make run
```

## Development

See [Makefile](Makefile) for available commands.

## Project Structure
```
my-notes-api/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── domain/           # Business entities
│   ├── dto/              # Data transfer objects
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middlewares
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── utils/            # Utility functions
├── pkg/                  # Reusable packages
├── configs/              # Configuration
├── migrations/           # Database migrations
└── test/                 # Tests
```

## License

MIT License
