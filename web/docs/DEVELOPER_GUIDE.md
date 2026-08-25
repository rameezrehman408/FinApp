# FinApp Frontend - Developer Documentation

## Design Patterns

### Component Architecture

#### Presentational vs Container Components
```tsx
// Presentational Component (ui)
<Button onClick={onClick} variant="primary">
  {children}
</Button>

// Container Component
<AuthenticatedButton onClick={handleAuthAction}>
  Delete Account
</AuthenticatedButton>
```

#### Compound Components Pattern
```tsx
<AccountList>
  <AccountList.Header>Add New Account</AccountList.Header>
  <AccountList.Body>
    <AccountCard />
  </AccountList.Body>
</AccountList>
```

### State Management (Zustand/React Context)

#### Feature-Sliced State
- **Shared**: Global state (auth, theme)
- **Entities**: Domain-specific state (accounts, transactions)
- **Processes**: Business logic (budget calculations)

```tsx
// store/features/budgets/index.ts
interface BudgetState {
  utilization: number;
  checkOverrun: (spent: number, limit: number) => boolean;
}

const useBudgetStore = create<BudgetState>()((set) => ({
  utilization: 0,
  checkOverrun: (spent, limit) => {
    return (spent / limit) * 100 > 100;
  }
}));
```

### Data Fetching (SWR Pattern)

#### API Service Layer
```typescript
// services/accounts.ts
export const accountService = {
  getAll: () => api.get<Account[]>('/accounts'),
  getById: (id: string) => api.get<Account>(`/accounts/${id}`),
  create: (data: CreateAccountDto) => api.post<Account>('/accounts', data)
};
```

#### Custom Hook with SWR
```typescript
// hooks/useAccounts.ts
export const useAccounts = () => {
  const { data, error, isLoading, mutate } = useSWR('/accounts', accountService.getAll);
  
  return {
    accounts: data || [],
    isLoading,
    error,
    refetch: mutate
  };
};
```

### Form Handling (React Hook Form)

#### Validation Schema
```typescript
// validators/accountSchema.ts
import { z } from 'zod';

export const accountSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  type: z.enum(['ASSET', 'LIABILITY', 'INCOME', 'EXPENSE', 'EQUITY']),
  currency: z.string().length(3, 'Currency must be 3 characters'),
  balance: z.number().optional()
});
```

#### Form Component
```typescript
const AccountForm = () => {
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: zodResolver(accountSchema)
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input {...register('name')} error={errors.name?.message} />
      <Select {...register('type')} />
      <SubmitButton>Save</SubmitButton>
    </form>
  );
};
```

### Authentication

#### Protected Routes
```tsx
// hooks/useAuth.ts
export const useAuth = () => {
  const { user, login, logout, isLoading } = useContext(AuthContext);
  
  const requireAuth = (callback: () => void) => {
    if (isLoading) return;
    if (!user) router.push('/login');
    else callback();
  };

  return { user, login, logout, requireAuth };
};
```

#### Protected Route Component
```tsx
// components/auth/ProtectedRoute.tsx
export const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  const { user, isLoading } = useAuth();
  
  if (isLoading) return <LoadingSpinner />;
  if (!user) return <LoginForm />;
  
  return <>{children}</>;
};
```

### Error Handling

#### API Error Interceptor
```typescript
// lib/axios.ts
axiosInstance.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      router.push('/login');
    }
    return Promise.reject(error);
  }
);
```

#### Error Boundary
```tsx
// components/ErrorBoundary.tsx
export class ErrorBoundary extends Component<Props, State> {
  componentDidCatch(error: Error) {
    Sentry.captureException(error);
    this.setState({ hasError: true });
  }

  render() {
    if (this.state.hasError) {
      return <ErrorPage />;
    }
    return this.props.children;
  }
}
```

## Folder Structure

```
src/
├── app/                 # Next.js App Router
│   ├── layout.tsx       # Root layout
│   ├── (auth)/          # Auth group (login, register)
│   ├── (dashboard)/     # Protected routes
│   └── api/             # API proxy (server actions)
├── components/          # Reusable UI components
│   ├── layout/          # Page layouts
│   ├── ui/              # Base UI elements
│   ├── auth/            # Authentication components
│   └── shared/          # Shared components
├── hooks/               # Custom React hooks
│   ├── useApi/          # API hooks
│   ├── useAuth/         # Auth hooks
│   └── useBudget/       # Budget hooks
├── store/               # State management
│   ├── features/        # Feature-specific stores
│   └── shared/          # Global stores
├── services/            # API services
│   ├── accounts.ts
│   ├── transactions.ts
│   └── budgets.ts
├── validators/          # Validation schemas
│   └── schemas.ts
├── types/               # TypeScript type definitions
└── lib/                 # Utility functions
```

## TypeScript Types

### API Response Types
```typescript
// types/entities.ts
interface Account {
  id: string;
  name: string;
  type: 'ASSET' | 'LIABILITY' | 'INCOME' | 'EXPENSE' | 'EQUITY';
  currency: string;
  balance: number;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

interface Transaction {
  id: string;
  accountId: string;
  type: 'INCOME' | 'EXPENSE' | 'TRANSFER' | 'INVESTMENT';
  amount: number;
  currency: string;
  description: string;
  transactionDate: string;
  status: 'PENDING' | 'COMPLETED' | 'CANCELLED';
  categoryId?: string;
  partyId?: string;
}

interface Budget {
  id: string;
  categoryId: string;
  householdId: string;
  period: 'WEEKLY' | 'MONTHLY' | 'YEARLY';
  limit: number;
  spent: number;
  utilization: number;
}
```

## Testing Strategy

### Unit Testing (Jest + React Testing Library)
```typescript
// services/__tests__/accounts.test.ts
describe('Account Service', () => {
  it('should create an account', async () => {
    const mockAccount = { name: 'Test Account', type: 'ASSET' };
    mockAxios.post.mockResolvedValue({ data: { id: '1', ...mockAccount } });
    
    const result = await accountService.create(mockAccount);
    expect(result.data.name).toBe('Test Account');
  });
});
```

### Component Testing
```typescript
// components/__tests__/AccountForm.test.tsx
describe('AccountForm', () => {
  it('renders all form fields', () => {
    render(<AccountForm />);
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/type/i)).toBeInTheDocument();
  });

  it('shows validation errors', async () => {
    render(<AccountForm />);
    fireEvent.click(screen.getByText(/submit/i));
    expect(await screen.findAllByRole('alert')).toHaveLength(2);
  });
});
```

## Performance Optimization

### React.memo for Expensive Components
```tsx
const AccountCard = React.memo(({ account }: { account: Account }) => {
  const formattedBalance = useFormatter.format(account.balance);
  return (
    <Card>
      <h3>{account.name}</h3>
      <p>{formattedBalance}</p>
    </Card>
  );
});
```

### Virtualized Lists
```tsx
<List
  height={600}
  itemCount={accounts.length}
  itemSize={100}
 .Item={AccountListItem}
/>
```

## CI/CD Pipeline

### GitHub Actions Workflow
```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - run: npm ci
      - run: npm test
      - run: npm run build
```