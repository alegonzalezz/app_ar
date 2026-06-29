CREATE TABLE IF NOT EXISTS administratives (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50) NOT NULL,
    role VARCHAR(100) NOT NULL,
    collective_agreement VARCHAR(255),
    work_schedule VARCHAR(255) NOT NULL,
    hire_date DATE NOT NULL,
    salary NUMERIC(12,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
