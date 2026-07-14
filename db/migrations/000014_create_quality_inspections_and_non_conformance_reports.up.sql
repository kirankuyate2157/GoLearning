CREATE TABLE IF NOT EXISTS quality_inspections (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE SET NULL,
    wip_item_id INTEGER REFERENCES wip_items(id) ON DELETE SET NULL,
    inspector_name VARCHAR(150),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    inspected_at TIMESTAMP WITH TIME ZONE,
    findings TEXT,
    result VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS non_conformance_reports (
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    inspection_id INTEGER NOT NULL REFERENCES quality_inspections(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE SET NULL,
    wip_item_id INTEGER REFERENCES wip_items(id) ON DELETE SET NULL,
    type VARCHAR(100),
    severity VARCHAR(50),
    description TEXT,
    corrective_action TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_quality_inspections_job_id ON quality_inspections(job_id);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_wip_item_id ON quality_inspections(wip_item_id);
CREATE INDEX IF NOT EXISTS idx_non_conformance_reports_inspection_id ON non_conformance_reports(inspection_id);
CREATE INDEX IF NOT EXISTS idx_non_conformance_reports_job_id ON non_conformance_reports(job_id);
CREATE INDEX IF NOT EXISTS idx_non_conformance_reports_wip_item_id ON non_conformance_reports(wip_item_id);
