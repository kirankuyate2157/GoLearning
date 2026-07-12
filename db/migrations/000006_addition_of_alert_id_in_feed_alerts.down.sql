
-- Re-add alert_type column (adjust datatype if needed)
ALTER TABLE feed_alerts
ADD COLUMN alert_type VARCHAR(255);

-- Drop foreign key constraint
ALTER TABLE feed_alerts
DROP CONSTRAINT fk_feed_alerts_alerts;

-- Drop alert_id column
ALTER TABLE feed_alerts
DROP COLUMN alert_id;

-- Drop foreign key constraint
ALTER TABLE feed_alerts
DROP CONSTRAINT fk_feed_alerts_feeds;

-- Remove feed_id column
ALTER TABLE feed_alerts
DROP COLUMN feed_id;
