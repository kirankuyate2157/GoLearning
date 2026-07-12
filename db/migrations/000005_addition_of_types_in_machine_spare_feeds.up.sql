-- Add type_id column to machines table
ALTER TABLE machines 
ADD COLUMN type_id INTEGER REFERENCES types(id);

-- Add operational and cost columns to machines table
ALTER TABLE machines 
ADD COLUMN min_operation_value NUMERIC,
ADD COLUMN min_operation_value_unit INTEGER REFERENCES types(id),
ADD COLUMN max_operation_value NUMERIC,
ADD COLUMN max_operation_value_unit INTEGER REFERENCES types(id),
ADD COLUMN electrity_uph NUMERIC,
ADD COLUMN rent_purchase_cost NUMERIC,
ADD COLUMN consumable_cost NUMERIC,
ADD COLUMN wages_per_hour NUMERIC;

-- Add type_id column to spare_parts table
ALTER TABLE spare_parts 
ADD COLUMN type_id INTEGER REFERENCES types(id);

-- Add type_id column to feeds table
ALTER TABLE feeds 
ADD COLUMN type_id INTEGER REFERENCES types(id);

-- Create indexes for better query performance
CREATE INDEX idx_machines_type ON machines(type_id);
CREATE INDEX idx_spare_parts_type ON spare_parts(type_id);
CREATE INDEX idx_feeds_type ON feeds(type_id);
