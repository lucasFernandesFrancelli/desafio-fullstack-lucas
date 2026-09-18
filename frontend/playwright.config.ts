import { defineConfig, devices } from "@playwright/test";

const FRONTEND_PORT = 5174;
const BACKEND_PORT = 8081;
const DATABASE_URL = "postgres://postgres:postgres@localhost:5432/ekaizen_e2e?sslmode=disable";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: [["list"], ["html", { open: "never" }]],
  use: {
    baseURL: `http://localhost:${FRONTEND_PORT}`,
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: [
    {
      command: "go run ./cmd/api",
      cwd: "../backend",
      url: `http://localhost:${BACKEND_PORT}/healthz`,
      env: {
        PORT: String(BACKEND_PORT),
        DATABASE_URL,
        JWT_SECRET: "e2e-secret",
        CORS_ORIGIN: `http://localhost:${FRONTEND_PORT}`,
        APP_ENV: "test",
      },
      reuseExistingServer: !process.env.CI,
      timeout: 60_000,
    },
    {
      command: `npm run dev -- --port ${FRONTEND_PORT} --strictPort`,
      url: `http://localhost:${FRONTEND_PORT}`,
      env: {
        VITE_API_BASE_URL: `http://localhost:${BACKEND_PORT}/api/v1`,
      },
      reuseExistingServer: !process.env.CI,
      timeout: 60_000,
    },
  ],
});
