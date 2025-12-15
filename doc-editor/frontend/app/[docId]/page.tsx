"use client"

import { useDocEditor, Block } from "@/hooks/use-doc-editor"
import { useParams } from "next/navigation"
import { useEffect, useRef } from "react"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import { Textarea } from "@/components/ui/textarea"

function BlockEditor({ 
  block, 
  onEdit, 
  onFocus, 
  cursors,
  onEnter
}: { 
  block: Block
  onEdit: (content: string) => void
  onFocus: () => void
  cursors: string[]
  onEnter: () => void
}) {
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto"
      textareaRef.current.style.height = textareaRef.current.scrollHeight + "px"
    }
  }, [block.content])

  return (
    <div className="relative group my-1">
      {cursors.length > 0 && (
        <div className="absolute -left-4 top-2 flex flex-col gap-1">
           {cursors.map(cid => (
             <div key={cid} className="w-2 h-2 rounded-full bg-blue-500" title={`User ${cid}`} />
           ))}
        </div>
      )}
      <Textarea
        ref={textareaRef}
        value={block.content}
        onChange={(e) => onEdit(e.target.value)}
        onFocus={onFocus}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            onEnter()
          }
        }}
        className={cn(
           "w-full border-none shadow-none focus-visible:ring-0 px-0 text-lg resize-none min-h-[1.5em] overflow-hidden leading-relaxed py-1 rounded-none bg-transparent hover:bg-muted/10 transition-colors",
           cursors.length > 0 ? "border-l-2 border-blue-500 pl-2 bg-blue-50/10" : ""
        )}
        placeholder="Type '/' for commands"
        rows={1}
      />
    </div>
  )
}

export default function DocPage() {
  const params = useParams()
  const docId = params.docId as string
  const { blocks, addBlock, editBlock, moveCursor, cursors, isConnected, clientId } = useDocEditor(docId)

  return (
    <div className="max-w-3xl mx-auto py-12 px-4">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">Document {docId.slice(0, 8)}</h1>
        <div className="flex items-center gap-2">
           <div className={`w-3 h-3 rounded-full ${isConnected ? "bg-green-500" : "bg-red-500"}`} />
           <span className="text-sm text-muted-foreground">{isConnected ? "Connected" : "Disconnected"}</span>
        </div>
      </div>
      
      <Card className="min-h-[500px] p-8 shadow-sm">
        {blocks.map((block) => {
           // Find cursors on this block
           const blockCursors = Object.entries(cursors)
             .filter(([cid, cursor]) => cursor.blockId === block.id && cid !== clientId)
             .map(([cid]) => cid)

           return (
             <BlockEditor 
               key={block.id} 
               block={block} 
               onEdit={(content) => editBlock({ ...block, content })}
               onFocus={() => moveCursor(block.id, 0)} // just tracking focus on block for now
               cursors={blockCursors}
               onEnter={() => addBlock("text")}
             />
           )
        })}

        <Button variant="ghost" className="w-full mt-4 justify-start text-muted-foreground" onClick={() => addBlock("text")}>
          + Add Block
        </Button>
      </Card>
    </div>
  )
}
