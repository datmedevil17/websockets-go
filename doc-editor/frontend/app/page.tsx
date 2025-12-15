"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useRouter } from "next/navigation"
import { useState } from "react"

export default function Home() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)

  const createDoc = async () => {
    setLoading(true)
    try {
      const res = await fetch("http://localhost:8080/docs", {
        method: "POST",
      })
      const data = await res.json()
      router.push(`/${data.docId}`)
    } catch (error) {
      console.error(error)
      alert("Failed to create document")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-muted/40">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Collab Docs</CardTitle>
          <CardDescription>Create a new collaborative document in seconds.</CardDescription>
        </CardHeader>
        <CardContent>
          <Button onClick={createDoc} disabled={loading} className="w-full">
            {loading ? "Creating..." : "Create New Document"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
