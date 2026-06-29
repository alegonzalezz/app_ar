DO $$
BEGIN
    -- Drop FK constraints before altering types
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'appointment_tasks_appointment_id_fkey') THEN
        ALTER TABLE appointment_tasks DROP CONSTRAINT appointment_tasks_appointment_id_fkey;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'appointment_tasks_task_id_fkey') THEN
        ALTER TABLE appointment_tasks DROP CONSTRAINT appointment_tasks_task_id_fkey;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'visit_tasks_visit_id_fkey') THEN
        ALTER TABLE visit_tasks DROP CONSTRAINT visit_tasks_visit_id_fkey;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'visit_tasks_task_id_fkey') THEN
        ALTER TABLE visit_tasks DROP CONSTRAINT visit_tasks_task_id_fkey;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'visits_worker_id_fkey') THEN
        ALTER TABLE visits DROP CONSTRAINT visits_worker_id_fkey;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'visits_appointment_id_fkey') THEN
        ALTER TABLE visits DROP CONSTRAINT visits_appointment_id_fkey;
    END IF;
END $$;

-- Convert PRIMARY KEY columns from VARCHAR(50) to UUID
ALTER TABLE greetings ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE users ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE workers ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE tasks ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE appointments ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE visits ALTER COLUMN id TYPE UUID USING id::uuid;

-- Convert FK columns from VARCHAR(50) to UUID
ALTER TABLE tasks ALTER COLUMN customer_id TYPE UUID USING customer_id::uuid;
ALTER TABLE tasks ALTER COLUMN worker_id TYPE UUID USING worker_id::uuid;
ALTER TABLE appointments ALTER COLUMN customer_id TYPE UUID USING customer_id::uuid;
ALTER TABLE appointments ALTER COLUMN worker_id TYPE UUID USING worker_id::uuid;
ALTER TABLE appointment_tasks ALTER COLUMN appointment_id TYPE UUID USING appointment_id::uuid;
ALTER TABLE appointment_tasks ALTER COLUMN task_id TYPE UUID USING task_id::uuid;
ALTER TABLE visits ALTER COLUMN appointment_id TYPE UUID USING appointment_id::uuid;
ALTER TABLE visits ALTER COLUMN worker_id TYPE UUID USING worker_id::uuid;
ALTER TABLE visit_tasks ALTER COLUMN visit_id TYPE UUID USING visit_id::uuid;
ALTER TABLE visit_tasks ALTER COLUMN task_id TYPE UUID USING task_id::uuid;

-- Re-add FK constraints
ALTER TABLE visits ADD CONSTRAINT visits_worker_id_fkey FOREIGN KEY (worker_id) REFERENCES workers(id) ON DELETE RESTRICT;
ALTER TABLE visits ADD CONSTRAINT visits_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE RESTRICT;
ALTER TABLE appointment_tasks ADD CONSTRAINT appointment_tasks_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE;
ALTER TABLE appointment_tasks ADD CONSTRAINT appointment_tasks_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;
ALTER TABLE visit_tasks ADD CONSTRAINT visit_tasks_visit_id_fkey FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE;
ALTER TABLE visit_tasks ADD CONSTRAINT visit_tasks_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- Restructure auth_users: add id, profile_id, profile_type
ALTER TABLE auth_users ADD COLUMN id UUID;
ALTER TABLE auth_users ADD COLUMN profile_id UUID;
ALTER TABLE auth_users ADD COLUMN profile_type VARCHAR(50);

-- Migrate existing data: all current auth_users are 'user' type
UPDATE auth_users SET
    id = gen_random_uuid(),
    profile_id = user_id::uuid,
    profile_type = 'user';

-- Make new columns NOT NULL
ALTER TABLE auth_users ALTER COLUMN id SET NOT NULL;
ALTER TABLE auth_users ALTER COLUMN profile_id SET NOT NULL;
ALTER TABLE auth_users ALTER COLUMN profile_type SET NOT NULL;

-- Drop old PK and column
ALTER TABLE auth_users DROP CONSTRAINT auth_users_pkey;
ALTER TABLE auth_users DROP COLUMN user_id;

-- Set new PK
ALTER TABLE auth_users ADD PRIMARY KEY (id);

-- Create index for profile lookups
CREATE INDEX IF NOT EXISTS idx_auth_users_profile ON auth_users(profile_id, profile_type);
