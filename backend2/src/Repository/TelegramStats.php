<?php

namespace App\Repository;

class TelegramStats
{
    public function __construct(
        public readonly int $sentCount,
        public readonly int $failedCount,
        public readonly ?\DateTimeInterface $lastSentAt
    ) {}
}
