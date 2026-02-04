# gRPC Fintech Project Documentation

Go that provides user authentication and wallet management services. The project uses Protocol Buffers for service definitions, PostgreSQL for data persistence, and Docker for containerization.

Architecture
Technology Stack

- Language: Go 1.25.5
- Framework: gRPC with Protocol Buffers
- Database: PostgreSQL 14
- ORM: GORM
- Authentication: JWT (golang-jwt)
- Hashing: Argon2 with base64 encoding
- Containerization: Docker & Docker Compose

```
grpc-fintech/
├── cmd/grpcapi/binaries/ # Server entry point
│ └── server.go
├── database/ # Database configuration
│ └── db.go
├── internal/
│ ├── api/
│ │ ├── handlers/ # gRPC service implementations
│ │ │ ├── users.go
│ │ │ ├── wallet.go
│ │ │ └── server_structs.go
│ │ └── interceptors/ # Middleware
│ │ ├── authentication.go
│ │ └── response_time.go
│ └── models/ # Data models
│ ├── users.go
│ ├── wallet.go
│ ├── wallet_history.go
│ └── transactions.go
├── proto/ # Protocol Buffer definitions
│ ├── users.proto
│ ├── wallet.proto
│ └── gen/ # Generated protobuf code
├── utils/ # Utility functions
│ ├── jwt.go
│ ├── password.go
│ └── error_handlers.go
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

Core Services

1. User Service (users.proto)
   Manages user authentication and account operations.

RPCs
| Method | Request | Response | Description |
| ----------------- | ------------------| ----------------- | ------------------|
| CreateUser | CreateUserRequest | UserResponse | Register a new user |
| LoginUser | LoginRequest | LoginResponse | Authenticate user and get JWT token |
| UserProfile | UserIdRequest | UserResponse | Fetch authenticated user profile |
| UpdateProfile | UserID | UserResponse | Update user profile |
| DeactivateAccount | UserID | DeactivateResponse| Deactivate user account |

Key Features

- Password hashing using Argon2 with salt (base64 encoded)
- JWT token generation and validation
- User status management

RPCs

| Method             | Request            | Response                   | Description                |
| ------------------ | ------------------ | -------------------------- | -------------------------- |
| Balance            | WalletIdRequest    | WalletResponse             | Get wallet balance         |
| Deposit            | DepositRequest     | DepositResponse            | Deposit money into wallet  |
| Withdrawl          | WithdrawlRequest   | WithdrawlResponse          | Withdraw money from wallet |
| TransactionHistory | TransactionRequest | TransactionHistoryResponse | Fetch transaction history  |

Key Features

- Real-time balance updates
- Transaction history tracking
- Wallet status management (active, locked, suspended)
- Atomic wallet balance updates with versioning

```
type User struct {
    ID int // Primary key
    UUID string // Unique identifier
    First_name string
    Last_name string
    Password string // Hashed
    PhoneNumber string
    Email string // Unique
    Status bool
    UpdatedAt time.Time
    CreatedAt time.Time
}
```

```
type Wallet struct {
    ID int
    UUID string // Unique
    UserID string // Foreign key to User
    Country string
    AvailableBalance float32
    Status string // active, locked, unlocked, suspended
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

WalletHistory Model
Tracks every balance change for audit purposes.

Transaction Model
Records all wallet transactions (deposits/withdrawals).

Authentication & Security
JWT Implementation
File: jwt.go

Algorithm: HS256 (HMAC SHA-256)
Lifespan: 24 hours
Claims: user_id, email, authorized
Secret: Read from JWT_SECRET environment variable
Password Hashing
File: password.go

Algorithm: Argon2ID
Parameters:
Time cost: 1
Memory: 64MB
Parallelism: 4
Hash length: 32 bytes
Storage Format: base64(salt).base64(hash)
Request Interceptors
File: authentication.go

Validates JWT tokens from authorization metadata header
Skips auth for public methods: Login, Register, ForgotPassword, ResetPassword
Adds user context to request
File: response_time.go

Measures RPC execution time
Returns X-Response-Time header
Database Configuration
File: db.go

Connection String Format
Environment Variables
Auto-Migration
Automatically creates/updates tables for:

users
wallets
wallet_histories
transactions
Deployment

```
# Start services
make up

# View logs
make logs

# Stop services
make down

# Check status
make ps
```

Docker Compose
File: docker-compose.yml

Services:

gRPC Server: Port 50052
PostgreSQL: Port 5432
Build & Run
Dockerfile
File: Dockerfile

Multi-stage build:

Builder stage: Compiles Go binary
Runtime stage: Alpine Linux for minimal image size
API Usage Examples
Create User
Login
Check Balance
Deposit
Error Handling
File: error_handlers.go

Global error logging and formatting:

gRPC error codes used:

Unauthenticated: Invalid/missing JWT

- PermissionDenied: Insufficient permissions
- Aborted: Transaction failed
- AlreadyExists: Duplicate resource
- Internal: Server error
- Key Implementation Notes
- Wallet Balance Updates
- Uses database transactions to ensure atomicity
- Creates WalletHistory records for audit trail
- Prevents concurrent balance inconsistencies
- Transaction History

Unimplemented Methods

- UpdateProfile: Returns placeholder response
- Interceptors: Not registered in main server setup
- Development Roadmap
- Register interceptors in gRPC server
- Implement UpdateProfile handler
- Fix transaction history query
- Add unit tests
- Implement rate limiting
- Add request validation middleware
- Implement transaction rollback logic
- Add comprehensive logging
- Dependencies
- See go.mod for complete list. Key dependencies:

Future Enhancements

- Multi-currency support: Add currency field to wallets
- Rate limiting: Prevent abuse
- Audit logging: Comprehensive transaction logging
- Transaction rollback: Handle failed operations
- Webhook notifications: Real-time balance updates
- Mobile app integration: Expose REST API gateway
- Enhanced security: 2FA, biometric auth
- Analytics: Transaction analytics dashboard
