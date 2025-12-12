-- +goose Up
-- Создаем тестовый магазин
INSERT INTO shops (id, name) VALUES (1, 'Тестовый магазин Posiflora')
ON CONFLICT (id) DO NOTHING;

-- Создаем тестовые заказы
INSERT INTO orders (shop_id, number, total, customer_name) VALUES
    (1, 'A-1001', 1500.00, 'Иван Петров'),
    (1, 'A-1002', 2300.50, 'Мария Сидорова'),
    (1, 'A-1003', 890.00, 'Алексей Иванов'),
    (1, 'A-1004', 3200.75, 'Елена Козлова'),
    (1, 'A-1005', 2490.00, 'Анна Смирнова'),
    (1, 'A-1006', 1750.25, 'Дмитрий Волков'),
    (1, 'A-1007', 2100.00, 'Ольга Новикова'),
    (1, 'A-1008', 980.50, 'Сергей Морозов')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM orders WHERE shop_id = 1;
DELETE FROM shops WHERE id = 1;


