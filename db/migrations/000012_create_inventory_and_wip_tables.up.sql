CREATE TABLE IF NOT EXISTS inventory_items (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    part_number VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity NUMERIC(14,4) NOT NULL DEFAULT 0,
    unit VARCHAR(50),
    location VARCHAR(255),
    lot_number VARCHAR(100),
    serial_number VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory_movements (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    inventory_item_id INTEGER NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    movement_type VARCHAR(50) NOT NULL,
    quantity NUMERIC(14,4) NOT NULL,
    reference_type VARCHAR(100),
    reference_uid VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wip_items (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    inventory_item_id INTEGER NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    quantity NUMERIC(14,4) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'started',
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inventory_items_name ON inventory_items(name);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_inventory_item_id ON inventory_movements(inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_wip_items_job_id ON wip_items(job_id);
