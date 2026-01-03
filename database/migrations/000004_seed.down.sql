DELETE FROM tasks
WHERE title IN ('Setup project', 'Create migrations', 'Implement API');

DELETE FROM projects
WHERE name = 'Demo Project';

DELETE FROM users
WHERE email IN ('admin@demo.com', 'user@demo.com');
