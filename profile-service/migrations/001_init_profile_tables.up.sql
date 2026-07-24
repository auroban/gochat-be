CREATE TABLE IF NOT EXISTS profile_service.profile (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by INT,
    updated_by INT,

    user_id INT NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    profile_photo_key VARCHAR(100) NOT NULL,
    bio VARCHAR(100) NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES auth_service.user(id),
    CONSTRAINT uk_user_id
        UNIQUE(user_id),
    CONSTRAINT uk_profile_photo_key
        UNIQUE(profile_photo_key)
);

CREATE INDEX idx_profile_user_id ON profile_service.profile(user_id);

CREATE TABLE IF NOT EXISTS profile_service.profile_history (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by INT,

    profile_id INT NOT NULL,
    action VARCHAR(50) NOT NULL,
    state_before JSONB NOT NULL,
    state_after JSONB NOT NULL,

    CONSTRAINT fk_profile
        FOREIGN KEY(profile_id)
            REFERENCES profile_service.profile(id)
    
);

CREATE INDEX idx_profile_history_profile_id ON profile_service.profile_history(profile_id);