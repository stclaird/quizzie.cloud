# Integration Testing Guide

This document describes the different approaches to integration testing in the Quizzie API project.

## Test Types

### 1. In-Memory Integration Tests (`api/integration_test.go`)

**Location**: `quizApp/api/integration_test.go`

**Purpose**: Fast integration tests that verify the full API workflow using an in-memory SQLite database.

**Run**:
```bash
# From quizApp directory
go test -tags=integration ./api/... -v

# Or using Make
make integration-test
```

**What it tests**:
- Full quiz workflow (get categories → get questions → answer questions)
- All API endpoints with real HTTP requests
- CORS headers and middleware
- Database interactions
- Invalid/error cases

**Benefits**:
- Fast execution (runs in <1 second)
- No external dependencies
- Runs in CI/CD easily
- Tests actual business logic integration

### 2. Docker-based Integration Tests (`tests/api_integration_test.go`)

**Location**: `quizApp/tests/api_integration_test.go`

**Purpose**: Test against a real running API instance (local Docker or deployed production).

**Run locally against Docker**:
```bash
# Start local Docker container
docker-compose -f docker-compose.test.yml up -d

# Run tests
export API_URL=http://localhost:8080
go test -tags=integration ./tests/... -v

# Cleanup
docker-compose -f docker-compose.test.yml down

# Or use automated script
./scripts/integration-test.sh

# Or using Make
make integration-test-docker
```

**Run against production**:
```bash
export API_URL=https://quizzierhg7bjpx-quizapp.functions.fnc.fr-par.scw.cloud
go test -tags=integration ./tests/... -v
```

**What it tests**:
- Actual HTTP calls to running API
- Full application startup and initialization
- Question loading from external sources (GitHub releases)
- Response times and performance
- CORS configuration
- Real-world deployment scenario

**Benefits**:
- Tests exactly what gets deployed
- Validates Docker configuration
- Tests with production data
- Performance benchmarking
- Post-deployment smoke tests

### 3. Unit Tests

**Location**: Various `*_test.go` files

**Run**:
```bash
# All unit tests
go test ./...

# With coverage
make test-coverage

# Verbose output
make test-verbose
```

## CI/CD Integration

The GitHub Actions workflow (`.github/workflows/BackEnd-Deploy.yaml`) runs tests in this order:

```yaml
# 1. Unit Tests (fast)
- name: Run Go Tests
  run: go test ./...

# 2. In-Memory Integration Tests (fast)
- name: Run Integration Tests
  run: go test -tags=integration ./api/... -v

# 3. Build & Security Scan
- name: Build Docker image for scanning
- name: Run Trivy vulnerability scanner

# 4. Deploy to Scaleway
- name: Deploy to Scaleway

# 5. Smoke Tests (production verification)
- name: Wait for deployment to stabilize
  run: sleep 30

- name: Run smoke tests against production
  env:
    API_URL: https://quizzierhg7bjpx-quizapp.functions.fnc.fr-par.scw.cloud
  run: go test -tags=integration ./tests/... -v
```

**Test Strategy**:
1. **Unit tests** - Fast feedback on individual functions
2. **In-memory integration** - Verify component integration (~1 second)
3. **Build & scan** - Only build if tests pass, check security
4. **Deploy** - Push to production
5. **Smoke tests** - Verify deployment works (continue on error)

## Test Structure

### In-Memory Integration Tests

```go
// +build integration

func setupIntegrationTest(t *testing.T) *gin.Engine {
    // Creates router with in-memory DB
    // Seeds test data
    // Returns configured router
}

func TestIntegration_FullQuizWorkflow(t *testing.T) {
    // Tests complete user journey
}
```

### Docker Integration Tests

```go
// +build integration

var apiURL string

func init() {
    apiURL = os.Getenv("API_URL")
    if apiURL == "" {
        apiURL = "http://localhost:8080"
    }
}

func TestAPI_GetCategories(t *testing.T) {
    waitForAPI(t)
    resp, err := http.Get(apiURL + "/categories/")
    // Makes real HTTP call to configured API
}
```

**Environment variable**: Set `API_URL` to test against different environments:
- Local Docker: `http://localhost:8080`
- Production: `https://quizzierhg7bjpx-quizapp.functions.fnc.fr-par.scw.cloud`

## Adding New Tests

### Adding In-Memory Integration Test

1. Add test function to `api/integration_test.go`
2. Use `// +build integration` tag
3. Call `setupIntegrationTest(t)` to get router
4. Make HTTP requests using `httptest`

Example:
```go
func TestIntegration_NewFeature(t *testing.T) {
    router := setupIntegrationTest(t)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/your/endpoint", nil)
    router.ServeHTTP(w, req)

    // Assertions...
}
```

### Adding Docker Integration Test

1. Add test to `tests/api_integration_test.go`
2. Use `// +build integration` tag
3. Use real HTTP calls to `apiURL`
4. Wait for API readiness with `waitForAPI(t)`

Example:
```go
func TestAPI_NewEndpoint(t *testing.T) {
    waitForAPI(t)

    resp, err := http.Get(apiURL + "/your/endpoint")
    // Handle response...
}
```

## Best Practices

1. **Use build tags**: Always use `// +build integration` to separate from unit tests
2. **Test workflows, not units**: Integration tests should test complete user journeys
3. **Seed realistic data**: Use data that represents actual use cases
4. **Test error cases**: Don't just test happy paths
5. **Keep tests independent**: Each test should be able to run standalone
6. **Clean up resources**: Close connections, remove test files, etc.

## Troubleshooting

### Tests timing out
- Check if the API is starting correctly
- Increase timeout in `waitForAPI()` function (default: 30 seconds)
- Check Docker container logs: `docker-compose -f docker-compose.test.yml logs`
- For production tests, verify the URL is accessible

### Database errors
- Ensure migrations are running
- Check that test data seeding is working
- Verify in-memory database is being used (for `api/integration_test.go`)

### Port conflicts
- Ensure port 8080 is not in use
- Check `docker-compose.test.yml` for port mappings
- Stop other running instances: `docker-compose down`

### Production smoke tests failing
- Check deployment succeeded: `scw container container get <container-id>`
- Verify API URL is correct in workflow
- Wait time may need adjustment (currently 30 seconds)
- Check if question pack is loading: visit `/categories/` endpoint

## Coverage

View test coverage including integration tests:

```bash
# Generate coverage report
make test-coverage

# View in browser (opens coverage.html)
open coverage.html
```

Integration tests improve overall coverage by testing:
- Middleware integration
- Route registration
- Database queries in context
- Full request/response cycle
- Error handling paths

## Quick Reference

```bash
# Run all tests (unit + integration)
make test

# Run only unit tests
go test ./...

# Run only in-memory integration tests
make integration-test
go test -tags=integration ./api/... -v

# Run Docker integration tests (local)
make integration-test-docker

# Run smoke tests against production
export API_URL=https://quizzierhg7bjpx-quizapp.functions.fnc.fr-par.scw.cloud
go test -tags=integration ./tests/... -v

# Generate coverage report
make test-coverage
```

## Test Results

**Latest production smoke test** (3 Dec 2025):
- ✅ 6 categories found
- ✅ 614 questions loaded
- ✅ Category filtering works
- ✅ Answer submission works
- ✅ CORS configured correctly
- ✅ Response times: ~100-200ms
- ✅ All tests passing
