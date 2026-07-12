-- Down migration for: create_workstation_vendors_and_related_tables
-- Drops tables created in the corresponding up migration in reverse order.
-- Uses IF EXISTS so the script is idempotent and safe to run multiple times.

-- DROP dependent tables first
DROP TABLE IF EXISTS vendor_contacts CASCADE;
DROP TABLE IF EXISTS vendor_types CASCADE;
DROP TABLE IF EXISTS vendors CASCADE;
DROP TABLE IF EXISTS workstation CASCADE;

-- DROP supporting tables
DROP TABLE IF EXISTS tax CASCADE;
DROP TABLE IF EXISTS warranty CASCADE;
DROP TABLE IF EXISTS vendor_contacts CASCADE; -- defensive: already dropped above
DROP TABLE IF EXISTS vendor_types CASCADE; -- defensive: already dropped above
DROP TABLE IF EXISTS person CASCADE;
DROP TABLE IF EXISTS address CASCADE;
DROP TABLE IF EXISTS types CASCADE;

-- Note: SERIAL and ENUM types (if created separately) are dropped automatically with CASCADE
-- Explicitly drop the enum type used by `types.category`
DROP TYPE IF EXISTS types_category_enum CASCADE;
