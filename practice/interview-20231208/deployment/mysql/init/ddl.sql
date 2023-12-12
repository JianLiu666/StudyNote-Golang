CREATE DATABASE IF NOT EXISTS `trading`;

DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (

    `id`           INT UNSIGNED     NOT NULL AUTO_INCREMENT COMMENT '交易單唯一識別碼',
    `userId`       INT UNSIGNED     NOT NULL                COMMENT '用戶唯一識別碼',
    `roleType`     TINYINT UNSIGNED NOT NULL                COMMENT '掛單角色(e.g. 買方/賣方)',
    `orderType`    TINYINT UNSIGNED NOT NULL                COMMENT '交易單類型(e.g. 市價單/限價單)',
    `durationType` TINYINT UNSIGNED NOT NULL                COMMENT '交易單期限(e.g. ROD/IOC/FOK)',
    `price`        INT UNSIGNED     NOT NULL                COMMENT '交易單價格',
    `quantity`     INT UNSIGNED     NOT NULL                COMMENT '交易單數量',
    `status`       TINYINT UNSIGNED NOT NULL                COMMENT '交易單狀態',
    `timestamp`    DATETIME         NOT NULL                COMMENT '交易單時間戳',
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `userId` (`userId`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易單';


DROP TABLE IF EXISTS `transactionLogs`;
CREATE TABLE `transactionLogs` (

    `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '成交紀錄唯一識別碼',
    `buyerOrderId`  INT UNSIGNED NOT NULL                COMMENT '買方唯一識別碼',
    `sellerOrderId` INT UNSIGNED NOT NULL                COMMENT '賣方唯一識別罵',
    `price`         INT UNSIGNED NOT NULL                COMMENT '成交價格',
    `quantity`      INT UNSIGNED NOT NULL                COMMENT '成交數量',
    `timstamp`      INT UNSIGNED NOT NULL                COMMENT '成交時間戳',
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `buyerOrderId` (`buyerOrderId`),
    UNIQUE KEY `sellerOrderId` (`sellerOrderId`),

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='成交紀錄';