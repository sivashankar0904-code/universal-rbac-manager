CREATE TABLE IF NOT EXISTS organisation (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    key             text UNIQUE NOT NULL,          -- 'acme', slug for URLs
    name            text NOT NULL,
    description     text,
    parent_id       uuid REFERENCES organisation (id),
    slugName        text,  
    is_active       boolean NOT NULL DEFAULT true,
    meta            jsonb NOT NULL,
    status          text,
    searchTags      TEXT[],
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_organisation_parent ON organisation (parent_id);
