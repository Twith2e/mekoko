INSERT INTO users (public_id, first_name, last_name, email, password_hash, role)
VALUES (
    gen_random_uuid(),
    'Mekoko',
    'Admin',
    'hellomekoko@gmail.com',
    '$2a$10$hognf5T7kh5WDohK1w1vROX16ouTzioLNL2XJ7QkQeE5HiJ4VGgIq',
    'admin'
);
