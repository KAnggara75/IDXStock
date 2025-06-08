-- CreateEnum
CREATE TYPE "Board" AS ENUM ('Watchlist', 'Main', 'Development', 'Acceleration', 'EkonomiBaru', 'A_SERIES', 'B_SERIES', 'C_SERIES', 'PREFEREN');

-- CreateTable
CREATE TABLE "daily" (
    "id" SERIAL NOT NULL,
    "code" VARCHAR(200) NOT NULL,
    "name" TEXT NOT NULL,
    "price" INTEGER,
    "priceopen" INTEGER,
    "high" INTEGER,
    "low" INTEGER,
    "volume" BIGINT,
    "marketcap" BIGINT,
    "tradetime" TIMESTAMP(3),
    "volumeavg" INTEGER,
    "pe" INTEGER,
    "eps" INTEGER,
    "high52" INTEGER,
    "low52" INTEGER,
    "change" INTEGER,
    "changepct" INTEGER,
    "closeyest" INTEGER,
    "shares" BIGINT,
    "insertBy" VARCHAR(100) NOT NULL,

    CONSTRAINT "daily_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "history" (
    "id" SERIAL NOT NULL,
    "code" VARCHAR(10) NOT NULL,
    "date" DATE NOT NULL,
    "previous" INTEGER,
    "open_price" INTEGER,
    "first_trade" INTEGER,
    "high" INTEGER,
    "low" INTEGER,
    "close" INTEGER,
    "change" INTEGER,
    "volume" BIGINT,
    "value" BIGINT,
    "frequency" BIGINT,
    "index_individual" DECIMAL(15,2),
    "offer" BIGINT,
    "offer_volume" BIGINT,
    "bid" BIGINT,
    "bid_volume" BIGINT,
    "listed_shares" BIGINT,
    "tradeble_shares" BIGINT,
    "weight_for_index" BIGINT,
    "foreign_sell" BIGINT,
    "foreign_buy" BIGINT,
    "delisting_date" DATE,
    "non_regular_volume" BIGINT,
    "non_regular_value" BIGINT,
    "non_regular_frequency" BIGINT,

    CONSTRAINT "history_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "stocks" (
    "code" VARCHAR(10) NOT NULL,
    "name" VARCHAR(200) NOT NULL,
    "listing_date" DATE,
    "delisting_date" DATE,
    "shares" BIGINT NOT NULL,
    "board" "Board" NOT NULL DEFAULT 'Main',

    CONSTRAINT "stocks_pkey" PRIMARY KEY ("code")
);

-- CreateTable
CREATE TABLE "under_writer" (
    "code" VARCHAR(10) NOT NULL,
    "name" VARCHAR(100) NOT NULL,

    CONSTRAINT "under_writer_pkey" PRIMARY KEY ("code")
);

-- CreateTable
CREATE TABLE "users" (
    "email" TEXT NOT NULL,
    "username" VARCHAR(100) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "isActive" BOOLEAN NOT NULL DEFAULT true,
    "logoutAt" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "users_pkey" PRIMARY KEY ("username")
);

-- CreateIndex
CREATE UNIQUE INDEX "daily_code_key" ON "daily"("code");

-- CreateIndex
CREATE UNIQUE INDEX "daily_name_key" ON "daily"("name");

-- CreateIndex
CREATE INDEX "daily_code_idx" ON "daily"("code");

-- CreateIndex
CREATE INDEX "history_code_date_idx" ON "history"("code", "date");

-- CreateIndex
CREATE UNIQUE INDEX "history_code_date_key" ON "history"("code", "date");

-- CreateIndex
CREATE UNIQUE INDEX "stocks_name_key" ON "stocks"("name");

-- CreateIndex
CREATE INDEX "stocks_code_idx" ON "stocks"("code");

-- CreateIndex
CREATE UNIQUE INDEX "stocks_code_key" ON "stocks"("code");

-- CreateIndex
CREATE INDEX "under_writer_code_idx" ON "under_writer"("code");

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- AddForeignKey
ALTER TABLE "daily" ADD CONSTRAINT "daily_insertBy_fkey" FOREIGN KEY ("insertBy") REFERENCES "users"("username") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "daily" ADD CONSTRAINT "daily_code_fkey" FOREIGN KEY ("code") REFERENCES "stocks"("code") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "history" ADD CONSTRAINT "history_code_fkey" FOREIGN KEY ("code") REFERENCES "stocks"("code") ON DELETE RESTRICT ON UPDATE CASCADE;
