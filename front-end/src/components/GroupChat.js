import { useEffect, useState, useRef } from "react";
import { useParams } from 'react-router-dom';
import { Link, useOutletContext } from "react-router-dom";
import Input from "./form/Input";
import ActiveChat from "./ActiveChat";
import EmojiPicker from 'emoji-picker-react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faFaceSmile } from '@fortawesome/free-solid-svg-icons'

const GroupChat = ({ members}) => {
const { userId } = useOutletContext();
let {id} = useParams();

const gcRef = useRef(null)
const [ws, setSocket] = useState(null)
const [user, setUser] = useState({})
const [onlineUsers, setOnlineUsers] = useState([])
const [message, setMessage] = useState("")
const [messages, setMessages] = useState([])
const [showEmojis, setShowEmojis] = useState(false)

const sendMessage = () => {
    if (message.trim() !== "") {
        const messageObj = {
            action: "message",
            message: message,
            groupId: parseInt(id),
        }
        console.log(messageObj)
        ws.send(JSON.stringify(messageObj))
        setMessage("")
    }
}

const handleKeyDown = (event) => {
    if (event.key === 'Enter') {
        sendMessage()
    }
}

const toggleEmojis = () => {
    setShowEmojis(!showEmojis)
}
const handleEmojiClick = (emoji) => {
    const newMessage = message + emoji.emoji
    setMessage(newMessage)
}

const loadMessages = (groupId) => {
    fetch("http://localhost:8080/groupChat/" + groupId, {
        method: "GET",
        credentials: "include"
    })
        .then((response) => response.json())
        .then((data) => {
            console.log(data)
            setMessages(data)
        })
        .catch((error) => {
            console.log(error);
        });
}

const handleMessage = (event) => {
    const data = JSON.parse(event.data)
    switch (data.action) {
        case "new_message":
            loadMessages(id)
            break
        case "online_group_members":
            console.log("online members arr: ",data.onlineGroupMemebers)

            setOnlineUsers(data.onlineGroupMemebers)
            break
    }
}

useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/groupChat/ws')

        socket.addEventListener('open', () => {
            console.log('WebSocket groupChat connection opened')
        })

        window.addEventListener('beforeunload', () => {
            socket.send(JSON.stringify({ action: 'offline' }))
        })

        setSocket(socket)

        return () => {
            socket.send(JSON.stringify({ action: 'offline' }));
            socket.close()
        }
}, [])

useEffect(() => {
    loadMessages(id)
    if (ws) {
        ws.addEventListener("message", handleMessage)
        return () => {
            ws.removeEventListener("message", handleMessage)
        }
    }
}, [ws])

const currentlyOnline = (onlineMemberIds, groupMemberId) => {
    if (onlineMemberIds == null) {
      return false;
    }
    return onlineMemberIds.includes(groupMemberId);
}


    return (
        <div className="d-flex mb-4">
            <div className="w-75 bg-light rounded-3">
                <ActiveChat messages={messages} userId={userId} user={user} type={"group"}/>
                <Input
                type="text"
                className="form-control mt-2"
                placeholder="Message..."
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                onKeyDown={(event) => handleKeyDown(event)}
                ref={gcRef}
            />
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
                <div className="d-flex align-items-center mt-2">
                <FontAwesomeIcon icon={faFaceSmile} style={{fontSize: "24px"}} className="text-primary px-2" onClick={toggleEmojis} />
                    <button className="btn btn-primary mt-2 w-100" onClick={sendMessage}>
                        Send Message
                    </button>
                </div>
            </div>
            <div style={{zIndex: -1}} className="w-25 bg-light rounded-3 mt-2">
                <ul className="list-group">
                {members.map((member) => (
                    <div key={member.id}>
                        {currentlyOnline(onlineUsers, member.id) 
                        ? <li style={{color: "green"}} className="list-group-item">{member.firstName} {member.lastName}</li> 
                        : <li className="list-group-item">{member.firstName} {member.lastName}</li>}
                    </div>
                ))}
                </ul>
            </div>
        </div>
    )
}

export default GroupChat