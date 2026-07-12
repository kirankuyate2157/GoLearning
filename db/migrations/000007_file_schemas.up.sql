CREATE TABLE
    files (
        id BIGSERIAL PRIMARY KEY,
        uid UUID NOT NULL DEFAULT gen_random_uuid(),
        original_name TEXT NOT NULL,
        stored_name TEXT NOT NULL,
        extension TEXT NOT NULL,
        mime_type TEXT NOT NULL,
        size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
        storage_provider TEXT NOT NULL, -- local | s3 | gcs | r2
        storage_bucket TEXT,
        storage_path TEXT NOT NULL,
        checksum_sha256 TEXT,
        uploaded_by BIGINT NOT NULL,
        uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (uid)
    );

CREATE TABLE
    entity_files (
        id BIGSERIAL PRIMARY KEY,
        file_id BIGINT NOT NULL,
        entity_type TEXT NOT NULL, -- machine, spare_part, vendor
        entity_id BIGINT NOT NULL, -- machine_id, spare_parts_id, vendor_id
        usage_type TEXT NOT NULL, -- thumbnail, attachment, report, image, diagram, document
        display_order INT DEFAULT 0,
        is_primary BOOLEAN DEFAULT FALSE,
        linked_by BIGINT NOT NULL,
        linked_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (file_id, entity_type, entity_id),
        FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
    );

-- Create indexes for better query performance
CREATE INDEX idx_files_uid ON files (uid);
CREATE INDEX idx_entity_files_entity ON entity_files (entity_type, entity_id);
CREATE INDEX idx_entity_files_file_id ON entity_files (file_id);