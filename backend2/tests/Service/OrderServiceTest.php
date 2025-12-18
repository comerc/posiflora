<?php

namespace App\Tests\Service;

use App\Entity\Order;
use App\Entity\TelegramIntegration;
use App\Entity\TelegramSendLog;
use App\Repository\ShopRepository;
use App\Repository\TelegramIntegrationRepository;
use App\Repository\TelegramSendLogRepository;
use App\Service\CreateOrderRequest;
use App\Service\CreateOrderResponse;
use App\Service\OrderService;
use App\Telegram\ClientInterface;
use Doctrine\ORM\EntityManagerInterface;
use PHPUnit\Framework\Attributes\AllowMockObjectsWithoutExpectations;
use PHPUnit\Framework\TestCase;

#[AllowMockObjectsWithoutExpectations]
class OrderServiceTest extends TestCase
{
    private TelegramIntegrationRepository $integrationRepo;
    private TelegramSendLogRepository $telegramSendLogRepo;
    private ClientInterface $telegramClient;
    private ShopRepository $shopRepo;
    private \App\Repository\OrderRepository $orderRepo;
    private EntityManagerInterface $entityManager;
    private OrderService $orderService;

    protected function setUp(): void
    {
        $this->integrationRepo = $this->createMock(TelegramIntegrationRepository::class);
        $this->telegramSendLogRepo = $this->createMock(TelegramSendLogRepository::class);
        $this->telegramClient = $this->createMock(ClientInterface::class);
        $this->shopRepo = $this->createMock(ShopRepository::class);
        $this->entityManager = $this->createMock(EntityManagerInterface::class);

        // Create a mock OrderRepository since we need it for the service
        $this->orderRepo = $this->createMock(\App\Repository\OrderRepository::class);

        $this->orderService = new OrderService(
            $this->orderRepo,
            $this->integrationRepo,
            $this->telegramSendLogRepo,
            $this->telegramClient,
            $this->shopRepo,
            $this->entityManager
        );
    }

    public function testCreateOrderWithEnabledIntegrationSendsMessage(): void
    {
        $shopId = 1;
        $request = new CreateOrderRequest('A-1005', 2490.0, 'Анна');

        $shop = $this->createMock(\App\Entity\Shop::class);
        $integration = $this->createMock(TelegramIntegration::class);

        $this->shopRepo->expects($this->once())
            ->method('getOrCreate')
            ->with($shopId)
            ->willReturn($shop);

        // Mock: заказ с таким номером еще не существует
        $this->orderRepo->expects($this->once())
            ->method('findOneBy')
            ->with(['shopId' => $shopId, 'number' => 'A-1005'])
            ->willReturn(null);

        $this->integrationRepo->expects($this->once())
            ->method('findOneBy')
            ->with(['shopId' => $shopId])
            ->willReturn($integration);

        // Mock: проверка идемпотентности не найдет существующий лог

        $integration->expects($this->once())
            ->method('isEnabled')
            ->willReturn(true);

        $integration->expects($this->once())
            ->method('getBotToken')
            ->willReturn('123456:ABC-DEF');

        $integration->expects($this->once())
            ->method('getChatId')
            ->willReturn('987654321');

        $this->telegramClient->expects($this->once())
            ->method('sendMessage')
            ->with('123456:ABC-DEF', '987654321', $this->callback(function($message) {
                return str_contains($message, 'Новый заказ A-1005');
            }))
            ->willReturn(true);

        // Mock EntityManager calls
        $this->entityManager->expects($this->exactly(2))
            ->method('persist')
            ->with($this->callback(function ($entity) {
                if ($entity instanceof \App\Entity\Order) {
                    // Присваиваем ID заказу при persist
                    $entity->setId(1);
                }
                return true;
            }));

        $this->entityManager->expects($this->exactly(2))
            ->method('flush');

        $response = $this->orderService->createOrder($shopId, $request);

        $this->assertInstanceOf(CreateOrderResponse::class, $response);
        $this->assertEquals('sent', $response->sendStatus);
    }

    public function testCreateOrderIdempotency(): void
    {
        $shopId = 1;
        $request = new CreateOrderRequest('A-1005', 2490.0, 'Анна');

        $shop = $this->createMock(\App\Entity\Shop::class);
        $existingOrder = $this->createMock(Order::class);
        $existingLog = $this->createMock(TelegramSendLog::class);

        $this->shopRepo->expects($this->once())
            ->method('getOrCreate')
            ->with($shopId)
            ->willReturn($shop);

        // Mock: заказ уже существует
        $this->orderRepo->expects($this->once())
            ->method('findOneBy')
            ->with(['shopId' => $shopId, 'number' => 'A-1005'])
            ->willReturn($existingOrder);

        // Заказ уже существует, проверяем идемпотентность отправки
        $existingOrder->expects($this->once())
            ->method('getId')
            ->willReturn(1);

        $this->telegramSendLogRepo->expects($this->once())
            ->method('getByShopIDAndOrderID')
            ->with($shopId, 1) // shopId и orderId
            ->willReturn($existingLog); // Return existing log

        $this->integrationRepo->expects($this->never())
            ->method('findOneBy'); // Не дойдет до проверки интеграции

        // Mock EntityManager calls - НЕ должны вызываться для нового заказа
        $this->entityManager->expects($this->never())
            ->method('persist');

        $this->entityManager->expects($this->never())
            ->method('flush');

        $response = $this->orderService->createOrder($shopId, $request);

        $this->assertInstanceOf(CreateOrderResponse::class, $response);
        $this->assertEquals('skipped', $response->sendStatus);
        $this->assertSame($existingOrder, $response->order);

        $this->telegramClient->expects($this->never())
            ->method('sendMessage');
    }

    public function testCreateOrderTelegramError(): void
    {
        $shopId = 1;
        $request = new CreateOrderRequest('A-1005', 2490.0, 'Анна');

        $shop = $this->createMock(\App\Entity\Shop::class);
        $integration = $this->createMock(TelegramIntegration::class);

        $this->shopRepo->expects($this->once())
            ->method('getOrCreate')
            ->with($shopId)
            ->willReturn($shop);

        // Mock: заказ с таким номером еще не существует
        $this->orderRepo->expects($this->once())
            ->method('findOneBy')
            ->with(['shopId' => $shopId, 'number' => 'A-1005'])
            ->willReturn(null);

        // Mock: проверка идемпотентности не найдет существующий лог
        $this->telegramSendLogRepo->expects($this->once())
            ->method('getByShopIDAndOrderID')
            ->with($shopId, 1) // shopId и orderId
            ->willReturn(null); // No existing log

        $this->integrationRepo->expects($this->once())
            ->method('findOneBy')
            ->with(['shopId' => $shopId])
            ->willReturn($integration);

        $integration->expects($this->once())
            ->method('isEnabled')
            ->willReturn(true);

        $integration->expects($this->once())
            ->method('getBotToken')
            ->willReturn('123456:ABC-DEF');

        $integration->expects($this->once())
            ->method('getChatId')
            ->willReturn('987654321');

        $this->telegramClient->expects($this->once())
            ->method('sendMessage')
            ->willReturn(false); // Telegram error

        // Mock EntityManager calls
        $this->entityManager->expects($this->exactly(2))
            ->method('persist')
            ->with($this->callback(function ($entity) {
                if ($entity instanceof \App\Entity\Order) {
                    // Присваиваем ID заказу при persist
                    $entity->setId(1);
                }
                return true;
            }));

        $this->entityManager->expects($this->exactly(2))
            ->method('flush');

        $response = $this->orderService->createOrder($shopId, $request);

        $this->assertInstanceOf(CreateOrderResponse::class, $response);
        $this->assertEquals('failed', $response->sendStatus);
    }
}
