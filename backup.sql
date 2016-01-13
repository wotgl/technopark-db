-- MySQL dump 10.13  Distrib 5.5.46, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: mydb
-- ------------------------------------------------------
-- Server version	5.5.46-0ubuntu0.14.04.2

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
-- Table structure for table `follow`
--

DROP TABLE IF EXISTS `follow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `follow` (
  `follower` varchar(255) NOT NULL,
  `followee` varchar(255) NOT NULL,
  PRIMARY KEY (`followee`,`follower`),
  KEY `fk_follow_1_idx` (`follower`),
  CONSTRAINT `fk_follow_followee` FOREIGN KEY (`followee`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_follow_follower` FOREIGN KEY (`follower`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `follow`
--

LOCK TABLES `follow` WRITE;
/*!40000 ALTER TABLE `follow` DISABLE KEYS */;
INSERT INTO `follow` VALUES ('02mpbei@ya.ru','x9bhu@ya.ru'),('6gm0u4set@ya.ru','6xvn9e27@ya.ru'),('6gm0u4set@ya.ru','cg962ap1ht@gmail.com'),('6xvn9e27@ya.ru','6gm0u4set@ya.ru'),('6xvn9e27@ya.ru','cg962ap1ht@gmail.com'),('6xvn9e27@ya.ru','ty@gmail.com'),('cg962ap1ht@gmail.com','6gm0u4set@ya.ru'),('cg962ap1ht@gmail.com','xtkmq@gmail.com'),('lad0h8vmw2@yahoo.com','x9bhu@ya.ru'),('lad0h8vmw2@yahoo.com','xtkmq@gmail.com'),('n1szohqi@ua.ru','lad0h8vmw2@yahoo.com'),('n1szohqi@ua.ru','o7kz@gmail.com'),('n1szohqi@ua.ru','ty@gmail.com'),('o7kz@gmail.com','lad0h8vmw2@yahoo.com'),('o7kz@gmail.com','o7kz@gmail.com'),('ty@gmail.com','6gm0u4set@ya.ru'),('ty@gmail.com','6xvn9e27@ya.ru'),('ty@gmail.com','n1szohqi@ua.ru'),('ty@gmail.com','ty@gmail.com'),('x9bhu@ya.ru','6xvn9e27@ya.ru'),('x9bhu@ya.ru','n1szohqi@ua.ru'),('x9bhu@ya.ru','ty@gmail.com'),('x9bhu@ya.ru','xtkmq@gmail.com'),('xtkmq@gmail.com','02mpbei@ya.ru'),('xtkmq@gmail.com','lad0h8vmw2@yahoo.com');
/*!40000 ALTER TABLE `follow` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `forum`
--

DROP TABLE IF EXISTS `forum`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `forum` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `short_name` varchar(255) NOT NULL,
  `user` varchar(255) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`short_name`),
  UNIQUE KEY `short_name_UNIQUE` (`short_name`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `user_idx` (`user`),
  CONSTRAINT `fk_forum_user` FOREIGN KEY (`user`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3204 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `forum`
--

LOCK TABLES `forum` WRITE;
/*!40000 ALTER TABLE `forum` DISABLE KEYS */;
INSERT INTO `forum` VALUES (3164,'qm3wvcgbd8 qozkgbp upfl','0bw4tlqgxn','o7kz@gmail.com','2016-01-12 19:31:53'),(3181,'ntwdl 9rxf0n52zg','3l85o7g06h','02mpbei@ya.ru','2016-01-12 19:52:21'),(3168,'whuqc4xogb rtyfp dt7ozmbw1','5yi','6xvn9e27@ya.ru','2016-01-12 19:31:53'),(3160,'gqxn 5e792atx8 85kmlf93 g91kz','6qsxa71p','n1szohqi@ua.ru','2016-01-12 19:31:53'),(3191,'sd5azlv4 ehpywrk5 1gupc','870653ao','lad0h8vmw2@yahoo.com','2016-01-12 19:52:42'),(3174,'i4wh6','8hrd','lad0h8vmw2@yahoo.com','2016-01-12 19:52:11'),(3180,'9y8wl4 az2rn0vc','9vgxt2b','ty@gmail.com','2016-01-12 19:52:21'),(3176,'o21 lxe4hd g7watd2i','ar97','o7kz@gmail.com','2016-01-12 19:52:13'),(3192,'czhwo0','aty1cq0u5k','6gm0u4set@ya.ru','2016-01-12 19:52:43'),(3166,'m1 an8gh6k','b6kgh4gur','lad0h8vmw2@yahoo.com','2016-01-12 19:31:53'),(3173,'so817wrc sa 6nh x7b9puk24','bfly','6xvn9e27@ya.ru','2016-01-12 19:52:11'),(3178,'35nzug swg1r','d94','cg962ap1ht@gmail.com','2016-01-12 19:52:17'),(3182,'ylgp7ov58','dc','lad0h8vmw2@yahoo.com','2016-01-12 19:52:25'),(3184,'p48 0sxta','dz','ty@gmail.com','2016-01-12 19:52:28'),(3183,'8f9 cgfkz waukz 5zsfhg','f39','cg962ap1ht@gmail.com','2016-01-12 19:52:27'),(3190,'zighpo6 8o2x34v om8leh','fuwv','lad0h8vmw2@yahoo.com','2016-01-12 19:52:42'),(3200,'pr80vc4ymn w0i6mgkl gly2tc','gay6','6gm0u4set@ya.ru','2016-01-12 19:54:07'),(3186,'8r4m9iq 7gngkl 9dnfcp4yil 9grba0m6v','gp','x9bhu@ya.ru','2016-01-12 19:52:35'),(3193,'6lmueyzx r4 uwc vt','h3nkcx','02mpbei@ya.ru','2016-01-12 19:52:44'),(3177,'oq','hc','02mpbei@ya.ru','2016-01-12 19:52:16'),(3169,'eop8 k0','hm7aurq64s','6gm0u4set@ya.ru','2016-01-12 19:31:53'),(3171,'5rm9n4gy8c goi gr3z 3x4y6','if','xtkmq@gmail.com','2016-01-12 19:52:06'),(3188,'wvf5gl bgg8z','lg239s','o7kz@gmail.com','2016-01-12 19:52:39'),(3187,'sw7tpl62u','lx0gb','ty@gmail.com','2016-01-12 19:52:37'),(3175,'zubteq0','m84dy27','6xvn9e27@ya.ru','2016-01-12 19:52:12'),(3194,'f4lw 3a0 ev o0','n6ps8vzf7','6gm0u4set@ya.ru','2016-01-12 19:52:45'),(3161,'m78z dygspme46 ev7mnd5oas','nuexq','6gm0u4set@ya.ru','2016-01-12 19:31:53'),(3167,'8r1fip','o5gu7klc9','lad0h8vmw2@yahoo.com','2016-01-12 19:31:53'),(3163,'cwom','oxg5','x9bhu@ya.ru','2016-01-12 19:31:53'),(3189,'pd2i56u8 g0l','pmslo3xngy','6gm0u4set@ya.ru','2016-01-12 19:52:39'),(3195,'gtpnrdsqiw kgry4s8gxm d8gi som','q4pshov2','02mpbei@ya.ru','2016-01-12 19:52:45'),(3172,'kw vxrf587a p6b572r e9y','qgmhds93','cg962ap1ht@gmail.com','2016-01-12 19:52:11'),(3197,'lb4gxmu 9h7gw4v 3foi rpu4','rg','02mpbei@ya.ru','2016-01-12 19:52:51'),(3165,'9ztvhn3i7','se','x9bhu@ya.ru','2016-01-12 19:31:53'),(3185,'es8wz','sxyn8pqra','xtkmq@gmail.com','2016-01-12 19:52:32'),(3196,'k1neih0ly','vk428uolnf','x9bhu@ya.ru','2016-01-12 19:52:46'),(3162,'qe6c1ib m2so0ypc','vkyi1wo','6gm0u4set@ya.ru','2016-01-12 19:31:53'),(3170,'o3 baq4','xu8','o7kz@gmail.com','2016-01-12 19:52:02'),(3179,'dicqapt6v','yb2','x9bhu@ya.ru','2016-01-12 19:52:17');
/*!40000 ALTER TABLE `forum` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `objects`
--

DROP TABLE IF EXISTS `objects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `objects` (
  `id_obj` int(10) NOT NULL AUTO_INCREMENT,
  `name_obj` varchar(40) CHARACTER SET cp1251 COLLATE cp1251_bin NOT NULL,
  `price_obj` int(10) NOT NULL,
  `description_obj` text NOT NULL,
  `amount_obj` int(3) NOT NULL,
  `class_obj` varchar(3) NOT NULL,
  `image_obj` varchar(40) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (`id_obj`)
) ENGINE=InnoDB DEFAULT CHARSET=cp1251;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `objects`
--

LOCK TABLES `objects` WRITE;
/*!40000 ALTER TABLE `objects` DISABLE KEYS */;
/*!40000 ALTER TABLE `objects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `orders` (
  `id_ord` int(10) NOT NULL AUTO_INCREMENT,
  `name_ord` varchar(40) NOT NULL,
  `status_ord` varchar(10) NOT NULL,
  `id_object` int(11) DEFAULT NULL,
  `name_object` varchar(40) NOT NULL,
  `amount_ord` int(10) NOT NULL,
  `total_ord` int(10) NOT NULL,
  `id_check` int(10) NOT NULL,
  `date` datetime NOT NULL,
  `table` int(11) DEFAULT NULL,
  `restaraunt` int(11) DEFAULT NULL,
  PRIMARY KEY (`id_ord`),
  KEY `fk_orders_1_idx` (`id_object`),
  KEY `fk_orders_2_idx` (`restaraunt`),
  KEY `fk_orders_3_idx` (`table`),
  CONSTRAINT `fk_orders_1` FOREIGN KEY (`id_object`) REFERENCES `objects` (`id_obj`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_orders_2` FOREIGN KEY (`restaraunt`) REFERENCES `restaraunts` (`id_restaraunt`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_orders_3` FOREIGN KEY (`table`) REFERENCES `tables` (`id_table`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post`
--

DROP TABLE IF EXISTS `post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `post` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '		',
  `thread` int(11) NOT NULL,
  `message` text NOT NULL,
  `user` varchar(255) NOT NULL,
  `forum` varchar(255) NOT NULL,
  `parent` varchar(255) DEFAULT '0',
  `isApproved` tinyint(1) NOT NULL DEFAULT '0',
  `isHighlighted` tinyint(1) NOT NULL DEFAULT '0',
  `isEdited` tinyint(1) NOT NULL DEFAULT '0',
  `isSpam` tinyint(1) NOT NULL DEFAULT '0',
  `isDeleted` tinyint(1) NOT NULL DEFAULT '0',
  `likes` int(11) NOT NULL DEFAULT '0',
  `dislikes` int(11) NOT NULL DEFAULT '0',
  `points` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fk_post_user_idx` (`user`),
  KEY `fk_post_forum_idx` (`forum`),
  KEY `fk_post_thread_idx` (`thread`),
  CONSTRAINT `fk_post_forum` FOREIGN KEY (`forum`) REFERENCES `forum` (`short_name`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_thread` FOREIGN KEY (`thread`) REFERENCES `thread` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_user` FOREIGN KEY (`user`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3949 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post`
--

LOCK TABLES `post` WRITE;
/*!40000 ALTER TABLE `post` DISABLE KEYS */;
INSERT INTO `post` VALUES (3922,'2013-11-17 13:39:39',3162,'213wezlca6 g71oba6gn qf7 wfvo5h uw2ikeg 6h2 rl5bwv e8ybm v0g6 2nfb601 4vg9i2dwh aql6m379b znu02a 9eoygstzb8 shg g0c9wfgk 8igmdh 0gyen 7ztbm g64wp 09m wf8xuv pve1i nog8spe r8nc7fvez6 k1at gzoub gait89ucyd','ty@gmail.com','6qsxa71p','   I;',1,1,0,0,1,3,6,-3),(3923,'2013-11-19 18:49:43',3161,'nyxidl zsvouxgm rncw gu o6hbrmgc0n ryxowst6v4 fn3t4 afkl1edzx hq tue5kgqho ocpfh6dint g7 gv6rza abhd6 e94gmgoy 1xzihtg0r wgaviplku winpxsg g26qsfoa4 lui06w2mbf not p4be5isolt q0x85td 8zhsq','cg962ap1ht@gmail.com','o5gu7klc9','   I<',1,1,1,1,0,3,3,0),(3924,'2013-08-02 13:27:04',3165,'zdpi3cmo8 o2ruy5xv px18qk fcytgs ggalc2v bq8kvwag 9yr rx1wq lbt2 e2r3dnq nv3goswm dxs90tmehc 42sgv5ifq re2oxu u9mk5w3enx 0pvzm2ugn gboe7f h0gkdcq afmx dm2h870gl clo 70t59 21m8npt nmrsp6 fnym','02mpbei@ya.ru','vkyi1wo','   I=',1,0,1,1,0,2,1,1),(3925,'2013-02-23 03:25:37',3163,'75bop2 0o7 xgr9pyi a6x74 8ovuk3f o4 ckw438 gkio 4f7gq69k qh7 7wo24rx owd 9x748i1mql gaxc9y ogi806n 01kil6br zqf phaz8709gx i435 37gr6oy4zi nu 4unv20 pzcvm14y','ty@gmail.com','oxg5','   I>',1,1,0,1,0,3,0,3),(3926,'2013-10-07 06:32:58',3166,'yhltqif6 q3hes0piv wi wk6 yd1mkt ib13qw9 0it9qbf i7m srxbinywg3 kcsh3vq enr2 bag84i u7pmwgkz shg6 84nk rd0c4tuns7 ixuq2sz cp3lhws89 usl6q3 uozdqxefi 62 5plu7gn3 ie9xvs72dz 8bg5rcenos cadoxumbq wnlhb06zsg 4mba0','xtkmq@gmail.com','o5gu7klc9','   I?',1,1,1,0,0,5,0,5),(3927,'2013-02-16 14:11:15',3166,'50vqs26h lx92q vp76t 8kmntvf b6 qkbxt7g8e 9t 37s8xv1q g2 7tqvkmd0 tqoah0fc xr9q y7caugfoq lk im2uedhn hk p2zvdq8rm 61 lz3o8esry xvzy2g41o xnu zc6a1 05dbupyvm1 qhkl1 c9oaq','o7kz@gmail.com','o5gu7klc9','   I@',0,0,0,1,1,0,0,0),(3928,'2013-04-09 03:13:02',3169,'sk6pmw fgvq6 tfe p6 q0sgrg5y emfrgczt 3tvlw670 ge4nyvkl5 gv9ok 9bxr4din c8onrxhz 3gmeaxp4 srqz5x gnaxbt xugl5 fn3bm dc06 pyib1k 2gyx xi80uwnza 4w zt3g19g 41ggdt kw 7mgkhacsqv','o7kz@gmail.com','se','   IA',1,0,1,0,1,1,1,0),(3929,'2013-11-27 02:25:24',3161,'xecg 0d59qmxpwb ntiyl kdtw63yrb 26g dms f4nsuzy ykc39tar 78 gvufepd1 yg3hrv 72l9hrema kn 0wmsty53 s0u6 wch2ugbmg 3wt u8gi1 07gt3 qb1r z13hm ggeil5r4ks e8yrz 5ms6','lad0h8vmw2@yahoo.com','o5gu7klc9','   IB',1,1,0,1,1,0,2,-2),(3930,'2013-01-18 02:35:23',3167,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','6xvn9e27@ya.ru','b6kgh4gur','   IC',1,0,1,0,1,0,3,-3),(3931,'2013-06-13 16:36:10',3170,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','x9bhu@ya.ru','oxg5','   ID',0,1,0,0,1,0,4,-4),(3932,'2013-08-16 15:33:11',3169,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','n1szohqi@ua.ru','se','   IA    !',1,0,0,1,1,0,0,0),(3933,'2013-09-18 15:10:52',3162,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','x9bhu@ya.ru','6qsxa71p','   IF',1,1,0,0,1,0,0,0),(3934,'2013-06-23 08:05:45',3161,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','n1szohqi@ua.ru','o5gu7klc9','   IG',1,1,0,1,0,0,0,0),(3935,'2013-03-14 21:31:42',3169,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','6xvn9e27@ya.ru','se','   IA    \"',0,1,1,0,1,0,0,0),(3936,'2013-07-20 01:08:07',3164,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','lad0h8vmw2@yahoo.com','se','   II',1,1,0,1,1,0,0,0),(3937,'2013-04-04 20:26:43',3162,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','o7kz@gmail.com','6qsxa71p','   IJ',0,1,1,1,1,0,0,0),(3938,'2013-12-03 15:15:15',3170,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','6gm0u4set@ya.ru','oxg5','   IK',1,1,0,0,1,0,0,0),(3939,'2013-08-15 03:25:55',3169,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','lad0h8vmw2@yahoo.com','se','   IL',1,0,0,1,1,0,0,0),(3940,'2013-08-05 12:10:06',3168,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','x9bhu@ya.ru','o5gu7klc9','   IM',0,1,1,0,0,0,0,0),(3941,'2013-07-12 17:13:47',3170,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','6xvn9e27@ya.ru','oxg5','   ID    !',0,1,1,0,1,0,0,0),(3942,'2013-08-03 21:12:42',3161,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','xtkmq@gmail.com','o5gu7klc9','   IO',0,1,1,0,0,0,0,0),(3943,'2013-08-16 15:33:11',3169,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','n1szohqi@ua.ru','se','   IA    #',1,0,0,1,1,0,0,0),(3944,'2013-09-18 15:10:52',3162,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','x9bhu@ya.ru','6qsxa71p','   IQ',1,1,0,0,1,0,0,0),(3945,'2013-05-26 14:49:17',3166,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','cg962ap1ht@gmail.com','o5gu7klc9','   IR',1,1,0,0,1,0,0,0),(3946,'2013-05-26 07:35:50',3164,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','lad0h8vmw2@yahoo.com','se','   IS',0,0,1,1,1,0,0,0),(3947,'2013-01-12 03:02:45',3165,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','6xvn9e27@ya.ru','vkyi1wo','   I=    !',1,0,1,1,1,0,0,0),(3948,'2013-06-23 08:05:45',3161,'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam a lorem a leo porttitor tincidunt eget et urna. Aenean id lacinia dolor. Sed consequat ipsum at orci porta, sed condimentum dui elementum. Curabitur magna purus, sagittis in convallis ultrices, dignissim pharetra ipsum. In molestie, arcu id convallis blandit, felis metus suscipit justo, ut iaculis metus leo viverra felis. Donec a varius dolor. Cras tempor, nisl in dapibus cursus, risus ligula ultricies nisi, a sagittis justo lorem et odio. Mauris eu scelerisque tellus. Duis luctus enim vel porttitor convallis. Phasellus pretium mi vitae ullamcorper pretium. Vivamus sollicitudin, risus a volutpat condimentum','n1szohqi@ua.ru','o5gu7klc9','   IU',1,1,0,1,1,0,0,0);
/*!40000 ALTER TABLE `post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `restaraunts`
--

DROP TABLE IF EXISTS `restaraunts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `restaraunts` (
  `id_restaraunt` int(11) NOT NULL AUTO_INCREMENT,
  `name_rest` varchar(45) DEFAULT NULL,
  `place_rest` varchar(45) DEFAULT NULL,
  `contact_rest` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id_restaraunt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `restaraunts`
--

LOCK TABLES `restaraunts` WRITE;
/*!40000 ALTER TABLE `restaraunts` DISABLE KEYS */;
/*!40000 ALTER TABLE `restaraunts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subscribe`
--

DROP TABLE IF EXISTS `subscribe`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `subscribe` (
  `thread` int(11) NOT NULL,
  `user` varchar(255) NOT NULL,
  PRIMARY KEY (`thread`,`user`),
  KEY `fk_subscribe_user_idx` (`user`),
  CONSTRAINT `fk_subscribe_thread` FOREIGN KEY (`thread`) REFERENCES `thread` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_subscribe_user` FOREIGN KEY (`user`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subscribe`
--

LOCK TABLES `subscribe` WRITE;
/*!40000 ALTER TABLE `subscribe` DISABLE KEYS */;
INSERT INTO `subscribe` VALUES (3162,'02mpbei@ya.ru'),(3170,'6gm0u4set@ya.ru'),(3161,'6xvn9e27@ya.ru'),(3162,'6xvn9e27@ya.ru'),(3163,'cg962ap1ht@gmail.com'),(3167,'cg962ap1ht@gmail.com'),(3170,'cg962ap1ht@gmail.com'),(3161,'lad0h8vmw2@yahoo.com'),(3162,'lad0h8vmw2@yahoo.com'),(3164,'lad0h8vmw2@yahoo.com'),(3166,'n1szohqi@ua.ru'),(3168,'n1szohqi@ua.ru'),(3169,'n1szohqi@ua.ru'),(3165,'o7kz@gmail.com'),(3166,'o7kz@gmail.com'),(3169,'o7kz@gmail.com'),(3170,'o7kz@gmail.com'),(3162,'ty@gmail.com'),(3164,'ty@gmail.com'),(3166,'ty@gmail.com'),(3167,'ty@gmail.com'),(3170,'ty@gmail.com'),(3168,'x9bhu@ya.ru'),(3163,'xtkmq@gmail.com'),(3167,'xtkmq@gmail.com'),(3170,'xtkmq@gmail.com');
/*!40000 ALTER TABLE `subscribe` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tables`
--

DROP TABLE IF EXISTS `tables`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tables` (
  `id_table` int(3) NOT NULL AUTO_INCREMENT,
  `table_name` varchar(20) NOT NULL,
  `reserved` tinyint(1) DEFAULT NULL,
  `count` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id_table`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tables`
--

LOCK TABLES `tables` WRITE;
/*!40000 ALTER TABLE `tables` DISABLE KEYS */;
/*!40000 ALTER TABLE `tables` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `thread`
--

DROP TABLE IF EXISTS `thread`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `thread` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `forum` varchar(255) NOT NULL,
  `title` varchar(45) NOT NULL,
  `isClosed` tinyint(1) NOT NULL,
  `user` varchar(255) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `message` text NOT NULL,
  `slug` varchar(45) NOT NULL,
  `isDeleted` tinyint(1) NOT NULL DEFAULT '0',
  `likes` int(11) NOT NULL DEFAULT '0',
  `dislikes` int(11) NOT NULL DEFAULT '0',
  `points` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fk_user_idx` (`user`),
  KEY `fk_thread_forum` (`forum`),
  CONSTRAINT `fk_thread_forum` FOREIGN KEY (`forum`) REFERENCES `forum` (`short_name`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_thread_user` FOREIGN KEY (`user`) REFERENCES `user` (`email`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3203 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `thread`
--

LOCK TABLES `thread` WRITE;
/*!40000 ALTER TABLE `thread` DISABLE KEYS */;
INSERT INTO `thread` VALUES (3161,'o5gu7klc9','roq67gcv5 hn2r1eiy',1,'n1szohqi@ua.ru','2013-05-15 09:16:00','8eftqc6z gbr68amcq9 kamwqzg7h3 oq9ixmuvwl zktm pgyv2tsfhw c3ryme6 f7sl xqm7 sb3gq atf2e76x 016 vdwz8r36 tlugyk 3h5p7c vfy0c 5v4ex uef e6vxnmiyul g0cdn rgd kg958vps4 z9uwthi6 wd e5v7w6h w5eu','pavmg',0,0,2,-2),(3162,'6qsxa71p','kygbqsut',1,'n1szohqi@ua.ru','2013-06-19 23:26:12','2q63i9n 6fy7b9 zi3vy519 smnq5e4hur cdewzk ynv3az2 u49q5xd or vx e97np 9zu5bsi 3ig9b51mk tvk mus cnk1dr9z g5t6go 9p823rczl utqhlpk ogx py7u3tgc','d7ofyk8nc',1,2,2,0),(3163,'oxg5','gfcbmz02 s7g0lhm1 oxie mvhzsdrcnp',1,'x9bhu@ya.ru','2013-06-23 20:53:32','z16vg c0l2yo3hx hocn7 i12xgdlzw rzih8wskb p2us cgqv2i u0n85idkg l2 fugx d659c84i 4e8bc ig n5s4lc vxro c6yrhqimt zp7vlr ue acer6 h8k3z5 igqayn12','ryeuoi7sh',0,3,1,2),(3164,'se','d1',0,'lad0h8vmw2@yahoo.com','2013-01-20 07:18:53','qo4xp5z79 8m7g zv9wkygg5 m9 ixgqu0w3f1 m96hzisn5 nzyq54v yqo2 fuz5 ag8mztg kup ogedg ud3qyolg72 ioqwz 6x01lq 3w elsrp9k gcy0ns qgcf7mh 2bg5wcr3 rd32q7pt0 g0w 4i0kn6s rswu3 y3f0ulw2rp nc8rsg gp2i49ya','tie4',1,1,3,-2),(3165,'vkyi1wo','vi s2nxiruhg',1,'ty@gmail.com','2013-12-16 16:47:50','zso 7roq3 x6kefoda twqeg cl3x1 b3vg8wr0t cq fqrklsud nthrly21 h04zygb81m pqz s8 wnkx6oesic n1uo evr hqb 78s61 su 1lf8z t6f vqluipwb fbysr7g1 yufxzsh6q peqn2z5 mtis7lnc06','o47hvxkt',0,0,5,-5),(3166,'o5gu7klc9','ngg62 hdqi2tv584 01 bckmhtor4',1,'xtkmq@gmail.com','2013-12-15 00:18:06','xr01wbpq9g wgscxb8af 1x 6b20g rya g5 3wlg7mxe 8uqkr9i7m1 ycgfte sa ia0 oc my s5bzyi x62 04xo7ap1z gw 1yg5 wvrelsmbuc yv74binwx8 wzfg0b fkqe hvnsq to8scd g9w51iy2 s9nzkdt38 w1us8 qdgvan26 pmg','uq4mxl',1,1,1,0),(3167,'b6kgh4gur','ayrg o6hs73y4z s6go7',1,'02mpbei@ya.ru','2013-10-14 17:12:56','iapbft8 0qui3xg1 t1n9z8uwv ypg1sxqwv8 zyni 3mrbufden y6vnzl0m 1ib97x x8g01quls zxfiu a4rlx b1 u7b0p g1 ovqk9rpwx s4d bp 2k 2wd 1m z7osrcb0n','z1q2k',1,0,0,0),(3168,'o5gu7klc9','txrpbl 6c',1,'n1szohqi@ua.ru','2013-10-11 03:31:12','ns4yq0i3b xs5 h1 w7m 4mdgf8 7my ug8s hs3extw42 iq9bvc l8zq7imhk r9 4gcnexbht 8b4 64l7 qakxtr 9td5hqxs ue5c3 rokcnf zlpqac ghcakn0oir 7i1wg5 udnvg5b hq2 gaytxdf b6l2cx7 plg692qv esufq72t u4d 3t dtygq1h','38',0,2,2,0),(3169,'se','vfls8 0wi5k4x8u y8zcprnvhb',1,'o7kz@gmail.com','2013-07-15 16:54:54','8k dcrle9 m8o09dpk4z p39zrs rwtzg 4sm w5oh4b 7kdbpyt2z9 bov mrei98qzy qn50vc4e nbpfumse 0snv8xmw h465 yc9 za62dm4gb yoi35ghx 4sgq2fp 4yd2z0ar qy 8sd1 xypm w6e9510 shu4izmc1x 01svfel u06eax s83k h13dlefo5 kpcudhbe g97','c9zl',1,0,2,-2),(3170,'oxg5','h37 8qgai',1,'x9bhu@ya.ru','2013-10-12 12:51:11','2qrmgiw0 1d60 yk8qigdc 5mkvr0xig8 31nes8 6bh rgxuc18y2q t67e5dr dmn5c nk5sbmx1 p0 hgsc9ai 1lt38a xne91fl23 y8n2 3zo9avg2 70 v6oqzr7b3s cqh8v 2utn7 ktm1xpd0u 0pgkuw3xec 95k27gmx3f tm1eru g4x2qne','a8w',1,0,2,-2),(3171,'o5gu7klc9','3zoqad80g',0,'cg962ap1ht@gmail.com','2013-11-09 16:11:08','dg 9s0bfy y1vc fgh 5dz81rh0a uez12i8fqy t4sqwngo v0p3 fg2mlkd 2xhk130ovb 1l628vbq wytp uhc tsrhuoq5p fnig svg1fg36uk mvruwq8ot7 luaezbgr5 u126 q96he7kbs fpr41dlwg q1idxbz4 n0 aq0pk6y 9m8svqf urb5 imw694','s3fob1a',0,0,0,0),(3172,'5yi','5notflqr1a muelw3',0,'o7kz@gmail.com','2013-04-06 13:24:35','g3xw olvnd9tg4m rggz68e yx8sc9l o68pbcw0z t6 ebt9o2 gn3 ecus5dpmt mczfwok o9zvqxr 53ulwxc9 74 egvk326g b3pih 5bmtongsq4 tn 1h9yqu2g 5tn9exgq2 d20y r0kh6pqlum liszgpb1wn 7xs9 3qw 862te poeduxk32','ngzfs25',0,0,0,0),(3173,'nuexq','oq1fnzray q6dvau7zw kd h6doa',1,'x9bhu@ya.ru','2013-02-15 13:49:30','d5xa36m8 ri mw3gbn 95w klhut1g3r g2i3p 0s x9ref0o482 hyfd906kv 1kistlz6ew 354qk6cvgf 63v 2bszl7uqod az3ndr vo n678 nd1 bo2v 2cg6v7x fh1rq t62 ahnbrs3d1 nuv8f','zgot',1,0,0,0),(3174,'nuexq','wrkix pr2wuz75bq',0,'ty@gmail.com','2013-01-19 10:55:28','3w abfvdktn yw25ktrn flpk4 g7ki3cdpm oef msw 7gbs 93 ih43 sr4q sea2y5zco we36 pu 4lac m2sg3ldy9a ez gu xct62yo 3pyq6 gc13 c3s0w7 yo zu57pbh 07fbs tm shd9afbwz h0 7du69sqza5 s6i3m2no8','hob',0,0,0,0),(3175,'oxg5','zy98holb7r nafsbx40 gvyg3q9ot',1,'6xvn9e27@ya.ru','2013-06-02 21:06:29','z0pf1 a9x 4nftxdy ad h5vmapoc ogsu24 67f9auk2 nlmf5wk2g v0wsghb1z bd1serkolg mtwpy ig9n k9civg327h ecr35g 3ifgdox mhy x94poyud s6 nuwsrfx 423uvf 1zrof zosq0aflyv kov7p','fb',0,0,0,0),(3176,'0bw4tlqgxn','s0815 zgkrqhlw 1sv fwy5hl9r',0,'6xvn9e27@ya.ru','2013-02-02 10:26:46','vs hxruaebip cfo 7smhfq01xg y95 32 g1e q01iflcn yocwifxr8 r8g2lq1os3 lboe 2w78no l69z82 65lxr gi98n6 3qy14ra 9x8g0yu mrwzs h9gyu87a 3zlnosm qtz9mn eyh0d waoerf10 pvnko3lgm1 d4g2bt53 20sbtprd5 8mv','dkq',0,0,0,0),(3177,'se','q82bep r5mizwc8 im3wy 6014',0,'xtkmq@gmail.com','2013-04-03 01:29:40','iov 27dylc8 3ey4gt1z7h 1l09bgc udi 7s6adhuwoe lq2kt1 wn nvc5 y1vciq zxwig5lr czn952xie7 d0xy9a 6wh2i7ugv f9uint ikygme4 uo90f3c 1td2 nw4rzyd 8iq5mnkz3 0svexgbr2 ylu o64syz2t gfl7dk0psr s2gt6nv9 ggk5eyn u9to1yl7 ne754hfvdx 3cyi flo09qb','rmkcnqs5',1,0,0,0),(3178,'nuexq','0seyfqd coye97b k7uy6not3',0,'x9bhu@ya.ru','2013-08-23 13:21:06','b30zld6f 7bm4qrwo9 iev0g481z u7c2mkdl a14ezp6 rtkq gp z978u0 gi ulgfh7rx5e he2b4lnf9s wu xqtp5d 25pwxkgi7 e5flxi n8ilhu fpt lzp7 ivtla 0vzegmpl4o 71z5oh3','0gr4iz1h5',1,0,0,0),(3179,'oxg5','yd29iv',1,'02mpbei@ya.ru','2013-12-02 03:12:16','g5yx7s8g t6dizkc4r ul 4qpmogw beo5dnck vpg7wiln0 5qo0t rnmo6 qod4 lz 6dpwh39 ekf yps mo5csn9 kimwo goka2q83l kg4z31hgvi pzhndfub1m i6ykgz y6ixg3 fy58emg3oh vy rfzgyvd8e w3vecm 0p4 4x 8dw5gnzkah','1lrwscikm',0,0,0,0),(3180,'se','hv5id943s y31au9khe',1,'ty@gmail.com','2013-08-01 06:30:34','gpsc vs8t srn3gv2p tgrcf1 ftsna4z5q 1c czhlfg kv v3l16a5p 60m ydg3itpbg2 xtlwv0f4 a4v7g530s xchydk vi3 z4gyvq x1pclw 9rbwc cuik x0k2g7mrep 34','0ubw',1,0,0,0),(3181,'nuexq','h5b81x6q 43mbp1o',0,'lad0h8vmw2@yahoo.com','2013-05-26 19:38:00','bgx 60pvf 95ecbd dbny at48yrfs2m 7zbw5 vfupacx y4tx870kua g6vw4 3kh2 ad8205 z52qgim4fs lm5egxuc0 34d17 zi 4z qv1b8 b1iwvarqm 2t 50z r6bg43tg ocb3znq5xg 1mcf68 93gwu','k3',0,0,0,0),(3182,'0bw4tlqgxn','2gne y8pguoik1t bgpy4eha mityxnr5qz',1,'02mpbei@ya.ru','2013-09-14 05:13:30','gn2ox4uy mzlhg2b amgolc71 1evl3op2t 6lyf731g4k 92dw7fv rs5gpi v0hqpcf 5f7hutda2i bcqudgnim9 ke91x 5h6x1cma lm2rezv lv mg8e istef2l4k rk401glxi3 dg iw9y 6blc psc sar4 1razkv p93coxdyu6 uw35g n0yeastwx','stauw',0,0,0,0),(3183,'se','qofuv',1,'ty@gmail.com','2013-08-03 12:15:08','8hy1 985kp4blsg io 8u20koxg5 hp g1ced9g eif107 f1 8y0mxgah7 7e94hg 0xw7mnla3 4f2pg oy9di in748hcakf d5ph6gg8qy i3kqod9 bwz tzkc coi5yt07dr 6t ytu 6e2b siphw ewxmp0fo 17be8rg','sqkmy3rdcv',0,0,0,0),(3184,'o5gu7klc9','avb59hu v8lz0fhk5',0,'ty@gmail.com','2013-03-16 01:01:48','x1q4ug tdvx7z emn gru8k9p5 nel357u80 ntk61 wuby r6gv v5nibdk 5yh nedb ibpakwf p4xtb5 i3 ubqr3let6 sxcgh5ign 1gg6cyz 4gf3pm1l 1rk dlsapbhg ns4c tgv 0w41ad','0epo',0,0,0,0),(3185,'0bw4tlqgxn','ygx7 a2wy le835bvocn qvoemu',0,'cg962ap1ht@gmail.com','2013-11-24 00:20:12','0br2 5ktlrqncg b250pdng saob5k 7foz2q z5gu6 879 6v 3b qyws 8w29mt3 5e 36kg421ruv gro7wvh3 i6g2051 yhqv90na b1kwg7qz mv74k2f3 3tgyh6eb0 s46urzywpv zm1qt8 wacpkosq vrutbcw af gs5y8bv yfa26n378','gvt918a',1,0,0,0),(3186,'6qsxa71p','g0ltpfk29w 6dolmp32',0,'02mpbei@ya.ru','2013-01-08 06:41:50','bgt8e9 62agzip gbzlye hw53f 20gb8 xfiqkca89w r87yvei2 cdxls ncid39 obg701l 3uyi y2izgxc fkmy0cng 3vkanpg vh 3eib1v7qc aifgd em3h8 krgu9da8g 51tflgxg9w nvq1 pf9dagbky h3f wnuy0q uecx8lai dy4 puesdl fd6cag1 py1wd2nz3h g0','b3d',1,0,0,0),(3187,'hm7aurq64s','76q1tfnu ah7pq0',1,'xtkmq@gmail.com','2013-11-19 14:45:36','pgwrg761v f7yxu 3r84yezvb9 t3n0 ke4onbsgt 1sg8hu 0ywimt95u 45duyf3ab 0y2 k26v i8g1 3ye1ztr twzv63gh14 h5sl1ityn3 kg 7ex0rzbn6g s6 ok xlu zo o5h fvxphcdrgw z9oiacvt5p dt5yb t950plef cuh lci9sm f69tw','gyg4p60597',1,0,0,0),(3188,'vkyi1wo','f5gedt1swy d1g',0,'lad0h8vmw2@yahoo.com','2013-09-07 20:24:57','baq3y2m 5b ugs yfneg e547o2a8mx godx694 bu9wgo6hfn s6gntfmh 2gnpbc0q 0ofx4 5mhz8 mr 4hpzlqygk 5o gwr46ap lh7mq3p12 khoage iw0toa13 2rg gx aeg imaw65gz','pzg',0,0,0,0),(3189,'hm7aurq64s','igc842a duoky26g9b dbg3 gs6cwb8d',1,'02mpbei@ya.ru','2013-06-13 00:36:47','ngc domqa egm rico0 f74x23vm hvus9e wiptgu tekfx1b30 h7gfo8 ag p4tl fu d9egfgq5 xpzrwq0 917kl 4vm2 6rgbfcum cva6 5ek0xi6 g6 8zbr4u9 ges 72gdwmi1k gz7 7fcrd0 argby85q g82xb xbour aox06p4','ubm7ap8re',1,0,0,0),(3190,'oxg5','cf7u8eq 5i 6izgd1 volf1',0,'x9bhu@ya.ru','2013-05-27 08:10:20','vwyfkl2gsx grugpynmq yzlcv uw2gx4 u9g2y51cn4 iqg9 1va 18zgf76 sbkl9g4 vyehlp th5svyw gcw qc79d g6 9h0r3yszb 9fag2u 1s u0kt7z k23isewf kio3g209zb','ugcd3am',0,0,0,0),(3191,'se','9ds u4sqkm3 ksqnvfgry',0,'6gm0u4set@ya.ru','2013-02-24 03:51:55','7ox q8n1 vlw tbm9cv6dk kuw75agql z1o9u0lwv 1tqn82yc3 iu89z uxtm9h oh5fqb 5c8sn9tvqb uy2wdizbf gc0y san4wk78mg z0uhyr6p lsxgcq63 cwfa 02m7i8cgyo hcvmw zs6n okyzarglv eo5q87w4 w3xn7mpla 4px67k81 48ra3f06s5 x5vd 7mtu8g ydo','9q',1,0,0,0),(3192,'5yi','1qby vrhok624 v9ux57',0,'o7kz@gmail.com','2013-11-08 14:14:49','dw fbuwmyhqg7 ynz9ofl wpoxfes pulm58f2 hbntw wlretaop46 tn glu 1loi aczk7pi 8xwto9um qchse137v dx1kg tk9vnb0 ckgl 9m3lit5x 9o1fyh3 sp53r4k8 4okxgfmh 5i3zlx ds7q2g 9sn6fv','wruf8',1,0,0,0),(3193,'5yi','wgcg5szq 7puy8',1,'lad0h8vmw2@yahoo.com','2013-04-11 18:39:36','o7dxsm8qn fbwx6928m gy mxc6 2risyk e7sgf2rnb wgaz tm ak61 bxog scn4fy x0yp6dmr y1cle2mk8g n9 7ck 8p3sdckgrq 13526qmoa8 2kl03m8 ggbfd pbfn3awtm 4shmzefn wnahg34l5i baroez2c1 wcik qx91va ue0os2t ueas iyz3 5g4bafld3','l0agbd2zo',0,0,0,0),(3194,'o5gu7klc9','3zoqad80g',0,'cg962ap1ht@gmail.com','2013-11-09 16:11:08','dg 9s0bfy y1vc fgh 5dz81rh0a uez12i8fqy t4sqwngo v0p3 fg2mlkd 2xhk130ovb 1l628vbq wytp uhc tsrhuoq5p fnig svg1fg36uk mvruwq8ot7 luaezbgr5 u126 q96he7kbs fpr41dlwg q1idxbz4 n0 aq0pk6y 9m8svqf urb5 imw694','s3fob1a',0,0,0,0),(3195,'5yi','5notflqr1a muelw3',0,'o7kz@gmail.com','2013-04-06 13:24:35','g3xw olvnd9tg4m rggz68e yx8sc9l o68pbcw0z t6 ebt9o2 gn3 ecus5dpmt mczfwok o9zvqxr 53ulwxc9 74 egvk326g b3pih 5bmtongsq4 tn 1h9yqu2g 5tn9exgq2 d20y r0kh6pqlum liszgpb1wn 7xs9 3qw 862te poeduxk32','ngzfs25',0,0,0,0),(3196,'b6kgh4gur','drhgy75w1a 0t snmghi rw5vit',1,'n1szohqi@ua.ru','2013-02-06 07:12:12','bkgqxwpr5 ulby 3tmwi15 8l2ogypiuw ryua l2 nypwibqd 8pza 4t5s de 4f0rgn8gw qp6t k9qbrgn o45k0lgd9 xmdr 46tupqz kglh8p ca5 uzf1c eio91bm x1g0fkzq2 8uo0ibm4y 9cby7gd aly01siz oqvwrhb','wnf0yi6',1,0,0,0),(3197,'hm7aurq64s','znegdku',1,'02mpbei@ya.ru','2013-10-13 16:26:12','gfld6 g1vy6t8 0ca2sbfo8 yte 7tap19vm 5wzqg1rnuc zo495gpf ai9g8e6ylq a836y9qk 4tk0g98 cggpoa9 wm1g6lqts0 i4 af7xcbv6r ozlf09 gtob93hk5 isglu1 qvw0g5o2 9udmhgb87 90x6g4 5gwcr13uxe k4 rv83bg2ch6','o5w3v',1,0,0,0),(3198,'vkyi1wo','ok8n4u6qyw 08ucd',0,'xtkmq@gmail.com','2013-12-15 19:23:18','ir4 a54ifcnq c6sal ctevws9 xtc5 yr adl64w3sv wxms 9mghr5 9r6db3 kxcb2 p860 g9ldch 9gqb2 3iuo09sr gnwlf4 aic0 wh lmqrghgu28 a82ovk30zq rct5pk3q g029 f03l','o2pd0',0,0,0,0),(3199,'nuexq','oq1fnzray q6dvau7zw kd h6doa',1,'x9bhu@ya.ru','2013-02-15 13:49:30','d5xa36m8 ri mw3gbn 95w klhut1g3r g2i3p 0s x9ref0o482 hyfd906kv 1kistlz6ew 354qk6cvgf 63v 2bszl7uqod az3ndr vo n678 nd1 bo2v 2cg6v7x fh1rq t62 ahnbrs3d1 nuv8f','zgot',1,0,0,0),(3200,'nuexq','wrkix pr2wuz75bq',0,'ty@gmail.com','2013-01-19 10:55:28','3w abfvdktn yw25ktrn flpk4 g7ki3cdpm oef msw 7gbs 93 ih43 sr4q sea2y5zco we36 pu 4lac m2sg3ldy9a ez gu xct62yo 3pyq6 gc13 c3s0w7 yo zu57pbh 07fbs tm shd9afbwz h0 7du69sqza5 s6i3m2no8','hob',0,0,0,0),(3201,'oxg5','zy98holb7r nafsbx40 gvyg3q9ot',1,'6xvn9e27@ya.ru','2013-06-02 21:06:29','z0pf1 a9x 4nftxdy ad h5vmapoc ogsu24 67f9auk2 nlmf5wk2g v0wsghb1z bd1serkolg mtwpy ig9n k9civg327h ecr35g 3ifgdox mhy x94poyud s6 nuwsrfx 423uvf 1zrof zosq0aflyv kov7p','fb',0,0,0,0),(3202,'0bw4tlqgxn','s0815 zgkrqhlw 1sv fwy5hl9r',0,'6xvn9e27@ya.ru','2013-02-02 10:26:46','vs hxruaebip cfo 7smhfq01xg y95 32 g1e q01iflcn yocwifxr8 r8g2lq1os3 lboe 2w78no l69z82 65lxr gi98n6 3qy14ra 9x8g0yu mrwzs h9gyu87a 3zlnosm qtz9mn eyh0d waoerf10 pvnko3lgm1 d4g2bt53 20sbtprd5 8mv','dkq',0,0,0,0);
/*!40000 ALTER TABLE `thread` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) DEFAULT NULL,
  `about` text,
  `name` varchar(32) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `isAnonymous` tinyint(1) NOT NULL DEFAULT '0',
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`email`),
  UNIQUE KEY `email_UNIQUE` (`email`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4769 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (4734,'ex2cfh5z8p','s68zd 9qzum1ta yx5 ykuzn 1g63melht z7r2 moygau 0rztk47fp flx h2 g6be4z gpltv 1y7qc zy6nog0d ggrpxq1d qumz58p4rl qrv 6dwf8blm2 60egyx4fw vzfmth 8kg y5eno1d','xg1w 58rfg c2ds8zygf k3g5tizyxc','02mpbei@ya.ru',0,'2016-01-12 19:31:52'),(4758,NULL,NULL,NULL,'04un@mail.ru',1,'2016-01-12 19:52:47'),(4759,'es72','8zmn9y5tkp yanzwuc7x3 bol9efgpg slwz4f xeka5mv2 8gq60 c9q6xd3 k7 bok8i1p lrwy 1iwck v62o 4m35d 0gl54ct3 dial7u porif 2oq43f qvk5ugn 8rw6pg 2kf7r6','ts51c','0bghz675fc@ya.ru',0,'2016-01-12 19:52:49'),(4748,'4a9w','egf2abw 3goz bp f7 l7g8vb niq 7dpmeg8ki ikr1lwxo 9r m0gc l06n9 np41vxhl 7olhd g5rxkfs g8dn2ytzqg be1oxrhka r2k tygae9fc0 pi w5az1qn gbx qggd9yf6be ify','kxg6rzpvc7 mp4','0yhgcs5@list.ru',0,'2016-01-12 19:52:31'),(4737,'g6x5eq','qdnfx sd8foca 7s0wn ph5 f62akzo sc86pvem m57csgf6z besaikoh 6d785m2 80ab1 bgderx5ho4 rybg6oxe5h 95ctsug2 rw5 xl9uh3zv 345yac86 fh0169sxb sao9xutf6g 2gtl ng5 98ce53 0uywapn4 pvs3n2l x6o 4pa gw ix8epmv 7l 365kq ux7qw2d1gb','4uxaies81 cmfu kmf5l pqucd','1k@mail.ru',0,'2016-01-12 19:52:04'),(4736,'hzisl7nv6','ox 5l pk13tbr7go a2h4vk8ng 0v5q1wu4b 3h9ilunsf 1qz 40uzh8mqp r9 g7u69ek80 z8 mph906i5 rgh 7sqp1dn hi6lvzm2q 3beroqc kvcfq hny2 6y xk65rup vyh5crf9 dgu3sc4q5i eq dp9q2vlt fpt 5d67 hq 039 csl2m3zwht xkgw','hmb 4a9lbd57m xqepvgby eg61u','4limneu@yahoo.com',0,'2016-01-12 19:52:03'),(4740,'ox6dbl','39isw5y hdy0rngu xtfpnc6m gzy5bv134r xuehv24 6g2o3ykq8g yg sc1xgwaq6 pwxtyrh 85gle tnud7l 9gbxc wmo 3lxiry eo x7feo z9hecp1q 9vo7l6 z84 um81p4hltk fzt74xyv 82xlofuc 9671f ig7q62t dqiw cx6a5wu2 9t20x qdn08r42','4e0g eax sq','6gcq@yahoo.com',0,'2016-01-12 19:52:15'),(4727,'6r','73ik2bgam0 kfp32u5rv ygqcdxb5 k361y aw xgu83t rse 483u06l1 x05mc6g 4azg imhgbg g7243tfcvr des2x txez vf o51prg p3 5q2e 81pu9gwda gw05 vue4 dnk 64x','g5v2a3rnk 5f8zi2','6gm0u4set@ya.ru',0,'2016-01-12 19:31:52'),(4726,'ugzo','8w1oyve5z fu76qg8rn3 z4y8ovli 3b0 0l iy05hn 1m92fhx g8qn4t6a f0 94elv dhgwzlsg fkbxg gvhurdtzi 3g20g9oc zxhkygt w15mlgi 4wi1 5oe g157 hngyu93e a14 wke3sxy','o0adeiwz7 n3st41 ogx3','6xvn9e27@ya.ru',0,'2016-01-12 19:31:52'),(4747,'6pwsd05','ouev30n sa xns5y2d 4bydxipv9 0qvpiclu8 h4l sv34b zf7liu1m 98e3cyniq4 ns kgsufa g20 lu2c laswmk ly5fvc c15xwn p3 eivg i4hepl2s5 sv7ga9 ucm utf7z96 q3hapwtk 6yr7ps qiyu 1tdyr8mif 2pa5 m3q6ls7eb 6357x41v mgfzei0y4','g2i17ckh','adikpgvwu@gmail.com',0,'2016-01-12 19:52:30'),(4735,'hf','qm8l3g06 zwyohs 4f932l zne g02vog hg r2t a94o ohvy67 avsm 8g7r5 nizwqc 4piom gd rok pvw ra rh si07rt q65 rg l3ougardwm 7g2s5 fd6 0mfw4qpacb i48ntqh5s ue05z71gg ghdl2o0c','w2ne','cg962ap1ht@gmail.com',0,'2016-01-12 19:31:52'),(4743,'8li0a','q98tzn1 pqn2g 0iky qi9 596 2apg08tm mbuplc 1hqcexo4gf pggm9 mhb1ypqfz 5l9f wi9oekz8f qwa9lbhx1 twcb s0qmgtz kr4c l28rfh 5ikghw yd 3gensckf tsvkg25cml k4fg3ztb hrofeu8 7t4o 9g3ig6z7 fgzgq7','yfei7nug6 53mzg','ch6f7wp@ya.ru',0,'2016-01-12 19:52:17'),(4744,NULL,NULL,NULL,'cxhg@bk.ru',1,'2016-01-12 19:52:18'),(4750,'lu','h671 kyma 9vlc2e pxhrd mg1 41setoq8 o1kg5 82zyrh1gg gt ie cedgns8yru vpywi5fd 0g viduoafe4 8anktx5czy 87fz0yp 8dalnugwr 7dfx4yme 5z yhov5 y63a 35cnag0ob crlgz8 2b4gh 9kc2yf40rb','po6cggbr3 rekh9l34','d5blro1k@bk.ru',0,'2016-01-12 19:52:33'),(4761,NULL,NULL,NULL,'g3obfc@ya.ru',1,'2016-01-12 19:52:51'),(4739,NULL,NULL,NULL,'ge2xn@ya.ru',1,'2016-01-12 19:52:11'),(4760,NULL,NULL,NULL,'gscr7dqgz@bk.ru',1,'2016-01-12 19:52:51'),(4755,NULL,NULL,NULL,'k5@ya.ru',1,'2016-01-12 19:52:44'),(4752,'6ebd2zh','80guzye176 p1w528gf 0e6lyw4 35 bg3ft2uw5 a0qo6gt3z z3s fpml we4 ows4f deo6buqi10 c4h6 ogm3q u4os 3t m4 xzw cmz gds1 qub20vdf twkl ysxv9gabn','qa1gwu q4u6z d9hg','kp@list.ru',0,'2016-01-12 19:52:37'),(4732,'civkw57y8b','sqyw scz 58guhwbmd 28zamigt fdsr7yeqw g2 5m4hix1 l2dsxqpkb vf0k dsyg6 1glri27e4 ltu0wyg nalez0q46 4bh7pa32te 8tiu3v5gpc cif0ve s2 fezlg0t7 dn32tw l9zgb se2x hpdba7qt3 0hiw12x93 mgkrca yb7zng','uzwk','lad0h8vmw2@yahoo.com',0,'2016-01-12 19:31:52'),(4738,NULL,NULL,NULL,'mosh8gnk0@bk.ru',1,'2016-01-12 19:52:10'),(4749,'z4bto','qp 7iwksbzo 4k fds nr7sfyih k23oqcxv9 box qoblg v6fawqc2 en gt 4ql0c8hg95 agr9t fe2c 2c7l9p1r 8s i2m mvu 346c2p soe08p','qk1','mq4cgk6p@bk.ru',0,'2016-01-12 19:52:32'),(4746,'5mg3gc8u','gerxli cb26qr1kgo rtdpwgfs 5p8a6 x3n7vowa 37ye l1naw9os amw9 bl3w tlehmzs qz k0 98todm kyz2rcum rxqo9yplcz nix31y 8gio1fe vfs mql3g zcp93vy2 9h6onx1c elfvzy309 6gpu0icx2 os 35gwb w6pd unyovfag vy vmwcpgz 6clv1pg2i','flashi','msz9v5c7@yahoo.com',0,'2016-01-12 19:52:23'),(4730,'ps2v40593i','5n4 48yqg0na omy8 9wargy 7kwo24z0v gwuslo vzyu6p8bm3 1ghnm abteo3 ck2wni3g f0 fy4bwsv 54h0om gqoxy 5xqswt9fhl wg27 m8x1 m0ya 06 28w4dyoa 7zoghm gix7 r71n5236ox s7zc gc734 bg9 z1g 3sht2edv','0yvh','n1szohqi@ua.ru',0,'2016-01-12 19:31:52'),(4745,'dmnat4w9u0','fk ted gd dks oc prfs1 9xdeowk3i lp8v3c 2g6stky qdpf tb25neuo nwg1 uam3qcg6tr kh g9xd7ar1q 3pnse41izu li0an46f2 w6isb8o xipg89oufq 6qgey cs5 3h02apvi g4 r3igw41az','5ru3l g5 vkrn4p','o07zbk@gmail.com',0,'2016-01-12 19:52:22'),(4728,'8xgcsyp','x14cgd pasg8 xn4y0kr u3wsa7d 7y 4w0mu3 5wy4mkbugn aqmop ux3vw1 ngg06 gmcuyw bndu 7cxpm wok a31w82neg vbg8et6g3x d3vn8e q4b6yv1 f708m2x awz3gy92t 23iz1ges4r hg9m 5nc78tu xf zn u6 7azstph9 ty l6s0g4','inuvcwxb7 ch9','o7kz@gmail.com',0,'2016-01-12 19:31:52'),(4757,NULL,NULL,NULL,'oq8b1r@gmail.com',1,'2016-01-12 19:52:46'),(4742,NULL,NULL,NULL,'phk04ngi@gmail.com',1,'2016-01-12 19:52:16'),(4762,NULL,NULL,NULL,'qiadn2g@yahoo.com',1,'2016-01-12 19:52:51'),(4765,'4e','or8lhqg 6onpdb 34qka 8xunywar x6 g4s3dx 8uv g36d9m5i io bwlpx4n a256vwmcb rc39 7o3as5vm gnidx23e6 nbw7s93 w6kub8 7nglutahw9 zxfsne placs7gidm 39ey4 dcr4o0u 8fvzuycdt 924te ggwhp h5 78gmi3kvx','x1mad m6x fngukl3r','ra32f0yn@gmail.com',0,'2016-01-12 19:54:09'),(4751,'bcvw2','pt0yvh s3c2 3segk4rn okg w4nqcok z16cuil 6r1ag4x hb78gugl ca nvo91f 35vz8tam a6i z0 5h h8zn2ol ro9 y64w7uhl8 x9wtp1yr xhitg as64gv 91s8g gy0n8b z17v8t 0o7t 0ri8n3cf7t nz09a4k 71hd8mt yqx w9t esggiha','xanu pz9yhciof 8ez7k 1w','s2c@ya.ru',0,'2016-01-12 19:52:34'),(4741,NULL,NULL,NULL,'su93a@gmail.com',1,'2016-01-12 19:52:15'),(4731,'gsqd8rp1f9','v31gixol a16x 8s756i x4 t3sdz5x ag5u 6f85gv7eh y918 9bcn kix9u67eq 2v u28e6ogd59 qb4 kdyilt49r fchk6dytl bm x6sqotr1w o0dlgm4 7nht5lyb gd5o4 cgqr gm','mw6hf 08','ty@gmail.com',0,'2016-01-12 19:31:52'),(4756,'pv7wb','bvoek0i 932sdowk nc9m0wht3x lh5zt8bk d91fb fp4v 90gnxq xpnqhm pwye upydv25 r46g2tmz yk6inmd2qv 4gib03f qn52w0 ov s0c 9bg01i3her aqzde3 97cytfd51 b3rg pmqo7n','zad g7w hl zxh','u0xvi@ya.ru',0,'2016-01-12 19:52:45'),(4729,'5pfotyhvl','1keamn yogdp09 s5 3x5r47medy uq6 2ig7g qfns829 t8ivm9a ptz pz4 i9no61y4 gf9qb zti z9 us840i f9d l4bcdgh8y zrg1e3sd z9 zy 2tlgeuf8cw l27ruzg3 2dh4pk t23g','hkoq2rd70z cl56g71 tlos8zd','x9bhu@ya.ru',0,'2016-01-12 19:31:52'),(4733,'p2nbm96l8','53eu42mh 2bz0nqym rn 87asq pycw91g ye2r pkmhv afcvzs4oh ewl 05743v l3wc korbcw4g5h veo 6f04 z9 oup8s kzt4s a8dz305g 5sc90gi6 diotarxbu tghwbcl7rv c8etka7s4f 0wv36o12mg 9gbhps51 go7694ha','gp ma8c ibop 8hgx7','xtkmq@gmail.com',0,'2016-01-12 19:31:52'),(4754,'kd6yq1b8','st7 li635u284y ugd70x n8 npkwlo d2f3gpsv90 yg36zpb q3dz4pf qopfm ziyw1c l20f9 fwg lgnr2zkq89 ih sxn85uwr luez5h9w mrgq4 uxac9bv2 qi43 0m6 k8y e9 u7ay3 8zrhu','tc 5diuek9cs erfm6a4','yghm1u@yahoo.com',0,'2016-01-12 19:52:43'),(4753,NULL,NULL,NULL,'z5@list.ru',1,'2016-01-12 19:52:40');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(3) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `password` varchar(20) NOT NULL,
  `first_name` varchar(20) DEFAULT NULL,
  `last_name` varchar(20) DEFAULT NULL,
  `sex` varchar(5) DEFAULT NULL,
  `position` varchar(20) DEFAULT NULL,
  `order` int(11) DEFAULT NULL,
  `restaraunt` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_users_1_idx` (`restaraunt`),
  KEY `fk_users_2_idx` (`order`),
  CONSTRAINT `fk_users_1` FOREIGN KEY (`restaraunt`) REFERENCES `restaraunts` (`id_restaraunt`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_users_2` FOREIGN KEY (`order`) REFERENCES `orders` (`id_ord`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-01-12 23:38:25
