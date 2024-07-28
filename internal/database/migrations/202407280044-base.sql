CREATE EXTENSION moddatetime;

-- Schemas

CREATE TABLE dp_schemas (
    schema_id VARCHAR(255) PRIMARY KEY,
    picture TEXT,
    description TEXT NOT NULL DEFAULT '',

    initial_sql TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TRIGGER dp_schemas_moddatetime
BEFORE UPDATE ON dp_schemas
FOR EACH ROW
EXECUTE PROCEDURE MODDATETIME(updated_at);

CREATE TYPE dp_difficulty AS ENUM ('easy', 'medium', 'hard');

CREATE TABLE dp_questions (
    question_id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    -- schema_id turns to NULL if the schema is deleted
    schema_id VARCHAR(255) REFERENCES dp_schemas ON DELETE SET NULL,

    schema_type VARCHAR(255) NOT NULL,
    difficulty DP_DIFFICULTY NOT NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    answer TEXT NOT NULL,
    solution_video TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TRIGGER dp_questions_moddatetime
BEFORE UPDATE ON dp_questions
FOR EACH ROW
EXECUTE PROCEDURE MODDATETIME(updated_at);

-- Users

CREATE TABLE dp_groups (
    group_id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TRIGGER dp_groups_moddatetime
BEFORE UPDATE ON dp_groups
FOR EACH ROW
EXECUTE PROCEDURE MODDATETIME(updated_at);

CREATE TABLE dp_users (
    logto_user_id TEXT PRIMARY KEY NOT NULL,
    group_id UUID REFERENCES dp_groups ON DELETE SET NULL
);
