-- Create company_invitations table
CREATE TABLE IF NOT EXISTS company_invitations (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email VARCHAR(150) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    position VARCHAR(100),
    role VARCHAR(30) NOT NULL DEFAULT 'recruiter' CHECK (role IN ('admin', 'recruiter', 'viewer')),
    token VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'expired')),
    invited_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    accepted_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    accepted_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_company_invitations_company_id ON company_invitations(company_id);
CREATE INDEX idx_company_invitations_email ON company_invitations(email);
CREATE INDEX idx_company_invitations_token ON company_invitations(token);
CREATE INDEX idx_company_invitations_status ON company_invitations(status);
CREATE INDEX idx_company_invitations_expires_at ON company_invitations(expires_at);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_company_invitations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_company_invitations_updated_at
    BEFORE UPDATE ON company_invitations
    FOR EACH ROW
    EXECUTE FUNCTION update_company_invitations_updated_at();

-- Add comment to table
COMMENT ON TABLE company_invitations IS 'Stores company employee invitation records with token-based acceptance system';
COMMENT ON COLUMN company_invitations.token IS 'Unique invitation token valid for 7 days';
COMMENT ON COLUMN company_invitations.status IS 'Invitation status: pending, accepted, rejected, or expired';
COMMENT ON COLUMN company_invitations.expires_at IS 'Token expiration timestamp (7 days from creation)';
