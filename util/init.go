/*
 *@author  chengkenli
 *@project srAnr
 *@package util
 *@file    init
 *@date    2025/6/12 15:47
 */

package util

const SchemaMetaCreate = `
CREATE TABLE SchemaMetaCreate (
  app                varchar(100)  NOT NULL COMMENT 'App StarRocks CN Name',
  nickname           varchar(100)  DEFAULT NULL COMMENT 'Alias',
  alias              varchar(100)  DEFAULT NULL COMMENT 'App Alias',
  feip               varchar(200)  NOT NULL COMMENT 'F5,VIP,CLB,FE',
  user               varchar(200)  NOT NULL COMMENT 'StarRocks Admin User',
  password           varchar(500)  NOT NULL COMMENT 'StarRocks Admin Password',
  feport             int           NOT NULL DEFAULT '9030' COMMENT 'FE Query Port',
  address            varchar(500)  DEFAULT NULL COMMENT 'Manager(Commercial)',
  expire             int           DEFAULT '30' COMMENT 'License expiration reminder countdown (Enterprise level), unit: day',
  status             int           NOT NULL DEFAULT '0' COMMENT 'License expiration (enterprise level) switch, 0 off, 1 on',
  fe_log_path        varchar(500)  DEFAULT NULL COMMENT 'FE LOGPath',
  be_log_path        varchar(500)  DEFAULT NULL COMMENT 'BE LOGPath',
  be_meta_log        varchar(200)  DEFAULT NULL COMMENT 'FE META LOGPath',
  java_udf_path      varchar(500)  DEFAULT NULL COMMENT 'UDF LOGPath',
  manager_access_key varchar(500)  DEFAULT NULL COMMENT 'Manager Access key',
  manager_secret_key varchar(500)  DEFAULT NULL COMMENT 'Manager Secret key',
  updated_at         timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'CURRENT_TIMESTAMP'
) ENGINE=InnoDB COMMENT='StarRocks信息表';
`

const SchemaMetaInsert = "INSERT INTO SchemaMetaInsert (app, nickname, alias, feip, `user`, password, feport, address, expire, status, fe_log_path, be_log_path, be_meta_log, java_udf_path, manager_access_key, manager_secret_key) VALUES('sr-test', 'StarRocks(Tencent Cloud) SR-TEST', NULL, '127.0.0.1', 'chengkenli', '**************1.', 9030, 'http://127.0.0.1:19321', 30, 0, NULL, NULL, NULL, NULL, NULL, NULL)"
