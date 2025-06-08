FROM oven/bun:alpine AS base

LABEL authors="KAnggara75"
LABEL maintainer="kanggara75@gmail.com"

WORKDIR /usr/app

# install dependencies into temp directory
# this will cache them and speed up future builds
FROM base AS install
WORKDIR /temp/prd

# install with --production (exclude devDependencies)
COPY . .

RUN bun install --frozen-lockfile --production

RUN bunx prisma generate

RUN bun build \
    --compile \
    --minify-whitespace \
    --minify-syntax \
    --target bun \
    --outfile app \
    ./src/index.ts


# copy production dependencies and source code into final image
FROM base AS release
COPY --from=install /temp/prd/node_modules/.prisma /temp/prd/node_modules/.prisma
COPY --from=install /temp/prd/app .

# [optional] tests & build
ENV NODE_ENV=production

# run the app
USER bun

EXPOSE 3000/tcp

CMD ["./app"]
