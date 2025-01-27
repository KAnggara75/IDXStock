-- CreateTable
CREATE TABLE `daily` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(200) NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `price` INTEGER NULL,
    `priceopen` INTEGER NULL,
    `high` INTEGER NULL,
    `low` INTEGER NULL,
    `volume` BIGINT NULL,
    `marketcap` BIGINT NULL,
    `tradetime` DATETIME(3) NULL,
    `volumeavg` INTEGER NULL,
    `pe` INTEGER NULL,
    `eps` INTEGER NULL,
    `high52` INTEGER NULL,
    `low52` INTEGER NULL,
    `change` INTEGER NULL,
    `changepct` INTEGER NULL,
    `closeyest` INTEGER NULL,
    `shares` BIGINT NULL,
    `insertBy` VARCHAR(100) NOT NULL,

    UNIQUE INDEX `daily_code_key`(`code`),
    UNIQUE INDEX `daily_name_key`(`name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `history` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(200) NOT NULL,
    `date` DATE NOT NULL,
    `previous` INTEGER NULL,
    `open_price` INTEGER NULL,
    `first_trade` INTEGER NULL,
    `high` INTEGER NULL,
    `low` INTEGER NULL,
    `close` INTEGER NULL,
    `change` INTEGER NULL,
    `volume` BIGINT NULL,
    `value` BIGINT NULL,
    `frequency` BIGINT NULL,
    `index_individual` DECIMAL(65, 30) NULL,
    `offer` BIGINT NULL,
    `offer_volume` BIGINT NULL,
    `bid` BIGINT NULL,
    `bid_volume` BIGINT NULL,
    `listed_shares` BIGINT NULL,
    `tradeble_shares` BIGINT NULL,
    `weight_for_index` BIGINT NULL,
    `foreign_sell` BIGINT NULL,
    `foreign_buy` BIGINT NULL,
    `delisting_date` DATE NULL,
    `non_regular_volume` BIGINT NULL,
    `non_regular_value` BIGINT NULL,
    `non_regular_frequency` BIGINT NULL,
    `insertBy` VARCHAR(100) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `stocks` (
    `code` VARCHAR(10) NOT NULL,
    `name` VARCHAR(200) NOT NULL,
    `listing_date` DATE NULL,
    `shares` BIGINT NOT NULL,
    `board` ENUM('Watchlist', 'Main', 'Development', 'Acceleration', 'EkonomiBaru') NOT NULL DEFAULT 'Main',

    UNIQUE INDEX `stocks_name_key`(`name`),
    INDEX `stocks_code_idx`(`code`),
    PRIMARY KEY (`code`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `under_writer` (
    `code` VARCHAR(10) NOT NULL,
    `name` VARCHAR(100) NOT NULL,

    INDEX `under_writer_code_idx`(`code`),
    PRIMARY KEY (`code`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `users` (
    `email` VARCHAR(191) NOT NULL,
    `username` VARCHAR(100) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `logoutAt` INTEGER NOT NULL DEFAULT 0,

    UNIQUE INDEX `users_email_key`(`email`),
    PRIMARY KEY (`username`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `daily` ADD CONSTRAINT `daily_insertBy_fkey` FOREIGN KEY (`insertBy`) REFERENCES `users`(`username`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `daily` ADD CONSTRAINT `daily_code_fkey` FOREIGN KEY (`code`) REFERENCES `stocks`(`code`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `history` ADD CONSTRAINT `history_insertBy_fkey` FOREIGN KEY (`insertBy`) REFERENCES `users`(`username`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `history` ADD CONSTRAINT `history_code_fkey` FOREIGN KEY (`code`) REFERENCES `stocks`(`code`) ON DELETE RESTRICT ON UPDATE CASCADE;
