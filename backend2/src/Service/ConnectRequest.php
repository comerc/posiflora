<?php

namespace App\Service;

class ConnectRequest
{
    public function __construct(
        public readonly string $botToken,
        public readonly string $chatId,
        public readonly bool $enabled
    ) {}
}
