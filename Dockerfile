FROM oven/bun:alpine AS base

LABEL authors="KAnggara75"
LABEL maintainer="kanggara75@gmail.com"

WORKDIR /usr/src/app

ENV NODE_ENV=production
COPY . .
RUN bun install --frozen-lockfile 
RUN bunx prisma generate
RUN bun build ./src/index.ts --compile --outfile=idxstock
RUN mkdir ./logs

# run the app
USER bun
EXPOSE 3000/tcp
CMD ["./idxstock"]
