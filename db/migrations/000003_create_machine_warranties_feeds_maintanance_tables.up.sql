
-- MACHINES

CREATE TABLE machines (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    description TEXT,
    model TEXT,
    vendor INTEGER REFERENCES vendors(id),
    serial_no TEXT,
    installation_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE machine_locations (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    plant TEXT,
    area TEXT,
    cell_or_line TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);


-- WARRANTIES

CREATE TABLE warranties (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    type TEXT CHECK (type IN ('machine', 'part')),
    ref_id INTEGER NOT NULL,
    provider INTEGER REFERENCES vendors(id),
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE INDEX idx_warranty_type_ref2 ON warranties(type, ref_id);
CREATE INDEX idx_warranties_provider ON warranties(provider);


-- SPARE PARTS

CREATE TABLE spare_parts (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE spare_part_vendors (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id) ON DELETE CASCADE,
    vendor_id INTEGER, 
    model_no TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE spare_parts_machines (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id) ON DELETE CASCADE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);


-- FEEDS

CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    description TEXT,
    low_value NUMERIC,
    high_value NUMERIC,
    display_unit TEXT,
    margin_of_error NUMERIC,
    connection_type TEXT, -- e.g. mqtt, plc
    ingestion_type TEXT CHECK (ingestion_type IN ('polling', 'subscribe')),
    polling_rate INTEGER,
    json JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE feed_connection_types (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    structure JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE feed_alerts (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    alert_type TEXT,
    low_alert_value NUMERIC,
    high_alert_value NUMERIC,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE feed_machine_mapping (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);


-- MAINTENANCE

CREATE TABLE maintenance_types (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE maintenance_templates (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    type_id INTEGER REFERENCES maintenance_types(id),
    name TEXT NOT NULL,
    description TEXT,
    machine_id INTEGER REFERENCES machines(id), --  remove keeep in separate table
    frequency TEXT, 
    verified_by INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE maintenance_template_procedures (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    maintenance_template_id INTEGER NOT NULL REFERENCES maintenance_templates(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE maintenance_procedure_steps (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    maintenance_template_procedure_id INTEGER NOT NULL REFERENCES maintenance_template_procedures(id) ON DELETE CASCADE,
    step_no INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE maintenance_procedure_steps_checklist (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    image TEXT,
    before_step_id INTEGER REFERENCES maintenance_procedure_steps(id),
    after_step_id INTEGER REFERENCES maintenance_procedure_steps(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

CREATE TABLE maintenance_schedules (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    date DATE NOT NULL,
    maintenance_template_id INTEGER REFERENCES maintenance_templates(id),
    assigned_to INTEGER ,
    machine_id INTEGER REFERENCES machines(id),
    start_timestamp TIMESTAMP,
    end_timestamp TIMESTAMP,
    findings TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
    -- sheduled maintainance start status
    -- status
    -- started_on
    -- ended on 
);


CREATE TABLE maintenance_checklists (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    maintenance_schedule_id INTEGER REFERENCES maintenance_schedules(id),
    checklist_id INTEGER REFERENCES maintenance_procedure_steps_checklist(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);


-- INDEX OPTIMIZATIONS

CREATE INDEX idx_machine_code ON machines(code);
CREATE INDEX idx_feed_code ON feeds(code);
CREATE INDEX idx_spare_code ON spare_parts(code);
CREATE INDEX idx_warranty_type_ref ON warranties(type, ref_id);
