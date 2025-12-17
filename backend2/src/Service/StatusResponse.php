<?php

namespace App\Service;

class StatusResponse
{
    public function __construct(
        public readonly bool $enabled,
        public readonly string $chat_id,
        public readonly ?string $last_sent_at,
        public readonly int $sent_count,
        public readonly int $failed_count
    ) {}

    public function toArray(): array
    {
        return [
            'enabled' => $this->enabled,
            'chat_id' => $this->chat_id,
            'last_sent_at' => $this->last_sent_at,
            'sent_count' => $this->sent_count,
            'failed_count' => $this->failed_count,
        ];
    }
}
