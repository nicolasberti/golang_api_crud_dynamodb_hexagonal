# API Documentation - Products Endpoint

## GET /api/v1/products

Enhanced endpoint for retrieving products with pagination, filtering, and sorting capabilities.

### Query Parameters

| Parameter | Type | Default | Description | Constraints |
|-----------|------|---------|-------------|-------------|
| `page` | integer | 1 | Page number for pagination | `min: 1`, `max: 1000` |
| `limit` | integer | 20 | Number of items per page | `min: 1`, `max: 100` |
| `name` | string | - | Filter products by name (partial match) | - |
| `min_price` | float | - | Minimum price filter | `min: 0` |
| `max_price` | float | - | Maximum price filter | `min: 0` |
| `sort_by` | string | `created_at` | Field to sort by | `name`, `price`, `created_at`, `updated_at` |
| `sort_order` | string | `desc` | Sort order | `asc`, `desc` |
| `fields` | string | - | Comma-separated list of fields to return | - |

### Response Structure

```json
{
  "products": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "price": "number",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ],
  "pagination": {
    "current_page": "integer",
    "per_page": "integer",
    "total_pages": "integer",
    "total_items": "integer",
    "has_next": "boolean",
    "has_prev": "boolean"
  },
  "filters_applied": {
    "name": "string",
    "min_price": "number",
    "max_price": "number"
  }
}
```

### Examples

#### 1. Basic Request (Default Parameters)
```bash
curl -X GET "http://localhost:8080/api/v1/products"
```

**Response:**
```json
{
  "products": [
    {
      "id": "prod-123",
      "name": "Laptop Pro",
      "description": "High-performance laptop",
      "price": 1299.99,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 1,
    "total_items": 1,
    "has_next": false,
    "has_prev": false
  }
}
```

#### 2. Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=2&limit=10"
```

#### 3. Filter by Name
```bash
curl -X GET "http://localhost:8080/api/v1/products?name=Laptop"
```

#### 4. Filter by Price Range
```bash
curl -X GET "http://localhost:8080/api/v1/products?min_price=500&max_price=1500"
```

#### 5. Sort by Price (Ascending)
```bash
curl -X GET "http://localhost:8080/api/v1/products?sort_by=price&sort_order=asc"
```

#### 6. Combined Filters and Sorting
```bash
curl -X GET "http://localhost:8080/api/v1/products?name=Pro&min_price=1000&sort_by=price&sort_order=desc&page=1&limit=5"
```

**Response with Filters:**
```json
{
  "products": [...],
  "pagination": {...},
  "filters_applied": {
    "name": "Pro",
    "min_price": 1000.0,
    "max_price": 0.0
  }
}
```

### Error Responses

#### 400 Bad Request - Invalid Parameters
```json
{
  "error": "invalid query parameters",
  "details": "Key: 'ListProductsRequest.SortBy' Error:Field validation for 'SortBy' failed on the 'oneof' tag"
}
```

#### 400 Bad Request - Page Limit Exceeded
```json
{
  "error": "page cannot exceed 1000"
}
```

#### 400 Bad Request - Invalid Price Range
```json
{
  "error": "min_price cannot be greater than max_price"
}
```

#### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

### Performance Considerations

1. **Pagination**: Always use pagination for large datasets to avoid memory issues
2. **Filtering**: Filters are applied at the database level for better performance
3. **Sorting**: Sorting is performed in-memory for DynamoDB Scan operations
4. **Limits**: Maximum page size is limited to 100 items to prevent large responses

### Best Practices

1. **Use appropriate page sizes**: 20-50 items per page is recommended
2. **Apply filters early**: Use specific filters to reduce dataset size
3. **Cache results**: Consider caching frequently accessed pages
4. **Monitor performance**: Track response times for large datasets

### Rate Limiting

- **Default**: 100 requests per minute per IP
- **Burst**: Up to 20 requests in a single burst
- **Headers**: Rate limit info is included in response headers:
  - `X-RateLimit-Limit`: Maximum requests per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when rate limit resets

### Security Notes

- All inputs are validated and sanitized
- SQL injection protection through parameterized queries
- Rate limiting prevents abuse
- Request size limits are enforced
- Sensitive data is never logged

### Monitoring and Metrics

The endpoint provides the following metrics:
- Request count and response time
- Filter usage statistics
- Pagination patterns
- Error rates by type

### SDK Examples

#### Go
```go
import "net/http"

func getProducts() (*http.Response, error) {
    url := "http://localhost:8080/api/v1/products?page=1&limit=20&sort_by=created_at"
    return http.Get(url)
}
```

#### JavaScript
```javascript
async function getProducts(page = 1, limit = 20, filters = {}) {
    const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
        ...filters
    });
    
    const response = await fetch(`http://localhost:8080/api/v1/products?${params}`);
    return response.json();
}

// Usage
getProducts(1, 20, { name: 'Laptop', min_price: 1000 });
```

#### Python
```python
import requests

def get_products(page=1, limit=20, **filters):
    params = {
        'page': page,
        'limit': limit,
        **filters
    }
    
    response = requests.get('http://localhost:8080/api/v1/products', params=params)
    response.raise_for_status()
    return response.json()

# Usage
products = get_products(page=1, limit=20, name='Laptop', min_price=1000)
```