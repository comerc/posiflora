<?php

namespace App\Service;

class CreateOrderRequest
{
    public function __construct(
        public readonly string $number,
        public readonly float $total,
        public readonly string $customerName
    ) {}
}
