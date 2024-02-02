<?php

/**
 * This class has been generated by dagger-php-sdk. DO NOT EDIT.
 */

declare(strict_types=1);

namespace Dagger;

/**
 * A port exposed by a container.
 */
class Port extends Client\AbstractObject implements Client\IdAble
{
    /**
     * The port description.
     */
    public function description(): string
    {
        $leafQueryBuilder = new \Dagger\Client\QueryBuilder('description');
        return (string)$this->queryLeaf($leafQueryBuilder, 'description');
    }

    /**
     * Skip the health check when run as a service.
     */
    public function experimentalSkipHealthcheck(): bool
    {
        $leafQueryBuilder = new \Dagger\Client\QueryBuilder('experimentalSkipHealthcheck');
        return (bool)$this->queryLeaf($leafQueryBuilder, 'experimentalSkipHealthcheck');
    }

    /**
     * A unique identifier for this Port.
     */
    public function id(): PortId
    {
        $leafQueryBuilder = new \Dagger\Client\QueryBuilder('id');
        return new \Dagger\PortId((string)$this->queryLeaf($leafQueryBuilder, 'id'));
    }

    /**
     * The port number.
     */
    public function port(): int
    {
        $leafQueryBuilder = new \Dagger\Client\QueryBuilder('port');
        return (int)$this->queryLeaf($leafQueryBuilder, 'port');
    }

    /**
     * The transport layer protocol.
     */
    public function protocol(): NetworkProtocol
    {
        $leafQueryBuilder = new \Dagger\Client\QueryBuilder('protocol');
        return \Dagger\NetworkProtocol::from((string)$this->queryLeaf($leafQueryBuilder, 'protocol'));
    }
}