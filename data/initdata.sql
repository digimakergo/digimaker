-- MySQL dump 10.13  Distrib 8.0.17, for osx10.15 (x86_64)
--
-- Host: 127.0.0.1    Database: dm_test
-- ------------------------------------------------------
-- Server version	5.7.30-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `dm_article`
--

DROP TABLE IF EXISTS `dm_article`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_article` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `status` int(11) NOT NULL DEFAULT '0',
  `version` int(11) NOT NULL DEFAULT '0',
  `author` int(11) NOT NULL DEFAULT '0',
  `editors` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `body` text NOT NULL,
  `languages` json DEFAULT NULL,
  `related_articles` json DEFAULT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_article`
--

LOCK TABLES `dm_article` WRITE;
/*!40000 ALTER TABLE `dm_article` DISABLE KEYS */;
/*!40000 ALTER TABLE `dm_article` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_draft`
--

DROP TABLE IF EXISTS `dm_draft`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_draft` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` int(11) NOT NULL,
  `parent_location_id` int(11) NOT NULL DEFAULT '0',
  `operation` varchar(150) NOT NULL,
  `time` datetime DEFAULT NULL,
  `data` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_draft`
--

LOCK TABLES `dm_draft` WRITE;
/*!40000 ALTER TABLE `dm_draft` DISABLE KEYS */;
/*!40000 ALTER TABLE `dm_draft` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_folder`
--

DROP TABLE IF EXISTS `dm_folder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_folder` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `folder_type` varchar(30) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `author` int(11) DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_folder`
--

LOCK TABLES `dm_folder` WRITE;
/*!40000 ALTER TABLE `dm_folder` DISABLE KEYS */;
INSERT INTO `dm_folder` VALUES (1,'','Content','',1560534265,1560534265,'bk1tpudi6ekibbmo2ch0',0,1,0),(2,'','Users','',1560534277,1560534277,'bk1tq1di6ekibbmo2ci0',0,1,0),(3,'','Demosite','',1560534450,1560534450,'bk1trcli6ekibbmo2cj0',0,1,0),(4,'','Organization','',1560534544,1560534544,'bk1ts45i6ekibbmo2ck0',0,1,0);
/*!40000 ALTER TABLE `dm_folder` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_image`
--

DROP TABLE IF EXISTS `dm_image`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` int(11) NOT NULL DEFAULT '0',
  `imagetype` varchar(30) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `path` varchar(255) NOT NULL DEFAULT '',
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_image`
--

LOCK TABLES `dm_image` WRITE;
/*!40000 ALTER TABLE `dm_image` DISABLE KEYS */;
/*!40000 ALTER TABLE `dm_image` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_location`
--

DROP TABLE IF EXISTS `dm_location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_location` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '-1',
  `main_id` int(11) NOT NULL,
  `hierarchy` varchar(255) NOT NULL DEFAULT '',
  `depth` int(11) NOT NULL DEFAULT '0',
  `identifier_path` varchar(255) NOT NULL DEFAULT '',
  `content_type` varchar(50) NOT NULL DEFAULT '',
  `content_id` int(11) NOT NULL,
  `language` varchar(20) NOT NULL DEFAULT '',
  `author` int(11) NOT NULL DEFAULT '0',
  `name` varchar(50) NOT NULL DEFAULT '',
  `is_hidden` tinyint(1) NOT NULL DEFAULT '0',
  `is_invisible` tinyint(1) NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `uid` varchar(50) NOT NULL DEFAULT '',
  `scope` varchar(30) NOT NULL DEFAULT '',
  `section` varchar(50) NOT NULL DEFAULT '',
  `p` varchar(30) NOT NULL DEFAULT 'c',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_location`
--

LOCK TABLES `dm_location` WRITE;
/*!40000 ALTER TABLE `dm_location` DISABLE KEYS */;
INSERT INTO `dm_location` VALUES (1,0,1,'1',1,'content','folder',1,'',0,'Content',0,0,0,'bk1tpudi6ekibbmo2chg','','public',''),(2,0,2,'2',1,'users','folder',2,'',0,'Users',0,0,0,'bk1tq1di6ekibbmo2cig','','',''),(3,1,3,'1/3',2,'content/demosite','folder',3,'',0,'Demosite',0,0,0,'bk1trcli6ekibbmo2cjg','','public',''),(4,2,4,'2/4',2,'users/organization','folder',4,'',0,'Organization',0,0,0,'bk1ts4di6ekibbmo2ckg','','',''),(5,4,5,'2/4/5',3,'users/organization/administrator-admin','user',1,'',0,'Administrator Admin',0,0,0,'bk1tsc5i6ekibbmo2clg','','',''),(6,4,6,'2/4/6',3,'users/organization/anonymous-user','user',2,'',0,'Anonymous User',0,0,0,'bk1tstli6ekibbmo2cmg','','',''),(7,2,7,'2/7',2,'users/anonymous','role',1,'',0,'Anonymous',0,0,0,'bk1vutti6ekij1eq9sgg','','',''),(8,7,8,'2/7/8',3,'users/anonymous/editor','role',2,'',0,'Editor',0,0,0,'bk2051di6ekislpnehu0','','','');
/*!40000 ALTER TABLE `dm_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_log_operation`
--

DROP TABLE IF EXISTS `dm_log_operation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_log_operation` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL DEFAULT '0',
  `category` varchar(50) NOT NULL DEFAULT '',
  `operation` varchar(255) NOT NULL DEFAULT '',
  `created` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_log_operation`
--

LOCK TABLES `dm_log_operation` WRITE;
/*!40000 ALTER TABLE `dm_log_operation` DISABLE KEYS */;
/*!40000 ALTER TABLE `dm_log_operation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_relation`
--

DROP TABLE IF EXISTS `dm_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_relation` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `to_content_id` int(11) NOT NULL DEFAULT '0',
  `to_type` varchar(50) NOT NULL DEFAULT '',
  `from_content_id` int(11) NOT NULL DEFAULT '0',
  `from_type` varchar(30) NOT NULL DEFAULT '',
  `from_location` int(11) NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `identifier` varchar(50) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `data` text NOT NULL,
  `uid` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_relation`
--

LOCK TABLES `dm_relation` WRITE;
/*!40000 ALTER TABLE `dm_relation` DISABLE KEYS */;
/*!40000 ALTER TABLE `dm_relation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_role`
--

DROP TABLE IF EXISTS `dm_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `author` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_role`
--

LOCK TABLES `dm_role` WRITE;
/*!40000 ALTER TABLE `dm_role` DISABLE KEYS */;
INSERT INTO `dm_role` VALUES (1,'Role Anonymous','',1560543095,1560543095,'bk1vutti6ekij1eq9sg0',0,1),(2,'Editor','',1560543877,1560543877,'bk2051di6ekislpnehtg',0,1);
/*!40000 ALTER TABLE `dm_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_role_policy`
--

DROP TABLE IF EXISTS `dm_role_policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_role_policy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `policy` varchar(100) NOT NULL DEFAULT '',
  `under` varchar(30) NOT NULL DEFAULT '',
  `scope` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_role_policy`
--

LOCK TABLES `dm_role_policy` WRITE;
/*!40000 ALTER TABLE `dm_role_policy` DISABLE KEYS */;
INSERT INTO `dm_role_policy` VALUES (1,8,'editor_visit','',''),(2,8,'edit','','');
/*!40000 ALTER TABLE `dm_role_policy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_user`
--

DROP TABLE IF EXISTS `dm_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `login` varchar(50) DEFAULT NULL,
  `firstname` varchar(50) NOT NULL,
  `lastname` varchar(50) DEFAULT NULL,
  `password` varchar(50) DEFAULT NULL,
  `mobile` varchar(50) DEFAULT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `version` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  `author` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_user`
--

LOCK TABLES `dm_user` WRITE;
/*!40000 ALTER TABLE `dm_user` DISABLE KEYS */;
INSERT INTO `dm_user` VALUES (1,'admin','Administrator','Admin','',NULL,1560534576,1560534576,1,'bk1tsc5i6ekibbmo2cl0',0),(2,'anonymous','Anonymous','User','',NULL,1560534645,1560534645,1,'bk1tstdi6ekibbmo2cm0',0);
/*!40000 ALTER TABLE `dm_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_user_role`
--

DROP TABLE IF EXISTS `dm_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_user_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_user_role`
--

LOCK TABLES `dm_user_role` WRITE;
/*!40000 ALTER TABLE `dm_user_role` DISABLE KEYS */;
INSERT INTO `dm_user_role` VALUES (1,1,8);
/*!40000 ALTER TABLE `dm_user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_user_usergroup`
--

DROP TABLE IF EXISTS `dm_user_usergroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_user_usergroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `usergroup_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_user_usergroup`
--

LOCK TABLES `dm_user_usergroup` WRITE;
/*!40000 ALTER TABLE `dm_user_usergroup` DISABLE KEYS */;
INSERT INTO `dm_user_usergroup` VALUES (4,2,7);
/*!40000 ALTER TABLE `dm_user_usergroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_usergroup`
--

DROP TABLE IF EXISTS `dm_usergroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_usergroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `cuid` varchar(30) NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `author` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_usergroup`
--

LOCK TABLES `dm_usergroup` WRITE;
/*!40000 ALTER TABLE `dm_usergroup` DISABLE KEYS */;
INSERT INTO `dm_usergroup` VALUES (1,'Anonymous','',1560543095,1560543095,'bk1vutti6ekij1eq9sg0',0,0),(2,'Editor','',1560543877,1560543877,'bk2051di6ekislpnehtg',0,0);
/*!40000 ALTER TABLE `dm_usergroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_usergroup_policy`
--

DROP TABLE IF EXISTS `dm_usergroup_policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_usergroup_policy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `usergroup_id` int(11) NOT NULL,
  `policy` varchar(100) NOT NULL DEFAULT '',
  `under` varchar(30) NOT NULL DEFAULT '',
  `scope` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_usergroup_policy`
--

LOCK TABLES `dm_usergroup_policy` WRITE;
/*!40000 ALTER TABLE `dm_usergroup_policy` DISABLE KEYS */;
INSERT INTO `dm_usergroup_policy` VALUES (1,7,'read','','');
/*!40000 ALTER TABLE `dm_usergroup_policy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_version`
--

DROP TABLE IF EXISTS `dm_version`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dm_version` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content_type` varchar(50) NOT NULL,
  `content_id` int(11) NOT NULL,
  `version` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `author` int(11) NOT NULL DEFAULT '0',
  `data` longtext NOT NULL,
  `location_id` int(11) NOT NULL DEFAULT '0',
  `created` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_version`
--

LOCK TABLES `dm_version` WRITE;
/*!40000 ALTER TABLE `dm_version` DISABLE KEYS */;
INSERT INTO `dm_version` VALUES (1,'user',9,1,0,0,'{\"id\":9,\"version\":1,\"published\":1560534576,\"modified\":1560534576,\"cuid\":\"bk1tsc5i6ekibbmo2cl0\",\"relations\":{\"list\":null},\"firstname\":{\"data\":\"Administrator\"},\"lastname\":{\"data\":\"Admin\"},\"login\":{\"data\":\"admin\"},\"password\":{\"data\":\"\"},\"location\":{\"id\":0,\"parent_id\":0,\"main_id\":0,\"identifier_path\":\"\",\"hierarchy\":\"\",\"depth\":0,\"content_type\":\"\",\"content_id\":0,\"language\":\"\",\"author\":0,\"name\":\"\",\"is_hidden\":false,\"is_invisible\":false,\"priority\":0,\"uid\":\"\",\"section\":\"\",\"p\":\"\"}}',0,0),(2,'user',10,1,0,0,'{\"id\":10,\"version\":1,\"published\":1560534645,\"modified\":1560534645,\"cuid\":\"bk1tstdi6ekibbmo2cm0\",\"relations\":{\"list\":null},\"firstname\":{\"data\":\"Anonymous\"},\"lastname\":{\"data\":\"User\"},\"login\":{\"data\":\"anonymous\"},\"password\":{\"data\":\"\"},\"location\":{\"id\":0,\"parent_id\":0,\"main_id\":0,\"identifier_path\":\"\",\"hierarchy\":\"\",\"depth\":0,\"content_type\":\"\",\"content_id\":0,\"language\":\"\",\"author\":0,\"name\":\"\",\"is_hidden\":false,\"is_invisible\":false,\"priority\":0,\"uid\":\"\",\"section\":\"\",\"p\":\"\"}}',0,0);
/*!40000 ALTER TABLE `dm_version` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-05-10 12:22:54
