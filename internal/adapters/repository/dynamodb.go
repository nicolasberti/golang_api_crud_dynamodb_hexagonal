package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/ports"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBRepository) Save(ctx context.Context, product domain.Product) error {
	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoDBRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return domain.Product{}, err
	}
	if result.Item == nil {
		return domain.Product{}, domain.ErrNotFound
	}

	var product domain.Product
	err = attributevalue.UnmarshalMap(result.Item, &product)
	return product, err
}

func (r *DynamoDBRepository) Update(ctx context.Context, product domain.Product) error {
	return r.Save(ctx, product) // PutItem overwrites
}

func (r *DynamoDBRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func (r *DynamoDBRepository) List(ctx context.Context) ([]domain.Product, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	err = attributevalue.UnmarshalListOfMaps(result.Items, &products)
	return products, err
}

func (r *DynamoDBRepository) ListWithFilters(ctx context.Context, filters ports.ProductFilters) (*ports.ProductListResult, error) {
	// Build scan input with filters
	scanInput := &dynamodb.ScanInput{
		TableName:         aws.String(r.tableName),
		Limit:             aws.Int32(int32(filters.Limit)),
		ExclusiveStartKey: nil, // Will be set for pagination
	}

	// Build filter expression if filters are applied
	var filterExpression strings.Builder
	var expressionAttributeNames map[string]string
	var expressionAttributeValues map[string]types.AttributeValue

	if filters.Name != "" || filters.MinPrice > 0 || filters.MaxPrice > 0 {
		expressionAttributeNames = make(map[string]string)
		expressionAttributeValues = make(map[string]types.AttributeValue)
		var conditions []string

		// Name filter (contains)
		if filters.Name != "" {
			conditions = append(conditions, "contains(#name, :name)")
			expressionAttributeNames["#name"] = "name"
			expressionAttributeValues[":name"] = &types.AttributeValueMemberS{Value: filters.Name}
		}

		// Price filters
		if filters.MinPrice > 0 {
			conditions = append(conditions, "price >= :min_price")
			expressionAttributeValues[":min_price"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", filters.MinPrice)}
		}

		if filters.MaxPrice > 0 {
			conditions = append(conditions, "price <= :max_price")
			expressionAttributeValues[":max_price"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", filters.MaxPrice)}
		}

		// Combine conditions
		filterExpression.WriteString(strings.Join(conditions, " AND "))
		scanInput.FilterExpression = aws.String(filterExpression.String())
		scanInput.ExpressionAttributeNames = expressionAttributeNames
		scanInput.ExpressionAttributeValues = expressionAttributeValues
	}

	// Execute scan
	result, err := r.client.Scan(ctx, scanInput)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products: %w", err)
	}

	// Unmarshal products
	var products []domain.Product
	err = attributevalue.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}

	// Get total count for pagination
	totalItems, err := r.getTotalCount(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Sort products in memory (DynamoDB Scan doesn't guarantee order)
	products = r.sortProducts(products, filters.SortBy, filters.SortOrder)

	// Apply offset for pagination
	if filters.Offset < len(products) {
		products = products[filters.Offset:]
	} else {
		products = []domain.Product{}
	}

	// Limit results
	if filters.Limit < len(products) {
		products = products[:filters.Limit]
	}

	return &ports.ProductListResult{
		Products:   products,
		TotalItems: totalItems,
	}, nil
}

func (r *DynamoDBRepository) getTotalCount(ctx context.Context, filters ports.ProductFilters) (int, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Select:    types.SelectCount,
	}

	// Apply same filters for count
	if filters.Name != "" || filters.MinPrice > 0 || filters.MaxPrice > 0 {
		var filterExpression strings.Builder
		var expressionAttributeNames map[string]string
		var expressionAttributeValues map[string]types.AttributeValue
		var conditions []string

		if filters.Name != "" {
			conditions = append(conditions, "contains(#name, :name)")
			if expressionAttributeNames == nil {
				expressionAttributeNames = make(map[string]string)
			}
			if expressionAttributeValues == nil {
				expressionAttributeValues = make(map[string]types.AttributeValue)
			}
			expressionAttributeNames["#name"] = "name"
			expressionAttributeValues[":name"] = &types.AttributeValueMemberS{Value: filters.Name}
		}

		if filters.MinPrice > 0 {
			conditions = append(conditions, "price >= :min_price")
			if expressionAttributeValues == nil {
				expressionAttributeValues = make(map[string]types.AttributeValue)
			}
			expressionAttributeValues[":min_price"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", filters.MinPrice)}
		}

		if filters.MaxPrice > 0 {
			conditions = append(conditions, "price <= :max_price")
			if expressionAttributeValues == nil {
				expressionAttributeValues = make(map[string]types.AttributeValue)
			}
			expressionAttributeValues[":max_price"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", filters.MaxPrice)}
		}

		if len(conditions) > 0 {
			filterExpression.WriteString(strings.Join(conditions, " AND "))
			scanInput.FilterExpression = aws.String(filterExpression.String())
			scanInput.ExpressionAttributeNames = expressionAttributeNames
			scanInput.ExpressionAttributeValues = expressionAttributeValues
		}
	}

	result, err := r.client.Scan(ctx, scanInput)
	if err != nil {
		return 0, err
	}

	return int(result.Count), nil
}

func (r *DynamoDBRepository) sortProducts(products []domain.Product, sortBy, sortOrder string) []domain.Product {
	if len(products) <= 1 {
		return products
	}

	// Simple bubble sort for demonstration - in production, consider more efficient sorting
	sorted := make([]domain.Product, len(products))
	copy(sorted, products)

	// Define comparison function based on sort field
	var compare func(i, j int) bool
	switch sortBy {
	case "name":
		compare = func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].Name > sorted[j].Name
			}
			return sorted[i].Name < sorted[j].Name
		}
	case "price":
		compare = func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].Price > sorted[j].Price
			}
			return sorted[i].Price < sorted[j].Price
		}
	case "created_at":
		fallthrough
	default:
		compare = func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
			}
			return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
		}
	}

	// Simple insertion sort
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && compare(j, j+1) {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	return sorted
}
