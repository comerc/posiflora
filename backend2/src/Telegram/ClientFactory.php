<?php

namespace App\Telegram;

use Psr\Log\LoggerInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;

class ClientFactory
{
    public function __construct(
        private readonly LoggerInterface $logger,
        private readonly HttpClientInterface $httpClient,
        private readonly bool $mockMode = true
    ) {}

    public function create(): ClientInterface
    {
        if ($this->mockMode) {
            return new MockClient($this->logger);
        }

        return new RealClient($this->logger, $this->httpClient);
    }
}
