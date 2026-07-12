-- Rename maintenance_procedure_steps_checklist table to maintenance_procedure_steps_questions
ALTER TABLE maintenance_procedure_steps_checklist RENAME TO maintenance_procedure_steps_questions;

-- Add display_order column to maintenance_procedure_steps_questions
ALTER TABLE maintenance_procedure_steps_questions ADD COLUMN display_order INTEGER;

-- Rename maintenance_checklists table to maintenance_schedule_questions
ALTER TABLE maintenance_checklists RENAME TO maintenance_schedule_questions;

-- Rename the foreign key column in maintenance_schedule_questions
ALTER TABLE maintenance_schedule_questions RENAME COLUMN checklist_id TO question_id;


-- Create the new answer list table
CREATE TABLE maintenance_question_answers (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    question_id INTEGER NOT NULL REFERENCES maintenance_procedure_steps_questions(id) ON DELETE CASCADE,
    answer_text TEXT NOT NULL,
    -- is_correct BOOLEAN DEFAULT FALSE,
    display_order INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INTEGER
);

-- Add answer_id column to maintenance_schedule_questions
ALTER TABLE maintenance_schedule_questions ADD COLUMN answer_id INTEGER REFERENCES maintenance_question_answers(id) ON DELETE SET NULL;


-- Add index for better query performance
CREATE INDEX idx_question_answers_question_id ON maintenance_question_answers(question_id);

UPDATE maintenance_procedure_steps_questions t
SET display_order = s.rn
FROM (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY step_id  ORDER BY id) AS rn
    FROM maintenance_procedure_steps_questions
) s
WHERE t.id = s.id;
