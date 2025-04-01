-- rooms
INSERT INTO rooms (name, type, price, tags, area, description, image_url) VALUES
    ('松湖智居 1 栋', 'b2l1', 1200.50, '海景, 阳台', 75.5, '宽敞明亮的海景房，带有独立阳台，适合度假和长期居住', '1.jpg'),
    ('松湖智居 2 栋', 'b1l1', 900.00, '市中心, 现代风格', 50.0, '位于市中心，装修现代，交通便利，适合商务人士', '2.jpg'),
    ('松湖智居 3 栋', 'b1',600.75, '温馨, 经济实惠', 30.2, '经济型公寓，空间紧凑但温馨，适合学生或短租客', '3.jpg');

-- billings
INSERT INTO billings (type, user_id, billing_id, price, paid, created_at, updated_at) VALUES
('rent-room', 1, 101, 1200.50, false, NOW(), NOW()),
('monthly-pay', 1, 102, 800.00, true, NOW(), NOW()),
('terminate-lease', 1, 103, 500.00, false, NOW(), NOW());

