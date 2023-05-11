import { useContext, useEffect, useState } from "react";
import { useParams, useOutletContext, useNavigate } from 'react-router-dom';
import Posts from "./PostsList";
import Followers from "./Followers";

const Profile = () => {
    let { id } = useParams();
    const { userId } = useOutletContext();
    const {setAlertClassName} = useOutletContext();

    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState({});
    const [isPublic, setIsPublic] = useState();
    const [isFollowed, setIsFollowed] = useState(false);
    const [isRequested, setIsRequested] = useState(false);
    const isCurrentUser = isNaN(id) || Number(id) === Number(userId);

    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"

    const changeProfileStatus = (status) => {
        fetch("http://localhost:8080/profile", {
            method: "PUT",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                profileP: status,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                setIsPublic(data.profileP == 1 ? true : false);
            })
            .catch((error) => {
                console.log(error);
            })
    }

    const fetchFollowStatus = () => {
        let profileEndpoint = 'http://localhost:8080/profile/followStatus/' + id

        fetch(profileEndpoint, {
            method: "GET",
            credentials: "include",
        })
        .then((response) => response.json())
        .then((data) => {
            setIsFollowed(data.isFollowing)
            setIsRequested(data.isRequested)
            // console.log(data);
        })
        .catch((error) => {
            console.log(error);
        })
    }

    const changeIsFollowed = () => {
        if (isPublic || isFollowed) {
            setIsFollowed(!isFollowed);
            fetchFollowData(!isFollowed, isRequested);
        } else {
            if (isRequested) {
                setIsRequested(false)
                setIsFollowed(false)
                fetchFollowData(false, false);
            } else {
                setIsRequested(true)
                setIsFollowed(false)
                fetchFollowData(false, true);
            }
        }

    }

    const fetchFollowData = (isFollowed, isRequested) => {
        fetch("http://localhost:8080/followUser", {
        method: "POST", 
        credentials: "include",
        headers: {
            "Content-Type": "application/json",     
        }, 
        body: JSON.stringify({
            currentUser: Number(userId),
            followedUser: Number(id),
            isRequested: isRequested,
            isFollowed: isFollowed,
        }),
    })
        .then(response => response.json())
        .then(data => {
            // console.log(data);
            setIsFollowed(data.isFollowed)
            setIsRequested(data.isRequested)
        })
        .catch(error => {
            console.log(error);
        })
    }

    const fetchData = () => {
        let profileEndpoint = 'http://localhost:8080/profile'
        if (id && !isCurrentUser) {
            profileEndpoint += `/${id}`
        } else {
            profileEndpoint += '/0'
        }

        fetch(profileEndpoint, {
            method: "GET",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                setUser(data);
                setIsPublic(data.profileP == 1 ? true : false);
                setLoading(false);

            })
            .catch((error) => {
                console.log(error);
            })
    }

    useEffect(() => {
        fetchData();
        if (!isCurrentUser) {
            fetchFollowStatus();
        }
        setAlertClassName('d-none');
    }, [id]);


    return (
        <>
            {loading ? "" : (
                <div className="container">
                    <div className="row">
                        <div className="profile-user-info">
                        <div className="col-md-6 profile-user-image-div">
                            <img src={`${profileImagesEndpoint}${user.image}`} alt="User Avatar" className="profile-user-image" />
                        </div>
                        <div className="col-md-6 profile-data">
                            {!isCurrentUser && (
                                <div className="d-flex justify-content-end">
                                    <button className="btn btn-warning me-0 btn-sm" onClick={() => { changeIsFollowed() }}>
                                        {isFollowed ? "Unfollow" : (isRequested ? "Requested" : "Follow")}
                                    </button>
                                </div>
                            )}
                            <h2>{user.firstName} {user.lastName}</h2>
                            <p className="profile-status">Profile: {isPublic ? "Public" : "Private"}</p>
                            <p>Email: {user.email}</p>
                            <p>Nickname: {user.nickname}</p>
                            <p>About Me: {user.about}</p>
                            <p>Created At: {new Date(user.createdAt).toLocaleDateString('en-GB')}</p>
                            {isCurrentUser && (
                                <div className="btn-group col-md-12 profile-status-buttons" role="group">
                                    <button
                                        type="button"
                                        className={`btn ${isPublic ? 'btn-warning' : 'btn-light profile-btn-border'} btn-sm` }
                                        onClick={() => {
                                            changeProfileStatus(1);
                                        }}
                                    >
                                        Public
                                    </button>
                                    <button
                                        type="button"
                                        className={`btn ${isPublic ? 'btn-light profile-btn-border' : 'btn-warning '} btn-sm profile-status-button`}
                                        onClick={() => {
                                            changeProfileStatus(0);
                                        }}
                                    >
                                        Private
                                    </button>
                                </div>
                            )}
                            </div>
                        </div>
                        <div className="col-md-12 followers-parent-div">
                            <hr />
                            {user.id && <Followers userId={user.id} />}
                        </div>
                            <hr />
                        <div className="col-md-12 posts-parent-div">
                            {user.id && <Posts userId={user.id} />}
                        </div>

                    </div>
                </div>
            )}
        </>
    )
}

export default Profile