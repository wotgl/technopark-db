SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

CREATE SCHEMA IF NOT EXISTS `mydb_s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `mydb_s` ;

-- -----------------------------------------------------
-- Table `mydb_s`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(32) NOT NULL,
  `about` TEXT NOT NULL,
  `name` VARCHAR(32) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `isAnonymous` TINYINT(1) NOT NULL DEFAULT 0,
  `date` TIMESTAMP NOT NULL,
  PRIMARY KEY (`email`),
  INDEX `user_id` (`id` ASC))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `mydb_s`.`forum`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`forum` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `short_name` VARCHAR(255) NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  `date` TIMESTAMP NOT NULL,
  INDEX `user_idx` (`user` ASC),
  PRIMARY KEY (`short_name`),
  INDEX `forum_id` (`id` ASC),
  CONSTRAINT `fk_forum_user`
    FOREIGN KEY (`user`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci
PACK_KEYS = DEFAULT;


-- -----------------------------------------------------
-- Table `mydb_s`.`follow`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`follow` (
  `follower` VARCHAR(255) NOT NULL,
  `followee` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`followee`, `follower`),
  INDEX `fk_follow_1_idx` (`follower` ASC),
  CONSTRAINT `fk_follow_follower`
    FOREIGN KEY (`follower`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_follow_followee`
    FOREIGN KEY (`followee`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `mydb_s`.`thread`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`thread` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `forum` VARCHAR(255) NOT NULL,
  `title` VARCHAR(45) NOT NULL,
  `isClosed` TINYINT(1) NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `message` TEXT NOT NULL,
  `slug` VARCHAR(45) NOT NULL,
  `isDeleted` TINYINT(1) NOT NULL DEFAULT 0,
  `likes` INT NOT NULL DEFAULT 0,
  `dislikes` INT NOT NULL DEFAULT 0,
  `points` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_user_idx` (`user` ASC),
  CONSTRAINT `fk_thread_forum`
    FOREIGN KEY (`forum`)
    REFERENCES `mydb_s`.`forum` (`short_name`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_thread_user`
    FOREIGN KEY (`user`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `mydb_s`.`subscribe`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`subscribe` (
  `thread` INT NOT NULL,
  `user` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`thread`, `user`),
  INDEX `fk_subscribe_user_idx` (`user` ASC),
  CONSTRAINT `fk_subscribe_thread`
    FOREIGN KEY (`thread`)
    REFERENCES `mydb_s`.`thread` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_subscribe_user`
    FOREIGN KEY (`user`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `mydb_s`.`post`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb_s`.`post` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `date` TIMESTAMP NOT NULL COMMENT '		',
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
  INDEX `fk_post_thread_idx` (`thread` ASC),
  CONSTRAINT `fk_post_thread`
    FOREIGN KEY (`thread`)
    REFERENCES `mydb_s`.`thread` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_user`
    FOREIGN KEY (`user`)
    REFERENCES `mydb_s`.`user` (`email`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_post_forum`
    FOREIGN KEY (`forum`)
    REFERENCES `mydb_s`.`forum` (`short_name`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
