CREATE TABLE IF NOT EXISTS mfa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    method TEXT NOT NULL,
    secret TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    last_used TIMESTAMP WITH TIME ZONE
);
