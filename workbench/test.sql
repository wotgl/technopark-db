SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

CREATE SCHEMA IF NOT EXISTS `test_db` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `test_db` ;

-- -----------------------------------------------------
-- Table `test_db`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(32) NULL,
  `about` TEXT NULL,
  `name` VARCHAR(32) NULL,
  `email` VARCHAR(255) NOT NULL,
  `isAnonymous` TINYINT(1) NOT NULL DEFAULT 0,
  `date` TIMESTAMP NOT NULL,
  UNIQUE INDEX `email_UNIQUE` (`email` ASC),
  PRIMARY KEY (`email`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `test_db`.`forum`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`forum` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `short_name` VARCHAR(255) NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `user_idx` (`user` ASC),
  PRIMARY KEY (`short_name`),
  UNIQUE INDEX `short_name_UNIQUE` (`short_name` ASC),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  CONSTRAINT `fk_forum_user`
    FOREIGN KEY (`user`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
PACK_KEYS = DEFAULT;


-- -----------------------------------------------------
-- Table `test_db`.`follow`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`follow` (
  `follower` VARCHAR(255) NOT NULL,
  `followee` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`followee`, `follower`),
  INDEX `fk_follow_1_idx` (`follower` ASC),
  CONSTRAINT `fk_follow_follower`
    FOREIGN KEY (`follower`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_follow_followee`
    FOREIGN KEY (`followee`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `test_db`.`thread`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`thread` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `forum` VARCHAR(255) NOT NULL,
  `title` VARCHAR(45) NOT NULL,
  `isClosed` TINYINT(1) NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `message` TEXT NOT NULL,
  `slug` VARCHAR(45) NOT NULL,
  `isDeleted` TINYINT(1) NOT NULL DEFAULT 0,
  `likes` INT NOT NULL DEFAULT 0,
  `dislikes` INT NOT NULL DEFAULT 0,
  `points` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_user_idx` (`user` ASC),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  CONSTRAINT `fk_thread_forum`
    FOREIGN KEY (`forum`)
    REFERENCES `test_db`.`forum` (`short_name`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_thread_user`
    FOREIGN KEY (`user`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `test_db`.`subscribe`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`subscribe` (
  `thread` INT NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`thread`, `user`),
  INDEX `fk_subscribe_user_idx` (`user` ASC),
  CONSTRAINT `fk_subscribe_thread`
    FOREIGN KEY (`thread`)
    REFERENCES `test_db`.`thread` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_subscribe_user`
    FOREIGN KEY (`user`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `test_db`.`post`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_db`.`post` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '		',
  `thread` INT NOT NULL,
  `message` TEXT NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  `forum` VARCHAR(255) NOT NULL,
  `parent` VARCHAR(255) NULL DEFAULT 0,
  `isApproved` TINYINT(1) NOT NULL DEFAULT 0,
  `isHighlighted` TINYINT(1) NOT NULL DEFAULT 0,
  `isEdited` TINYINT(1) NOT NULL DEFAULT 0,
  `isSpam` TINYINT(1) NOT NULL DEFAULT 0,
  `isDeleted` TINYINT(1) NOT NULL DEFAULT 0,
  `likes` INT NOT NULL DEFAULT 0,
  `dislikes` INT NOT NULL DEFAULT 0,
  `points` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_post_user_idx` (`user` ASC),
  INDEX `fk_post_forum_idx` (`forum` ASC),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  INDEX `fk_post_thread_idx` (`thread` ASC),
  CONSTRAINT `fk_post_thread`
    FOREIGN KEY (`thread`)
    REFERENCES `test_db`.`thread` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_user`
    FOREIGN KEY (`user`)
    REFERENCES `test_db`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_forum`
    FOREIGN KEY (`forum`)
    REFERENCES `test_db`.`forum` (`short_name`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
PACK_KEYS = DEFAULT;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
