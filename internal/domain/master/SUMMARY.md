# Master Domain

## Overview

The Master domain manages reference/master data for the Keerja job portal, specifically skills and benefits that are used throughout the system. This domain provides standardized data for job postings, user profiles, and matching algorithms.

## Entities (2)

### 1. BenefitsMaster

Master data for job benefits (e.g., health insurance, remote work, etc.).

**Fields:**

- `ID` (int64): Primary key
- `Code` (string): Unique benefit code (max 50 chars)
- `Name` (string): Unique benefit name (max 150 chars)
- `Category` (string): Enum: financial, health, career, lifestyle, flexibility, other
- `Description` (text): Detailed description
- `Icon` (string): Icon identifier for UI (max 100 chars)
- `IsActive` (bool): Active status (default: true)
- `PopularityScore` (numeric 5,2): Score 0-100, indexed DESC
- `CreatedAt`, `UpdatedAt`: Timestamps
- `DeletedAt`: Soft delete support

**Helper Methods:**

- `IsFinancial()`, `IsHealth()`, `IsCareer()`, `IsLifestyle()`, `IsFlexibility()`: Category checks
- `IsPopular()`: Checks if popularity >= 70
- `IncrementPopularity(amount)`: Increase score (capped at 100)

**Categories:**

- **Financial**: Salary, bonus, stock options, retirement plans
- **Health**: Medical insurance, dental, mental health support
- **Career**: Training, certifications, career development
- **Lifestyle**: Gym, meals, transportation, company events
- **Flexibility**: Remote work, flexible hours, work-life balance
- **Other**: Miscellaneous benefits

**Usage:**

- Referenced by JobBenefit (job benefits)
- Used in job search filters
- Used in company profiles
- Popularity tracked for trending benefits

### 2. SkillsMaster

Master data for skills (technical, soft, languages, tools).

**Fields:**

- `ID` (int64): Primary key
- `Code` (string): Unique skill code (max 50 chars)
- `Name` (string): Unique skill name (max 150 chars)
- `NormalizedName` (string): Lowercase normalized name for matching
- `CategoryID` (int64): Reference to JobCategory (optional)
- `Description` (text): Detailed description
- `SkillType` (string): Enum: technical, soft, language, tool
- `DifficultyLevel` (string): Enum: beginner, intermediate, advanced
- `PopularityScore` (numeric 5,2): Score 0-100, indexed DESC
- `Aliases` (text[]): Alternative names for the skill
- `ParentID` (int64): Self-referential for skill hierarchy
- `IsActive` (bool): Active status (default: true)
- `CreatedAt`, `UpdatedAt`: Timestamps
- `DeletedAt`: Soft delete support

**Relationships:**

- `Parent` (SkillsMaster): Parent skill in hierarchy
- `Children` ([]SkillsMaster): Child skills in hierarchy
- Note: CategoryID references job_categories (cross-domain reference)

**Helper Methods:**

- `IsTechnical()`, `IsSoft()`, `IsLanguage()`, `IsTool()`: Type checks
- `IsBeginner()`, `IsIntermediate()`, `IsAdvanced()`: Difficulty checks
- `IsPopular()`: Checks if popularity >= 70
- `HasParent()`, `HasChildren()`: Hierarchy checks
- `IncrementPopularity(amount)`: Increase score (capped at 100)
- `AddAlias(alias)`, `RemoveAlias(alias)`, `HasAlias(alias)`: Alias management

**Skill Types:**

- **Technical**: Programming languages, frameworks, databases
- **Soft**: Communication, leadership, problem-solving
- **Language**: English, Spanish, Mandarin, etc.
- **Tool**: Software, platforms, IDE, design tools

**Difficulty Levels:**

- **Beginner**: 0-2 years experience
- **Intermediate**: 2-5 years experience
- **Advanced**: 5+ years experience

**Hierarchy Example:**

```
JavaScript (root)
├── React (child)
│   ├── Next.js (grandchild)
│   └── Redux (grandchild)
├── Node.js (child)
│   └── Express.js (grandchild)
└── TypeScript (child)
```

**Usage:**

- Referenced by JobSkill (job requirements)
- Referenced by UserSkill (user profiles)
- Used in job matching algorithms
- Used in skill-based search and recommendations
- Popularity tracked for trending skills

## Repository Layer (2 Interfaces)

### BenefitsMasterRepository (26 methods)

**Basic CRUD (6):**

- `Create`, `FindByID`, `FindByCode`, `FindByName`, `Update`, `Delete`

**Listing & Search (4):**

- `List(filter)`: Paginated list with filtering
- `ListActive`: All active benefits
- `ListByCategory(category)`: Benefits in category
- `SearchBenefits(query, page, pageSize)`: Full-text search

**Category Operations (2):**

- `GetCategories`: List all unique categories
- `CountByCategory(category)`: Count benefits in category

**Popularity Operations (4):**

- `UpdatePopularity(id, score)`: Set popularity score
- `IncrementPopularity(id, amount)`: Increase popularity
- `GetMostPopular(limit)`: Top benefits by popularity
- `GetByPopularityRange(min, max)`: Benefits in range

**Status Operations (2):**

- `Activate(id)`: Set IsActive = true
- `Deactivate(id)`: Set IsActive = false

**Bulk Operations (2):**

- `BulkCreate(benefits)`: Create multiple benefits
- `BulkUpdatePopularity(updates)`: Update multiple scores

**Statistics (4):**

- `Count`: Total benefits count
- `CountActive`: Active benefits count
- `GetBenefitStats`: Comprehensive statistics

### SkillsMasterRepository (48 methods)

**Basic CRUD (6):**

- `Create`, `FindByID`, `FindByCode`, `FindByName`, `FindByNormalizedName`, `Update`, `Delete`

**Listing & Search (7):**

- `List(filter)`: Paginated list with filtering
- `ListActive`: All active skills
- `ListByType(type)`: Skills by type
- `ListByDifficulty(level)`: Skills by difficulty
- `ListByCategory(categoryID)`: Skills in category
- `SearchSkills(query, page, pageSize)`: Full-text search
- `SearchByAlias(alias)`: Find skills with specific alias

**Hierarchy Operations (5):**

- `GetRootSkills`: Skills without parent
- `GetChildren(parentID)`: Direct children of skill
- `GetParent(childID)`: Parent of skill
- `GetSkillTree(rootID)`: Complete tree with children
- `HasChildren(id)`: Check if skill has children

**Popularity Operations (5):**

- `UpdatePopularity(id, score)`: Set popularity score
- `IncrementPopularity(id, amount)`: Increase popularity
- `GetMostPopular(limit)`: Top skills by popularity
- `GetByPopularityRange(min, max)`: Skills in range
- `GetPopularByType(type, limit)`: Top skills by type

**Alias Operations (4):**

- `AddAlias(id, alias)`: Add new alias
- `RemoveAlias(id, alias)`: Remove alias
- `UpdateAliases(id, aliases)`: Replace all aliases
- `FindByAliases(aliases)`: Find skills with aliases

**Status Operations (2):**

- `Activate(id)`: Set IsActive = true
- `Deactivate(id)`: Set IsActive = false

**Bulk Operations (2):**

- `BulkCreate(skills)`: Create multiple skills
- `BulkUpdatePopularity(updates)`: Update multiple scores

**Statistics (6):**

- `Count`, `CountActive`, `CountByType(type)`, `CountByDifficulty(level)`
- `GetSkillStats`: Comprehensive statistics

**Recommendations (2):**

- `GetRelatedSkills(skillID, limit)`: Related skills (siblings, similar)
- `GetComplementarySkills(skillIDs, limit)`: Skills that complement given set

## Service Layer (2 Interfaces)

### BenefitsMasterService (20 methods)

**Benefit Management (7):**

- `CreateBenefit(req)`: Create new benefit
- `UpdateBenefit(id, req)`: Update existing benefit
- `DeleteBenefit(id)`: Delete benefit (soft delete)
- `GetBenefit(id)`: Get benefit details with usage count
- `GetBenefitByCode(code)`: Get by unique code
- `GetBenefits(filter)`: Paginated list with filtering
- `SearchBenefits(query, page, pageSize)`: Full-text search

**Category Operations (3):**

- `GetBenefitsByCategory(category)`: List benefits in category
- `GetCategories`: List all categories with counts
- `GetCategoryStats`: Statistics by category

**Popularity Management (3):**

- `UpdatePopularity(id, score)`: Set popularity score
- `IncrementPopularity(id)`: Increase by default amount
- `GetMostPopular(limit)`: Top benefits
- `GetPopularByCategory(category, limit)`: Top in category

**Status Management (2):**

- `ActivateBenefit(id)`: Activate benefit
- `DeactivateBenefit(id)`: Deactivate benefit
- `GetActiveBenefits`: List all active

**Bulk Operations (4):**

- `BulkCreateBenefits(benefits)`: Create multiple
- `BulkUpdatePopularity(updates)`: Update multiple scores
- `ImportBenefits(data)`: Import from CSV/JSON with validation
- `ExportBenefits(filter)`: Export to CSV/JSON

**Statistics (2):**

- `GetBenefitStats`: Overall statistics
- `GetTrendingBenefits(limit)`: Benefits with rising popularity

**Validation (2):**

- `ValidateBenefit(req)`: Validate benefit data
- `CheckDuplicateBenefit(name, code)`: Check for duplicates

### SkillsMasterService (35 methods)

**Skill Management (7):**

- `CreateSkill(req)`: Create new skill
- `UpdateSkill(id, req)`: Update existing skill
- `DeleteSkill(id)`: Delete skill (soft delete)
- `GetSkill(id)`: Get skill details with usage count
- `GetSkillByCode(code)`: Get by unique code
- `GetSkills(filter)`: Paginated list with filtering
- `SearchSkills(query, page, pageSize)`: Full-text search

**Type & Difficulty Operations (5):**

- `GetSkillsByType(type)`: List skills by type
- `GetSkillsByDifficulty(level)`: List skills by difficulty
- `GetSkillsByCategory(categoryID)`: List skills in category
- `GetTypeStats`: Statistics by type
- `GetDifficultyStats`: Statistics by difficulty

**Hierarchy Operations (6):**

- `GetRootSkills`: Top-level skills
- `GetChildSkills(parentID)`: Direct children
- `GetParentSkill(childID)`: Parent skill
- `GetSkillTree(rootID)`: Complete tree structure
- `SetParentSkill(skillID, parentID)`: Set parent relationship
- `RemoveParentSkill(skillID)`: Remove from hierarchy

**Popularity Management (3):**

- `UpdatePopularity(id, score)`: Set popularity score
- `IncrementPopularity(id)`: Increase by default amount
- `GetMostPopular(limit)`: Top skills
- `GetPopularByType(type, limit)`: Top by type

**Alias Management (4):**

- `AddAlias(skillID, alias)`: Add new alias
- `RemoveAlias(skillID, alias)`: Remove alias
- `UpdateAliases(skillID, aliases)`: Replace all aliases
- `SearchByAlias(alias)`: Find skills with alias

**Status Management (2):**

- `ActivateSkill(id)`: Activate skill
- `DeactivateSkill(id)`: Deactivate skill
- `GetActiveSkills`: List all active

**Bulk Operations (4):**

- `BulkCreateSkills(skills)`: Create multiple
- `BulkUpdatePopularity(updates)`: Update multiple scores
- `ImportSkills(data)`: Import from CSV/JSON with validation
- `ExportSkills(filter)`: Export to CSV/JSON

**Recommendations (3):**

- `GetRelatedSkills(skillID, limit)`: Related skills
- `GetComplementarySkills(skillIDs, limit)`: Complementary skills
- `GetSkillSuggestions(userSkills, limit)`: Skill recommendations for user

**Statistics (2):**

- `GetSkillStats`: Overall statistics
- `GetTrendingSkills(limit)`: Skills with rising popularity

**Validation (3):**

- `ValidateSkill(req)`: Validate skill data
- `CheckDuplicateSkill(name, code)`: Check for duplicates
- `ValidateHierarchy(skillID, parentID)`: Validate parent-child relationship

## DTOs

### Request DTOs (4)

1. **CreateBenefitRequest**: Code, Name, Category, Description, Icon, IsActive
2. **UpdateBenefitRequest**: Code, Name, Category, Description, Icon, IsActive
3. **CreateSkillRequest**: Code, Name, NormalizedName, CategoryID, Description, SkillType, DifficultyLevel, Aliases, ParentID, IsActive
4. **UpdateSkillRequest**: Code, Name, NormalizedName, CategoryID, Description, SkillType, DifficultyLevel, Aliases, ParentID, IsActive

### Response DTOs (12)

1. **BenefitResponse**: Benefit details with usage count
2. **BenefitListResponse**: Paginated benefit list
3. **SkillResponse**: Skill details with usage count, parent, children count
4. **SkillListResponse**: Paginated skill list
5. **SkillTreeResponse**: Hierarchical skill tree
6. **BenefitStatsResponse**: Benefit statistics with distribution
7. **SkillStatsResponse**: Skill statistics with distribution
8. **CategoryInfo**: Category with label and count
9. **CategoryStatResponse**: Category statistics
10. **TypeStatResponse**: Skill type statistics
11. **DifficultyStatResponse**: Difficulty level statistics
12. **ImportResult**: Import operation result with errors

### Supporting Types (6)

1. **BenefitImportData**: Data for importing benefits
2. **BenefitExportData**: Data for exporting benefits
3. **SkillImportData**: Data for importing skills
4. **SkillExportData**: Data for exporting skills
5. **BenefitsFilter**: Search, Category, IsActive, MinPopularity, MaxPopularity, Page, PageSize, SortBy, SortOrder
6. **SkillsFilter**: Search, SkillType, DifficultyLevel, CategoryID, ParentID, IsActive, MinPopularity, MaxPopularity, HasParent, HasChildren, Page, PageSize, SortBy, SortOrder

## Business Features

### 1. Reference Data Management

- Centralized master data for skills and benefits
- Standardized naming and categorization
- Active/inactive status management
- Soft delete support

### 2. Skill Hierarchy

- Parent-child relationships for skill organization
- Multi-level hierarchy support (root → child → grandchild)
- Tree structure navigation
- Related skills discovery

### 3. Skill Aliases

- Multiple names for the same skill
- Support for abbreviations and variations
- Search by alias
- Alias management (add, remove, update)

### 4. Popularity Tracking

- Popularity score (0-100) for trending analysis
- Dynamic score updates based on usage
- Most popular items discovery
- Popularity-based recommendations

### 5. Categorization

- Benefits: 6 categories (financial, health, career, lifestyle, flexibility, other)
- Skills: 4 types (technical, soft, language, tool)
- Skills: 3 difficulty levels (beginner, intermediate, advanced)
- Category-based filtering and statistics

### 6. Import/Export

- Bulk data import from CSV/JSON
- Validation during import
- Error reporting for failed imports
- Export to CSV/JSON for backup/sharing

### 7. Search & Discovery

- Full-text search
- Filter by category/type/difficulty
- Search by alias
- Popularity-based ranking

### 8. Recommendations

- Related skills (siblings, similar categories)
- Complementary skills (often used together)
- Skill suggestions based on user profile
- Trending skills and benefits

## Technical Features

### 1. Data Persistence

- GORM models with proper relationships
- Soft delete support (DeletedAt)
- Unique constraints on code and name
- Indexed fields for performance (popularity DESC)

### 2. Validation

- Struct validation tags
- Enum validation for categories/types/levels
- Code and name uniqueness
- Hierarchy validation (no circular references)

### 3. Relationships

- Self-referential (SkillsMaster.Parent/Children)
- Cross-domain (SkillsMaster.CategoryID → JobCategory)
- One-to-many cascading

### 4. Filtering & Search

- Comprehensive filter types
- Full-text search
- Popularity range filtering
- Hierarchy filtering (HasParent, HasChildren)
- Pagination support

### 5. Statistics

- Count by category/type/difficulty
- Average popularity
- Most popular items
- Distribution statistics

### 6. PostgreSQL Features

- Array type for aliases (text[])
- Numeric type for popularity (5,2)
- Indexes for performance
- Foreign key constraints

## Key Workflows

### Skill Creation Flow

1. Admin creates skill with name, type, difficulty
2. System generates normalized name (lowercase)
3. Optionally set parent skill for hierarchy
4. Add aliases for search optimization
5. Set initial popularity score (default 0)
6. Skill is now available for job postings and user profiles

### Popularity Update Flow

1. User adds skill to profile or job requires skill
2. System increments skill popularity
3. Popularity score updated (capped at 100)
4. Trending skills list refreshed
5. Recommendations updated based on new scores

### Skill Hierarchy Navigation

1. User views root skill (e.g., JavaScript)
2. System loads immediate children (React, Node.js, TypeScript)
3. User expands React to see Next.js, Redux
4. System supports unlimited depth
5. Breadcrumb navigation from child to root

### Benefit Import Flow

1. Admin uploads CSV/JSON file with benefits
2. System validates each entry (name, code, category)
3. Check for duplicates by code/name
4. Create valid entries, report errors
5. Return ImportResult with success/failure counts
6. Admin reviews errors and fixes data

### Skill Recommendation Flow

1. User has skills: [JavaScript, React, Node.js]
2. System finds complementary skills (TypeScript, Redux, Express.js)
3. Filter by popularity and relevance
4. Exclude skills user already has
5. Return top N recommendations ordered by score

## Statistics

### Master Domain Summary

- **Entities**: 2 (BenefitsMaster, SkillsMaster)
- **Repository Methods**: 74 total (26 + 48)
- **Service Methods**: 55 total (20 + 35)
- **Request DTOs**: 4
- **Response DTOs**: 12
- **Supporting Types**: 6
- **Total Lines**: ~580 (entity: ~200, repository: ~140, service: ~240)

### Complexity Analysis

- **Hierarchy**: Self-referential relationships, tree navigation
- **Aliases**: Array field management, search optimization
- **Popularity**: Dynamic scoring, trending analysis
- **Import/Export**: Bulk operations, validation, error handling
- **Recommendations**: Related items, complementary items, suggestions

### Integration Points

- **Job Domain**: JobBenefit references BenefitsMaster, JobSkill references SkillsMaster
- **User Domain**: UserSkill references SkillsMaster
- **Company Domain**: Company benefits reference BenefitsMaster
- **Application Domain**: Matching algorithms use skills and benefits

## Data Examples

### Sample Benefits

```
Code: HEALTH_INS, Name: Health Insurance, Category: health, Popularity: 95
Code: REMOTE_WORK, Name: Remote Work Option, Category: flexibility, Popularity: 88
Code: STOCK_OPT, Name: Stock Options, Category: financial, Popularity: 72
Code: GYM_MEMBER, Name: Gym Membership, Category: lifestyle, Popularity: 65
Code: TRAINING, Name: Training Budget, Category: career, Popularity: 78
```

### Sample Skills

```
Name: JavaScript, Type: technical, Difficulty: intermediate, Popularity: 98
  ├── Name: React, Type: technical, Difficulty: intermediate, Popularity: 95
  │   ├── Name: Next.js, Type: tool, Difficulty: intermediate, Popularity: 85
  │   └── Name: Redux, Type: tool, Difficulty: intermediate, Popularity: 80
  ├── Name: Node.js, Type: technical, Difficulty: intermediate, Popularity: 92
  │   └── Name: Express.js, Type: tool, Difficulty: beginner, Popularity: 88
  └── Name: TypeScript, Type: technical, Difficulty: advanced, Popularity: 90

Name: Communication, Type: soft, Difficulty: beginner, Popularity: 87
Name: Leadership, Type: soft, Difficulty: advanced, Popularity: 82
Name: English, Type: language, Difficulty: intermediate, Popularity: 95
```

## Next Steps

1. **Implementation Phase**:

   - Implement BenefitsMasterRepository with GORM
   - Implement SkillsMasterRepository with GORM
   - Implement BenefitsMasterService with business logic
   - Implement SkillsMasterService with recommendations

2. **Data Seeding**:

   - Create comprehensive skills seed data (programming languages, frameworks, tools, soft skills, languages)
   - Create comprehensive benefits seed data (common job benefits across categories)
   - Set initial popularity scores based on industry standards
   - Establish skill hierarchies

3. **API Layer**:

   - Benefits CRUD endpoints
   - Skills CRUD endpoints
   - Search and filter endpoints
   - Import/export endpoints
   - Statistics endpoints

4. **Integration**:

   - Connect JobBenefit to BenefitsMaster
   - Connect JobSkill to SkillsMaster
   - Connect UserSkill to SkillsMaster
   - Implement skill matching algorithm
   - Implement popularity tracking on usage

5. **Testing**:

   - Unit tests for repositories
   - Unit tests for services
   - Integration tests for hierarchy operations
   - Integration tests for recommendations
   - Performance tests for search and tree operations

6. **Documentation**:
   - API documentation (Swagger)
   - Skill taxonomy documentation
   - Benefit categories guide
   - Data import guide
   - Recommendation algorithm documentation
