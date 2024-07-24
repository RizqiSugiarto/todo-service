CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    invitation_label_id UUID,
    invitation_category_id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_invitation_label
        FOREIGN KEY(invitation_label_id) 
        REFERENCES invitation_labels(id),
    CONSTRAINT fk_invitation_category
        FOREIGN KEY(invitation_category_id) 
        REFERENCES invitation_categories(id)
);

-- Add the index for user_id
CREATE INDEX idx_user_id_invitations ON invitations(user_id);