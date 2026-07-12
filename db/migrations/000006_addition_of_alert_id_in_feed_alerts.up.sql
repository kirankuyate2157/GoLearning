-- Add alert_id column to feed_alerts table
alter table feed_alerts
add column alert_id int default 0 ;

-- Add Foreign Key Constraint 
ALTER TABLE feed_alerts
ADD CONSTRAINT fk_feed_alerts_alerts
FOREIGN KEY (alert_id)
REFERENCES alerts (id);

-- Drop alert_type column from feed_alerts table
ALTER table feed_alerts
Drop column alert_type ;

-- Add feed_id column to feed_alerts table
alter table feed_alerts 
add column feed_id int null;

-- Add Foreign Key Constraint 
ALTER TABLE feed_alerts
ADD CONSTRAINT fk_feed_alerts_feeds
FOREIGN KEY (feed_id)
REFERENCES feeds (id);