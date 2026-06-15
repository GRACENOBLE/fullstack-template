import { vi, describe, it, expect } from "vitest";

// Mock @sentry/nextjs so tests don't try to connect
vi.mock("@sentry/nextjs", () => ({
  init: vi.fn(),
}));

describe("Sentry config", () => {
  it("client config imports without error", async () => {
    await expect(import("../sentry.client.config")).resolves.toBeDefined();
  });

  it("server config imports without error", async () => {
    await expect(import("../sentry.server.config")).resolves.toBeDefined();
  });

  it("edge config imports without error", async () => {
    await expect(import("../sentry.edge.config")).resolves.toBeDefined();
  });
});
