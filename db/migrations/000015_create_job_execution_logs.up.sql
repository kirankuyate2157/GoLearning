CREATE TABLE IF NOT EXISTS job_execution_logs (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    machine_id INTEGER REFERENCES machines(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    status VARCHAR(50),
    operator_name VARCHAR(150),
    note TEXT,
    logged_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_execution_logs_job_id ON job_execution_logs(job_id);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_machine_id ON job_execution_logs(machine_id);
