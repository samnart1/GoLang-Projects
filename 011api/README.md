```
blog-api/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go               # Authentication handlers
│   │   │   ├── posts.go              # Blog post handlers
│   │   │   ├── users.go              # User handlers
│   │   │   └── health.go             # Health check handlers
│   │   ├── middleware/
│   │   │   ├── auth.go               # Authentication middleware
│   │   │   ├── cors.go               # CORS middleware
│   │   │   ├── logging.go            # Request logging middleware
│   │   │   └── ratelimit.go          # Rate limiting middleware
│   │   └── routes/
│   │       └── routes.go             # Route definitions
│   ├── config/
│   │   └── config.go                 # Configuration management
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── user.go               # User entity
│   │   │   ├── post.go               # Post entity
│   │   │   └── comment.go            # Comment entity
│   │   ├── repositories/
│   │   │   ├── user_repository.go    # User repository interface
│   │   │   └── post_repository.go    # Post repository interface
│   │   └── services/
│   │       ├── user_service.go       # User business logic
│   │       ├── post_service.go       # Post business logic
│   │       └── auth_service.go       # Authentication business logic
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── migrations/
│   │   │   │   ├── 001_create_users_table.up.sql
│   │   │   │   ├── 001_create_users_table.down.sql
│   │   │   │   ├── 002_create_posts_table.up.sql
│   │   │   │   └── 002_create_posts_table.down.sql
│   │   │   ├── connection.go         # Database connection setup
│   │   │   └── postgres.go           # PostgreSQL specific implementation
│   │   ├── repositories/
│   │   │   ├── postgres_user_repo.go # PostgreSQL user repository
│   │   │   └── postgres_post_repo.go # PostgreSQL post repository
│   │   └── cache/
│   │       └── redis.go              # Redis cache implementation
│   └── utils/
│       ├── validator.go              # Input validation utilities
│       ├── jwt.go                    # JWT token utilities
│       ├── password.go               # Password hashing utilities
│       └── response.go               # Standard response utilities
├── pkg/
│   ├── logger/
│   │   └── logger.go                 # Logging package
│   ├── errors/
│   │   └── errors.go                 # Custom error types
│   └── pagination/
│       └── pagination.go             # Pagination utilities
├── api/
│   └── openapi/
│       └── swagger.yaml              # OpenAPI/Swagger specification
├── docs/
│   ├── README.md                     # Project documentation
│   ├── API.md                        # API documentation
│   └── DEPLOYMENT.md                 # Deployment instructions
├── scripts/
│   ├── migrate.sh                    # Database migration script
│   ├── seed.sh                       # Database seeding script
│   └── build.sh                      # Build script
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── tests/
│   ├── integration/
│   │   ├── auth_test.go
│   │   └── posts_test.go
│   ├── unit/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── repositories/
│   └── fixtures/
│       └── test_data.json
├── configs/
│   ├── config.yaml                   # Default configuration
│   ├── config.dev.yaml               # Development configuration
│   └── config.prod.yaml              # Production configuration
├── .env.example                      # Environment variables example
├── .gitignore                        # Git ignore file
├── go.mod                            # Go module file
├── go.sum                            # Go module checksums
├── Makefile                          # Build and development commands
└── README.md                         # Project README
```