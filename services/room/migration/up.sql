-- Ensure a clean slate by dropping the tables if they already exist.
DROP TABLE IF EXISTS "scheduled"; -- Removes the 'scheduled' table to avoid conflicts during creation.
DROP TABLE IF EXISTS "room_expiry"; -- Removes the 'room_expiry' table to avoid conflicts during creation.
DROP TABLE IF EXISTS "room_result"; -- Removes the 'room_result' table to avoid conflicts during creation.
DROP TABLE IF EXISTS "room_log"; -- Removes the 'room_log' table to avoid conflicts during creation, ensuring no leftover data affects the new schema.
DROP TABLE IF EXISTS "rooms"; -- Removes the 'rooms' table to avoid conflicts during creation. This is done last to maintain referential integrity for tables with foreign keys.
DROP TABLE IF EXISTS "bans";
-- DROP TABLE IF EXISTS "archived_rooms"; -- Option to remove the 'archived_rooms' table to avoid conflicts during creation. This is commented out as an optional step.

-- Define the "rooms" table to store detailed information about individual chat or meeting rooms.
CREATE TABLE "rooms"(
    room_id VARCHAR(14), -- Unique identifier for each room, with a maximum length of 14 characters.
    creator TEXT NOT NULL, -- ID or name of the room's creator, ensuring every room has an associated creator.
    users_length INTEGER NOT NULL, -- Tracks the current number of users in the room, ensuring this value is always provided.
    closed BOOLEAN NOT NULL DEFAULT false, -- Flags whether the room is closed to new entries, defaults to open.
    closed_at TIMESTAMP, -- Records the exact time the room was closed, nullable if the room remains open.
    created_at TIMESTAMP NOT NULL, -- Captures the creation time of the room, essential for tracking room age and cleanup.
    PRIMARY KEY (room_id) -- Ensures room_id is unique and serves as the primary key for identification.
);

-- Define the "scheduled" table to manage scheduling for events in rooms.
CREATE TABLE "scheduled"(
    room_id VARCHAR(14) PRIMARY KEY REFERENCES "rooms" (room_id), -- Links to the "rooms" table, ensuring scheduled events are valid rooms.
    starts_at TIMESTAMP NOT NULL -- Specifies the start time of the scheduled event, mandatory for event timing.
);

-- Define the "room_expiry" table for auto-expiration of rooms.
CREATE TABLE "room_expiry"(
    room_id VARCHAR(14) PRIMARY KEY REFERENCES "rooms" (room_id), -- Links to the "rooms" table, marking rooms for automatic expiration.
    ends_at TIMESTAMP NOT NULL -- Determines when the room is considered expired and potentially auto-deleted or archived.
);

-- Define the "room_result" table to track decisions made in rooms.
CREATE TABLE "room_result"(
    room_id VARCHAR(14) PRIMARY KEY REFERENCES "rooms" (room_id), -- Links to the "rooms" table, associating decisions with specific rooms.
    approvers TEXT[], -- Array of text strings, each representing an approver's ID or name, to track approvals within the room.
    deniers TEXT[],    -- Array of text strings, each representing a denier's ID or name, to track denials within the room.
    created_at TIMESTAMP NOT NULL -- Captures the creation time of the room result, essential for tracking decisions' timing.
);

-- Define the "room_log" table to keep track of user events in rooms.
CREATE TABLE "room_log"(
    room_id VARCHAR(14) REFERENCES "rooms" (room_id), -- Links to the "rooms" table, ensuring logs are associated with valid rooms.
    user_id TEXT NOT NULL, -- ID of the user related to the log event.
    user_event TEXT NOT NULL, -- Describes the event type (e.g., 'joined', 'left').
    created_at TIMESTAMP NOT NULL, -- Timestamp when the event occurred.
    PRIMARY KEY (room_id, created_at)
);

CREATE TABLE "bans"(
    user_id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL
);

-- Define the "archived_rooms" table for historical record keeping of rooms.
-- CREATE TABLE "archived_rooms"(
--     room_id VARCHAR(14), -- Ensures reference to unique room identifiers, linking to the "rooms" table.
--     archived_at TIMESTAMP NOT NULL, -- Records the time when the room was moved to the archive.
--     creator TEXT NOT NULL, -- Maintains record of the room's original creator in the archive.
--     users_length INTEGER NOT NULL, -- Preserves the final count of users in the room upon archiving.
--     closed BOOLEAN NOT NULL DEFAULT false, -- Indicates the room's status at the time of archiving.
--     created_at TIMESTAMP NOT NULL, -- Keeps track of the original creation timestamp for historical reference.
--     PRIMARY KEY (room_id, archived_at) -- Composite key ensuring uniqueness for each archived instance based on room and time.
-- );
