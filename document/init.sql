CREATE DATABASE IF NOT EXISTS `gfdemo` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

use gfdemo;

# drop table if exists user;
create table if not exists `user`
(
  `id`       bigint unsigned   not null primary key auto_increment,
  `name`     varchar(21)       not null comment 'user name',
  `nick`     varchar(20)       not null comment 'user nickname',
  `password` varchar(31)       not null comment 'user password',
  `email`    varchar(255)      not null comment 'user email',
  `age`      smallint unsigned not null comment 'user age',
  `head_img` varchar(255)      not null default '' comment 'user head img url'
  );

# drop table if exists equipment;
create table if not exists `equipment`
(
  id   bigint unsigned primary key auto_increment,
  name varchar(255)     not null comment 'equipment name',
  type tinyint unsigned not null comment 'equipment type',
  atk  int              not null comment 'attack damage',
  mag  int              not null comment 'magical damage',
  def  int              not null comment 'physical defense',
  res  int              not null comment 'magical defense',
  spd  int              not null comment 'speed'
  );
