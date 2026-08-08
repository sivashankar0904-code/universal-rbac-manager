CREATE TABLE IF NOT EXISTS Address (
    id              TEXT uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organizationID  TEXT,
    userID          TEXT,
    email           TEXT,
    mobile          TEXT,
    landline        TEXT,
    othernumber     TEXT
    is_active       boolean NOT NULL DEFAULT true,
    status          jsonb,
    searchTags      TEXT[],    
    createdAt       timestamptz DEFAULT now(),
    updatedAt       timestamptz DEFAULT now()    
)