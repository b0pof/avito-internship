-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);


CREATE TYPE tender_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status tender_status NOT NULL DEFAULT 'Created',
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tender_version (
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type tender_service_type NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK ( version >= 1 )
);

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

CREATE TYPE bid_author_type AS ENUM (
    'Organization',
    'User'
);

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status bid_status NOT NULL DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type NOT NULL,
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid_version (
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK ( version >= 1 )
);

-- Reset updated_at function
CREATE OR REPLACE FUNCTION update_time()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- trigger for employee updated_at
CREATE TRIGGER employee_updator
    BEFORE UPDATE ON employee
    FOR EACH ROW EXECUTE FUNCTION update_time();

-- trigger for organization updated_at
CREATE TRIGGER organization_updator
    BEFORE UPDATE ON organization
    FOR EACH ROW EXECUTE FUNCTION update_time();

-- indexes

CREATE INDEX employee_id_idx ON employee(username);
CREATE INDEX organization_id_idx ON organization(id);
CREATE INDEX bid_id_idx ON bid(id);
CREATE INDEX tender_id_idx ON tender(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS bid_version;
DROP TABLE IF EXISTS bid;
DROP TYPE IF EXISTS bid_author_type;
DROP TYPE IF EXISTS bid_status;
DROP TABLE IF EXISTS tender_version;
DROP TABLE IF EXISTS tender;
DROP TYPE IF EXISTS tender_status;
DROP TYPE IF EXISTS tender_service_type;
DROP TABLE IF EXISTS organization_responsible;
DROP TABLE IF EXISTS organization;
DROP TYPE IF EXISTS organization_type;
DROP TABLE IF EXISTS employee;
DROP FUNCTION IF EXISTS update_time;
DROP TRIGGER IF EXISTS employee_updator ON employee;
DROP TRIGGER IF EXISTS organization_updator ON organization;
DROP INDEX IF EXISTS employee_username_idx;
DROP INDEX IF EXISTS organization_id_idx;
DROP INDEX IF EXISTS bid_id_idx;
DROP INDEX IF EXISTS tender_id_idx;

-- +goose StatementEnd
