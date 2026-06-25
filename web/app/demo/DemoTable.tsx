"use client"

import { DataTable } from "@/components/data-table"
import { userColumns, type User } from "@/components/data-table"

const users: User[] = [
  { id: "u-001", name: "Alice Martin", email: "alice.martin@example.com", createdAt: "2024-01-05T10:00:00Z" },
  { id: "u-002", name: "Bob Chen", email: "bob.chen@example.com", createdAt: "2024-01-12T09:15:00Z" },
  { id: "u-003", name: "Clara Osei", email: "clara.osei@example.com", createdAt: "2024-01-20T14:30:00Z" },
  { id: "u-004", name: "David Kim", email: "david.kim@example.com", createdAt: "2024-02-03T08:45:00Z" },
  { id: "u-005", name: "Evelyn Brooks", email: "evelyn.brooks@example.com", createdAt: "2024-02-14T11:00:00Z" },
  { id: "u-006", name: "Frank Müller", email: "frank.muller@example.com", createdAt: "2024-02-28T16:20:00Z" },
  { id: "u-007", name: "Grace Nakamura", email: "grace.nakamura@example.com", createdAt: "2024-03-07T13:10:00Z" },
  { id: "u-008", name: "Henry Okafor", email: "henry.okafor@example.com", createdAt: "2024-03-15T10:55:00Z" },
  { id: "u-009", name: "Isabel Santos", email: "isabel.santos@example.com", createdAt: "2024-03-22T09:30:00Z" },
  { id: "u-010", name: "James Patel", email: "james.patel@example.com", createdAt: "2024-04-01T12:00:00Z" },
  { id: "u-011", name: "Karen Liu", email: "karen.liu@example.com", createdAt: "2024-04-10T15:45:00Z" },
  { id: "u-012", name: "Liam Johansson", email: "liam.johansson@example.com", createdAt: "2024-04-18T08:00:00Z" },
  { id: "u-013", name: "Maya Rossi", email: "maya.rossi@example.com", createdAt: "2024-04-25T17:30:00Z" },
  { id: "u-014", name: "Nolan Wright", email: "nolan.wright@example.com", createdAt: "2024-05-02T10:15:00Z" },
  { id: "u-015", name: "Olivia Diallo", email: "olivia.diallo@example.com", createdAt: "2024-05-09T11:40:00Z" },
  { id: "u-016", name: "Peter Andersen", email: "peter.andersen@example.com", createdAt: "2024-05-17T14:00:00Z" },
  { id: "u-017", name: "Quinn Yamamoto", email: "quinn.yamamoto@example.com", createdAt: "2024-05-24T09:00:00Z" },
  { id: "u-018", name: "Rachel Nkosi", email: "rachel.nkosi@example.com", createdAt: "2024-06-01T13:25:00Z" },
  { id: "u-019", name: "Samuel Torres", email: "samuel.torres@example.com", createdAt: "2024-06-10T10:50:00Z" },
  { id: "u-020", name: "Tanya Kowalski", email: "tanya.kowalski@example.com", createdAt: "2024-06-18T16:10:00Z" },
  { id: "u-021", name: "Uma Fitzgerald", email: "uma.fitzgerald@example.com", createdAt: "2024-06-25T08:30:00Z" },
  { id: "u-022", name: "Victor Mensah", email: "victor.mensah@example.com", createdAt: "2024-07-03T11:20:00Z" },
  { id: "u-023", name: "Wendy Larsson", email: "wendy.larsson@example.com", createdAt: "2024-07-11T14:45:00Z" },
  { id: "u-024", name: "Xavier Dubois", email: "xavier.dubois@example.com", createdAt: "2024-07-19T09:05:00Z" },
  { id: "u-025", name: "Yuki Tanaka", email: "yuki.tanaka@example.com", createdAt: "2024-07-28T12:35:00Z" },
]

export function DemoTable() {
  return (
    <DataTable
      columns={userColumns}
      data={users}
      filterColumn="name"
      filterPlaceholder="Search by name…"
    />
  )
}
