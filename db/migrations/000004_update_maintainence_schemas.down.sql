-- Remove remark column from maintenance_checklists
ALTER TABLE maintenance_checklists DROP COLUMN remark;


-- Remove status, started_on, ended_on columns from maintenance_schedules
ALTER TABLE maintenance_schedules 
DROP COLUMN status,
DROP COLUMN started_on,
DROP COLUMN ended_on;


-- Add back machine_id column to maintenance_templates
ALTER TABLE maintenance_templates 
ADD COLUMN machine_id INTEGER REFERENCES machines(id);



-- Drop the mapping table
DROP TABLE maintenance_template_machines;
