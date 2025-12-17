<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Seed data migration
 */
final class Version20251217153502 extends AbstractMigration
{
    public function getDescription(): string
    {
        return 'Seed test data';
    }

    public function up(Schema $schema): void
    {
        // Создаем тестовый магазин
        $this->addSql("INSERT OR IGNORE INTO shops (id, name, created_at, updated_at) VALUES (1, 'Тестовый магазин Posiflora', datetime('now'), datetime('now'))");

        // Создаем тестовые заказы
        $this->addSql("INSERT OR IGNORE INTO orders (shop_id, number, total, customer_name, created_at) VALUES
            (1, 'A-1001', 1500.00, 'Иван Петров', datetime('now', '-7 days')),
            (1, 'A-1002', 2300.50, 'Мария Сидорова', datetime('now', '-6 days')),
            (1, 'A-1003', 890.00, 'Алексей Иванов', datetime('now', '-5 days')),
            (1, 'A-1004', 3200.75, 'Елена Козлова', datetime('now', '-4 days')),
            (1, 'A-1005', 2490.00, 'Анна Смирнова', datetime('now', '-3 days')),
            (1, 'A-1006', 1750.25, 'Дмитрий Волков', datetime('now', '-2 days')),
            (1, 'A-1007', 2100.00, 'Ольга Новикова', datetime('now', '-1 day')),
            (1, 'A-1008', 980.50, 'Сергей Морозов', datetime('now'))");
    }

    public function down(Schema $schema): void
    {
        $this->addSql('DELETE FROM orders WHERE shop_id = 1');
        $this->addSql('DELETE FROM shops WHERE id = 1');
    }
}
