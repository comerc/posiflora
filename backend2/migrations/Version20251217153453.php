<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20251217153453 extends AbstractMigration
{
    public function getDescription(): string
    {
        return '';
    }

    public function up(Schema $schema): void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->addSql('CREATE TABLE orders (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, shop_id INTEGER NOT NULL, number VARCHAR(255) NOT NULL, total DOUBLE PRECISION NOT NULL, customer_name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL)');
        $this->addSql('CREATE TABLE shops (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL)');
        $this->addSql('CREATE TABLE telegram_integrations (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, shop_id INTEGER NOT NULL, bot_token VARCHAR(255) NOT NULL, chat_id VARCHAR(255) NOT NULL, enabled BOOLEAN NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL)');
        $this->addSql('CREATE UNIQUE INDEX shop_id_unique ON telegram_integrations (shop_id)');
        $this->addSql('CREATE TABLE telegram_send_log (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, shop_id INTEGER NOT NULL, order_id INTEGER NOT NULL, message CLOB NOT NULL, status VARCHAR(10) NOT NULL, error CLOB DEFAULT NULL, sent_at DATETIME NOT NULL)');
    }

    public function down(Schema $schema): void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->addSql('DROP TABLE orders');
        $this->addSql('DROP TABLE shops');
        $this->addSql('DROP TABLE telegram_integrations');
        $this->addSql('DROP TABLE telegram_send_log');
    }
}
