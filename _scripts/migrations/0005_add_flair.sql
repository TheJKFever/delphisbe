/* Script to create the tables and constraints for Flair.
 */
CREATE TABLE IF NOT EXISTS flair_templates (
    id varchar(36) PRIMARY KEY,
    display_name varchar(128),
    image_url text,
    source varchar(128) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

/* Create the user/flair join table.
 *
 * user_id has a foreign key constraint to ensure the referenced user exists.
 * If the user is deleted, delete the user's available flair also.
 * flair_id has a foreign key constraint to ensure the referenced flair exists.
 * If the flair is deleted, delete the user's available flair also.
 */
CREATE TABLE IF NOT EXISTS flairs (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) REFERENCES users(id) ON DELETE CASCADE,
    template_id varchar(36) REFERENCES flair_templates(id) ON DELETE CASCADE,
    -- Uncomment this if you want to track when/if flair is verified
    -- verified_at timestamp with time zone
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

-- Add optional flair_id column and foreign key to participants table.
ALTER TABLE participants
ADD COLUMN flair_id varchar(36)
REFERENCES flairs(id) ON DELETE SET NULL;
