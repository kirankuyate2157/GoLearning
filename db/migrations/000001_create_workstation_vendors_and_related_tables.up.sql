



DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'types_category_enum') THEN
        CREATE TYPE types_category_enum AS ENUM ('vendor', 'workstation', 'sparepart', 'unit', 'feed' , 'warranty','maintenancefrequency','feedsdisplayunit');
    END IF;
END
$$ LANGUAGE plpgsql;

CREATE TABLE types (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    category types_category_enum,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE warranty (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    provider TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE person (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    email text UNIQUE,
    mobile_number VARCHAR(15),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE tax (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    type TEXT NOT NULL,
    percent NUMERIC(5,2) CHECK (percent >= 0 AND percent <= 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE address (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    address_line_1 TEXT NOT NULL,
    address_line_2 TEXT,
    address_line_3 TEXT,
    landmark TEXT,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL,
    pin_code VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name TEXT NOT NULL,
    website TEXT,
    email TEXT,
    phone_no VARCHAR(20),
    join_date DATE,
    gst_no VARCHAR(20) UNIQUE,
    gst_type INTEGER REFERENCES tax(id),
    address_id INTEGER REFERENCES address(id),
    contact_person_id INTEGER REFERENCES person(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE TABLE vendor_types (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE,
    type_id INTEGER REFERENCES types(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);
ALTER TABLE vendor_types ADD CONSTRAINT vendor_types_vendor_type_unique UNIQUE (vendor_id, type_id);





CREATE TABLE vendor_contacts (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE,
    contact_id INTEGER REFERENCES person(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);
ALTER TABLE vendor_contacts ADD CONSTRAINT vendor_contacts_vendor_contact_unique UNIQUE (vendor_id, contact_id);





CREATE TABLE workstation (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    code TEXT UNIQUE,
    type_id INTEGER REFERENCES types(id),
    manufacturer_id INTEGER REFERENCES vendors(id),
    warranty_id INTEGER REFERENCES warranty(id),
    model TEXT,
    area TEXT,
    installation_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER
);





CREATE INDEX idx_vendor_name ON vendors(name);
CREATE INDEX idx_vendor_gst_no ON vendors(gst_no);
CREATE INDEX idx_vendors_gst_type ON vendors(gst_type);
CREATE INDEX idx_vendors_address ON vendors(address_id);
CREATE INDEX idx_vendors_contact_person ON vendors(contact_person_id);


CREATE INDEX idx_vendor_types_vendor ON vendor_types(vendor_id);
CREATE INDEX idx_vendor_types_type ON vendor_types(type_id);


CREATE INDEX idx_vendor_contacts_vendor ON vendor_contacts(vendor_id);
CREATE INDEX idx_vendor_contacts_contact ON vendor_contacts(contact_id);


CREATE INDEX idx_workstation_name ON workstation(name);
CREATE INDEX idx_workstation_code ON workstation(code);
CREATE INDEX idx_workstation_type ON workstation(type_id);
CREATE INDEX idx_workstation_warranty ON workstation(warranty_id);
CREATE INDEX idx_workstation_manufacturer ON workstation(manufacturer_id);


