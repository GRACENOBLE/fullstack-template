import { type ColumnDef } from "@tanstack/react-table"

export interface User {
  id: string
  name: string
  email: string
  createdAt: string
}

export const userColumns: ColumnDef<User, string>[] = [
  {
    accessorKey: "id",
    header: "ID",
    enableSorting: false,
  },
  {
    accessorKey: "name",
    header: "Name",
    enableSorting: true,
  },
  {
    accessorKey: "email",
    header: "Email",
    enableSorting: true,
  },
  {
    accessorKey: "createdAt",
    header: "Created At",
    enableSorting: false,
    cell: ({ getValue }) => {
      const raw = getValue()
      const date = new Date(raw)
      if (isNaN(date.getTime())) return raw
      const yyyy = date.getUTCFullYear()
      const mm = String(date.getUTCMonth() + 1).padStart(2, "0")
      const dd = String(date.getUTCDate()).padStart(2, "0")
      return `${yyyy}-${mm}-${dd}`
    },
  },
]
