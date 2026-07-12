CREATE TABLE maintenance_schedule_observations_answers (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    maintenance_schedule_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by INTEGER,
    CONSTRAINT uq_checked_checkbox
        UNIQUE (maintenance_schedule_id, question_id, answer_id)
);
