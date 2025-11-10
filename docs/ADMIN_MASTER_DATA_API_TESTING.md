# Admin Master Data API - Testing Guide

## Quick Start Testing

### Prerequisites
1. Server berjalan di `http://localhost:8080`
2. Admin token yang valid
3. Postman atau cURL installed

---

## 1. Setup Authentication

### Get Admin Token

**Login sebagai Admin:**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin_password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "..."
  }
}
```

**Simpan token untuk digunakan di semua request berikutnya.**

---

## 2. Test Province CRUD

### 2.1 Create Province

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jawa Barat",
    "code": "32",
    "is_active": true
  }'
```

**Expected Response (201):**
```json
{
  "success": true,
  "message": "Province created successfully",
  "data": {
    "id": 1,
    "name": "Jawa Barat",
    "code": "32",
    "is_active": true
  }
}
```

**Test Duplicate Code (Should Fail):**
```bash
curl -X POST http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jawa Timur",
    "code": "32",
    "is_active": true
  }'
```

**Expected Response (409):**
```json
{
  "success": false,
  "message": "Province with this code already exists"
}
```

### 2.2 Get All Provinces

```bash
curl -X GET "http://localhost:8080/api/v1/admin/master/provinces?active=true" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 2.3 Get Province by ID

```bash
curl -X GET http://localhost:8080/api/v1/admin/master/provinces/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 2.4 Update Province

```bash
curl -X PUT http://localhost:8080/api/v1/admin/master/provinces/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jawa Barat Updated",
    "is_active": false
  }'
```

### 2.5 Delete Province (Test Relational Constraint)

**Step 1: Create City dengan Province ID 1**
```bash
curl -X POST http://localhost:8080/api/v1/admin/master/cities \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bandung",
    "type": "Kota",
    "code": "3273",
    "province_id": 1,
    "is_active": true
  }'
```

**Step 2: Try Delete Province (Should Fail)**
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/master/provinces/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Response (409):**
```json
{
  "success": false,
  "message": "Cannot delete province: it is still referenced by cities or companies"
}
```

---

## 3. Test City CRUD

### 3.1 Create City

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/cities \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bandung",
    "type": "Kota",
    "code": "3273",
    "province_id": 1,
    "is_active": true
  }'
```

### 3.2 Get All Cities

```bash
curl -X GET "http://localhost:8080/api/v1/admin/master/cities?province_id=1&active=true" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 3.3 Update City

```bash
curl -X PUT http://localhost:8080/api/v1/admin/master/cities/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bandung Updated",
    "type": "Kota"
  }'
```

### 3.4 Delete City

**Test dengan districts:**
1. Create district dengan city_id
2. Try delete city → Should fail (409)
3. Delete district
4. Delete city → Should success (200)

---

## 4. Test District CRUD

### 4.1 Create District

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/districts \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Batujajar",
    "code": "3273010",
    "postal_code": "40561",
    "city_id": 1,
    "is_active": true
  }'
```

### 4.2 Get All Districts

```bash
curl -X GET "http://localhost:8080/api/v1/admin/master/districts?city_id=1" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 4.3 Get District by ID (Full Hierarchy)

```bash
curl -X GET http://localhost:8080/api/v1/admin/master/districts/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Batujajar",
    "code": "3273010",
    "postal_code": "40561",
    "city_id": 1,
    "city": {
      "id": 1,
      "name": "Bandung",
      "full_name": "Kota Bandung",
      "province": {
        "id": 1,
        "name": "Jawa Barat"
      }
    },
    "full_location_path": "Batujajar, Kota Bandung, Jawa Barat",
    "is_active": true
  }
}
```

---

## 5. Test Industry CRUD

### 5.1 Create Industry

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/industries \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technology",
    "slug": "technology",
    "description": "Technology and IT industry",
    "icon_url": "https://example.com/icon.png",
    "is_active": true
  }'
```

### 5.2 Create Industry dengan Auto-Generated Slug

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/industries \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Financial Services",
    "description": "Banking and finance industry",
    "is_active": true
  }'
```

**Slug akan auto-generated:** `financial-services`

### 5.3 Test Duplicate Name

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/industries \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technology",
    "slug": "technology-2",
    "is_active": true
  }'
```

**Expected Response (409):**
```json
{
  "success": false,
  "message": "Industry with this name already exists"
}
```

---

## 6. Test Job Type CRUD

### 6.1 Create Job Type

```bash
curl -X POST http://localhost:8080/api/v1/admin/master/job-types \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Full-Time",
    "code": "full_time",
    "order": 1
  }'
```

### 6.2 Get All Job Types

```bash
curl -X GET http://localhost:8080/api/v1/admin/master/job-types \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 6.3 Update Job Type

```bash
curl -X PUT http://localhost:8080/api/v1/admin/master/job-types/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Full-Time Updated",
    "order": 1
  }'
```

### 6.4 Test Delete dengan Job References

**Precondition:** Pastikan ada job yang menggunakan job_type_id = 1

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/master/job-types/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Response (409):**
```json
{
  "success": false,
  "message": "Cannot delete job type: it is still referenced by jobs"
}
```

---

## 7. Test Company Size CRUD

### 7.1 Create Company Size

```bash
curl -X POST http://localhost:8080/api/v1/admin/meta/company-sizes \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "1-10 employees",
    "min_employees": 1,
    "max_employees": 10,
    "is_active": true
  }'
```

### 7.2 Create Company Size dengan Unlimited Max

```bash
curl -X POST http://localhost:8080/api/v1/admin/meta/company-sizes \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "1000+ employees",
    "min_employees": 1000,
    "max_employees": null,
    "is_active": true
  }'
```

### 7.3 Get All Company Sizes

```bash
curl -X GET "http://localhost:8080/api/v1/admin/meta/company-sizes?active=true" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## 8. Complete Test Flow

### Test Flow: Location Hierarchy

```bash
# 1. Create Province
PROVINCE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Jawa Barat", "code": "32", "is_active": true}')

PROVINCE_ID=$(echo $PROVINCE_RESPONSE | jq -r '.data.id')
echo "Created Province ID: $PROVINCE_ID"

# 2. Create City
CITY_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/admin/master/cities \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Bandung\", \"type\": \"Kota\", \"code\": \"3273\", \"province_id\": $PROVINCE_ID, \"is_active\": true}")

CITY_ID=$(echo $CITY_RESPONSE | jq -r '.data.id')
echo "Created City ID: $CITY_ID"

# 3. Create District
DISTRICT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/admin/master/districts \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Batujajar\", \"code\": \"3273010\", \"postal_code\": \"40561\", \"city_id\": $CITY_ID, \"is_active\": true}")

DISTRICT_ID=$(echo $DISTRICT_RESPONSE | jq -r '.data.id')
echo "Created District ID: $DISTRICT_ID"

# 4. Try Delete Province (Should Fail)
curl -X DELETE http://localhost:8080/api/v1/admin/master/provinces/$PROVINCE_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
# Expected: 409 Conflict

# 5. Delete District
curl -X DELETE http://localhost:8080/api/v1/admin/master/districts/$DISTRICT_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 6. Delete City
curl -X DELETE http://localhost:8080/api/v1/admin/master/cities/$CITY_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 7. Delete Province (Should Success)
curl -X DELETE http://localhost:8080/api/v1/admin/master/provinces/$PROVINCE_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
# Expected: 200 OK
```

---

## 9. Postman Collection

### Import Collection

1. Buka Postman
2. Klik **Import**
3. Pilih file collection JSON (dibuat manual atau dari export)

### Environment Variables

Buat environment dengan variables:
```
base_url: http://localhost:8080
admin_token: YOUR_ADMIN_TOKEN
```

### Collection Structure

```
Admin Master Data CRUD
├── Province
│   ├── Create Province
│   ├── Get All Provinces
│   ├── Get Province by ID
│   ├── Update Province
│   └── Delete Province
├── City
│   ├── Create City
│   ├── Get All Cities
│   ├── Get City by ID
│   ├── Update City
│   └── Delete City
├── District
│   ├── Create District
│   ├── Get All Districts
│   ├── Get District by ID
│   ├── Update District
│   └── Delete District
├── Industry
│   ├── Create Industry
│   ├── Get All Industries
│   ├── Get Industry by ID
│   ├── Update Industry
│   └── Delete Industry
├── Job Type
│   ├── Create Job Type
│   ├── Get All Job Types
│   ├── Get Job Type by ID
│   ├── Update Job Type
│   └── Delete Job Type
└── Company Size
    ├── Create Company Size
    ├── Get All Company Sizes
    ├── Get Company Size by ID
    ├── Update Company Size
    └── Delete Company Size
```

---

## 10. Test Scenarios Checklist

### Province Tests
- [ ] ✅ Create province dengan data valid
- [ ] ✅ Create province dengan code duplicate → 409
- [ ] ✅ Get all provinces → 200
- [ ] ✅ Get province by ID → 200
- [ ] ✅ Update province → 200
- [ ] ✅ Delete province tanpa references → 200
- [ ] ✅ Delete province dengan cities → 409
- [ ] ✅ Delete province digunakan companies → 409

### City Tests
- [ ] ✅ Create city dengan data valid
- [ ] ✅ Create city dengan name duplicate dalam province → 409
- [ ] ✅ Get all cities by province_id → 200
- [ ] ✅ Get city by ID → 200
- [ ] ✅ Update city → 200
- [ ] ✅ Delete city tanpa references → 200
- [ ] ✅ Delete city dengan districts → 409
- [ ] ✅ Delete city digunakan companies → 409

### District Tests
- [ ] ✅ Create district dengan data valid
- [ ] ✅ Create district dengan name duplicate dalam city → 409
- [ ] ✅ Get all districts by city_id → 200
- [ ] ✅ Get district by ID dengan full hierarchy → 200
- [ ] ✅ Update district → 200
- [ ] ✅ Delete district tanpa references → 200
- [ ] ✅ Delete district digunakan companies → 409

### Industry Tests
- [ ] ✅ Create industry dengan data valid
- [ ] ✅ Create industry dengan name duplicate → 409
- [ ] ✅ Create industry dengan slug auto-generated
- [ ] ✅ Get all industries → 200
- [ ] ✅ Get industry by ID → 200
- [ ] ✅ Update industry → 200
- [ ] ✅ Delete industry tanpa references → 200
- [ ] ✅ Delete industry digunakan companies → 409

### Job Type Tests
- [ ] ✅ Create job type dengan data valid
- [ ] ✅ Create job type dengan code duplicate → 409
- [ ] ✅ Get all job types → 200
- [ ] ✅ Get job type by ID → 200
- [ ] ✅ Update job type → 200
- [ ] ✅ Delete job type tanpa references → 200
- [ ] ✅ Delete job type digunakan jobs → 409

### Company Size Tests
- [ ] ✅ Create company size dengan data valid
- [ ] ✅ Create company size dengan label duplicate → 409
- [ ] ✅ Create company size dengan max_employees unlimited (null)
- [ ] ✅ Get all company sizes → 200
- [ ] ✅ Get company size by ID → 200
- [ ] ✅ Update company size → 200
- [ ] ✅ Delete company size tanpa references → 200
- [ ] ✅ Delete company size digunakan companies → 409

### Authentication Tests
- [ ] ✅ Request tanpa token → 401
- [ ] ✅ Request dengan token non-admin → 403
- [ ] ✅ Request dengan admin token → 200/201

---

## 11. Error Testing

### Test Validation Errors

**Missing Required Field:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "32"
  }'
```

**Expected Response (400):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "name": "name is required"
  }
}
```

**Invalid Data Type:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "A",
    "code": "32",
    "is_active": "not_boolean"
  }'
```

**Expected Response (400):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "name": "name must be at least 2 characters"
  }
}
```

---

## 12. Performance Testing

### Test dengan Banyak Data

```bash
# Create 100 provinces
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/admin/master/provinces \
    -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"Province $i\",
      \"code\": \"P$i\",
      \"is_active\": true
    }"
done

# Get all provinces and measure response time
time curl -X GET http://localhost:8080/api/v1/admin/master/provinces \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## 13. Tips Testing

1. **Gunakan Environment Variables**: Simpan token dan base_url di environment variables
2. **Test Relational Constraints**: Pastikan test delete dengan references untuk memastikan validasi bekerja
3. **Test Duplicate Validation**: Test semua skenario duplikasi untuk memastikan validasi bekerja
4. **Test Error Handling**: Test berbagai error scenarios untuk memastikan error handling bekerja
5. **Clean Up**: Hapus test data setelah testing selesai

---

## 14. Common Issues & Solutions

### Issue: "Unauthorized" Error
**Solution:** 
- Pastikan token valid dan belum expired
- Pastikan header Authorization format benar: `Bearer <token>`

### Issue: "Forbidden" Error
**Solution:**
- Pastikan user memiliki role admin
- Check user permissions di database

### Issue: "Cannot delete: still referenced"
**Solution:**
- Hapus semua referensi terlebih dahulu
- Atau update referensi ke master data lain

### Issue: "Duplicate entry" Error
**Solution:**
- Gunakan name/code/label yang berbeda
- Atau update entitas yang sudah ada

---

## 15. API Response Examples

### Success Response (200/201)
```json
{
  "success": true,
  "message": "Province created successfully",
  "data": {
    "id": 1,
    "name": "Jawa Barat",
    "code": "32",
    "is_active": true
  }
}
```

### Error Response (400)
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "name": "name is required",
    "code": "code must be at least 2 characters"
  }
}
```

### Error Response (404)
```json
{
  "success": false,
  "message": "Province not found"
}
```

### Error Response (409)
```json
{
  "success": false,
  "message": "Province with this code already exists"
}
```

---

## Summary

Semua endpoint CRUD untuk Master Data sudah diimplementasi dengan:
- ✅ Validasi input lengkap
- ✅ Validasi duplikasi
- ✅ Validasi relasional pada DELETE
- ✅ Error handling yang jelas
- ✅ Authentication & Authorization
- ✅ Cache invalidation
- ✅ Response format konsisten

Endpoint siap digunakan dan di-test!

