# DTO Mappers

This package provides conversion functions between domain entities and DTOs (Data Transfer Objects).

## Structure

- `user_mapper.go` - User entity ↔ DTO conversions
- `company_mapper.go` - Company entity ↔ DTO conversions
- `job_mapper.go` - Job entity ↔ DTO conversions
- `application_mapper.go` - Application entity ↔ DTO conversions
- `helpers.go` - Helper functions for pointer conversions

## Usage

### Entity to Response DTO

```go
import "keerja-backend/internal/dto/mapper"

// Convert user entity to response DTO
userResp := mapper.ToUserResponse(userEntity)

// Convert job entity with relations to detail response
jobDetail := mapper.ToJobDetailResponse(jobEntity)
```

### Request DTO to Entity

```go
// Update entity from request DTO
mapper.UpdateProfileRequestToEntity(req, profileEntity)
```

## Notes

- These mappers provide basic field mapping
- Handlers should add computed fields (counts, relationships, etc.)
- Some mappings require database queries for related data
- Validation should be done before mapping

## TODO

- Add auth mapper for token generation
- Add comprehensive tests
- Add batch conversion functions
- Optimize for performance with goroutines
