<?php

namespace App\Controller;

use App\Service\CreateOrderRequest;
use App\Service\OrderService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;

class OrderController extends AbstractController
{
    public function __construct(
        private readonly OrderService $orderService
    ) {}

    #[Route('/shops/{shopId}/orders', name: 'create_order', methods: ['POST'])]
    public function createOrder(Request $request, int $shopId): JsonResponse
    {
        $data = json_decode($request->getContent(), true);

        if (json_last_error() !== JSON_ERROR_NONE) {
            return $this->json(['error' => 'Invalid JSON'], 400);
        }

        $createOrderRequest = new CreateOrderRequest(
            number: $data['number'] ?? '',
            total: (float) ($data['total'] ?? 0),
            customerName: $data['customerName'] ?? ''
        );

        try {
            $response = $this->orderService->createOrder($shopId, $createOrderRequest);
            return $this->json($response->toArray());
        } catch (\Exception $e) {
            return $this->json(['error' => 'Internal server error'], 500);
        }
    }
}
