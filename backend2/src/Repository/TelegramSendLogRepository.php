<?php

namespace App\Repository;

use App\Entity\TelegramSendLog;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * @extends ServiceEntityRepository<TelegramSendLog>
 */
class TelegramSendLogRepository extends ServiceEntityRepository
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, TelegramSendLog::class);
    }

    //    /**
    //     * @return TelegramSendLog[] Returns an array of TelegramSendLog objects
    //     */
    //    public function findByExampleField($value): array
    //    {
    //        return $this->createQueryBuilder('t')
    //            ->andWhere('t.exampleField = :val')
    //            ->setParameter('val', $value)
    //            ->orderBy('t.id', 'ASC')
    //            ->setMaxResults(10)
    //            ->getQuery()
    //            ->getResult()
    //        ;
    //    }

    //    public function findOneBySomeField($value): ?TelegramSendLog
    //    {
    //        return $this->createQueryBuilder('t')
    //            ->andWhere('t.exampleField = :val')
    //            ->setParameter('val', $value)
    //            ->getQuery()
    //            ->getOneOrNullResult()
    //        ;
    //    }

    public function getByShopIDAndOrderID(int $shopId, int $orderId): ?TelegramSendLog
    {
        return $this->createQueryBuilder('t')
            ->andWhere('t.shopId = :shopId')
            ->andWhere('t.orderId = :orderId')
            ->setParameter('shopId', $shopId)
            ->setParameter('orderId', $orderId)
            ->getQuery()
            ->getOneOrNullResult();
    }

    public function getStatsForLast7Days(int $shopId): TelegramStats
    {
        $sevenDaysAgo = new \DateTime('-7 days');

        $qb = $this->createQueryBuilder('t')
            ->select('COUNT(t.id) as totalCount')
            ->addSelect('SUM(CASE WHEN t.status = :sentStatus THEN 1 ELSE 0 END) as sentCount')
            ->addSelect('SUM(CASE WHEN t.status = :failedStatus THEN 1 ELSE 0 END) as failedCount')
            ->addSelect('MAX(t.sentAt) as lastSentAt')
            ->andWhere('t.shopId = :shopId')
            ->andWhere('t.sentAt >= :sevenDaysAgo')
            ->setParameter('shopId', $shopId)
            ->setParameter('sevenDaysAgo', $sevenDaysAgo)
            ->setParameter('sentStatus', TelegramSendLog::STATUS_SENT)
            ->setParameter('failedStatus', TelegramSendLog::STATUS_FAILED);

        $result = $qb->getQuery()
            ->setCacheMode(\Doctrine\ORM\Cache::MODE_REFRESH) // Отключаем кэш для этого запроса
            ->getSingleResult();

        return new TelegramStats(
            (int) $result['sentCount'] ?: 0,
            (int) $result['failedCount'] ?: 0,
            $result['lastSentAt'] ? new \DateTime($result['lastSentAt']) : null
        );
    }
}
