<?php

namespace App\Service;

use App\Entity\Order;
use App\Entity\TelegramIntegration;
use App\Entity\TelegramSendLog;
use App\Repository\OrderRepository;
use App\Repository\ShopRepository;
use App\Repository\TelegramIntegrationRepository;
use App\Repository\TelegramSendLogRepository;
use App\Telegram\ClientInterface;
use Doctrine\ORM\EntityManagerInterface;

class OrderService
{
    public function __construct(
        private readonly OrderRepository $orderRepo,
        private readonly TelegramIntegrationRepository $integrationRepo,
        private readonly TelegramSendLogRepository $telegramSendLogRepo,
        private readonly ClientInterface $telegramClient,
        private readonly ShopRepository $shopRepo,
        private readonly EntityManagerInterface $entityManager
    ) {}

    public function createOrder(int $shopId, CreateOrderRequest $request): CreateOrderResponse
    {
        // Создаем магазин, если его нет
        $this->shopRepo->getOrCreate($shopId);

        // 1. Проверяем, существует ли уже заказ с таким номером для этого магазина
        $existingOrder = $this->orderRepo->findOneBy(['shopId' => $shopId, 'number' => $request->number]);
        if ($existingOrder !== null) {
            // Заказ уже существует, проверяем идемпотентность отправки
            $existingLog = $this->telegramSendLogRepo->getByShopIDAndOrderID($shopId, $existingOrder->getId());
            if ($existingLog !== null) {
                // Уже отправляли, пропускаем
                return new CreateOrderResponse($existingOrder, 'skipped');
            }
            // Заказ существует, но отправки не было - используем существующий заказ
            $order = $existingOrder;
        } else {
            // Создаем новый заказ
            $order = new Order();
            $order->setShopId($shopId);
            $order->setNumber($request->number);
            $order->setTotal($request->total);
            $order->setCustomerName($request->customerName);
            $order->setCreatedAt(new \DateTime());

            $this->entityManager->persist($order);
            $this->entityManager->flush(); // Flush to get the order ID
        }

        // 2. Проверяем идемпотентность отправки (для нового или существующего заказа)
        $existingLog = $this->telegramSendLogRepo->getByShopIDAndOrderID($shopId, $order->getId());
        if ($existingLog !== null) {
            // Уже отправляли, пропускаем
            return new CreateOrderResponse($order, 'skipped');
        }

        // 3. Проверяем наличие активной Telegram-интеграции
        $integration = $this->integrationRepo->findOneBy(['shopId' => $shopId]);

        // 4. Если интеграция включена - отправляем уведомление
        if ($integration !== null && $integration->isEnabled()) {
            $message = sprintf(
                'Новый заказ %s на сумму %.2f ₽, клиент %s',
                $order->getNumber(),
                $order->getTotal(),
                $order->getCustomerName()
            );

            $sendLog = new TelegramSendLog();
            $sendLog->setShopId($shopId);
            $sendLog->setOrderId($order->getId());
            $sendLog->setMessage($message);
            $sendLog->setStatus(TelegramSendLog::STATUS_SENT);

            $success = $this->telegramClient->sendMessage(
                $integration->getBotToken(),
                $integration->getChatId(),
                $message
            );

            if (!$success) {
                $sendLog->setStatus(TelegramSendLog::STATUS_FAILED);
                $sendLog->setError('Telegram API error');
            }

            $this->entityManager->persist($sendLog);
            $this->entityManager->flush();

            $status = $sendLog->getStatus() === TelegramSendLog::STATUS_SENT ? 'sent' : 'failed';

            return new CreateOrderResponse($order, $status);
        }

        // Интеграция не настроена или отключена
        return new CreateOrderResponse($order, 'skipped');
    }
}
