import { Link, Outlet, useNavigate } from "react-router-dom";
import React, { useContext, useEffect, useState } from "react";
import Alert from "./components/form/Alert";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBell, faArrowRightFromBracket, faHouse, faAddressCard, faUsersLine, faUsersViewfinder } from '@fortawesome/free-solid-svg-icons'
import Messenger from "./components/Messenger";
import Welcome from "./components/Welcome";

function App() {

  const [loggedIn, setLoggedIn] = useState(null);
  const [userId, setUserId] = useState(null);
  const [alertMessage, setAlertMessage] = useState("");
  const [alertClassName, setAlertClassName] = useState("d-none");
  const [notifications, setNotifications] = useState([]);
  const [lastTimeStamp, setLastTimeStamp] = useState(null);
  const [data, setData] = useState(null);
  const navigate = useNavigate();

  const getNotifications = () => {
    fetch(`http://localhost:8080/notifications`, {
      method: "GET",
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        if (data && data.length > 0) {
          const sortedNotifications = data.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
          setNotifications(sortedNotifications);
          setLastTimeStamp(sortedNotifications[0].createdAt);
        }
      })
      .catch((error) => {
        console.log(error);
      })
  }

  useEffect(() => {
    if (loggedIn === false || loggedIn === null) {
      return;
    }
    regularFetchNotifications();

    const interval = setInterval(() => {
      regularFetchNotifications();
    }, 10000)
    return () => clearInterval(interval);
  }, [lastTimeStamp, loggedIn])

  const regularFetchNotifications = () => {
    fetch(`http://localhost:8080/notifications?lastTimeStamp=${lastTimeStamp}`, {
      method: "GET",
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        if (data && data.length > 0) {
          const sortedNotifications = data.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
          setNotifications(prevNotifications => [...sortedNotifications, ...prevNotifications]);
          setLastTimeStamp(sortedNotifications[0].createdAt);
          sortedNotifications.map((notification) => {
            if (notification.seen === 0) {
              setData("New notification")
            }
          })
        }
      })
      .catch((error) => { console.log(error) })
  }

  const handleNotificationsClick = () => {
    getNotifications();
    navigate("/notifications");
  }

  const isLoggedIn = () => {
    const cookies = document.cookie.split(";");
    for (const cookie of cookies) {
      const [name, value] = cookie.split("=");
      if (name === "SN-Session" && value) {
        fetch("http://localhost:8080/profile/0", {
          method: "GET",
          credentials: "include",
        })
          .then((response) => response.json())
          .then((data) => {
            setUserId(data.id);
            setLoggedIn(true);
          })
          .catch((error) => {
            console.log(error);
          });
      } else {
        setLoggedIn(false);
      }
    }
  }

  const logOut = () => {
    fetch("http://localhost:8080/logout", {
      method: "POST",
      credentials: "include",
    })
      .then(response => response.json())
    console.log("user logged out")
    setUserId(null);
    setLoggedIn(false);
    setAlertClassName('d-none');
    setNotifications([])
  }

  useEffect(() => {
    isLoggedIn()
  }, []);

  return (
    <div className="container pb-3">
      {loggedIn === null ? (
        ""
      ) : loggedIn ? (
        <>
          <div className="row">
            <div className="col mt-2 ">
              <div className=" header-instruments-left" >
                <Link to="/" className="users-users-link">
                <img src={"http://localhost:8080/media/various/logo.png"} alt="Social Network Logo" className="welcome-logo-header" />
                </Link>
                <h1 className="mt-4 header-name">dunder mingle</h1>
              </div>
            </div>
            <div className="col text-end">
              <div className="header-instruments-right mt-5">
                {data === null ? (
                  <FontAwesomeIcon icon={faBell} onClick={handleNotificationsClick} className="notifications-icon" />
                ) : (
                  <FontAwesomeIcon icon={faBell} style={{ color: "#ff0000", }} onClick={handleNotificationsClick} className="notifications-icon" />
                )}
                <Link to="/#" onClick={logOut}>
                  <FontAwesomeIcon icon={faArrowRightFromBracket} className="log-out" />
                </Link>
              </div>
            </div>
            <hr className="mb-3" />
          </div>

          <div className="row">
            <div className="col-md-2">
              <nav>
                <div className="list-group text-center">
                  <Link to="/" className="list-group-item list-group-item-action"> <FontAwesomeIcon icon={faHouse} /> Home </Link>
                  <Link to="/profile/me" className="list-group-item list-group-item-action"> <FontAwesomeIcon icon={faAddressCard} /> Profile</Link>
                  <Link to="/users" className="list-group-item list-group-item-action"> <FontAwesomeIcon icon={faUsersLine} /> Users</Link>
                  <Link to="/groups" className="list-group-item list-group-item-action"> <FontAwesomeIcon icon={faUsersViewfinder} /> Groups</Link>
                </div>
              </nav>
            </div>
            <div className="col-md-7">
              <Alert
                message={alertMessage}
                className={alertClassName}
              />
              <Outlet context={{
                loggedIn,
                setLoggedIn,
                userId,
                setUserId,
                setData,
                setAlertMessage,
                setAlertClassName,
                notifications,
                setNotifications,
                getNotifications,
              }} />
            </div>
            <div className="col-md-3">
              <div className="rounded-3">
                <Messenger userId={userId} />
              </div>
            </div>
          </div>
        </>
      ) : (
        <>
          <Welcome
            setLoggedIn={setLoggedIn}
            setAlertClassName={setAlertClassName}
            setAlertMessage={setAlertMessage}
            setUserId={setUserId}
            alertMessage={alertMessage}
            alertClassName={alertClassName}
          />
        </>
      )}
    </div>
  );
}
export default App;