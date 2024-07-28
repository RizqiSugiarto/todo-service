CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE texts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    text TEXT NOT NULL,
    activity_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_activity_id
        FOREIGN KEY(activity_id) 
        REFERENCES activities(id)
);