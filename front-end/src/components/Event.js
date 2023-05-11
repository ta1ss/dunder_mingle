import { useEffect, useState } from "react";
import { useParams } from 'react-router-dom';
import { Link, useOutletContext } from "react-router-dom";

const Event = () => {
let {id} = useParams();
const { userId } = useOutletContext();
const [event, setEvent] = useState([])
const { setAlertMessage } = useOutletContext();
const { setAlertClassName } = useOutletContext();
const [showModal, setShowModal] = useState(false);
const [eventView, setEventView] = useState("going");

const handleOpenModal = () => {
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEventView("going");
  };


const viewGoing = () => {
    setEventView("going");
}

const viewNotGoing = () => {
    setEventView("notGoing");
}

const fetchEvent = () => {
    const options = {
        method: 'GET',
        credentials: 'include',
    }
    let getGroupEventEndpoint =  `http://localhost:8080/group/event?eventId=${id}`
    fetch(getGroupEventEndpoint, options)
        .then(response => response.json())
        .then(event => setEvent(event))
        .catch(error => console.log(error))
}

const goingStatus = (eventId, status) => {
    const options = {
        method: 'POST',
        credentials: 'include',
    }
    let goingStatusEndpoint =  `http://localhost:8080/group/eventStatus?eventId=${eventId}&status=${status}`
    fetch(goingStatusEndpoint, options)
    .then(response => response.json())
    .then(() => {
        setAlertMessage("");
        setAlertClassName("d-none");
        fetchEvent()
    })
    .catch(error => console.log(error));
}

const userInArr = (arr) => {
    if (arr != null){
        for (let i =0 ; i < arr.length; i++) {
            if (userId === arr[i].id){
                return true
            }
        }
    }
    return false
}

useEffect(() => {
    fetchEvent();
}, []);

    return (
        <div>
            <div className='post border rounded-3'>
                <h6 className='event-time'>{new Date(event.eventStart).toLocaleString('en-GB', {day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute:'2-digit'})} - {new Date(event.eventEnd).toLocaleString('en-GB', {day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute:'2-digit'})}</h6>
                <h2 className='m-0'>{event.title}</h2>
                <input id='groupId' type='hidden' value={event.id || ""} />
                <p className='m-0'>{event.description}</p>
                <p className='postInfo mb-1'>
                    Created by: <Link to={`/profile/${event.createdBy}`}>{event.creatorName}</Link>
                </p>
                <div className="goingDiv">
                    {userInArr(event.goingUsers) && <div className="btn-group " role="group">
                        <button type="button" className={`btn btn-primary`}>Going</button>
                        <button type="button" className={`btn btn-light`} onClick={() => {goingStatus(event.id, 0)}}>Not going</button>
                    </div>}
                    {userInArr(event.notGoingUsers) && <div className="btn-group " role="group">
                        <button type="button" className={`btn btn-light`} onClick={() => {goingStatus(event.id, 1)}}>Going</button>
                        <button type="button" className={`btn btn-primary`}>Not going</button>
                    </div>}
                    {!userInArr(event.goingUsers) && !userInArr(event.notGoingUsers) && <div className="btn-group " role="group">
                        <button type="button" className={`btn btn-light`} onClick={() => {goingStatus(event.id, 1)}}>Going</button>
                        <button type="button" className={`btn btn-light`} onClick={() => {goingStatus(event.id, 0)}}>Not going</button>
                    </div>}
                    <div className="modal-link" onClick={handleOpenModal}>
                    {event.goingUsers == null && event.notGoingUsers == null && <h6>Going: 0 | Not going: 0</h6>}
                    {event.notGoingUsers == null && event.goingUsers != null && <h6>Going: {event.goingUsers.length} | Not going: 0</h6>}
                    {event.goingUsers == null  && event.notGoingUsers != null && <h6>Going: 0 | Not going: {event.notGoingUsers.length}</h6>}
                    {event.goingUsers != null && event.notGoingUsers != null && <h6>Going: {event.goingUsers.length} | Not going: {event.notGoingUsers.length}</h6>}
                    </div>
                </div>
            </div>
            <div className={showModal ? 'modal show' : 'modal'}>
                <div className="modal-content">
                    <div className="modal-header">
                        <h2>Guests</h2>
                        <button style={{zIndex: 1}} onClick={handleCloseModal} type="button" className="btn-close" aria-label="Close"></button>
                    </div>
                    <div style={{ display: 'flex', justifyContent: "space-evenly" }} className="mt-3">
                    {event.goingUsers == null 
                    ? <div><h4 onClick={viewGoing} className={eventView == "going" ?  'active group-page-links' :  'group-page-links'}>Going (0)</h4></div>
                    : <div><h4 onClick={viewGoing} className={eventView == "going" ?  'active group-page-links' :  'group-page-links'}>Going ({event.goingUsers.length})</h4></div>}
                    {event.notGoingUsers == null 
                    ? <div><h4 onClick={viewNotGoing} className={eventView == "notGoing" ?  'active group-page-links' :  'group-page-links'}>Not going (0)</h4></div>
                    : <div><h4 onClick={viewNotGoing} className={eventView == "notGoing" ?  'active group-page-links' :  'group-page-links'}>Not going ({event.notGoingUsers.length})</h4></div>}
                    </div>
                    {eventView == "going" && event.goingUsers != null && <div>
                        {event.goingUsers.map((user) => (
                            <div key={user.id}>
                                <Link to={`/profile/${user.id}`} className="users-users-link">
                                    <div className='d-flex modal-user-container' >
                                        <img src={`http://localhost:8080/profile/image/${user.id}`} alt="profile image" className="users-user-image img-fluid" />
                                        <p className="users-user-name">{user.firstName} {user.lastName}</p>
                                    </div>
                                </Link>
                            </div>
                        ))}
                    </div>
                    }
                    {eventView == "notGoing" && event.notGoingUsers != null && <div>
                        {event.notGoingUsers.map((user) => (
                            <div key={user.id}>
                                <Link to={`/profile/${user.id}`} className="users-users-link">
                                    <div className='d-flex modal-user-container' >
                                        <img src={`http://localhost:8080/profile/image/${user.id}`} alt="profile image" className="users-user-image img-fluid" />
                                        <p className="users-user-name">{user.firstName} {user.lastName}</p>
                                    </div>
                                </Link>
                            </div>
                        ))}
                    </div>
                    }
                </div>
            </div>
        </div>
    )
}

export default Event