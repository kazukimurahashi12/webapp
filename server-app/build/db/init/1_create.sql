CREATE DATABASE IF NOT EXISTS user_info;
USE user_info;
INSERT INTO USERS (id, user_id, password)
SELECT 1, 'root', '$2a$10$m5QnQxuZYnwSQ/.937pl/ulKB7ljIviSrhu52P/eNctOhukfdKwz'
WHERE NOT EXISTS (
    SELECT 1 FROM USERS WHERE user_id = 'root'
);