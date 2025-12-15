import { useState, useCallback, useEffect, useRef } from "react"
import useWebSocket, { ReadyState } from "react-use-websocket"
import { debounce } from "lodash"

export type BlockType = "text" | "image"

export interface Block {
  id: string
  type: BlockType
  content: string
}

export interface Cursor {
  blockId: string
  offset: number
}

interface Message {
  type: string
  docId: string
  clientId: string
  block?: Block
  blocks?: Block[]
  cursor?: Cursor
}

export function useDocEditor(docId: string) {
  const [blocks, setBlocks] = useState<Block[]>([])
  const [clientId, setClientId] = useState<string>("")
  const [cursors, setCursors] = useState<Record<string, Cursor>>({})
  
  const { sendMessage, lastMessage, readyState } = useWebSocket(
    docId ? `ws://localhost:8080/ws?doc=${docId}` : null, 
    {
      shouldReconnect: () => true,
    }
  )

  useEffect(() => {
    if (lastMessage !== null) {
      const msg = JSON.parse(lastMessage.data) as Message

      switch (msg.type) {
        case "identity":
          setClientId(msg.clientId)
          break
        case "sync_document":
          if (msg.blocks) setBlocks(msg.blocks)
          break
        case "add_block":
          if (msg.block) {
            setBlocks((prev) => {
              if (prev.some((b) => b.id === msg.block!.id)) {
                return prev
              }
              return [...prev, msg.block!]
            })
          }
          break
        case "edit_block":
           if (msg.block) {
             setBlocks((prev) => 
               prev.map((b) => (b.id === msg.block!.id ? msg.block! : b))
             )
           }
           break
        case "cursor_move":
           if (msg.clientId && msg.cursor) {
             setCursors((prev) => ({
               ...prev,
               [msg.clientId]: msg.cursor!
             }))
           }
           break
        case "client_disconnected":
           if (msg.clientId) {
             setCursors((prev) => {
                const next = { ...prev }
                delete next[msg.clientId]
                return next
             })
           }
           break
      }
    }
  }, [lastMessage])

  const addBlock = useCallback((type: BlockType = "text") => {
    const newBlock: Block = {
      id: crypto.randomUUID(),
      type,
      content: "",
    }
    // Optimistic update
    setBlocks((prev) => [...prev, newBlock])
    
    sendMessage(JSON.stringify({
      type: "add_block",
      docId,
      block: newBlock
    }))
  }, [docId, sendMessage])

  const editBlock = useCallback((block: Block) => {
    // Optimistic update
    setBlocks((prev) => 
       prev.map((b) => (b.id === block.id ? block : b))
    )

    sendMessage(JSON.stringify({
      type: "edit_block",
      docId,
      block
    }))
  }, [docId, sendMessage])

  // Debounced move cursor
  const debouncedMoveCursor = useRef(
    debounce((blockId: string, offset: number) => {
        sendMessage(JSON.stringify({
            type: "cursor_move",
            docId,
            cursor: { blockId, offset }
        }))
    }, 100)
  ).current

  return {
    blocks,
    clientId,
    cursors,
    addBlock,
    editBlock,
    moveCursor: debouncedMoveCursor,
    isConnected: readyState === ReadyState.OPEN
  }
}
