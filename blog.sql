  
CREATE TABLE `article` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL DEFAULT '',
  `content` text,
  `author_id` int(11) NOT NULL,
  `subtract` varchar(500) NOT NULL DEFAULT '',
  `click_number` int(11) NOT NULL DEFAULT '0',
  `praised_number` int(11) NOT NULL DEFAULT '0',
  `create_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  `is_show` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_authorId` (`author_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
 
CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '',
  `salt` varchar(32) NOT NULL DEFAULT '',
  `password` varchar(64) NOT NULL DEFAULT '',
  `email` varchar(64) DEFAULT NULL,
  `del_flag` tinyint(1) NOT NULL,
  `create_time` int(11) NOT NULL,
  `avatar` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
 
CREATE TABLE `visit_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_agent` varchar(512) DEFAULT NULL,
  `full_ip` varchar(512) DEFAULT NULL,
  `trim_ip` varchar(64) DEFAULT NULL,
  `visit_path` varchar(512) NOT NULL DEFAULT '',
  `user_hash` bigint(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
