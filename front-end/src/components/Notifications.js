import { useEffect, useState } from "react";
import { Link, useOutletContext } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faUserPlus, faCalendar, faUsers, faPersonRays } from '@fortawesome/free-solid-svg-icons'
const Notifications = () => {
    const { notifications, setNotifications, setData } = useOutletContext();
    const [loading, setLoading] = useState(false);
    const [actions, setActions] = useState({});
    const {setAlertClassName} = useOutletContext();
    const displayNotifications = (notification) => {
        switch (notification.type) {
            case "follow_request":
                return (
                    <div className={`notifications-message notification-container ${notification.seen === 0 ? "notification-unseen" : ""}`}>
                        <FontAwesomeIcon icon={faUserPlus} className="notification-icon"/>
                        <div className="notification-image-container">
                            <img src={`http://localhost:8080/profile/image/${notification.data.id}`} alt="profile image" className="notification-image img-fluid" />
                        </div>
                        <div className="notification-message-container">
                            <span className="notification-generic-message">
                                <Link to={`/profile/${notification.data.id}`} className="notification-id-link">
                                    {notification.data.firstName} {notification.data.lastName}
                                </Link> requested to follow you.</span>
                        </div>
                        {actions[notification.id] || notification.actioned == 1 ? (
                        ""
                    ) : (
                        <div className="d-flex justify-content-end notification-buttons">
                            <button onClick={() => {
                                notificationAction(notification, "accept")
                            }} className="btn btn-outline-success btn-sm notification-accept-button">Accept</button>
                            <button onClick={() => {
                                notificationAction(notification, "decline")
                            }} className="btn btn-outline-danger btn-sm notification-decline-button">Decline</button>
                        </div>
                    )}
                    </div>
                )
            case "group_event":
                return (
                    <div className={`notification-container ${notification.seen === 0 ? "notification-unseen" : ""} `}>
                        <div className="notification-image-container ">
                            <img src={`/media/group_images/${notification.data.img}`} alt="group image" className="notification-image img-fluid" />
                        </div>
                        <div className="notification-message-container">
                        <FontAwesomeIcon icon={faCalendar} className="notification-icon" />
                            <span className="notification-generic-message">New event in the group:</span>
                            <div className="notification-link-container">
                                <Link to={`/group/event/${notification.data.id}`} className="notification-id-link">
                                    {notification.data.groupName}
                                </Link>
                            </div>
                        </div>
                    </div>
                )
            case "group_join_request":
                return (
                    <div className={`notification-container ${notification.seen === 0 ? "notification-unseen" : ""} `}>
                        <div className="notification-image-container ">
                            <img src={`/media/group_images/${notification.data.img}`} alt="group image" className="notification-image img-fluid" />
                        </div>
                        <div className="notification-message-container">
                            <FontAwesomeIcon icon={faPersonRays} className="notification-icon" />
                            <span className="notification-generic-message">
                                <Link to={`/profile/${notification.optionalData.id}`} className="notification-id-link">
                                    {notification.optionalData.firstName} {notification.optionalData.lastName}
                                </Link> requested to join your group:</span>
                            <div className="notification-link-container">
                                <Link to={`/group/${notification.data.id}`} className="notification-id-link">
                                    {notification.data.title}
                                </Link>
                            </div>
                        </div>
                        {actions[notification.id] || notification.actioned == 1 ? (
                        ""
                    ) : (
                        <div className="d-flex justify-content-end notification-buttons">
                            <button onClick={() => {
                                notificationAction(notification, "accept")
                            }} className="btn btn-outline-success btn-sm notification-accept-button">Accept</button>
                            <button onClick={() => {
                                notificationAction(notification, "decline")
                            }} className="btn btn-outline-danger btn-sm notification-decline-button">Decline</button>
                        </div>
                    )}
                    </div>
                )
            case "group_invitation":
                return (
                    <div className={`notification-container ${notification.seen === 0 ? "notification-unseen" : ""} `}>
                        <div className="notification-image-container ">
                            <img src={`/media/group_images/${notification.data.img}`} alt="group image" className="notification-image img-fluid" />
                        </div>
                        <div className="notification-message-container">
                        <FontAwesomeIcon icon={faUsers} className="notification-icon" />
                            <span className="notification-generic-message">
                                <Link to={`/profile/${notification.optionalData.id}`} className="notification-id-link">
                                    {notification.optionalData.firstName} {notification.optionalData.lastName}
                                </Link> invited you to join the group: </span>
                            <div className="notification-link-container">
                                <Link to={`/group/${notification.data.id}`} className="notification-id-link">
                                    {notification.data.title}
                                </Link>
                            </div>
                        </div>
                        {actions[notification.id] || notification.actioned == 1 ? (
                        ""
                    ) : (
                        <div className="d-flex justify-content-end notification-buttons">
                            <button onClick={() => {
                                notificationAction(notification, "accept")
                            }} className="btn btn-outline-success btn-sm notification-accept-button">Accept</button>
                            <button onClick={() => {
                                notificationAction(notification, "decline")
                            }} className="btn btn-outline-danger btn-sm notification-decline-button">Decline</button>
                        </div>
                    )}
                    </div>
                )
        }
    }
    const notificationAction = (notification, action) => {
        fetch("http://localhost:8080/notificationAction", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ notification, action }),
        })
            .then((response) => response.json())
            .then(() => {
                setActions((prevState) => ({
                    ...prevState,
                    [notification.id]: true,
                }));
             })
            .catch((error) => { console.log(error) })
    }
    const markAsSeen = () => {
        fetch("http://localhost:8080/notificationsSeen", {
            method: "PUT",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => { })
            .catch((error) => { console.log(error) })
    }
    const clearAllNotifications = () => {
        fetch("http://localhost:8080/notifications", {
            method: "DELETE",
            credentials: "include",
        })
            .then((response) => response.json())
            .then(() => {
                setNotifications([]);
            })
            .catch((error) => { console.log(error) })
    } 
    useEffect(() => {
        setData(null);
        markAsSeen();
        setAlertClassName('d-none');
    }, [])
    return (
        <>
            {loading ? (<div>Loading...</div>) : (
                <div className="container">
                    <div className="row">
                        <div className="notifications-main-container">
                            <h2>Notifications</h2>
                            <hr />
                            <button onClick={clearAllNotifications} className="btn btn-outline-info btn-sm notifications-clear">Clear all</button>
                            { notifications.length === 0  ? (
                                <div className="notifications-none">
                                    <img src={"http://localhost:8080/media/various/noInfo.gif"} alt="no notifications" className="notifications-none-image" loop />
                                </div>
                            ) : (
                                notifications.map((notification, index) => (
                                    <div key={index}>
                                        {displayNotifications(notification)}
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}
export default Notifications;