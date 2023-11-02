CREATE DATABASE IF NOT EXISTS `blog`;

DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'UUID',
  `name` varchar(100) DEFAULT '' COMMENT '標籤名稱',
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '創建時間',
  `created_by` varchar(100) DEFAULT '' COMMENT '創建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改時間',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '狀態: 0禁用, 1:啟用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章標籤管理';

-- Prometheus metrics
CREATE USER 'exporter'@'%' IDENTIFIED BY '123456';
GRANT PROCESS, REPLICATION CLIENT ON *.* TO 'exporter'@'%';
GRANT SELECT ON *.* TO 'exporter'@'%';