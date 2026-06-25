import { describe, it, expect } from "vitest"
import { render, screen, within } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

import { DataTable } from "../DataTable"
import { userColumns, type User } from "../columns"

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function buildUsers(count: number): User[] {
  return Array.from({ length: count }, (_, i) => ({
    id: `u-${String(i + 1).padStart(3, "0")}`,
    name: `User ${String(i + 1).padStart(3, "0")}`,
    email: `user${i + 1}@example.com`,
    createdAt: `2024-01-${String((i % 28) + 1).padStart(2, "0")}T00:00:00Z`,
  }))
}

// ---------------------------------------------------------------------------
// Test 1 — initial render shows exactly 10 rows for a 25-item dataset
// ---------------------------------------------------------------------------

describe("DataTable — initial render", () => {
  it("shows 10 data rows by default for a 25-item dataset", () => {
    const data = buildUsers(25)
    render(<DataTable columns={userColumns} data={data} />)

    // tbody rows only (exclude header row)
    const tbody = screen.getByRole("table").querySelector("tbody")!
    const rows = within(tbody).getAllByRole("row")
    expect(rows).toHaveLength(10)
  })
})

// ---------------------------------------------------------------------------
// Test 2 — filter input narrows visible rows
// ---------------------------------------------------------------------------

describe("DataTable — filtering", () => {
  it("shows only rows matching the filter value", async () => {
    const user = userEvent.setup()
    const data = buildUsers(25)
    render(
      <DataTable
        columns={userColumns}
        data={data}
        filterColumn="name"
        filterPlaceholder="Search by name…"
      />
    )

    const input = screen.getByRole("textbox", { name: "Search by name…" })
    // "User 001" is the exact name; typing "001" should match only that row
    await user.type(input, "001")

    const tbody = screen.getByRole("table").querySelector("tbody")!
    const rows = within(tbody).getAllByRole("row")
    expect(rows).toHaveLength(1)
    expect(within(rows[0]).getByText("User 001")).toBeInTheDocument()
  })
})

// ---------------------------------------------------------------------------
// Test 3 — pagination advances pages; Previous is disabled on page 1
// ---------------------------------------------------------------------------

describe("DataTable — pagination", () => {
  it("Previous button is disabled on page 1", () => {
    const data = buildUsers(25)
    render(<DataTable columns={userColumns} data={data} />)

    const prevButton = screen.getByRole("button", { name: /previous page/i })
    expect(prevButton).toBeDisabled()
  })

  it("Next button advances to page 2 and Previous becomes enabled", async () => {
    const user = userEvent.setup()
    const data = buildUsers(25)
    render(<DataTable columns={userColumns} data={data} />)

    const nextButton = screen.getByRole("button", { name: /next page/i })
    await user.click(nextButton)

    expect(screen.getByText(/page 2 of/i)).toBeInTheDocument()

    const prevButton = screen.getByRole("button", { name: /previous page/i })
    expect(prevButton).not.toBeDisabled()
  })

  it("page 2 shows remaining rows (15 items → 5 on page 2)", async () => {
    const user = userEvent.setup()
    const data = buildUsers(15)
    render(<DataTable columns={userColumns} data={data} />)

    await user.click(screen.getByRole("button", { name: /next page/i }))

    const tbody = screen.getByRole("table").querySelector("tbody")!
    const rows = within(tbody).getAllByRole("row")
    expect(rows).toHaveLength(5)
  })
})

// ---------------------------------------------------------------------------
// Test 4 — sorting reorders rows when a sortable header is clicked
// ---------------------------------------------------------------------------

describe("DataTable — sorting", () => {
  it("sorts Name column ascending on first click", async () => {
    const user = userEvent.setup()
    // Build 5 users with names out of natural order to make sorting visible
    const data: User[] = [
      { id: "1", name: "Zara Adams", email: "z@example.com", createdAt: "2024-01-01T00:00:00Z" },
      { id: "2", name: "Alice Brooks", email: "a@example.com", createdAt: "2024-01-02T00:00:00Z" },
      { id: "3", name: "Mike Chen", email: "m@example.com", createdAt: "2024-01-03T00:00:00Z" },
    ]
    render(<DataTable columns={userColumns} data={data} />)

    const nameHeader = screen.getByRole("columnheader", { name: /name/i })
    await user.click(nameHeader)

    const tbody = screen.getByRole("table").querySelector("tbody")!
    const rows = within(tbody).getAllByRole("row")
    const firstCellText = within(rows[0]).getAllByRole("cell")[1].textContent
    expect(firstCellText).toBe("Alice Brooks")
  })

  it("sorts Name column descending on second click", async () => {
    const user = userEvent.setup()
    const data: User[] = [
      { id: "1", name: "Zara Adams", email: "z@example.com", createdAt: "2024-01-01T00:00:00Z" },
      { id: "2", name: "Alice Brooks", email: "a@example.com", createdAt: "2024-01-02T00:00:00Z" },
      { id: "3", name: "Mike Chen", email: "m@example.com", createdAt: "2024-01-03T00:00:00Z" },
    ]
    render(<DataTable columns={userColumns} data={data} />)

    const nameHeader = screen.getByRole("columnheader", { name: /name/i })
    await user.click(nameHeader) // asc
    await user.click(nameHeader) // desc

    const tbody = screen.getByRole("table").querySelector("tbody")!
    const rows = within(tbody).getAllByRole("row")
    const firstCellText = within(rows[0]).getAllByRole("cell")[1].textContent
    expect(firstCellText).toBe("Zara Adams")
  })
})
