<?php

namespace App\Service;

use App\Entity\TelegramIntegration;
use App\Repository\ShopRepository;
use App\Repository\TelegramIntegrationRepository;
use Doctrine\ORM\EntityManagerInterface;

class TelegramIntegrationService
{
    public function __construct(
        private readonly TelegramIntegrationRepository $integrationRepo,
        private readonly ShopRepository $shopRepo,
        private readonly EntityManagerInterface $entityManager
    ) {}

    public function connect(int $shopId, ConnectRequest $request): TelegramIntegration
    {
        if (empty($request->botToken) || empty($request->chatId)) {
            throw new \InvalidArgumentException('botToken and chatId are required');
        }

        // Создаем магазин, если его нет
        $this->shopRepo->getOrCreate($shopId);

        $now = new \DateTime();
        $integration = new TelegramIntegration();
        $integration->setShopId($shopId);
        $integration->setBotToken($request->botToken);
        $integration->setChatId($request->chatId);
        $integration->setEnabled($request->enabled);
        $integration->setUpdatedAt($now);

        $existing = $this->integrationRepo->findOneBy(['shopId' => $shopId]);
        if ($existing !== null) {
            $integration = $existing;
            $integration->setBotToken($request->botToken);
            $integration->setChatId($request->chatId);
            $integration->setEnabled($request->enabled);
            $integration->setUpdatedAt($now);
        } else {
            $integration->setCreatedAt($now);
        }

        $this->entityManager->persist($integration);
        $this->entityManager->flush();

        return $integration;
    }

    public function getStatus(int $shopId): ?TelegramIntegration
    {
        return $this->integrationRepo->findOneBy(['shopId' => $shopId]);
    }
}
