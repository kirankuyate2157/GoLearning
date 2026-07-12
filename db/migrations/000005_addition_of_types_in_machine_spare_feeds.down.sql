-- Drop indexes
DROP INDEX IF EXISTS idx_feeds_type;
DROP INDEX IF EXISTS idx_spare_parts_type;
DROP INDEX IF EXISTS idx_machines_type;

-- Remove type_id column from feeds table
ALTER TABLE feeds 
DROP COLUMN IF EXISTS type_id;

-- Remove type_id column from spare_parts table
ALTER TABLE spare_parts 
DROP COLUMN IF EXISTS type_id;

-- Remove operational and cost columns from machines table
ALTER TABLE machines 
DROP COLUMN IF EXISTS wages_per_hour,
DROP COLUMN IF EXISTS consumable_cost,
DROP COLUMN IF EXISTS rent_purchase_cost,
DROP COLUMN IF EXISTS electrity_uph,
DROP COLUMN IF EXISTS max_operation_value_unit,
DROP COLUMN IF EXISTS max_operation_value,
DROP COLUMN IF EXISTS min_operation_value_unit,
DROP COLUMN IF EXISTS min_operation_value;

-- Remove type_id column from machines table
ALTER TABLE machines 
DROP COLUMN IF EXISTS type_id;
