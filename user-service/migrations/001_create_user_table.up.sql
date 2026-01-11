CREATE TABLE IF NOT EXISTS user_service.user (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by INT,
    updated_by INT,

    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) DEFAULT NULL,
    email VARCHAR(255) DEFAULT NULL,
    display_name VARCHAR(100) NOT NULL,
    auth_provider VARCHAR(50) DEFAULT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES user_service.user(id) ON DELETE SET NULL,
    CONSTRAINT fk_updated_by FOREIGN KEY (updated_by) REFERENCES user_service.user(id) ON DELETE SET NULL
);