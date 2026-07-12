-- Create the new mapping table
CREATE TABLE maintenance_template_machines (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    maintenance_template_id INTEGER NOT NULL REFERENCES maintenance_templates(id) ON DELETE CASCADE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);


-- Remove the column from maintenance_templates
ALTER TABLE maintenance_templates DROP COLUMN machine_id;


ALTER TABLE maintenance_schedules 
ADD COLUMN status TEXT,
ADD COLUMN started_on TIMESTAMPTZ,
ADD COLUMN ended_on TIMESTAMPTZ;


ALTER TABLE maintenance_checklists 
ADD COLUMN remark TEXT;