import { useState, useEffect, useRef } from 'react'
import './App.css'

interface Message {
  type: string
  content?: string
  from?: string
  to?: string
  timestamp?: number
}

function App() {
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)
  const [username, setUsername] = useState('')
  const [room, setRoom] = useState('general')
  const [messages, setMessages] = useState<Message[]>([])
  const [messageInput, setMessageInput] = useState('')
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const connect = () => {
    if (!username.trim()) {
      alert('Please enter a username')
      return
    }

    const wsUrl = `ws://localhost:8080/ws?user_id=${encodeURIComponent(username)}&room=${encodeURIComponent(room)}`
    const ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('Connected to WebSocket')
      setConnected(true)
      setIsLoggedIn(true)
      setSocket(ws)
      
      // Send join message
      const joinMsg: Message = {
        type: 'join',
        from: username,
        to: room,
        timestamp: Date.now()
      }
      ws.send(JSON.stringify(joinMsg))
    }

    ws.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data)
        setMessages(prev => [...prev, message])
      } catch (err) {
        console.error('Error parsing message:', err)
      }
    }

    ws.onclose = () => {
      console.log('Disconnected from WebSocket')
      setConnected(false)
      setSocket(null)
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  const disconnect = () => {
    if (socket) {
      const leaveMsg: Message = {
        type: 'leave',
        from: username,
        to: room,
        timestamp: Date.now()
      }
      socket.send(JSON.stringify(leaveMsg))
      socket.close()
    }
    setIsLoggedIn(false)
    setMessages([])
  }

  const sendMessage = () => {
    if (!socket || !messageInput.trim()) return

    const message: Message = {
      type: 'chat',
      from: username,
      content: messageInput.trim(),
      to: room,
      timestamp: Date.now()
    }

    socket.send(JSON.stringify(message))
    setMessageInput('')
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  const formatTimestamp = (timestamp?: number) => {
    if (!timestamp) return ''
    return new Date(timestamp).toLocaleTimeString()
  }

  if (!isLoggedIn) {
    return (
      <div className="login-container">
        <div className="login-form">
          <h1>🦍 Gorilla Chat</h1>
          <div className="input-group">
            <label>Username:</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              onKeyPress={(e) => e.key === 'Enter' && connect()}
            />
          </div>
          <div className="input-group">
            <label>Room:</label>
            <input
              type="text"
              value={room}
              onChange={(e) => setRoom(e.target.value)}
              placeholder="Enter room name"
              onKeyPress={(e) => e.key === 'Enter' && connect()}
            />
          </div>
          <button onClick={connect} disabled={!username.trim()}>
            Join Chat
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h2>🦍 Gorilla Chat - Room: {room}</h2>
        <div className="user-info">
          <span className={`status ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? '🟢 Connected' : '🔴 Disconnected'}
          </span>
          <span>👤 {username}</span>
          <button onClick={disconnect}>Disconnect</button>
        </div>
      </div>
      
      <div className="messages-container">
        {messages.length === 0 ? (
          <div className="no-messages">No messages yet. Start the conversation!</div>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.type}`}>
              <div className="message-header">
                {msg.type === 'chat' && (
                  <>
                    <span className="sender">{msg.from}</span>
                    <span className="timestamp">{formatTimestamp(msg.timestamp)}</span>
                  </>
                )}
                {msg.type === 'join' && (
                  <span className="system-msg">👋 {msg.from} joined the room</span>
                )}
                {msg.type === 'leave' && (
                  <span className="system-msg">👋 {msg.from} left the room</span>
                )}
              </div>
              {msg.type === 'chat' && (
                <div className="message-content">{msg.content}</div>
              )}
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>
      
      <div className="message-input-container">
        <input
          type="text"
          value={messageInput}
          onChange={(e) => setMessageInput(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type a message..."
          disabled={!connected}
        />
        <button onClick={sendMessage} disabled={!connected || !messageInput.trim()}>
          Send
        </button>
      </div>
    </div>
  )
}

export default App
