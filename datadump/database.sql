-- MySQL dump 10.13  Distrib 5.7.25, for Linux (x86_64)
--
-- Host: 185.35.187.91    Database: dev_emf
-- ------------------------------------------------------
-- Server version	5.7.25-0ubuntu0.16.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
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
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_article` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `status` int(11) NOT NULL DEFAULT '0',
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
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_article`
--

LOCK TABLES `dm_article` WRITE;
/*!40000 ALTER TABLE `dm_article` DISABLE KEYS */;
INSERT INTO `dm_article` VALUES (1,1,1,'','Welcome','','test2019-04-28 23:36:26.242035 +0200 CEST m=+0.116724089',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(2,0,0,'','Welcome','','',NULL,NULL,0,0,'bisdku5i6ekklfcg6sf0'),(3,1,1,'','','','',NULL,NULL,0,0,''),(4,1,1,'','','','',NULL,NULL,0,231213,''),(5,1,1,'','','','',NULL,NULL,0,0,''),(6,0,0,'','','','',NULL,NULL,0,0,'5555555'),(7,0,1,'','','','',NULL,NULL,1555160703,1555160703,''),(8,0,1,'','','','',NULL,NULL,1555160780,1555160780,''),(9,0,1,'','','','',NULL,NULL,1555163008,1555163008,''),(10,0,1,'','','','',NULL,NULL,1555167726,1555167726,''),(11,0,1,'','','','',NULL,NULL,1555172650,1555172650,''),(12,0,1,'','','','',NULL,NULL,1555172678,1555172678,''),(13,0,0,'','','','',NULL,NULL,0,0,'bis7ehti6ekh9l3ahps0'),(14,0,0,'','','','',NULL,NULL,0,0,'bisdi55i6ekkgnv0v0mg'),(15,0,0,'','Welcome','','test2019-04-21 21:11:36.811602 +0200 CEST m=+0.295915534',NULL,NULL,111111,12313,'biuc2dti6ekjbfn7atk0'),(16,0,0,'','','','',NULL,NULL,0,5555555,''),(17,0,0,'','','','',NULL,NULL,1555876380,1555876380,'biuck75i6ekkfim1fpa0'),(21,0,0,'','Test','','',NULL,NULL,1555877658,1555877658,'biucu6li6ekl32m1osp0'),(22,0,0,'','Test21.04.2019 22:25','','',NULL,NULL,1555878336,1555878336,'biud3g5i6eklll6a37og'),(23,0,0,'','Test 21.04.2019 22:25','','',NULL,NULL,1555878350,1555878350,'biud3jli6eklm8ldajr0'),(24,0,0,'','Test 21.04.2019 22:28','','Hello',NULL,NULL,1555878483,1555878483,'biud4kti6eklqplqrd6g'),(25,0,0,'','Test 22.04.2019 11:56','','Hello',NULL,NULL,1555927010,1555927010,'biuovoli6eko05jjm38g'),(26,0,0,'','Test 22.04.2019 12:04','','Hello',NULL,NULL,1555927488,1555927488,'biup3g5i6eko66e9n4mg'),(30,0,0,'','Test 22.04.2019 14:11','','Hello',NULL,NULL,1555935086,1555935086,'biuqurli6ekgrts6st2g'),(31,0,0,'','Test 22.04.2019 14:25','','Hello',NULL,NULL,1555935929,1555935929,'biur5edi6ekh8f72bo7g'),(32,0,0,'','Title1','','Hello world',NULL,NULL,1555937684,1555937684,'biurj55i6ekhu6c75ddg'),(33,0,0,'','Good1','','Test',NULL,NULL,1555937783,1555937783,'biurjtti6ekhu6c75deg'),(34,0,0,'','This is good','','This is nice. :)',NULL,NULL,1555938679,1555938679,'biurqtti6eki3if9odf0'),(35,0,0,'','This is nice news','','Hello world!',NULL,NULL,1555938869,1555938869,'biursda23akqmj2f50jg'),(36,0,0,'','Hello world','','Hello world body....',NULL,NULL,1555939459,1555939459,'bius10q23akqkp7g55l0'),(37,0,0,'','This is created from mobile','','Hjnvhknssw',NULL,NULL,1555941765,1555941765,'biusj1a23akqkp7g55m0'),(38,0,0,'{\"name\":\"chen\", \"email\":\"chen@digimaker.no\"}','This is nice','','This is nice...',NULL,NULL,1555956628,1555956628,'biv075223akkdeldlaa0'),(39,0,0,'','Test','','This is the body text',NULL,NULL,1556010895,1556010895,'bivdf3q23akl47h33nj0'),(40,0,0,'','Welcome','','test2019-04-25 21:17:54.254423 +0200 CEST m=+0.114920807',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(41,0,0,'','Welcome','','test2019-04-25 21:19:41.86809 +0200 CEST m=+0.115214177',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(42,0,0,'','Welcome','','test2019-04-25 21:27:05.720371 +0200 CEST m=+0.118570830',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(43,0,0,'','Welcome','','test2019-04-25 21:37:36.70084 +0200 CEST m=+0.125834685',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(44,0,0,'','Welcome','','test2019-04-25 21:42:00.007948 +0200 CEST m=+0.150465957',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(45,0,0,'','Welcome','','test2019-04-26 20:22:58.858131 +0200 CEST m=+0.124329785',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(46,0,0,'','Welcome','','test2019-04-26 20:39:29.818206 +0200 CEST m=+0.121351436',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(47,0,0,'','Welcome','','test2019-04-26 20:48:49.455475 +0200 CEST m=+0.152287431',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0'),(48,0,0,'','Welcome','','test2019-04-26 20:49:24.199451 +0200 CEST m=+0.147711459',NULL,NULL,111111,12313,'biud4cti6eklq51ea0p0');
/*!40000 ALTER TABLE `dm_article` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_folder`
--

DROP TABLE IF EXISTS `dm_folder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_folder` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `folder_type` varchar(30) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `remote_id` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_folder`
--

LOCK TABLES `dm_folder` WRITE;
/*!40000 ALTER TABLE `dm_folder` DISABLE KEYS */;
INSERT INTO `dm_folder` VALUES (1,'','Home','',0,0,''),(2,'','Blog','',0,0,''),(3,'','News','',0,0,''),(4,'','Contact Us','',0,0,''),(5,'','Users','',0,0,''),(6,'','Share Content','',0,0,'master-50'),(7,'','Test 22.04.2019 14:28','Hello',1555936081,1555936081,'biur6kdi6ekh9pi7ku3g'),(8,'','Test 22.04.2019 14:53','Hello',1555937595,1555937595,'biurieti6ekhrsfq478g'),(9,'','Test 22.04.2019 15:43','Hello',1555940637,1555940637,'biusa7di6eki608cromg'),(10,'','Test folder','Folder',1555960873,1555960873,'biv18adi6ekji63i8nh0'),(11,'','Test','Test1',1555960929,1555960929,'biv18odi6ekjir18f450');
/*!40000 ALTER TABLE `dm_folder` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_location`
--

DROP TABLE IF EXISTS `dm_location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_location` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '-1',
  `main_id` int(11) NOT NULL,
  `hierarchy` varchar(255) NOT NULL DEFAULT '',
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
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_location`
--

LOCK TABLES `dm_location` WRITE;
/*!40000 ALTER TABLE `dm_location` DISABLE KEYS */;
INSERT INTO `dm_location` VALUES (1,0,1,'1','','folder',1,'eng-gb',0,'Home',0,0,0,'','','content','c'),(2,1,2,'1/2','','article',1,'eng-gb',0,'Welcome',0,0,0,'bin1oj5i6ekglgsliq4g','','content','c'),(3,1,3,'1/3','','folder',2,'eng-gb',0,'Blog',0,0,0,'','','content','c'),(4,1,4,'1/4','','folder',3,'eng-gb',0,'News',0,0,0,'','','content','recent'),(5,1,5,'1/5','','folder',4,'eng-gb',0,'Contact Us',0,0,0,'','','content','archive'),(6,4,6,'1/4/6','','article',2,'eng-gb',0,'EM EMF is released',0,0,0,'','','content','c'),(7,0,7,'7','','folder',5,'eng-gb',0,'Users',0,0,0,'','','user','c'),(8,7,8,'7/8','','user',1,'eng-gb',0,'Chen Xiongjie',0,0,0,'','','user','c'),(11,-1,11,'','','folder',0,'eng-gb',0,'Test folder',0,0,0,'','','','c'),(12,-1,12,'','','folder',0,'eng-gb',0,'Test folder',0,0,0,'','','','c'),(14,1,0,'','','',0,'',0,'',0,0,0,'','','',''),(15,1,0,'','','',0,'',0,'',0,0,0,'biottj5i6ekipupsp7b0','','',''),(16,1,0,'','','',0,'',0,'',0,0,0,'biouf05i6ekj5fs07k20','','',''),(17,1,0,'','','',0,'',0,'',0,0,0,'biovjrli6ekjf1sm1n6g','','',''),(18,1,0,'','','',0,'',0,'',0,0,0,'bip0qali6ekk630i8slg','','',''),(19,-1,0,'','','article',0,'',0,'',0,0,0,'bip0qhli6ekk6rst2630','','',''),(20,4,0,'','','article',0,'',0,'',0,0,0,'biuck75i6ekkfim1fpag','','',''),(24,4,0,'','','article',21,'',0,'Test',0,0,0,'biucu6li6ekl32m1ospg','','',''),(25,4,0,'','','article',22,'',0,'Test21.04.2019 22:25',0,0,0,'biud3g5i6eklll6a37p0','','',''),(26,4,0,'','','article',23,'',0,'Test 21.04.2019 22:25',0,0,0,'biud3jli6eklm8ldajrg','','',''),(27,4,0,'','','article',24,'',0,'Test 21.04.2019 22:28',0,0,0,'biud4kti6eklqplqrd70','','',''),(28,4,0,'','','article',0,'',0,'',0,0,0,'biuovoli6eko05jjm390','','',''),(33,4,0,'','','article',31,'',0,'Test 22.04.2019 14:25',0,0,0,'biur5edi6ekh8f72bo80','','',''),(34,4,0,'','','folder',7,'',0,'Test 22.04.2019 14:28',0,0,0,'biur6kdi6ekh9pi7ku40','','',''),(35,4,0,'','','folder',8,'',0,'Test 22.04.2019 14:53',0,0,0,'biurieti6ekhrsfq4790','','',''),(36,1,0,'','','article',32,'',0,'Title1',0,0,0,'biurj55i6ekhu6c75de0','','',''),(37,1,0,'','','article',33,'',0,'Good1',0,0,0,'biurjtti6ekhu6c75df0','','',''),(38,1,0,'','','article',34,'',0,'This is good',0,0,0,'biurqu5i6eki3if9odfg','','',''),(39,4,0,'','','article',35,'',0,'This is nice news',0,0,0,'biursda23akqmj2f50k0','','',''),(40,1,0,'','','article',36,'',0,'Hello world',0,0,0,'bius10q23akqkp7g55lg','','',''),(41,4,0,'','','folder',9,'',0,'Test 22.04.2019 15:43',0,0,0,'biusa7di6eki608cron0','','',''),(42,1,0,'','','article',37,'',0,'This is created from mobile',0,0,0,'biusj1a23akqkp7g55mg','','',''),(43,1,0,'','','article',38,'',0,'This is nice',0,0,0,'biv075223akkdeldlaag','','',''),(44,5,0,'','','folder',10,'',0,'Test folder',0,0,0,'biv18adi6ekji63i8nhg','','',''),(45,5,0,'','','folder',11,'',0,'Test',0,0,0,'biv18odi6ekjir18f45g','','',''),(46,4,0,'','home/news/test','article',39,'',0,'Test',0,0,0,'bivdf3q23akl47h33njg','','','');
/*!40000 ALTER TABLE `dm_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_relation`
--

DROP TABLE IF EXISTS `dm_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
  `remote_id` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_relation`
--

LOCK TABLES `dm_relation` WRITE;
/*!40000 ALTER TABLE `dm_relation` DISABLE KEYS */;
INSERT INTO `dm_relation` VALUES (1,10,'article',5,'',0,1,'cover_image','Profile picture','web/images/22x/3ee22/ddxexxxxx.png',''),(3,37,'article',1,'article',2,0,'related_articles','','Welcome',''),(4,37,'article',1,'',0,0,'translation','nor-no','',''),(8,2,'article',0,'',37,10,'related_articles','f','',''),(9,37,'article',0,'',40,10,'related_articles','hello:)','Hello world',''),(10,2,'article',0,'',37,0,'related_articles','Test2019-04-25 22:08:24.937221 +0200 CEST m=+0.077377383','',''),(13,2,'article',0,'',42,0,'related_articles','','37 {This is created from mobile <nil>}',''),(14,2,'article',0,'',42,0,'related_articles','','37 {This is created from mobile <nil>}',''),(15,2,'article',0,'',42,427123440,'related_articles','2019-04-27 12:34:40.415358 +0200 CEST m=+0.079602793','37 {This is created from mobile <nil>}',''),(16,2,'article',0,'',42,427123629,'related_articles','2019-04-27 12:36:29.588874 +0200 CEST m=+0.075406532','37 {This is created from mobile <nil>}',''),(19,0,'',0,'',0,0,'','Test1','',''),(21,0,'',0,'',0,0,'','Test2','',''),(22,0,'',0,'',0,0,'','Test2','',''),(25,0,'',0,'',0,0,'','Test1','',''),(26,0,'',0,'',0,0,'','Test2','','');
/*!40000 ALTER TABLE `dm_relation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_user`
--

DROP TABLE IF EXISTS `dm_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `login` varchar(50) DEFAULT NULL,
  `firstname` varchar(50) NOT NULL,
  `lastname` varchar(50) DEFAULT NULL,
  `password` varchar(50) DEFAULT NULL,
  `mobile` varchar(50) DEFAULT NULL,
  `remote_id` varchar(30) DEFAULT '',
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_user`
--

LOCK TABLES `dm_user` WRITE;
/*!40000 ALTER TABLE `dm_user` DISABLE KEYS */;
INSERT INTO `dm_user` VALUES (1,'chen','Chen','Xiongjie','fdsafasfiifhsdf23131','+4796888261','',0,0);
/*!40000 ALTER TABLE `dm_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_version`
--

DROP TABLE IF EXISTS `dm_version`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_version` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL,
  `content_id` int(11) NOT NULL,
  `version` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `author` int(11) NOT NULL DEFAULT '0',
  `data` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_version`
--

LOCK TABLES `dm_version` WRITE;
/*!40000 ALTER TABLE `dm_version` DISABLE KEYS */;
INSERT INTO `dm_version` VALUES (1,'article',1,1,0,0,'');
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

-- Dump completed on 2019-04-29 21:50:31
