# Product CRUD Hexagonal

Una aplicación CRUD de productos en Go siguiendo arquitectura hexagonal y buenas prácticas.

## Arquitectura

- **Core Domain**: Lógica de negocio pura
- **Ports**: Interfaces para desacoplar
- **Adapters**: Implementaciones concretas (HTTP, DynamoDB)
- **Platform**: Configuración y logging

## Requisitos

- Go 1.21+
- AWS CLI configurado
- Terraform 1.0+

## Infraestructura

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

## Ejecución Local

```bash
# Copiar variables de entorno
cp .env.example .env

# Ejecutar aplicación
go run cmd/api/main.go
```

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/products` - Crear producto
- `GET /api/v1/products` - Listar productos
- `GET /api/v1/products/:id` - Obtener producto
- `PUT /api/v1/products/:id` - Actualizar producto
- `DELETE /api/v1/products/:id` - Eliminar producto

## Ejemplo de Uso

```bash
# Crear producto
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":1299.99}'

# Listar productos
curl http://localhost:8080/api/v1/products
```