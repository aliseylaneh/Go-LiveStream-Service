-- Drop the table 'files' if it already exists to avoid conflicts
DROP TABLE IF EXISTS "files";

-- Create the 'files' table with specific column definitions
CREATE TABLE "files" (
    file_id VARCHAR(30) NOT NULL,    -- Unique identifier for the file
    file_type VARCHAR(20) NOT NULL,  -- Type of the file
    room_id VARCHAR(14) NOT NULL,    -- Identifier for the room associated with the file
    user_id TEXT,                    -- identifier for the user (assuming PostgreSQL syntax)
    created_at TIMESTAMP NOT NULL,   -- Timestamp indicating when the file was created
    PRIMARY KEY (file_id, file_type) -- Composite primary key composed of 'file_id' and 'file_type'
);
