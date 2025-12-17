<?php

namespace App\Telegram;

use Psr\Log\LoggerInterface;

class MockClient implements ClientInterface
{
    public function __construct(
        private readonly LoggerInterface $logger,
        private bool $enabled = true
    ) {}

    public function sendMessage(string $token, string $chatId, string $message): bool
    {
        if (!$this->enabled) {
            $this->logger->error('Mock client is disabled');
            return false;
        }

        // Логируем вместо реальной отправки
        $this->logger->info('[MOCK TELEGRAM] Sending message', [
            'chat_id' => $chatId,
            'message' => $message,
            'token' => substr($token, 0, 10) . '...',
        ]);

        return true;
    }

    public function setEnabled(bool $enabled): void
    {
        $this->enabled = $enabled;
    }
}
