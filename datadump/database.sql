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
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` text NOT NULL,
  `body` text NOT NULL,
  `published` int(11) NOT NULL DEFAULT '0',
  `modified` int(11) NOT NULL DEFAULT '0',
  `remote_id` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_article`
--

LOCK TABLES `dm_article` WRITE;
/*!40000 ALTER TABLE `dm_article` DISABLE KEYS */;
INSERT INTO `dm_article` VALUES (1,1,1,'Welcome','','test2019-04-21 21:15:03.379424 +0200 CEST m=+0.119375313',111111,12313,'bit31n5i6eko5oe5fo9g'),(2,0,0,'Welcome','','',0,0,'bisdku5i6ekklfcg6sf0'),(3,1,1,'','','',0,0,''),(4,1,1,'','','',0,231213,''),(5,1,1,'','','',0,0,''),(6,0,0,'','','',0,0,'5555555'),(7,0,1,'','','',1555160703,1555160703,''),(8,0,1,'','','',1555160780,1555160780,''),(9,0,1,'','','',1555163008,1555163008,''),(10,0,1,'','','',1555167726,1555167726,''),(11,0,1,'','','',1555172650,1555172650,''),(12,0,1,'','','',1555172678,1555172678,''),(13,0,0,'','','',0,0,'bis7ehti6ekh9l3ahps0'),(14,0,0,'','','',0,0,'bisdi55i6ekkgnv0v0mg'),(15,0,0,'Welcome','','test2019-04-21 21:11:36.811602 +0200 CEST m=+0.295915534',111111,12313,'biuc2dti6ekjbfn7atk0'),(16,0,0,'','','',0,5555555,'');
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
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_folder`
--

LOCK TABLES `dm_folder` WRITE;
/*!40000 ALTER TABLE `dm_folder` DISABLE KEYS */;
INSERT INTO `dm_folder` VALUES (1,'','Home','',0,0,''),(2,'','Blog','',0,0,''),(3,'','News','',0,0,''),(4,'','Contact Us','',0,0,''),(5,'','Users','',0,0,''),(6,'','Share Content','',0,0,'master-50');
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
  `content_type` varchar(50) NOT NULL DEFAULT '',
  `content_id` int(11) NOT NULL,
  `language` varchar(20) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL DEFAULT '',
  `is_hidden` tinyint(1) NOT NULL DEFAULT '0',
  `is_invisible` tinyint(1) NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `uid` varchar(50) NOT NULL DEFAULT '',
  `section` varchar(50) NOT NULL DEFAULT '',
  `p` varchar(30) NOT NULL DEFAULT 'c',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_location`
--

LOCK TABLES `dm_location` WRITE;
/*!40000 ALTER TABLE `dm_location` DISABLE KEYS */;
INSERT INTO `dm_location` VALUES (1,0,1,'1','folder',1,'eng-gb','Home',0,0,0,'','content','c'),(2,1,2,'1/2','article',1,'eng-gb','Welcome',0,0,0,'bin1oj5i6ekglgsliq4g','content','c'),(3,1,3,'1/3','folder',2,'eng-gb','Blog',0,0,0,'','content','c'),(4,1,4,'1/4','folder',3,'eng-gb','News',0,0,0,'','content','recent'),(5,1,5,'1/5','folder',4,'eng-gb','Contact Us',0,0,0,'','content','archive'),(6,4,6,'1/4/6','article',2,'eng-gb','EM EMF is released',0,0,0,'','content','c'),(7,0,7,'7','folder',5,'eng-gb','Users',0,0,0,'','user','c'),(8,7,8,'7/8','user',1,'eng-gb','Chen Xiongjie',0,0,0,'','user','c'),(9,0,9,'9','folder',9,'eng-gb','Shared Contents',0,0,0,'','',''),(10,0,10,'','folder',9,'eng-gb','Test folder',0,0,0,'','',''),(11,-1,11,'','folder',0,'eng-gb','Test folder',0,0,0,'','','c'),(12,-1,12,'','folder',0,'eng-gb','Test folder',0,0,0,'','','c'),(14,1,0,'','',0,'','',0,0,0,'','',''),(15,1,0,'','',0,'','',0,0,0,'biottj5i6ekipupsp7b0','',''),(16,1,0,'','',0,'','',0,0,0,'biouf05i6ekj5fs07k20','',''),(17,1,0,'','',0,'','',0,0,0,'biovjrli6ekjf1sm1n6g','',''),(18,1,0,'','',0,'','',0,0,0,'bip0qali6ekk630i8slg','',''),(19,-1,0,'','article',0,'','',0,0,0,'bip0qhli6ekk6rst2630','','');
/*!40000 ALTER TABLE `dm_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dm_relation`
--

DROP TABLE IF EXISTS `dm_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dm_relation` (
  `from_content` int(11) DEFAULT NULL,
  `to_content` int(11) DEFAULT NULL,
  `relation_type` varchar(50) DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `identifier` varchar(50) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `remote_id` varchar(30) DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dm_relation`
--

LOCK TABLES `dm_relation` WRITE;
/*!40000 ALTER TABLE `dm_relation` DISABLE KEYS */;
INSERT INTO `dm_relation` VALUES (5,10,'image',1,'cover_image','Profile picture','');
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

-- Dump completed on 2019-04-21 21:16:49
