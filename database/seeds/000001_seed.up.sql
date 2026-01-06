WITH admin_user AS (
  INSERT INTO users (email, password_hash, role)
  VALUES ('admin@demo.com', '$2a$10$1blEC2WDt7LhFcVNLVHSv.IK4bs2Yeo2Lu9QbD0zK3/LWvbA0d2mm', 'admin')
  ON CONFLICT (email) DO NOTHING
  RETURNING id
),
normal_user AS (
  INSERT INTO users (email, password_hash, role)
  VALUES ('user@demo.com', '$2a$10$1blEC2WDt7LhFcVNLVHSv.IK4bs2Yeo2Lu9QbD0zK3/LWvbA0d2mm', 'user')
  ON CONFLICT (email) DO NOTHING
  RETURNING id
),
project AS (
  INSERT INTO projects (name, owner_id)
  SELECT
    'Demo Project',
    id
  FROM admin_user
  UNION
  SELECT
    'Demo Project',
    id
  FROM users
  WHERE email = 'admin@demo.com'
  LIMIT 1
  RETURNING id
)
INSERT INTO tasks (project_id, title, description, status, assignee_id)
SELECT
  project.id,
  t.title,
  t.description,
  t.status,
  u.id
FROM project
CROSS JOIN (
  VALUES
    ('Setup project', 'Initial setup of the project', 'done'),
    ('Create migrations', 'Add database migrations', 'in_progress'),
    ('Implement API', 'Start building the API', 'todo')
) AS t(title, description, status)
LEFT JOIN users u ON u.email = 'user@demo.com';
