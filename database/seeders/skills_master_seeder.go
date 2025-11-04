package seeders

import (
	"keerja-backend/internal/domain/master"
	"log"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SkillsMasterSeeder seeds the skills_master table
func SkillsMasterSeeder(db *gorm.DB) error {
	log.Println("Seeding skills_master table...")

	skills := []master.SkillsMaster{
		// Programming Languages
		{Code: "GO", Name: "Go", NormalizedName: "go", Description: "Programming language developed by Google", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 95, IsActive: true},
		{Code: "JAVASCRIPT", Name: "JavaScript", NormalizedName: "javascript", Description: "Popular scripting language for web development", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 98, IsActive: true},
		{Code: "TYPESCRIPT", Name: "TypeScript", NormalizedName: "typescript", Description: "Typed superset of JavaScript", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 92, IsActive: true},
		{Code: "PYTHON", Name: "Python", NormalizedName: "python", Description: "High-level general-purpose programming language", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 97, IsActive: true},
		{Code: "JAVA", Name: "Java", NormalizedName: "java", Description: "Object-oriented programming language", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 95, IsActive: true},
		{Code: "PHP", Name: "PHP", NormalizedName: "php", Description: "Server-side scripting language", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 85, IsActive: true},
		{Code: "RUBY", Name: "Ruby", NormalizedName: "ruby", Description: "Dynamic programming language", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 75, IsActive: true},
		{Code: "CPP", Name: "C++", NormalizedName: "cpp", Description: "General-purpose programming language", SkillType: "technical", DifficultyLevel: "advanced", PopularityScore: 80, IsActive: true},
		{Code: "CSHARP", Name: "C#", NormalizedName: "csharp", Description: "Microsoft .NET language", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},
		{Code: "KOTLIN", Name: "Kotlin", NormalizedName: "kotlin", Description: "Modern Android development language", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 82, IsActive: true},
		{Code: "SWIFT", Name: "Swift", NormalizedName: "swift", Description: "iOS development language", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 80, IsActive: true},
		{Code: "RUST", Name: "Rust", NormalizedName: "rust", Description: "Systems programming language", SkillType: "technical", DifficultyLevel: "advanced", PopularityScore: 78, IsActive: true},

		// Frontend Frameworks
		{Code: "REACT", Name: "React", NormalizedName: "react", Description: "JavaScript library for building user interfaces", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 96, IsActive: true},
		{Code: "VUEJS", Name: "Vue.js", NormalizedName: "vuejs", Description: "Progressive JavaScript framework", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 88, IsActive: true},
		{Code: "ANGULAR", Name: "Angular", NormalizedName: "angular", Description: "TypeScript-based web framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},
		{Code: "NEXTJS", Name: "Next.js", NormalizedName: "nextjs", Description: "React framework for production", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 90, IsActive: true},
		{Code: "SVELTE", Name: "Svelte", NormalizedName: "svelte", Description: "Compile-time framework", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 75, IsActive: true},

		// Backend Frameworks
		{Code: "NODEJS", Name: "Node.js", NormalizedName: "nodejs", Description: "JavaScript runtime for server-side", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 95, IsActive: true},
		{Code: "EXPRESS", Name: "Express.js", NormalizedName: "expressjs", Description: "Node.js web framework", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 92, IsActive: true},
		{Code: "DJANGO", Name: "Django", NormalizedName: "django", Description: "Python web framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 90, IsActive: true},
		{Code: "FLASK", Name: "Flask", NormalizedName: "flask", Description: "Lightweight Python web framework", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 85, IsActive: true},
		{Code: "LARAVEL", Name: "Laravel", NormalizedName: "laravel", Description: "PHP web framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 88, IsActive: true},
		{Code: "SPRING", Name: "Spring Boot", NormalizedName: "springboot", Description: "Java application framework", SkillType: "technical", DifficultyLevel: "advanced", PopularityScore: 90, IsActive: true},
		{Code: "RAILS", Name: "Ruby on Rails", NormalizedName: "rails", Description: "Ruby web framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 78, IsActive: true},
		{Code: "ASPNET", Name: "ASP.NET Core", NormalizedName: "aspnetcore", Description: "Microsoft web framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},

		// Databases
		{Code: "POSTGRESQL", Name: "PostgreSQL", NormalizedName: "postgresql", Description: "Advanced open-source database", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 92, IsActive: true},
		{Code: "MYSQL", Name: "MySQL", NormalizedName: "mysql", Description: "Popular open-source database", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 90, IsActive: true},
		{Code: "MONGODB", Name: "MongoDB", NormalizedName: "mongodb", Description: "NoSQL document database", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 88, IsActive: true},
		{Code: "REDIS", Name: "Redis", NormalizedName: "redis", Description: "In-memory data structure store", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},
		{Code: "ELASTICSEARCH", Name: "Elasticsearch", NormalizedName: "elasticsearch", Description: "Search and analytics engine", SkillType: "technical", DifficultyLevel: "advanced", PopularityScore: 82, IsActive: true},

		// Cloud & DevOps
		{Code: "AWS", Name: "AWS", NormalizedName: "aws", Description: "Amazon Web Services", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 95, IsActive: true},
		{Code: "AZURE", Name: "Azure", NormalizedName: "azure", Description: "Microsoft cloud platform", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 88, IsActive: true},
		{Code: "GCP", Name: "Google Cloud Platform", NormalizedName: "gcp", Description: "Google cloud services", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},
		{Code: "DOCKER", Name: "Docker", NormalizedName: "docker", Description: "Containerization platform", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 93, IsActive: true},
		{Code: "KUBERNETES", Name: "Kubernetes", NormalizedName: "kubernetes", Description: "Container orchestration", SkillType: "technical", DifficultyLevel: "advanced", PopularityScore: 90, IsActive: true},
		{Code: "JENKINS", Name: "Jenkins", NormalizedName: "jenkins", Description: "CI/CD automation server", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 80, IsActive: true},
		{Code: "GIT", Name: "Git", NormalizedName: "git", Description: "Version control system", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 98, IsActive: true},
		{Code: "GITLAB", Name: "GitLab", NormalizedName: "gitlab", Description: "DevOps platform", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 82, IsActive: true},
		{Code: "GITHUB", Name: "GitHub", NormalizedName: "github", Description: "Code hosting platform", SkillType: "technical", DifficultyLevel: "beginner", PopularityScore: 95, IsActive: true},

		// Mobile Development
		{Code: "REACT_NATIVE", Name: "React Native", NormalizedName: "reactnative", Description: "Cross-platform mobile framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 88, IsActive: true},
		{Code: "FLUTTER", Name: "Flutter", NormalizedName: "flutter", Description: "Google's mobile UI framework", SkillType: "technical", DifficultyLevel: "intermediate", PopularityScore: 90, IsActive: true},

		// Soft Skills
		{Code: "COMMUNICATION", Name: "Communication", NormalizedName: "communication", Description: "Effective interpersonal communication", SkillType: "soft", DifficultyLevel: "intermediate", PopularityScore: 95, IsActive: true},
		{Code: "TEAMWORK", Name: "Teamwork", NormalizedName: "teamwork", Description: "Collaborative working ability", SkillType: "soft", DifficultyLevel: "beginner", PopularityScore: 93, IsActive: true},
		{Code: "PROBLEM_SOLVING", Name: "Problem Solving", NormalizedName: "problemsolving", Description: "Analytical thinking and solutions", SkillType: "soft", DifficultyLevel: "intermediate", PopularityScore: 98, IsActive: true},
		{Code: "TIME_MANAGEMENT", Name: "Time Management", NormalizedName: "timemanagement", Description: "Efficient task prioritization", SkillType: "soft", DifficultyLevel: "intermediate", PopularityScore: 90, IsActive: true},
		{Code: "LEADERSHIP", Name: "Leadership", NormalizedName: "leadership", Description: "Team guidance and motivation", SkillType: "soft", DifficultyLevel: "advanced", PopularityScore: 92, IsActive: true},
		{Code: "ADAPTABILITY", Name: "Adaptability", NormalizedName: "adaptability", Description: "Flexibility in changing environments", SkillType: "soft", DifficultyLevel: "intermediate", PopularityScore: 88, IsActive: true},
		{Code: "CRITICAL_THINKING", Name: "Critical Thinking", NormalizedName: "criticalthinking", Description: "Objective analysis and evaluation", SkillType: "soft", DifficultyLevel: "advanced", PopularityScore: 90, IsActive: true},
		{Code: "CREATIVITY", Name: "Creativity", NormalizedName: "creativity", Description: "Innovative and original thinking", SkillType: "soft", DifficultyLevel: "intermediate", PopularityScore: 85, IsActive: true},
	}

	// Use OnConflict to update existing skills or create new ones
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"code", "normalized_name", "description", "skill_type", "difficulty_level", "popularity_score", "is_active", "updated_at"}),
	}).Create(&skills)

	if result.Error != nil {
		log.Printf("Failed to seed skills: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d skills", len(skills))
	return nil
}

// Helper function to normalize skill names
func normalizeSkillName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}
