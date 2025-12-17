<?php

namespace App\Controller;

use App\Service\ConnectRequest;
use App\Service\StatusService;
use App\Service\TelegramIntegrationService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;

class TelegramController extends AbstractController
{
    public function __construct(
        private readonly TelegramIntegrationService $integrationService,
        private readonly StatusService $statusService
    ) {}

    #[Route('/shops/{shopId}/telegram/connect', name: 'telegram_connect', methods: ['POST'])]
    public function connect(Request $request, int $shopId): JsonResponse
    {
        $data = json_decode($request->getContent(), true);

        if (json_last_error() !== JSON_ERROR_NONE) {
            return $this->json(['error' => 'Invalid JSON'], 400);
        }

        $connectRequest = new ConnectRequest(
            botToken: $data['botToken'] ?? '',
            chatId: $data['chatId'] ?? '',
            enabled: $data['enabled'] ?? false
        );

        try {
            $integration = $this->integrationService->connect($shopId, $connectRequest);
            return $this->json([
                'id' => $integration->getId(),
                'shop_id' => $integration->getShopId(),
                'bot_token' => $integration->getBotToken(),
                'chat_id' => $integration->getChatId(),
                'enabled' => $integration->isEnabled(),
                'created_at' => $integration->getCreatedAt()?->format('c'),
                'updated_at' => $integration->getUpdatedAt()?->format('c'),
            ]);
        } catch (\InvalidArgumentException $e) {
            return $this->json(['error' => $e->getMessage()], 400);
        } catch (\Exception $e) {
            return $this->json(['error' => 'Internal server error'], 500);
        }
    }

    #[Route('/shops/{shopId}/telegram/status', name: 'telegram_status', methods: ['GET'])]
    public function getStatus(Request $request, int $shopId): JsonResponse
    {
        // Проверяем параметр mask=disabled для отключения маскирования
        $maskDisabled = $request->query->get('mask') === 'disabled';

        try {
            $status = $this->statusService->getStatus($shopId, $maskDisabled);
            return $this->json($status->toArray());
        } catch (\Exception $e) {
            return $this->json(['error' => 'Internal server error'], 500);
        }
    }
}
