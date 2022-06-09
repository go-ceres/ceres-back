--  //   Copyright 2021. Go-Ceres
--  //   Author https://github.com/go-ceres/go-ceres
--  //
--  //   Licensed under the Apache License, Version 2.0 (the "License");
--  //   you may not use this file except in compliance with the License.
--  //   You may obtain a copy of the License at
--  //
--  //       http://www.apache.org/licenses/LICENSE-2.0
--  //
--  //   Unless required by applicable law or agreed to in writing, software
--  //   distributed under the License is distributed on an "AS IS" BASIS,
--  //   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
--  //   See the License for the specific language governing permissions and
--  //   limitations under the License.

CREATE TABLE `t_store_industry` (
                                    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                    `industryNo` varchar(50) DEFAULT NULL COMMENT '行业编码',
                                    `name` varchar(20) DEFAULT NULL COMMENT '行业名称',
                                    `parentedNo` varchar(50) DEFAULT NULL COMMENT '父级编号',
                                    `imageUrl` varchar(50) DEFAULT NULL COMMENT '图片',
                                    `updateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    `creationDate` datetime DEFAULT CURRENT_TIMESTAMP,
                                    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='店铺经营行业标签';
