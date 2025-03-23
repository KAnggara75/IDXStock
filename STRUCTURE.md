# Folder Structure

```md
📂 root/
├── 📂 src/
│   ├── 📂 config/         # Configuration (env, database, etc.)
│   ├── 📂 controllers/    # Handles HTTP requests
│   ├── 📂 routes/         # API route definitions
│   ├── 📂 middleware/     # Authentication, logging, etc.
│   ├── 📂 services/       # Business logic layer
│   ├── 📂 repositories/   # Database access layer using Prisma
│   ├── 📂 models/         # DTOs, Zod schemas, and TypeScript types
│   ├── 📂 utils/          # Helper functions (date formatting, hashing, etc.)
│   ├── app.ts            # Initializes Hono app
│   ├── server.ts         # Starts the Bun server
│   ├── prisma/           # Prisma schema and migrations
│   ├── tests/            # Unit and integration tests
│   ├── types.d.ts        # Additional TypeScript types
├── 📂 scripts/           # CLI scripts for seeding, migrations, etc.
├── .env                  # Environment variables
├── .gitignore            # Ignore unnecessary files in Git
├── bun.lockb             # Bun's lock file
├── package.json          # Project metadata and dependencies
├── tsconfig.json         # TypeScript configuration
└── README.md             # Project documentation
```