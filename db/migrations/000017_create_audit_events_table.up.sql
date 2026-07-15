-- Audit events table to record critical mutations and access patterns

CREATE TABLE IF NOT EXISTS audit_events (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL DEFAULT gen_random_uuid(),
  user_uid TEXT NOT NULL,
  action TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_uid UUID,
  old_value TEXT,
  new_value TEXT,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_events_uid ON audit_events(uid);
CREATE INDEX IF NOT EXISTS idx_audit_events_user_uid ON audit_events(user_uid);
CREATE INDEX IF NOT EXISTS idx_audit_events_resource ON audit_events(resource_type, resource_uid);
CREATE INDEX IF NOT EXISTS idx_audit_events_created_at ON audit_events(created_at DESC);
