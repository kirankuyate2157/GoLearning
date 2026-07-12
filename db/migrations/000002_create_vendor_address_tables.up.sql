-- Create mapping table vendor_addresses and remove address/contact refs from vendors

CREATE TABLE vendor_addresses (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid(),
    vendor_id INTEGER NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES address(id) ON DELETE CASCADE,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);

-- Drop vendor columns that are now modeled via vendor_addresses
ALTER TABLE vendors DROP COLUMN IF EXISTS address_id;
ALTER TABLE vendors DROP COLUMN IF EXISTS contact_person_id;

-- Drop related indexes if they exist
DROP INDEX IF EXISTS idx_vendors_address;
DROP INDEX IF EXISTS idx_vendors_contact_person;

CREATE INDEX idx_vendor_addresses_vendor ON vendor_addresses(vendor_id);
CREATE INDEX idx_vendor_addresses_address ON vendor_addresses(address_id);
