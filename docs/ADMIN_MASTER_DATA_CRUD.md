# Dokumentasi Admin Master Data CRUD

## Overview

Dokumentasi ini menjelaskan implementasi CRUD (Create, Read, Update, Delete) lengkap untuk Master Data yang dapat dikelola oleh Admin. Semua endpoint memerlukan autentikasi admin.

## Daftar Isi

1. [Endpoint Overview](#endpoint-overview)
2. [Authentication](#authentication)
3. [Province CRUD](#1-province-crud)
4. [City CRUD](#2-city-crud)
5. [District CRUD](#3-district-crud)
6. [Industry CRUD](#4-industry-crud)
7. [Job Type CRUD](#5-job-type-crud)
8. [Company Size CRUD](#6-company-size-crud)
9. [Validasi dan Error Handling](#validasi-dan-error-handling)
10. [Testing API](#testing-api)

---

## Endpoint Overview

### Base URL
```
/api/v1/admin/master
```

### Endpoints Summary

| Entity | POST | GET | GET by ID | PUT | DELETE |
|--------|------|-----|-----------|-----|--------|
| Provinces | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cities | ✅ | ✅ | ✅ | ✅ | ✅ |
| Districts | ✅ | ✅ | ✅ | ✅ | ✅ |
| Industries | ✅ | ✅ | ✅ | ✅ | ✅ |
| Job Types | ✅ | ✅ | ✅ | ✅ | ✅ |
| Company Sizes | ✅ | ✅ | ✅ | ✅ | ✅ |

**Catatan:** Company Sizes menggunakan base path `/api/v1/admin/meta/company-sizes`

---

## Authentication

Semua endpoint memerlukan:
1. **Authentication Token** (Bearer Token)
2. **Admin Role** (user harus memiliki role admin)

### Header yang Diperlukan
```
Authorization: Bearer <token>
Content-Type: application/json
```

---

## 1. Province CRUD

### 1.1 Create Province
**POST** `/api/v1/admin/master/provinces`

**Request Body:**
```json
{
  "name": "Jawa Barat",
  "code": "32",
  "is_active": true
}
```

**Response (201 Created):**
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

**Validasi:**
- `name`: required, min 2, max 255 characters
- `code`: required, min 2, max 50 characters, unique
- `is_active`: boolean (default: true)

**Error Responses:**
- `400`: Invalid request body / Validation failed
- `409`: Province with this code already exists
- `401`: Unauthorized
- `403`: Forbidden (not admin)

---

### 1.2 Get All Provinces
**GET** `/api/v1/admin/master/provinces`

**Query Parameters:**
- `search` (optional): Search by name
- `active` (optional): Filter by active status (true/false)

**Example:**
```
GET /api/v1/admin/master/provinces?search=jawa&active=true
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Provinces retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Jawa Barat",
      "code": "32",
      "is_active": true
    },
    {
      "id": 2,
      "name": "Jawa Timur",
      "code": "35",
      "is_active": true
    }
  ]
}
```

---

### 1.3 Get Province by ID
**GET** `/api/v1/admin/master/provinces/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Province retrieved successfully",
  "data": {
    "id": 1,
    "name": "Jawa Barat",
    "code": "32",
    "is_active": true
  }
}
```

**Error Responses:**
- `400`: Invalid province ID
- `404`: Province not found

---

### 1.4 Update Province
**PUT** `/api/v1/admin/master/provinces/:id`

**Request Body:**
```json
{
  "name": "Jawa Barat Updated",
  "code": "32",
  "is_active": false
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Province updated successfully",
  "data": {
    "id": 1,
    "name": "Jawa Barat Updated",
    "code": "32",
    "is_active": false
  }
}
```

**Validasi:**
- Semua field optional
- Jika `code` diubah, akan dicek duplikasi
- `is_active` dapat diubah untuk mengaktifkan/nonaktifkan

---

### 1.5 Delete Province
**DELETE** `/api/v1/admin/master/provinces/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Province deleted successfully",
  "data": null
}
```

**Validasi Relasional:**
- Province **TIDAK BISA** dihapus jika:
  - Masih memiliki cities (anak kota/kabupaten)
  - Masih digunakan oleh companies

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete province: it is still referenced by cities or companies"
}
```

---

## 2. City CRUD

### 2.1 Create City
**POST** `/api/v1/admin/master/cities`

**Request Body:**
```json
{
  "name": "Bandung",
  "type": "Kota",
  "code": "3273",
  "province_id": 1,
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "City created successfully",
  "data": {
    "id": 1,
    "name": "Bandung",
    "full_name": "Kota Bandung",
    "type": "Kota",
    "code": "3273",
    "province_id": 1,
    "province": {
      "id": 1,
      "name": "Jawa Barat",
      "code": "32",
      "is_active": true
    },
    "is_active": true
  }
}
```

**Validasi:**
- `name`: required, min 2, max 255 characters
- `type`: required, must be "Kota" or "Kabupaten"
- `code`: required, min 2, max 50 characters, unique
- `province_id`: required, min 1, must exist
- `is_active`: boolean (default: true)
- Name harus unique dalam satu province

---

### 2.2 Get All Cities
**GET** `/api/v1/admin/master/cities?province_id=1&search=bandung&active=true`

**Query Parameters:**
- `province_id` (required): Filter by province ID
- `search` (optional): Search by name
- `active` (optional): Filter by active status

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Cities retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Bandung",
      "full_name": "Kota Bandung",
      "type": "Kota",
      "code": "3273",
      "province_id": 1,
      "is_active": true
    }
  ]
}
```

---

### 2.3 Get City by ID
**GET** `/api/v1/admin/master/cities/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "City retrieved successfully",
  "data": {
    "id": 1,
    "name": "Bandung",
    "full_name": "Kota Bandung",
    "type": "Kota",
    "code": "3273",
    "province_id": 1,
    "province": {
      "id": 1,
      "name": "Jawa Barat",
      "code": "32",
      "is_active": true
    },
    "is_active": true
  }
}
```

---

### 2.4 Update City
**PUT** `/api/v1/admin/master/cities/:id`

**Request Body:**
```json
{
  "name": "Bandung Updated",
  "type": "Kota",
  "code": "3273",
  "province_id": 1,
  "is_active": true
}
```

**Validasi:**
- Semua field optional
- Jika `name` atau `province_id` diubah, akan dicek duplikasi name dalam province baru

---

### 2.5 Delete City
**DELETE** `/api/v1/admin/master/cities/:id`

**Validasi Relasional:**
- City **TIDAK BISA** dihapus jika:
  - Masih memiliki districts (kecamatan)
  - Masih digunakan oleh companies

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete city: it is still referenced by districts or companies"
}
```

---

## 3. District CRUD

### 3.1 Create District
**POST** `/api/v1/admin/master/districts`

**Request Body:**
```json
{
  "name": "Batujajar",
  "code": "3273010",
  "postal_code": "40561",
  "city_id": 1,
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "District created successfully",
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
      "type": "Kota",
      "code": "3273",
      "province_id": 1,
      "province": {
        "id": 1,
        "name": "Jawa Barat",
        "code": "32",
        "is_active": true
      },
      "is_active": true
    },
    "full_location_path": "Batujajar, Kota Bandung, Jawa Barat",
    "is_active": true
  }
}
```

**Validasi:**
- `name`: required, min 2, max 255 characters
- `code`: required, min 2, max 50 characters, unique
- `postal_code`: optional, must be 5 characters if provided
- `city_id`: required, min 1, must exist
- `is_active`: boolean (default: true)
- Name harus unique dalam satu city

---

### 3.2 Get All Districts
**GET** `/api/v1/admin/master/districts?city_id=1&search=batu&active=true`

**Query Parameters:**
- `city_id` (required): Filter by city ID
- `search` (optional): Search by name
- `active` (optional): Filter by active status

---

### 3.3 Get District by ID
**GET** `/api/v1/admin/master/districts/:id`

**Response:** Mengembalikan district dengan full location hierarchy (city dan province)

---

### 3.4 Update District
**PUT** `/api/v1/admin/master/districts/:id`

**Request Body:**
```json
{
  "name": "Batujajar Updated",
  "code": "3273010",
  "postal_code": "40561",
  "city_id": 1,
  "is_active": true
}
```

---

### 3.5 Delete District
**DELETE** `/api/v1/admin/master/districts/:id`

**Validasi Relasional:**
- District **TIDAK BISA** dihapus jika:
  - Masih digunakan oleh companies

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete district: it is still referenced by companies"
}
```

---

## 4. Industry CRUD

### 4.1 Create Industry
**POST** `/api/v1/admin/master/industries`

**Request Body:**
```json
{
  "name": "Technology",
  "slug": "technology",
  "description": "Technology industry",
  "icon_url": "https://example.com/icon.png",
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Industry created successfully",
  "data": {
    "id": 1,
    "name": "Technology",
    "slug": "technology",
    "description": "Technology industry",
    "icon_url": "https://example.com/icon.png",
    "is_active": true
  }
}
```

**Validasi:**
- `name`: required, min 2, max 255 characters, unique
- `slug`: required, min 2, max 255 characters, unique (auto-generated jika tidak provided)
- `description`: optional
- `icon_url`: optional, must be valid URL if provided
- `is_active`: boolean (default: true)

**Catatan:** Jika `slug` tidak provided, akan auto-generated dari `name`. Jika slug sudah ada, akan ditambahkan timestamp.

---

### 4.2 Get All Industries
**GET** `/api/v1/admin/master/industries?search=tech&active=true`

**Query Parameters:**
- `search` (optional): Search by name
- `active` (optional): Filter by active status

---

### 4.3 Get Industry by ID
**GET** `/api/v1/admin/master/industries/:id`

---

### 4.4 Update Industry
**PUT** `/api/v1/admin/master/industries/:id`

**Request Body:**
```json
{
  "name": "Technology Updated",
  "slug": "technology-updated",
  "description": "Updated description",
  "icon_url": "https://example.com/new-icon.png",
  "is_active": true
}
```

**Validasi:**
- Semua field optional
- Jika `name` diubah, akan dicek duplikasi
- Jika `slug` diubah, akan dicek duplikasi

---

### 4.5 Delete Industry
**DELETE** `/api/v1/admin/master/industries/:id`

**Validasi Relasional:**
- Industry **TIDAK BISA** dihapus jika:
  - Masih digunakan oleh companies

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete industry: it is still referenced by companies"
}
```

---

## 5. Job Type CRUD

### 5.1 Create Job Type
**POST** `/api/v1/admin/master/job-types`

**Request Body:**
```json
{
  "name": "Full-Time",
  "code": "full_time",
  "order": 1
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Job type created successfully",
  "data": {
    "id": 1,
    "name": "Full-Time",
    "code": "full_time",
    "order": 1
  }
}
```

**Validasi:**
- `name`: required, min 2, max 100 characters
- `code`: required, min 2, max 30 characters, unique
- `order`: optional integer (default: 0)

---

### 5.2 Get All Job Types
**GET** `/api/v1/admin/master/job-types`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Job types retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Full-Time",
      "code": "full_time",
      "order": 1
    },
    {
      "id": 2,
      "name": "Part-Time",
      "code": "part_time",
      "order": 2
    }
  ]
}
```

---

### 5.3 Get Job Type by ID
**GET** `/api/v1/admin/master/job-types/:id`

---

### 5.4 Update Job Type
**PUT** `/api/v1/admin/master/job-types/:id`

**Request Body:**
```json
{
  "name": "Full-Time Updated",
  "code": "full_time",
  "order": 1
}
```

**Validasi:**
- Semua field optional
- Jika `code` diubah, akan dicek duplikasi (exclude current ID)

---

### 5.5 Delete Job Type
**DELETE** `/api/v1/admin/master/job-types/:id`

**Validasi Relasional:**
- Job Type **TIDAK BISA** dihapus jika:
  - Masih digunakan oleh jobs

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete job type: it is still referenced by jobs"
}
```

---

## 6. Company Size CRUD

### 6.1 Create Company Size
**POST** `/api/v1/admin/meta/company-sizes`

**Request Body:**
```json
{
  "label": "1-10 employees",
  "min_employees": 1,
  "max_employees": 10,
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Company size created successfully",
  "data": {
    "id": 1,
    "label": "1-10 employees",
    "min_employees": 1,
    "max_employees": 10,
    "is_active": true
  }
}
```

**Validasi:**
- `label`: required, min 2, max 255 characters, unique
- `min_employees`: required, min 0
- `max_employees`: optional, must be >= min_employees (null = unlimited)
- `is_active`: boolean (default: true)

**Catatan:** `max_employees` dapat null/unlimited. Contoh: untuk "1000+" set `max_employees` ke null.

---

### 6.2 Get All Company Sizes
**GET** `/api/v1/admin/meta/company-sizes?active=true`

**Query Parameters:**
- `active` (optional): Filter by active status

---

### 6.3 Get Company Size by ID
**GET** `/api/v1/admin/meta/company-sizes/:id`

---

### 6.4 Update Company Size
**PUT** `/api/v1/admin/meta/company-sizes/:id`

**Request Body:**
```json
{
  "label": "1-10 employees Updated",
  "min_employees": 1,
  "max_employees": 10,
  "is_active": true
}
```

**Validasi:**
- Semua field optional
- Jika `label` diubah, akan dicek duplikasi
- `max_employees` harus >= `min_employees` jika keduanya provided

---

### 6.5 Delete Company Size
**DELETE** `/api/v1/admin/meta/company-sizes/:id`

**Validasi Relasional:**
- Company Size **TIDAK BISA** dihapus jika:
  - Masih digunakan oleh companies

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Cannot delete company size: it is still referenced by companies"
}
```

---

## Validasi dan Error Handling

### Validasi Duplikat

Semua entitas memiliki validasi duplikat:

1. **Province**: Code harus unique
2. **City**: Name harus unique dalam satu province
3. **District**: Name harus unique dalam satu city
4. **Industry**: Name dan slug harus unique
5. **Job Type**: Code harus unique
6. **Company Size**: Label harus unique

### Validasi Relasional (DELETE)

Semua entitas dicek apakah masih digunakan sebelum dihapus:

1. **Province**: 
   - Tidak bisa dihapus jika masih ada cities
   - Tidak bisa dihapus jika masih digunakan oleh companies

2. **City**: 
   - Tidak bisa dihapus jika masih ada districts
   - Tidak bisa dihapus jika masih digunakan oleh companies

3. **District**: 
   - Tidak bisa dihapus jika masih digunakan oleh companies

4. **Industry**: 
   - Tidak bisa dihapus jika masih digunakan oleh companies

5. **Job Type**: 
   - Tidak bisa dihapus jika masih digunakan oleh jobs

6. **Company Size**: 
   - Tidak bisa dihapus jika masih digunakan oleh companies

### Error Responses

**400 Bad Request:**
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

**404 Not Found:**
```json
{
  "success": false,
  "message": "Province not found"
}
```

**409 Conflict:**
```json
{
  "success": false,
  "message": "Province with this code already exists"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Forbidden"
}
```

---

## Testing API

### 1. Menggunakan cURL

#### Create Province
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

#### Get All Provinces
```bash
curl -X GET "http://localhost:8080/api/v1/admin/master/provinces?active=true" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### Update Province
```bash
curl -X PUT http://localhost:8080/api/v1/admin/master/provinces/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jawa Barat Updated",
    "is_active": false
  }'
```

#### Delete Province
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/master/provinces/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

### 2. Menggunakan Postman

#### Setup Collection

1. Buat collection baru: "Admin Master Data CRUD"
2. Set environment variable:
   - `base_url`: `http://localhost:8080`
   - `admin_token`: `<your_admin_token>`

#### Create Request Template

**Headers:**
```
Authorization: Bearer {{admin_token}}
Content-Type: application/json
```

**Body (for POST/PUT):**
```json
{
  "name": "Example",
  "code": "example_code",
  "is_active": true
}
```

---

### 3. Test Scenarios

#### Test Case 1: Create Province
1. **Request:** POST `/api/v1/admin/master/provinces`
2. **Body:** Valid province data
3. **Expected:** 201 Created dengan data province baru

#### Test Case 2: Create Duplicate Province Code
1. **Request:** POST `/api/v1/admin/master/provinces`
2. **Body:** Province dengan code yang sudah ada
3. **Expected:** 409 Conflict dengan message "Province with this code already exists"

#### Test Case 3: Get All Provinces
1. **Request:** GET `/api/v1/admin/master/provinces`
2. **Expected:** 200 OK dengan array provinces

#### Test Case 4: Update Province
1. **Request:** PUT `/api/v1/admin/master/provinces/:id`
2. **Body:** Updated province data
3. **Expected:** 200 OK dengan data province yang diupdate

#### Test Case 5: Delete Province dengan References
1. **Request:** DELETE `/api/v1/admin/master/provinces/:id`
2. **Condition:** Province masih memiliki cities atau digunakan oleh companies
3. **Expected:** 409 Conflict dengan message "Cannot delete province: it is still referenced by cities or companies"

#### Test Case 6: Delete Province tanpa References
1. **Request:** DELETE `/api/v1/admin/master/provinces/:id`
2. **Condition:** Province tidak memiliki cities dan tidak digunakan oleh companies
3. **Expected:** 200 OK dengan message "Province deleted successfully"

---

### 4. Test Script untuk Automation

#### Using JavaScript (Postman Tests)

```javascript
// Test Create Province
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Response has success true", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.eql(true);
});

pm.test("Province has id", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.id).to.be.a('number');
});

// Test Get Provinces
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response is array", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.be.an('array');
});
```

---

### 5. Integration Test Flow

#### Complete CRUD Flow untuk Province

1. **Create Province**
   ```
   POST /api/v1/admin/master/provinces
   → Simpan province_id dari response
   ```

2. **Get Province**
   ```
   GET /api/v1/admin/master/provinces/{province_id}
   → Verify data sesuai dengan yang dibuat
   ```

3. **Update Province**
   ```
   PUT /api/v1/admin/master/provinces/{province_id}
   → Verify data terupdate
   ```

4. **Create City dengan Province ID**
   ```
   POST /api/v1/admin/master/cities
   → Simpan city_id dari response
   ```

5. **Try Delete Province (Should Fail)**
   ```
   DELETE /api/v1/admin/master/provinces/{province_id}
   → Should return 409 Conflict (still has cities)
   ```

6. **Delete City**
   ```
   DELETE /api/v1/admin/master/cities/{city_id}
   → Verify city deleted
   ```

7. **Delete Province (Should Success)**
   ```
   DELETE /api/v1/admin/master/provinces/{province_id}
   → Should return 200 OK
   ```

---

## File Structure

### Handler
- `internal/handler/http/admin/masterdata_handler.go` - Main handler untuk semua CRUD operations

### Services
- `internal/service/master/admin_industry_service.go` - Admin industry service
- `internal/service/master/admin_province_service.go` - Admin province service
- `internal/service/master/admin_city_service.go` - Admin city service
- `internal/service/master/admin_district_service.go` - Admin district service
- `internal/service/master/admin_company_size_service.go` - Admin company size service
- `internal/service/master/admin_job_type_service.go` - Admin job type service

### Domain Models
- `internal/domain/master/models.go` - Request/Response DTOs
- `internal/domain/master/admin_service.go` - Admin service interfaces

### Routes
- `internal/routes/admin_routes.go` - Route definitions

---

## Checklist Testing

### Province CRUD
- [ ] Create province dengan data valid
- [ ] Create province dengan code duplicate (should fail)
- [ ] Get all provinces
- [ ] Get province by ID
- [ ] Update province
- [ ] Delete province tanpa references (should success)
- [ ] Delete province dengan cities (should fail)
- [ ] Delete province digunakan companies (should fail)

### City CRUD
- [ ] Create city dengan data valid
- [ ] Create city dengan name duplicate dalam province sama (should fail)
- [ ] Get all cities by province_id
- [ ] Get city by ID
- [ ] Update city
- [ ] Delete city tanpa references (should success)
- [ ] Delete city dengan districts (should fail)
- [ ] Delete city digunakan companies (should fail)

### District CRUD
- [ ] Create district dengan data valid
- [ ] Create district dengan name duplicate dalam city sama (should fail)
- [ ] Get all districts by city_id
- [ ] Get district by ID (dengan full location hierarchy)
- [ ] Update district
- [ ] Delete district tanpa references (should success)
- [ ] Delete district digunakan companies (should fail)

### Industry CRUD
- [ ] Create industry dengan data valid
- [ ] Create industry dengan name duplicate (should fail)
- [ ] Create industry dengan slug auto-generated
- [ ] Get all industries
- [ ] Get industry by ID
- [ ] Update industry
- [ ] Delete industry tanpa references (should success)
- [ ] Delete industry digunakan companies (should fail)

### Job Type CRUD
- [ ] Create job type dengan data valid
- [ ] Create job type dengan code duplicate (should fail)
- [ ] Get all job types
- [ ] Get job type by ID
- [ ] Update job type
- [ ] Delete job type tanpa references (should success)
- [ ] Delete job type digunakan jobs (should fail)

### Company Size CRUD
- [ ] Create company size dengan data valid
- [ ] Create company size dengan label duplicate (should fail)
- [ ] Create company size dengan max_employees unlimited (null)
- [ ] Get all company sizes
- [ ] Get company size by ID
- [ ] Update company size
- [ ] Delete company size tanpa references (should success)
- [ ] Delete company size digunakan companies (should fail)

### Authentication & Authorization
- [ ] Request tanpa token (should return 401)
- [ ] Request dengan token non-admin (should return 403)
- [ ] Request dengan admin token (should success)

---

## Catatan Penting

1. **Soft Delete**: Beberapa entitas menggunakan soft delete (deleted_at), bukan hard delete
2. **Cache Invalidation**: Cache akan di-invalidate otomatis setelah create/update/delete
3. **Transactional**: Semua operasi DELETE melakukan pengecekan references terlebih dahulu
4. **Error Messages**: Semua error messages dalam Bahasa Indonesia untuk konsistensi
5. **Validation**: Semua input divalidasi menggunakan struct tags (validate)

---

## Troubleshooting

### Error: "Cannot delete: still referenced"
**Solusi:** Hapus atau update semua referensi terlebih dahulu sebelum menghapus master data

### Error: "Duplicate entry"
**Solusi:** Gunakan name/code/label yang berbeda, atau update entitas yang sudah ada

### Error: "Unauthorized" / "Forbidden"
**Solusi:** Pastikan menggunakan admin token yang valid dan user memiliki role admin

---

## Update Log

- **2024-01-XX**: Initial implementation
  - Created all CRUD endpoints for Provinces, Cities, Districts, Industries, Job Types, and Company Sizes
  - Implemented relational constraint validation
  - Added duplicate validation for all entities
  - Integrated with authentication and authorization middleware

