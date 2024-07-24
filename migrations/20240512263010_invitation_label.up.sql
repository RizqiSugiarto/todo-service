CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE invitation_labels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Add the index for user_id
CREATE INDEX idx_user_id_invitation_labels ON invitation_labels(user_id);