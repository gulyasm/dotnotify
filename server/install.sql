DROP DATABASE IF EXISTS Dotnotify;
CREATE DATABASE Dotnotify;
USE Dotnotify;

CREATE TABLE `Files` (
    `Name` varchar(256) NOT NULL,
    `User` varchar(256) NOT NULL,
    `Hash` bigint NOT NULL,
    `CreatedAt` datetime NOT NULL,
    `ModifiedAt` datetime NOT NULL,
    PRIMARY KEY(`Name`, `User`)
);

CREATE TABLE `Users` (
    `Name` varchar(256) NOT NULL,
    PRIMARY KEY(`Name`)
);
