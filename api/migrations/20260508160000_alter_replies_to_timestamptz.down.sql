ALTER TABLE replies
    ALTER COLUMN created_at TYPE TIMESTAMP
        USING created_at AT TIME ZONE current_setting('TimeZone'),
    ALTER COLUMN updated_at TYPE TIMESTAMP
        USING updated_at AT TIME ZONE current_setting('TimeZone');
