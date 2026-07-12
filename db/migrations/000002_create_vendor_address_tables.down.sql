-- Revert vendor_addresses mapping: recreate columns on vendors and drop mapping table

-- Add columns back to vendors (nullable to avoid data issues)
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS address_id INTEGER REFERENCES address(id);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS contact_person_id INTEGER REFERENCES person(id);

-- Recreate indexes on vendors
CREATE INDEX IF NOT EXISTS idx_vendors_address ON vendors(address_id);
CREATE INDEX IF NOT EXISTS idx_vendors_contact_person ON vendors(contact_person_id);

-- Drop mapping table
DROP TABLE IF EXISTS vendor_addresses CASCADE;
