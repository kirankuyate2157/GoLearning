-- Remove answer_id column from maintenance_schedule_questions (drops FK first)
ALTER TABLE maintenance_schedule_questions DROP COLUMN IF EXISTS answer_id;

-- Drop the answer list table
DROP INDEX IF EXISTS idx_question_answers_question_id;
DROP TABLE IF EXISTS maintenance_question_answers;

-- Rename the foreign key column back to checklist_id
ALTER TABLE maintenance_schedule_questions RENAME COLUMN question_id TO checklist_id;

-- Rename maintenance_schedule_questions table back to maintenance_checklists
ALTER TABLE maintenance_schedule_questions RENAME TO maintenance_checklists;

-- Rename maintenance_procedure_steps_questions table back to maintenance_procedure_steps_checklist
ALTER TABLE maintenance_procedure_steps_questions RENAME TO maintenance_procedure_steps_checklist;

-- Remove display_order column from maintenance_procedure_steps_checklist
ALTER TABLE maintenance_procedure_steps_checklist DROP COLUMN IF EXISTS display_order;
