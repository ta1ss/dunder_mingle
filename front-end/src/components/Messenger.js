import Input from "./form/Input";
import ActiveChat from "./ActiveChat";
import EmojiPicker from 'emoji-picker-react';
import { useEffect, useState, useRef } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faFaceSmile } from '@fortawesome/free-solid-svg-icons'
import { Link, useOutletContext } from "react-router-dom";

const Messenger = ({ userId }) => {

    const msgRef = useRef(null)
    const [user, setUser] = useState({})
    const [ws, setSocket] = useState(null)
    const [users, setUsers] = useState([])
    const [message, setMessage] = useState("")
    const [messages, setMessages] = useState([])
    const [recipient, setRecipient] = useState(null)
    const [showEmojis, setShowEmojis] = useState(false)

    const sendMessage = () => {
        if (recipient && message.trim() !== "") {
            ws.send(JSON.stringify({ action: "message", message: message, targetId: recipient }))
            setMessage("")
        }
    }

    const fetchUsers = () => {
        fetch("http://localhost:8080/messenger/users", {
            method: "GET",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                setUsers(data)
            })
            .catch((error) => {
                console.log(error)
            })
    }

    const handleClick = (index) => {
        loadMessages(users[index].id)
        users[index].messageRead = 1
        ws.send(JSON.stringify({ action: "read_messages", senderId: users[index].id }))
    }

    const handleKeyDown = (event) => {
        if (event.key === 'Enter') {
            sendMessage()
        }
    }

    const loadMessages = (targetId) => {
        if (targetId !== recipient) {
            fetch("http://localhost:8080/messenger/messages/" + targetId, {
                method: "GET",
                credentials: "include"
            })
                .then((response) => response.json())
                .then((data) => {
                    setMessages(data)
                    setRecipient(targetId)
                })
                .catch((error) => {
                    console.log(error);
                });
        }
    }

    const receiveMessage = (msg) => {
        setUsers((prevUsers) => {
            if (recipient === userId) { // from current user to current user
                const currentUser = prevUsers.find((user) => user.id === recipient)
                const otherUsers = prevUsers.filter((user) => user.id !== recipient)
                return [currentUser, ...otherUsers]
            } else if (msg.targetId === userId) { // from other user to current user
                const sender = prevUsers.find((user) => user.id === msg.senderId)
                const otherUsers = prevUsers.filter((user) => user.id !== msg.senderId)
                if (msg.senderId === recipient) {
                    sender.messageRead = 1
                    ws.send(JSON.stringify({ action: "read_messages", senderId: msg.senderId }))
                } else {
                    sender.messageRead = 0
                }
                return [sender, ...otherUsers]
            } else if (msg.senderId === userId) { // from current user to other user
                const target = prevUsers.find((user) => user.id === msg.targetId)
                const otherUsers = prevUsers.filter((user) => user.id !== msg.targetId)
                if (msg.targetId === recipient) {
                    target.messageRead = 1
                    ws.send(JSON.stringify({ action: "read_messages", senderId: msg.senderId }))
                }
                return [target, ...otherUsers]
            }
        })
        if (msg.senderId === recipient || msg.targetId === recipient) {
            setMessages((prevMessages) => {
                if (prevMessages) {
                    return [msg, ...prevMessages]
                } else {
                    return [msg]
                }
            })
        }
    }

    const toggleEmojis = () => {
        setShowEmojis(!showEmojis)
    }

    const updateOnlineUsers = (onlineUsers) => {
        setUsers((prevUsers) =>
            prevUsers.map((user) => {
                return { ...user, online: onlineUsers.includes(user.id) ? 1 : 0 }
            })
        )
    }

    const handleEmojiClick = (emoji) => {
        const newMessage = message + emoji.emoji
        setMessage(newMessage)
    }

    const handleMessage = (event) => {
        const data = JSON.parse(event.data)
        switch (data.action) {
            case "new_message":
                receiveMessage(data.userMessage)
                break
            case "online_users":
                updateOnlineUsers(data.onlineUsers)
                break
        }
    }

    useEffect(() => {
        if (userId) {
            fetchUsers()

            const socket = new WebSocket('ws://localhost:8080/messenger/ws')

            socket.addEventListener('open', () => {
                console.log('WebSocket connection opened')
            })

            window.addEventListener('beforeunload', () => {
                socket.send(JSON.stringify({ action: 'offline' }))
            })

            setSocket(socket)

            return () => {
                socket.send(JSON.stringify({ action: 'offline' }));
                socket.close()
            }
        }
    }, [userId])

    useEffect(() => {
        if (ws) {
            ws.addEventListener("message", handleMessage)
            return () => {
                ws.removeEventListener("message", handleMessage)
            }
        }
    }, [ws, recipient])

    return (
        <>
            <ul id="messengerUsersList" className="list-group rounded-3 border">
                {users && users.length > 0 && users.map((user, index) => (
                    <li
                        key={index}
                        className={`list-group-item pointer border-0 border-bottom ${user.online && user.online === 1 ? 'text-success' : ''}`}
                        onClick={() => {
                            if (recipient !== user.id) {
                                handleClick(index)
                                setUser(user)
                            } else {
                                setRecipient(null)
                                setShowEmojis(false)
                            }
                        }}
                    >
                        {user.firstName} {user.lastName} <span className="text-primary">{user.messageRead === 0 ? "(new)" : ""}</span>
                    </li>
                ))}
            </ul>
            {recipient &&
                <>
                    <ActiveChat messages={messages} userId={userId} user={user} type={"user"} />
                    <Input
                        type="text"
                        name="message"
                        className="form-control mt-2"
                        placeholder="Message..."
                        value={message}
                        onChange={(event) => setMessage(event.target.value)}
                        onKeyDown={(event) => handleKeyDown(event)}
                        ref={msgRef}
                    />
                    <div className="d-flex align-items-center mt-2">
                        <FontAwesomeIcon icon={faFaceSmile} style={{fontSize: "24px"}} className="text-primary px-2" onClick={toggleEmojis} />
                        <button className="btn btn-primary w-100" onClick={sendMessage}>
                            Send Message
                        </button>
                    </div>
                    {showEmojis &&
                        <div className="mt-2">
                            <EmojiPicker
                                searchDisabled={true}
                                skinTonesDisabled={true}
                                previewConfig={{
                                    showPreview: false
                                }}
                                width="100%"
                                height={200}
                                onEmojiClick={handleEmojiClick}
                            />
                        </div>
                    }
                </>
            }
        </>
    )
}

export default Messenger;