<?php

namespace App\Service;

use App\Entity\Order;

class CreateOrderResponse
{
    public function __construct(
        public readonly Order $order,
        public readonly string $sendStatus // sent, failed, skipped
    ) {}

    public function toArray(): array
    {
        return [
            'order' => [
                'id' => $this->order->getId(),
                'shop_id' => $this->order->getShopId(),
                'number' => $this->order->getNumber(),
                'total' => $this->order->getTotal(),
                'customer_name' => $this->order->getCustomerName(),
                'created_at' => $this->order->getCreatedAt()?->format('c'),
            ],
            'send_status' => $this->sendStatus,
        ];
    }
}
