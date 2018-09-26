--
-- CREATE TABLE Users
--
  CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `username` varchar(32) NOT NULL,
    `email` varchar(255) NOT NULL UNIQUE,
    `password` varchar(255) NOT NULL,
    `createdAt` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
  )
