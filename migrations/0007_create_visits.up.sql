CREATE TABLE IF NOT EXISTS visits (
    id VARCHAR(50) PRIMARY KEY,
    appointment_id VARCHAR(50) NOT NULL REFERENCES appointments(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    worker_id VARCHAR(50) NOT NULL REFERENCES workers(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress',
    notes TEXT,
    deleted_reason VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_visits_appointment_id ON visits(appointment_id);
CREATE INDEX IF NOT EXISTS idx_visits_customer_id ON visits(customer_id);
CREATE INDEX IF NOT EXISTS idx_visits_worker_id ON visits(worker_id);
CREATE INDEX IF NOT EXISTS idx_visits_created_at ON visits(created_at DESC);

CREATE TABLE IF NOT EXISTS visit_tasks (
    visit_id VARCHAR(50) NOT NULL REFERENCES visits(id) ON DELETE CASCADE,
    task_id VARCHAR(50) NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    notes TEXT,
    PRIMARY KEY (visit_id, task_id)
);

CREATE INDEX IF NOT EXISTS idx_visit_tasks_visit_id ON visit_tasks(visit_id);
CREATE INDEX IF NOT EXISTS idx_visit_tasks_task_id ON visit_tasks(task_id);
