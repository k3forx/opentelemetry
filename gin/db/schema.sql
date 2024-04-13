CREATE DATABASE IF NOT EXISTS `app`;

use `app`;

CREATE TABLE `authors` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(255) NOT NULL,
  `bio`  TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `books` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `author_id` BIGINT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY fkey_author_id (author_id) REFERENCES authors (id)
);

INSERT INTO `authors`
  (name, bio)
VALUES
  ("東野圭吾", "1985年、『放課後』で第31回江戸川乱歩賞を受賞し、作家デビュー。1999年に『秘密』で日本推理作家協会賞を受賞し、直木賞候補になってからは毎年のように作品が直木賞候補に挙がり、2006年に『容疑者Xの献身』で直木賞や本格ミステリ大賞を受賞する。"),
  ("湊かなえ", "広島県因島市中庄町（現・尾道市因島中庄町）生まれ。武庫川女子大学家政学部被服学科卒業。2007年には金戸 美苗（かなと みなえ）の名義で第35回創作ラジオドラマ大賞を受賞した。")
;
