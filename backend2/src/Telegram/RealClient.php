<?php

namespace App\Telegram;

use Psr\Log\LoggerInterface;

class RealClient implements ClientInterface
{
    public function __construct(
        private readonly LoggerInterface $logger,
        private readonly \Symfony\Contracts\HttpClient\HttpClientInterface $httpClient
    ) {}

    public function sendMessage(string $token, string $chatId, string $message): bool
    {
        $url = sprintf('https://api.telegram.org/bot%s/sendMessage', $token);

        $payload = [
            'chat_id' => $chatId,
            'text' => $message,
        ];

        try {
            $response = $this->httpClient->request('POST', $url, [
                'json' => $payload,
                'headers' => [
                    'Content-Type' => 'application/json',
                ],
            ]);

            $statusCode = $response->getStatusCode();
            $responseData = $response->toArray(false);

            if ($statusCode !== 200) {
                $this->logger->error('Telegram API error', [
                    'status' => $statusCode,
                    'response' => $responseData,
                    'token' => substr($token, 0, 10) . '...',
                    'chat_id' => $chatId,
                ]);
                return false;
            }

            return true;
        } catch (\Exception $e) {
            $this->logger->error('Failed to send Telegram message', [
                'error' => $e->getMessage(),
                'token' => substr($token, 0, 10) . '...',
                'chat_id' => $chatId,
            ]);
            return false;
        }
    }
}
