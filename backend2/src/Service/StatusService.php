<?php

namespace App\Service;

use App\Repository\TelegramIntegrationRepository;
use App\Repository\TelegramSendLogRepository;

class StatusService
{
    public function __construct(
        private readonly TelegramIntegrationRepository $integrationRepo,
        private readonly TelegramSendLogRepository $telegramSendLogRepo
    ) {}

    public function getStatus(int $shopId, bool $maskDisabled = false): StatusResponse
    {
        $integration = $this->integrationRepo->findOneBy(['shopId' => $shopId]);
        $stats = $this->telegramSendLogRepo->getStatsForLast7Days($shopId);

        $response = new StatusResponse(
            enabled: $integration?->isEnabled() ?? false,
            chat_id: $this->maskChatId($integration?->getChatId(), $maskDisabled),
            last_sent_at: $stats->lastSentAt?->format('c'),
            sent_count: $stats->sentCount,
            failed_count: $stats->failedCount
        );

        return $response;
    }

    private function maskChatId(?string $chatId, bool $maskDisabled): string
    {
        if ($chatId === null) {
            return '';
        }

        // Если maskDisabled=true, то маскируем
        // Если maskDisabled=false, то показываем полный Chat ID
        if (!$maskDisabled) {
            return $chatId;
        }

        // Маскируем chatId (показываем первые 3 и последние 3 символа)
        if (strlen($chatId) > 6) {
            return substr($chatId, 0, 3) . str_repeat('*', strlen($chatId) - 6) . substr($chatId, -3);
        }

        return str_repeat('*', strlen($chatId));
    }
}
