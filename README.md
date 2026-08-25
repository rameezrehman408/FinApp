# FinApp - Financial Application

## Overview

FinApp is a web-based financial management application built with Next.js (React) and a Go backend. The application enables users to manage their finances, track assets, investments, budgets, and transactions.

## Architecture

### Tech Stack
- **Frontend**: Next.js 14, React 18, TypeScript
- **Backend**: Go (Gin framework), PostgreSQL, Redis, MinIO
- **Deployment**: Docker Compose

### Design Patterns

#### 1. Component-Based Architecture
- Each UI component is self-contained and reusable
- Clear separation between presentational and container components
- Example: `AccountForm` vs `AccountFormContainer`

#### 2. Repository Pattern
- Data access is decoupled from business logic
- Interfaces for data operations (AccountRepository, TransactionRepository, etc.)
- Allows swapping persistence layers (DB → cache → file)

#### 3. Service Layer
- Business logic is encapsulated in service classes
- Handles complex operations like budget calculations
- Example: `BudgetService.calculateUtilization()`

#### 4. CQRS (Command Query Responsibility Segregation)
- Separates read (Query) and write (Command) operations
- API endpoints differentiate between GET (queries) and POST (commands)
- Enables potential scaling of reads vs writes

#### 5. Dependency Injection
- Configuration via environment variables
- Services injected through constructors
- Testable and maintainable code structure

#### 6. State Management
- React Context API for global state
- Feature-specific contexts (auth, transactions, budgets)
- Type-safe state with TypeScript interfaces

## Project Structure

```
~/repos/FinApp/
├── web/                          # Frontend (Next.js)
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx         # Root layout
│   │   │   ├── (auth)/           # Authentication components
│   │   │   ├── (dashboard)/      # Main dashboard UI
│   │   │   ├── api/             # API service layer
│   │   │   │   ├── accounts/
│   │   │   │   ├── transactions/
│   │   │   │   ├── budgets/
│   │   │   │   └── investments/
│   │   │   ├── components/       # Reusable UI components
│   │   │   ├── hooks/           # Custom hooks
│   │   │   └── store/           # State management
│   │   └── types/               # TypeScript interfaces
│   ├── package.json
│   └── tsconfig.json
├── api/                          # Backend (Go)
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   └── Dockerfile
└── docker-compose.yml
```

## Getting Started

1. **Start the backend**:
   ```bash
   cd ~/repos/FinApp/api
   docker-compose up -d postgres redis minio finapp-api
   ```

2. **Start the frontend**:
   ```bash
   cd ~/repos/FinApp/web
   npm run dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/accounts` | List all accounts |
| POST | `/api/v1/accounts` | Create a new account |
| GET | `/api/v1/accounts/:id` | Get account by ID |
| PUT | `/api/v1/accounts/:id` | Update account |
| DELETE | `/api/v1/accounts/:id` | Delete account |

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/transactions` | List transactions |
| POST | `/api/v1/transactions` | Create a transaction |
| GET | `/api/v1/transactions/:id` | Get transaction by ID |

## Development

### Running Tests
```bash
# Backend tests
cd ~/repos/FinApp/api
npm test

# Frontend tests
cd ~/repos/FinApp/web
npm test
```

### Building
```bash
# Backend
cd ~/repos/FinApp/api
npm run build

# Frontend
cd ~/repos/FinApp/web
npm run build
```

## Security Considerations

- **Authentication**: JWT-based authentication with refresh tokens
- **Authorization**: Role-based access control (user vs admin)
- **Data Protection**: Encrypted storage for sensitive data
- **Input Validation**: All API inputs validated and sanitized
- **Rate Limiting**: API rate limiting to prevent abuse

## Future Enhancements

- Real-time updates with WebSocket connections
- Export/Import functionality (PDF, CSV)
- Multi-currency support
- Advanced reporting and analytics
- Mobile-responsive design
- Integration with payment gateways (Stripe, PayPal)
