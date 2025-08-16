CREATE TABLE jwks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kid TEXT UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    algorithm TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE
);
