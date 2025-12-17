<?php

namespace App\Telegram;

interface ClientInterface
{
    public function sendMessage(string $token, string $chatId, string $message): bool;
}
