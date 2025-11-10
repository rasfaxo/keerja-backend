$files = @(
    "create_oauth_providers_table.sql",
    "create_otp_codes_table.sql",
    "create_refresh_tokens_table.sql",
    "create_email_logs_table.sql",
    "create_notifications_table.sql",
    "000015_create_company_invitations_table.up.sql",
    "000016_create_device_tokens_table.up.sql",
    "000017_create_push_notification_logs_table.up.sql",
    "000018_create_industries_table.up.sql",
    "000019_create_company_sizes_table.up.sql",
    "000020_create_provinces_table.up.sql",
    "000021_create_cities_table.up.sql",
    "000022_create_districts_table.up.sql",
    "000023_alter_companies_add_master_relations.up.sql",
    "000024_alter_users_add_company_relation.up.sql",
    "000025_alter_user_profiles_add_master_relations.up.sql",
    "000026_create_job_master_tables.up.sql",
    "000027_add_job_master_data_fks.up.sql",
    "000028_remove_legacy_employment_type_constraints.up.sql",
    "20251029_add_missing_profile_fields.sql"
)

foreach ($file in $files) {
    Write-Host "Migrating: $file"
    psql -U postgres -h localhost -d keerja_db -f "database/migrations/$file"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error migrating $file"
        exit 1
    }
}