CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255), -- Nullable for GitHub users
    github_id VARCHAR(255) UNIQUE,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bugs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    description TEXT,
    url TEXT,
    user_agent TEXT,
    viewport JSONB,
    console_logs JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
