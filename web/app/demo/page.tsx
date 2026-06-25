import { DemoTable } from "./DemoTable"

export default function DemoPage() {
  return (
    <main className="mx-auto max-w-5xl px-4 py-10">
      <h1 className="text-2xl font-semibold text-foreground mb-2">
        Data Table Demo
      </h1>
      <p className="text-muted-foreground mb-8">
        A reusable table with sorting, filtering, and pagination. Showing 25
        users — 10 per page.
      </p>
      <DemoTable />
    </main>
  )
}
