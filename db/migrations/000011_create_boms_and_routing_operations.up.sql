CREATE TABLE IF NOT EXISTS boms (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    production_order_id INTEGER REFERENCES production_orders(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    part_number VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bom_lines (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    bom_id INTEGER NOT NULL REFERENCES boms(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    part_number VARCHAR(100),
    quantity NUMERIC(12,2) NOT NULL DEFAULT 1,
    unit VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS routing_operations (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    production_order_id INTEGER REFERENCES production_orders(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    sequence_number INTEGER NOT NULL DEFAULT 1,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    estimated_duration_hours NUMERIC(10,2),
    machine_id INTEGER REFERENCES machines(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_boms_name ON boms(name);
CREATE INDEX IF NOT EXISTS idx_bom_lines_bom_id ON bom_lines(bom_id);
CREATE INDEX IF NOT EXISTS idx_routing_operations_sequence ON routing_operations(sequence_number);
