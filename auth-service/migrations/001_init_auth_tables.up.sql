CREATE SCHEMA IF NOT EXISTS auth_service;

-- Reference table for the valid account lifecycle states. Adding a new status
-- later is a one-row INSERT (no ALTER on the user table, independent of its
-- size). Value validity is enforced via the FK from user.status; transition
-- rules (which state can move to which) live in the application — is_terminal
-- is a convenience flag the app reads, not a state machine the DB enforces.
CREATE TABLE IF NOT EXISTS auth_service.account_status (
    code        VARCHAR(30) PRIMARY KEY,
    label       VARCHAR(50) NOT NULL,
);

INSERT INTO auth_service.account_status (code, label) VALUES
    ('PENDING_VERIFICATION', 'Pending Verification'),
    ('ACTIVE',               'Active'),
    ('LOCKED',               'Locked'),
    ('DISABLED',             'Disabled'),
    ('SUSPENDED',            'Suspended'),
    ('PENDING_DELETION',     'Pending Deletion'),
    ('DELETED',              'Deleted')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE IF NOT EXISTS auth_service.user (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by INT,
    updated_by INT,

    -- username is the always-present public identity: chosen at signup or
    -- auto-generated. Never null, so it needs no "at least one identifier" check.
    username VARCHAR(50) NOT NULL,
    -- email is optional, private (login / search / recovery), never displayed
    -- as a handle.
    email VARCHAR(255),
    -- nullable: OAuth-only accounts have no local password.
    password_hash VARCHAR(255),

    auth_provider VARCHAR(30) NOT NULL DEFAULT 'local',
    type VARCHAR(30) NOT NULL DEFAULT 'user',
    -- status validity enforced by FK into account_status (status owned by auth,
    -- not profile).
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_VERIFICATION'
        REFERENCES auth_service.account_status(code),

    email_verified BOOLEAN NOT NULL DEFAULT false,
    -- true when the handle was auto-generated at signup; flips false once
    -- the user chooses their own.
    username_is_generated BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP
);

-- Case-insensitive uniqueness. username is NOT NULL so a plain unique index;
-- email is nullable, so a partial index (multiple NULLs allowed for accounts
-- without an email set).
CREATE UNIQUE INDEX IF NOT EXISTS uk_user_username_lower
    ON auth_service.user (lower(username));

CREATE UNIQUE INDEX IF NOT EXISTS uk_user_email_lower
    ON auth_service.user (lower(email))
    WHERE email IS NOT NULL;
