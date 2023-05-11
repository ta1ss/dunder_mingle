import { useEffect, useState } from "react";
import { Link, useOutletContext } from "react-router-dom";

const ActiveChat = ({ messages, userId, user, type }) => {
    const [displayedMessages, setDisplayedMessages] = useState([])
    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"

    function formatTime(timeString) {
        const date = new Date(timeString)
        const now = new Date()

        if (
            date.getFullYear() === now.getFullYear() &&
            date.getMonth() === now.getMonth() &&
            date.getDate() === now.getDate()
        ) {
            const hours = ('0' + date.getHours()).slice(-2)
            const minutes = ('0' + date.getMinutes()).slice(-2)
            return `${hours}:${minutes}`
        } else {
            const month = ('0' + (date.getMonth() + 1)).slice(-2)
            const day = ('0' + date.getDate()).slice(-2)
            const year = date.getFullYear().toString().slice(-2)
            return `${day}.${month}.${year}`
        }
    }

    useEffect(() => {
        if (messages && userId) {
            setDisplayedMessages(messages)
        }
    }, [messages]);

    return (
        <>
            {user &&
                <div className="rounded-3 border mt-2">
                     {type === "user" && <div className="px-3 py-2 bg-white border-bottom" style={{borderTopLeftRadius: "7px", borderTopRightRadius: "7px"}}>
                        <img src={`${profileImagesEndpoint}${user.image}`} className="messengerUserImg" />
                        {user.firstName} {user.lastName}
                    </div>}
                    <div id="activeChat" className="p-3 bg-white" style={{borderBottomLeftRadius: "7px", borderBottomRightRadius: "7px"}}>
                        {messages && displayedMessages.length > 0 && displayedMessages.map((message, index) => (
                            (userId === message.senderId)
                                ? <div key={index} className="msg-right">
                                    <p className="msg-body m-0">{message.body}</p>
                                    <p className="msg-date m-0">{formatTime(message.CreatedAt)}</p>
                                </div>
                                : <div key={index} className="msg-left">
                                    {type === "group" && <p style={{fontWeight: "bold"}} className="m-0">{message.userName}</p>}
                                    <p className="msg-body m-0">{message.body}</p>
                                    <p className="msg-date m-0">{formatTime(message.CreatedAt)}</p>
                                </div>
                        ))}
                    </div>
                </div>
            }
        </>
    )
}

export default ActiveChat;