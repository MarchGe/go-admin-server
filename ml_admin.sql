/*
 Navicat Premium Data Transfer

 Source Server         : MySQL（101.201.29.174）
 Source Server Type    : MySQL
 Source Server Version : 80024
 Source Host           : 101.201.29.174:3306
 Source Schema         : ml_admin

 Target Server Type    : MySQL
 Target Server Version : 80024
 File Encoding         : 65001

 Date: 11/05/2024 01:29:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_casbin_rule`(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2629 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_app
-- ----------------------------
DROP TABLE IF EXISTS `dv_app`;
CREATE TABLE `dv_app`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '版本',
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '部署包的路劲',
  `file_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '部署包的文件名',
  `port` smallint(0) NOT NULL COMMENT '应用端口',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_n_v`(`name`, `version`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_group
-- ----------------------------
DROP TABLE IF EXISTS `dv_group`;
CREATE TABLE `dv_group`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分组名称',
  `sort_num` int(0) NULL DEFAULT 0 COMMENT '顺序',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_host
-- ----------------------------
DROP TABLE IF EXISTS `dv_host`;
CREATE TABLE `dv_host`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'IP',
  `port` smallint(0) NOT NULL COMMENT 'ssh端口',
  `user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'ssh用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'ssh密码',
  `sort_num` int(0) NULL DEFAULT 0 COMMENT '顺序',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ip`(`ip`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_host_group
-- ----------------------------
DROP TABLE IF EXISTS `dv_host_group`;
CREATE TABLE `dv_host_group`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `host_id` bigint(0) NOT NULL COMMENT '主机ID',
  `group_id` bigint(0) NOT NULL COMMENT '分组ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_h_g_id`(`group_id`, `host_id`) USING BTREE,
  INDEX `idx_h_id`(`host_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 90 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_script
-- ----------------------------
DROP TABLE IF EXISTS `dv_script`;
CREATE TABLE `dv_script`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '版本',
  `content` varchar(10000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '脚本内容',
  `description` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '使用说明',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_n_v`(`name`, `version`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_task
-- ----------------------------
DROP TABLE IF EXISTS `dv_task`;
CREATE TABLE `dv_task`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `type` tinyint(0) NOT NULL COMMENT '任务类型',
  `status` tinyint(0) NULL DEFAULT 0 COMMENT '任务状态',
  `association_id` bigint(0) NOT NULL COMMENT '关联的具体任务的ID',
  `association_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '关联的具体任务表类型',
  `cron` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'cron表达式',
  `execute_type` tinyint(0) NULL DEFAULT 0 COMMENT '执行方式',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dv_task_deploy
-- ----------------------------
DROP TABLE IF EXISTS `dv_task_deploy`;
CREATE TABLE `dv_task_deploy`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `upload_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '部署包的上传路劲',
  `app_id` bigint(0) NOT NULL COMMENT '关联的应用ID',
  `script_id` bigint(0) NOT NULL COMMENT '关联的部署脚本ID',
  `host_group_id` bigint(0) NOT NULL COMMENT '关联的服务器分组ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_dept
-- ----------------------------
DROP TABLE IF EXISTS `ops_dept`;
CREATE TABLE `ops_dept`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '部门名称',
  `sort_num` int(0) NULL DEFAULT 0 COMMENT '部门顺序',
  `parent_id` bigint(0) NULL DEFAULT 0 COMMENT '父级部门ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE,
  INDEX `idx_pid`(`parent_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 52 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_exception_log
-- ----------------------------
DROP TABLE IF EXISTS `ops_exception_log`;
CREATE TABLE `ops_exception_log`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `user_id` bigint(0) NULL DEFAULT NULL COMMENT '登录用户的ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户昵称',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录IP',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '浏览器的userAgent',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求url',
  `query` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求url',
  `body` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'body参数信息',
  `error` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误内容',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_uname`(`nickname`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 84 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_icon
-- ----------------------------
DROP TABLE IF EXISTS `ops_icon`;
CREATE TABLE `ops_icon`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `value` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '图标',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 114 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ops_icon
-- ----------------------------
INSERT INTO `ops_icon` VALUES (1, 'zhangdan');
INSERT INTO `ops_icon` VALUES (2, 'tijikongjian');
INSERT INTO `ops_icon` VALUES (3, 'yewu');
INSERT INTO `ops_icon` VALUES (4, 'yingyongchengxu');
INSERT INTO `ops_icon` VALUES (5, 'quanxianyuechi');
INSERT INTO `ops_icon` VALUES (6, 'ziyuan');
INSERT INTO `ops_icon` VALUES (7, 'xinwenzixun');
INSERT INTO `ops_icon` VALUES (8, 'hezuoguanxi');
INSERT INTO `ops_icon` VALUES (9, '-fuwu');
INSERT INTO `ops_icon` VALUES (10, '-kefu');
INSERT INTO `ops_icon` VALUES (11, '-guoji');
INSERT INTO `ops_icon` VALUES (12, 'haiguan');
INSERT INTO `ops_icon` VALUES (13, 'touchengkongyun');
INSERT INTO `ops_icon` VALUES (14, 'caiwu');
INSERT INTO `ops_icon` VALUES (15, 'mianfei');
INSERT INTO `ops_icon` VALUES (16, 'jisuanqilishuai');
INSERT INTO `ops_icon` VALUES (17, 'checkbox-xuanzhong');
INSERT INTO `ops_icon` VALUES (18, 'Raidobox-xuanzhong');
INSERT INTO `ops_icon` VALUES (19, 'checkbox-xuanzhongbufen');
INSERT INTO `ops_icon` VALUES (20, 'youxiajiaogouxuan');
INSERT INTO `ops_icon` VALUES (21, 'shouye');
INSERT INTO `ops_icon` VALUES (22, 'wenti');
INSERT INTO `ops_icon` VALUES (23, 'liaotianduihua');
INSERT INTO `ops_icon` VALUES (24, 'dianhua');
INSERT INTO `ops_icon` VALUES (25, 'dianhua-yuankuang');
INSERT INTO `ops_icon` VALUES (26, 'lingdang');
INSERT INTO `ops_icon` VALUES (27, 'laba');
INSERT INTO `ops_icon` VALUES (28, 'shoucang');
INSERT INTO `ops_icon` VALUES (29, 'maikefeng');
INSERT INTO `ops_icon` VALUES (30, 'xihuan');
INSERT INTO `ops_icon` VALUES (31, 'shijian');
INSERT INTO `ops_icon` VALUES (32, 'shanguangdeng-zidong');
INSERT INTO `ops_icon` VALUES (33, 'shanguangdeng-guanbi');
INSERT INTO `ops_icon` VALUES (34, 'baocun');
INSERT INTO `ops_icon` VALUES (35, 'morentouxiang');
INSERT INTO `ops_icon` VALUES (36, 'zhucetianjiahaoyou');
INSERT INTO `ops_icon` VALUES (37, 'renwu');
INSERT INTO `ops_icon` VALUES (38, 'bianjishuru');
INSERT INTO `ops_icon` VALUES (39, 'yingwenmoshi');
INSERT INTO `ops_icon` VALUES (40, 'jianpan');
INSERT INTO `ops_icon` VALUES (41, 'paizhao');
INSERT INTO `ops_icon` VALUES (42, 'zhongwenmoshi');
INSERT INTO `ops_icon` VALUES (43, 'tupian');
INSERT INTO `ops_icon` VALUES (44, 'xianshikejian');
INSERT INTO `ops_icon` VALUES (45, 'suoding');
INSERT INTO `ops_icon` VALUES (46, 'yincangbukejian');
INSERT INTO `ops_icon` VALUES (47, 'jiesuo');
INSERT INTO `ops_icon` VALUES (48, 'shaixuanguolv');
INSERT INTO `ops_icon` VALUES (49, 'anzhuangshigong');
INSERT INTO `ops_icon` VALUES (50, 'zhuxiaoguanji');
INSERT INTO `ops_icon` VALUES (51, 'chaping');
INSERT INTO `ops_icon` VALUES (52, 'haoping');
INSERT INTO `ops_icon` VALUES (53, 'liebiaoshitucaidan');
INSERT INTO `ops_icon` VALUES (54, 'gonggeshitu');
INSERT INTO `ops_icon` VALUES (55, 'jia-fangkuang');
INSERT INTO `ops_icon` VALUES (56, 'jia-yuankuang');
INSERT INTO `ops_icon` VALUES (57, 'jian-fangkuang');
INSERT INTO `ops_icon` VALUES (58, 'jian-yuankuang');
INSERT INTO `ops_icon` VALUES (59, 'zhengquewancheng-yuankuang');
INSERT INTO `ops_icon` VALUES (60, 'cuowuguanbiquxiao-yuankuang');
INSERT INTO `ops_icon` VALUES (61, 'cuowuguanbiquxiao-fangkuang');
INSERT INTO `ops_icon` VALUES (62, 'wenhao-yuankuang');
INSERT INTO `ops_icon` VALUES (63, 'xinxi-yuankuang');
INSERT INTO `ops_icon` VALUES (64, 'gantanhao-sanjiaokuang');
INSERT INTO `ops_icon` VALUES (65, 'gantanhao-yuankuang');
INSERT INTO `ops_icon` VALUES (66, 'shangyiyehoutuifanhui-yuankuang');
INSERT INTO `ops_icon` VALUES (67, 'xiayiyeqianjinchakangengduo-yuankuang');
INSERT INTO `ops_icon` VALUES (68, 'xiangxiazhankai-yuankuang');
INSERT INTO `ops_icon` VALUES (69, 'xiangshangshouqi-yuankuang');
INSERT INTO `ops_icon` VALUES (70, 'weizhi');
INSERT INTO `ops_icon` VALUES (71, 'daohang');
INSERT INTO `ops_icon` VALUES (72, 'jiankongshexiangtou');
INSERT INTO `ops_icon` VALUES (73, 'baobiao');
INSERT INTO `ops_icon` VALUES (74, 'bingtu');
INSERT INTO `ops_icon` VALUES (75, 'tiaoxingtu');
INSERT INTO `ops_icon` VALUES (76, 'zhexiantu');
INSERT INTO `ops_icon` VALUES (77, 'zhinanzhidao');
INSERT INTO `ops_icon` VALUES (78, 'dianpu');
INSERT INTO `ops_icon` VALUES (79, 'yonghuziliao');
INSERT INTO `ops_icon` VALUES (80, 'pifuzhuti');
INSERT INTO `ops_icon` VALUES (81, 'diamond');
INSERT INTO `ops_icon` VALUES (82, 'yinhangqia');
INSERT INTO `ops_icon` VALUES (83, 'yunshuzhongwuliu');
INSERT INTO `ops_icon` VALUES (84, 'baoguofahuo');
INSERT INTO `ops_icon` VALUES (85, 'chaibaoguoqujian');
INSERT INTO `ops_icon` VALUES (86, 'zitigui');
INSERT INTO `ops_icon` VALUES (87, 'caigou');
INSERT INTO `ops_icon` VALUES (88, 'shangpin');
INSERT INTO `ops_icon` VALUES (89, 'peizaizhuangche');
INSERT INTO `ops_icon` VALUES (90, 'zhiliang');
INSERT INTO `ops_icon` VALUES (91, 'anquanbaozhang');
INSERT INTO `ops_icon` VALUES (92, 'cangkucangchu');
INSERT INTO `ops_icon` VALUES (93, 'zhongzhuanzhan');
INSERT INTO `ops_icon` VALUES (94, 'kucun');
INSERT INTO `ops_icon` VALUES (95, 'moduanwangdian');
INSERT INTO `ops_icon` VALUES (96, 'qianshoushenpitongguo');
INSERT INTO `ops_icon` VALUES (97, 'juqianshou');
INSERT INTO `ops_icon` VALUES (98, 'jijianfasong');
INSERT INTO `ops_icon` VALUES (99, 'qiyeyuanquwuye');
INSERT INTO `ops_icon` VALUES (100, 'jiesuan');
INSERT INTO `ops_icon` VALUES (101, 'jifen');
INSERT INTO `ops_icon` VALUES (102, 'ziliaoshouce');
INSERT INTO `ops_icon` VALUES (103, 'youhuijuan');
INSERT INTO `ops_icon` VALUES (104, 'danju');
INSERT INTO `ops_icon` VALUES (105, 'chuangjiandanju');
INSERT INTO `ops_icon` VALUES (106, 'shanchu');
INSERT INTO `ops_icon` VALUES (107, 'tubiao');
INSERT INTO `ops_icon` VALUES (108, 'yonghuguanli');
INSERT INTO `ops_icon` VALUES (109, 'xitongguanli');
INSERT INTO `ops_icon` VALUES (110, 'xitongguanli-caidanguanli');
INSERT INTO `ops_icon` VALUES (111, 'jiaoseguanli');
INSERT INTO `ops_icon` VALUES (112, 'suofang');
INSERT INTO `ops_icon` VALUES (113, 'shousuo');

-- ----------------------------
-- Table structure for ops_job
-- ----------------------------
DROP TABLE IF EXISTS `ops_job`;
CREATE TABLE `ops_job`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '岗位名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '岗位描述',
  `sort_num` bigint(0) NOT NULL DEFAULT 0 COMMENT '排序',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_login_log
-- ----------------------------
DROP TABLE IF EXISTS `ops_login_log`;
CREATE TABLE `ops_login_log`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `user_id` bigint(0) NULL DEFAULT NULL COMMENT '登录用户的ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户昵称',
  `real_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户真名',
  `dept_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '部门名称',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录IP',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '浏览器的userAgent',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录地点',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_uname`(`nickname`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_menu
-- ----------------------------
DROP TABLE IF EXISTS `ops_menu`;
CREATE TABLE `ops_menu`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `symbol` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '权限标识',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '菜单名称',
  `icon` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '菜单图标',
  `sort_num` int(0) NOT NULL DEFAULT 0 COMMENT '菜单顺序',
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '路由路径',
  `display` tinyint(0) NOT NULL DEFAULT 0 COMMENT '是否显示，0-否，1-是',
  `external` tinyint(0) NOT NULL DEFAULT 0 COMMENT '是否外链，0-否，1-是',
  `parent_id` bigint(0) NOT NULL DEFAULT 0 COMMENT '父菜单ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `external_way` tinyint(0) NULL DEFAULT 0 COMMENT '外链打开方式（仅外链有效），0-外联，1-内嵌',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pid`(`parent_id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE,
  INDEX `idx_symbol`(`symbol`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 137 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ops_menu
-- ----------------------------
INSERT INTO `ops_menu` VALUES (1, '', '系统管理', 'xitongguanli', 100, '', 1, 0, 0, '2018-10-22 21:53:02', '2024-05-09 14:42:35', 0);
INSERT INTO `ops_menu` VALUES (2, '', '用户管理', 'yonghuguanli', 200, '/sys/user', 1, 0, 1, '2018-10-22 22:44:08', '2024-03-12 23:51:41', 0);
INSERT INTO `ops_menu` VALUES (3, '', '角色管理', 'jiaoseguanli', 300, '/sys/role', 1, 0, 1, '2018-10-22 22:44:32', '2024-03-12 23:53:21', 0);
INSERT INTO `ops_menu` VALUES (4, '', '菜单管理', 'xitongguanli-caidanguanli', 400, '/sys/menu', 1, 0, 1, '2018-10-22 22:51:33', '2024-03-12 23:53:47', 0);
INSERT INTO `ops_menu` VALUES (27, 'user:update', '编辑', '', 200, '', 0, 0, 2, '2024-03-08 15:59:32', '2024-03-12 23:52:26', 0);
INSERT INTO `ops_menu` VALUES (30, 'user:disable', '禁用', '', 300, '', 0, 0, 2, '2024-03-09 01:45:19', '2024-03-12 23:52:18', 0);
INSERT INTO `ops_menu` VALUES (31, 'user:enable', '启用', '', 400, '', 0, 0, 2, '2024-03-09 01:45:47', '2024-03-12 23:52:15', 0);
INSERT INTO `ops_menu` VALUES (34, '', '部门管理', 'ziyuan', 600, '/sys/dept', 1, 0, 1, '2024-03-10 01:21:21', '2024-03-16 00:59:47', 0);
INSERT INTO `ops_menu` VALUES (35, '', '岗位管理', 'xinwenzixun', 700, '/sys/job', 1, 0, 1, '2024-03-10 01:24:15', '2024-03-15 13:43:32', 0);
INSERT INTO `ops_menu` VALUES (38, '', '系统日志', 'danju', 200, '/monitor', 1, 0, 0, '2024-03-11 19:19:34', '2024-04-06 06:19:52', 0);
INSERT INTO `ops_menu` VALUES (39, '', '操作日志', 'danju', 1, '/sys/log/op', 1, 0, 38, '2024-03-11 19:21:55', '2024-05-09 14:58:52', 0);
INSERT INTO `ops_menu` VALUES (40, '', '异常日志', 'chuangjiandanju', 2, '/sys/log/exception', 1, 0, 38, '2024-03-11 19:22:37', '2024-05-09 14:59:32', 0);
INSERT INTO `ops_menu` VALUES (41, 'user:changePassword', '修改密码', '', 600, '', 0, 0, 2, '2024-03-12 17:51:11', '2024-03-12 23:52:09', 0);
INSERT INTO `ops_menu` VALUES (42, 'user:resetPassword', '重置密码', '', 700, '', 0, 0, 2, '2024-03-12 17:51:35', '2024-03-12 23:52:06', 0);
INSERT INTO `ops_menu` VALUES (43, 'user:delete', '删除用户', '', 800, '', 0, 0, 2, '2024-03-12 17:51:54', '2024-03-12 23:52:02', 0);
INSERT INTO `ops_menu` VALUES (44, 'user:menus', '分配权限', '', 900, '', 0, 0, 2, '2024-03-12 17:55:27', '2024-03-12 23:51:58', 0);
INSERT INTO `ops_menu` VALUES (45, 'user:add', '新增', '', 100, '', 0, 0, 2, '2024-03-12 19:56:30', '2024-03-12 23:52:31', 0);
INSERT INTO `ops_menu` VALUES (46, 'role:add', '新增', '', 200, '', 0, 0, 3, '2024-03-12 19:56:49', '2024-03-12 23:53:29', 0);
INSERT INTO `ops_menu` VALUES (47, 'role:update', '编辑', '', 300, '', 0, 0, 3, '2024-03-12 19:57:10', '2024-03-12 23:53:34', 0);
INSERT INTO `ops_menu` VALUES (48, 'role:delete', '删除', '', 400, '', 0, 0, 3, '2024-03-12 19:57:25', '2024-03-12 23:53:39', 0);
INSERT INTO `ops_menu` VALUES (52, 'user:list', '列表', '', 50, '', 0, 0, 2, '2024-03-12 23:52:52', '2024-04-14 23:11:18', 0);
INSERT INTO `ops_menu` VALUES (53, 'role:list', '列表', '', 100, '', 0, 0, 3, '2024-03-12 23:53:17', '2024-03-13 00:20:44', 0);
INSERT INTO `ops_menu` VALUES (54, 'menu:tree', '查询菜单树', '', 100, '', 0, 0, 4, '2024-03-12 23:56:40', '2024-03-12 23:56:48', 0);
INSERT INTO `ops_menu` VALUES (55, 'role:menus', '分配权限', '', 500, '', 0, 0, 3, '2024-03-12 23:58:29', '2024-03-12 23:58:29', 0);
INSERT INTO `ops_menu` VALUES (56, 'menu:add', '新增', '', 200, '', 0, 0, 4, '2024-03-13 00:00:02', '2024-03-13 00:00:02', 0);
INSERT INTO `ops_menu` VALUES (57, 'menu:update', '编辑', '', 300, '', 0, 0, 4, '2024-03-13 00:00:22', '2024-03-13 00:00:22', 0);
INSERT INTO `ops_menu` VALUES (58, 'menu:delete', '删除', '', 400, '', 0, 0, 4, '2024-03-13 00:00:36', '2024-03-13 00:00:36', 0);
INSERT INTO `ops_menu` VALUES (59, '', '常用网站', 'tijikongjian', 300, '', 1, 0, 0, '2024-03-13 23:05:56', '2024-04-14 01:25:07', 0);
INSERT INTO `ops_menu` VALUES (60, '', 'Element-Plus', 'zhinanzhidao', 1000, 'https://element-plus.org/zh-CN/', 1, 1, 59, '2024-03-13 23:07:07', '2024-04-04 04:01:24', 1);
INSERT INTO `ops_menu` VALUES (62, 'icon:all', '查询全部图标', '', 500, '', 0, 0, 4, '2024-03-14 01:49:58', '2024-04-22 01:13:34', 0);
INSERT INTO `ops_menu` VALUES (70, 'job:list', '列表', '', 100, '', 0, 0, 35, '2024-03-15 02:41:43', '2024-03-15 02:41:53', 0);
INSERT INTO `ops_menu` VALUES (71, 'job:add', '新增', '', 200, '', 0, 0, 35, '2024-03-15 02:42:14', '2024-03-15 02:42:14', 0);
INSERT INTO `ops_menu` VALUES (72, 'job:update', '编辑', '', 300, '', 0, 0, 35, '2024-03-15 02:42:31', '2024-03-15 02:42:31', 0);
INSERT INTO `ops_menu` VALUES (73, 'job:delete', '删除', '', 400, '', 0, 0, 35, '2024-03-15 02:42:47', '2024-03-15 02:42:47', 0);
INSERT INTO `ops_menu` VALUES (75, 'dept:tree', '查询部门树', '', 100, '', 0, 0, 34, '2024-03-16 01:01:05', '2024-03-16 01:01:05', 0);
INSERT INTO `ops_menu` VALUES (76, 'dept:add', '新增', '', 200, '', 0, 0, 34, '2024-03-16 01:01:27', '2024-03-16 01:01:27', 0);
INSERT INTO `ops_menu` VALUES (77, 'dept:update', '编辑', '', 300, '', 0, 0, 34, '2024-03-16 01:01:44', '2024-03-16 01:01:44', 0);
INSERT INTO `ops_menu` VALUES (78, 'dept:delete', '删除', '', 400, '', 0, 0, 34, '2024-03-16 01:02:01', '2024-03-16 01:02:01', 0);
INSERT INTO `ops_menu` VALUES (82, '', '登录日志', 'qianshoushenpitongguo', 0, '/sys/log/login', 1, 0, 38, '2024-03-17 03:25:36', '2024-05-09 14:58:01', 0);
INSERT INTO `ops_menu` VALUES (83, '', '运维管理', 'zhiliang', 150, '', 1, 0, 0, '2024-03-17 03:28:03', '2024-04-21 18:02:42', 0);
INSERT INTO `ops_menu` VALUES (84, 'terminal:connect', 'Web Shell', 'shanguangdeng-zidong', 600, '/devops/terminal', 1, 0, 112, '2024-03-17 03:31:41', '2024-04-21 18:02:52', 0);
INSERT INTO `ops_menu` VALUES (93, '', '服务器监控', 'baobiao', 500, '/devops/monitor/performance', 1, 0, 83, '2024-04-06 05:10:22', '2024-04-22 00:24:54', 0);
INSERT INTO `ops_menu` VALUES (94, 'monitor:detail', '实时状况', '', 200, '', 0, 0, 93, '2024-04-08 10:50:42', '2024-04-22 00:24:51', 0);
INSERT INTO `ops_menu` VALUES (95, 'monitor:delete', '删除', '', 300, '', 0, 0, 93, '2024-04-08 10:51:10', '2024-04-22 00:24:47', 0);
INSERT INTO `ops_menu` VALUES (97, '', '应用管理', 'tubiao', 300, '/devops/app', 1, 0, 83, '2024-04-16 16:12:40', '2024-04-21 21:57:04', 0);
INSERT INTO `ops_menu` VALUES (98, '', '服务器组', 'xinwenzixun', 200, '/devops/group', 1, 0, 83, '2024-04-16 16:12:59', '2024-04-21 18:05:56', 0);
INSERT INTO `ops_menu` VALUES (99, '', '服务器', 'renwu', 100, '/devops/host', 1, 0, 83, '2024-04-20 02:13:25', '2024-04-21 18:06:00', 0);
INSERT INTO `ops_menu` VALUES (100, 'host:list', '列表', '', 100, '', 0, 0, 99, '2024-04-20 02:45:24', '2024-04-20 17:22:32', 0);
INSERT INTO `ops_menu` VALUES (101, 'host:add', '新增', '', 200, '', 0, 0, 99, '2024-04-20 02:46:08', '2024-04-20 02:46:08', 0);
INSERT INTO `ops_menu` VALUES (102, 'host:update', '编辑', '', 300, '', 0, 0, 99, '2024-04-20 02:46:33', '2024-04-20 02:46:33', 0);
INSERT INTO `ops_menu` VALUES (103, 'host:delete', '删除', '', 400, '', 0, 0, 99, '2024-04-20 02:46:49', '2024-04-20 02:46:49', 0);
INSERT INTO `ops_menu` VALUES (104, 'group:list', '列表', '', 100, '', 0, 0, 98, '2024-04-20 02:47:19', '2024-04-20 02:47:19', 0);
INSERT INTO `ops_menu` VALUES (105, 'group:add', '新增', '', 200, '', 0, 0, 98, '2024-04-20 02:47:33', '2024-04-20 02:47:33', 0);
INSERT INTO `ops_menu` VALUES (106, 'group:update', '编辑', '', 300, '', 0, 0, 98, '2024-04-20 02:47:50', '2024-04-20 02:47:50', 0);
INSERT INTO `ops_menu` VALUES (107, 'group:delete', '删除', '', 400, '', 0, 0, 98, '2024-04-20 02:48:07', '2024-04-20 02:48:07', 0);
INSERT INTO `ops_menu` VALUES (108, 'host:connectTest', '连接测试', '', 500, '', 0, 0, 99, '2024-04-20 03:49:45', '2024-04-20 03:49:45', 0);
INSERT INTO `ops_menu` VALUES (109, 'host:connect', 'Shell', '', 600, '', 0, 0, 99, '2024-04-20 23:54:15', '2024-04-21 02:10:17', 0);
INSERT INTO `ops_menu` VALUES (110, '', '任务管理', 'checkbox-xuanzhong', 420, '/devops/task', 1, 0, 83, '2024-04-21 16:01:40', '2024-04-21 21:57:24', 0);
INSERT INTO `ops_menu` VALUES (111, '', '脚本管理', 'jijianfasong', 400, '/devops/script', 1, 0, 83, '2024-04-21 17:58:59', '2024-04-21 21:57:18', 0);
INSERT INTO `ops_menu` VALUES (112, '', '系统工具', 'anzhuangshigong', 170, '', 1, 0, 0, '2024-04-21 18:02:08', '2024-05-09 14:42:32', 0);
INSERT INTO `ops_menu` VALUES (113, 'app:add', '新增', '', 200, '', 0, 0, 97, '2024-04-21 22:08:04', '2024-04-21 22:08:04', 0);
INSERT INTO `ops_menu` VALUES (114, 'app:list', '列表', '', 100, '', 0, 0, 97, '2024-04-21 22:08:24', '2024-04-21 22:08:24', 0);
INSERT INTO `ops_menu` VALUES (115, 'app:update', '编辑', '', 300, '', 0, 0, 97, '2024-04-21 22:08:43', '2024-04-21 22:08:43', 0);
INSERT INTO `ops_menu` VALUES (116, 'app:delete', '删除', '', 400, '', 0, 0, 97, '2024-04-21 22:09:04', '2024-04-21 22:09:04', 0);
INSERT INTO `ops_menu` VALUES (117, 'script:list', '列表', '', 100, '', 0, 0, 111, '2024-04-21 23:27:39', '2024-04-21 23:27:39', 0);
INSERT INTO `ops_menu` VALUES (118, 'script:add', '新增', '', 200, '', 0, 0, 111, '2024-04-21 23:28:02', '2024-04-21 23:28:02', 0);
INSERT INTO `ops_menu` VALUES (119, 'script:update', '编辑', '', 300, '', 0, 0, 111, '2024-04-21 23:28:16', '2024-04-21 23:28:16', 0);
INSERT INTO `ops_menu` VALUES (120, 'script:delete', '删除', '', 400, '', 0, 0, 111, '2024-04-21 23:28:37', '2024-04-21 23:28:37', 0);
INSERT INTO `ops_menu` VALUES (121, 'monitor:list', '列表', '', 100, '', 0, 0, 93, '2024-04-22 00:24:37', '2024-04-22 00:24:59', 0);
INSERT INTO `ops_menu` VALUES (122, 'task:list', '列表', '', 100, '', 0, 0, 110, '2024-04-22 00:49:15', '2024-04-22 00:49:15', 0);
INSERT INTO `ops_menu` VALUES (123, 'task:add', '新增', '', 200, '', 0, 0, 110, '2024-04-22 00:49:29', '2024-04-22 00:49:29', 0);
INSERT INTO `ops_menu` VALUES (124, 'task:update', '编辑', '', 300, '', 0, 0, 110, '2024-04-22 00:49:46', '2024-04-22 00:49:46', 0);
INSERT INTO `ops_menu` VALUES (125, 'task:delete', '删除', '', 400, '', 0, 0, 110, '2024-04-22 00:50:00', '2024-04-22 00:50:00', 0);
INSERT INTO `ops_menu` VALUES (126, 'app:upload', '上传', '', 500, '', 0, 0, 97, '2024-04-22 20:21:20', '2024-04-22 20:21:20', 0);
INSERT INTO `ops_menu` VALUES (127, 'app:download', '下载', '', 600, '', 0, 0, 97, '2024-04-23 04:13:11', '2024-04-23 04:13:34', 0);
INSERT INTO `ops_menu` VALUES (128, 'task:start', '启动', '', 250, '', 0, 0, 110, '2024-04-28 23:47:40', '2024-05-03 12:13:38', 0);
INSERT INTO `ops_menu` VALUES (129, 'task:stop', '停止', '', 260, '', 0, 0, 110, '2024-04-30 23:17:27', '2024-05-03 12:13:43', 0);
INSERT INTO `ops_menu` VALUES (130, 'task:log', '日志', '', 700, '', 0, 0, 110, '2024-05-03 12:14:09', '2024-05-03 16:13:27', 0);
INSERT INTO `ops_menu` VALUES (131, 'loginLog:list', '列表', '', 100, '', 0, 0, 82, '2024-05-09 14:58:31', '2024-05-09 14:58:31', 0);
INSERT INTO `ops_menu` VALUES (132, 'loginLog:delete', '清空', '', 200, '', 0, 0, 82, '2024-05-09 14:58:46', '2024-05-09 14:58:46', 0);
INSERT INTO `ops_menu` VALUES (133, 'opLog:list', '列表', '', 100, '', 0, 0, 39, '2024-05-09 14:59:09', '2024-05-09 14:59:09', 0);
INSERT INTO `ops_menu` VALUES (134, 'opLog:delete', '清空', '', 200, '', 0, 0, 39, '2024-05-09 14:59:25', '2024-05-09 14:59:25', 0);
INSERT INTO `ops_menu` VALUES (135, 'exceptionLog:list', '列表', '', 100, '', 0, 0, 40, '2024-05-09 14:59:44', '2024-05-09 14:59:44', 0);
INSERT INTO `ops_menu` VALUES (136, 'exceptionLog:delete', '清空', '', 200, '', 0, 0, 40, '2024-05-09 14:59:58', '2024-05-09 14:59:58', 0);

-- ----------------------------
-- Table structure for ops_op_log
-- ----------------------------
DROP TABLE IF EXISTS `ops_op_log`;
CREATE TABLE `ops_op_log`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `user_id` bigint(0) NULL DEFAULT NULL COMMENT '登录用户的ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户昵称',
  `real_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户真名',
  `dept_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '部门名称',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录IP',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '浏览器的userAgent',
  `action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作名称',
  `target` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作对象',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求url',
  `query` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求url',
  `body` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'body参数信息',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录地点',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_uname`(`nickname`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1011 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_role
-- ----------------------------
DROP TABLE IF EXISTS `ops_role`;
CREATE TABLE `ops_role`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '角色名称',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `sort_num` bigint(0) NOT NULL DEFAULT 0 COMMENT '排序',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 40 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `ops_role_menu`;
CREATE TABLE `ops_role_menu`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `role_id` bigint(0) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(0) NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_r_m_id`(`role_id`, `menu_id`) USING BTREE,
  INDEX `idx_m_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3118 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色与菜单关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_role_user
-- ----------------------------
DROP TABLE IF EXISTS `ops_role_user`;
CREATE TABLE `ops_role_user`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `role_id` bigint(0) NOT NULL COMMENT '角色ID',
  `user_id` bigint(0) NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_r_u_id`(`role_id`, `user_id`) USING BTREE,
  INDEX `idx_uid`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 231 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色与用户关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_settings
-- ----------------------------
DROP TABLE IF EXISTS `ops_settings`;
CREATE TABLE `ops_settings`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(0) NOT NULL COMMENT '用户ID',
  `key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '键名',
  `value` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '值',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_uid_key`(`user_id`, `key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_user
-- ----------------------------
DROP TABLE IF EXISTS `ops_user`;
CREATE TABLE `ops_user`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `cellphone` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '邮箱',
  `sex` tinyint(0) NOT NULL DEFAULT 0 COMMENT '性别，0-男，1-女',
  `birthday` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '生日',
  `status` tinyint(0) NOT NULL DEFAULT 0 COMMENT '账号状态，0-正常，1-禁用',
  `root` tinyint(0) NOT NULL DEFAULT 0 COMMENT '是否超级用户，0-否，1-是',
  `create_time` datetime(0) NOT NULL COMMENT '创建时间',
  `update_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `dept_id` bigint(0) NULL DEFAULT 0 COMMENT '部门ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_email`(`email`) USING BTREE,
  INDEX `idx_name`(`name`) USING BTREE,
  INDEX `idx_cellphone`(`cellphone`) USING BTREE,
  INDEX `idx_nickname`(`nickname`) USING BTREE,
  INDEX `fk_ops_user_dept`(`dept_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 78 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '系统用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_user_job
-- ----------------------------
DROP TABLE IF EXISTS `ops_user_job`;
CREATE TABLE `ops_user_job`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(0) NOT NULL COMMENT '用户ID',
  `job_id` bigint(0) NOT NULL COMMENT '岗位ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_u_j_id`(`user_id`, `job_id`) USING BTREE,
  INDEX `idx_j_id`(`job_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 105 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_user_menu
-- ----------------------------
DROP TABLE IF EXISTS `ops_user_menu`;
CREATE TABLE `ops_user_menu`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `menu_id` bigint(0) NOT NULL COMMENT '菜单ID',
  `user_id` bigint(0) NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_u_m_id`(`user_id`, `menu_id`) USING BTREE,
  INDEX `idx_m_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 390 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_user_password
-- ----------------------------
DROP TABLE IF EXISTS `ops_user_password`;
CREATE TABLE `ops_user_password`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(0) NOT NULL COMMENT '用户ID',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户密码',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_uid`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 72 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户密码表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
