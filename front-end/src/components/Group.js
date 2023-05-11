import { useEffect, useState } from "react";
import { useParams } from 'react-router-dom';
import { Link, useOutletContext } from "react-router-dom";
import Input from "./form/Input";
import PostsList from "./PostsList";
import Textarea from "./form/Textarea";
import Checkbox from "./form/Checkbox";
import Events from "./Events";
import GroupChat from "./GroupChat";

const Group = () => {
let {id} = useParams();
const { userId } = useOutletContext();

const [loading, setLoading] = useState(true);
const [group, setGroup] = useState({});
const [title, setTitle] = useState('');
const [body, setBody] = useState('');
const [image, setImage] = useState(null)
const [imageName, setImageName] = useState("Add Image")

const [groupId, setGroupId] = useState(null);
const [viewPage, setViewPage] = useState("posts");

const [showPostForm, setShowPostForm] = useState(false);
const [showInviteForm, setShowInviteForm] = useState(false);

const [followers, setFollowers] = useState([]);
const [followersToggle, toggleFollowers] = useState("d-none");

const [confirmLeaveGroup, setConfirmLeaveGroup] = useState({});
const [confirmRemoveUser, setConfirmRemoveUser] = useState({});

const { setAlertMessage } = useOutletContext();
const { setAlertClassName } = useOutletContext();

const toggleLeaveGroupBtn = (groupId) => {
    setConfirmLeaveGroup((prevState) => ({
        ...prevState,
        [groupId]: !prevState[groupId],
    }));
};

const toggleRemoveUserBtn = (userId) => {
    setConfirmRemoveUser((prevState) => ({
        ...prevState,
        [userId]: !prevState[userId],
    }));
};

const clearGroupPostForm = () => {
    setTitle('')
    setBody('')
    setImage(null)
    document.getElementById("postImgUpload").value = ""
}

const fetchGroup = () => {
    const options = {
        method: 'GET',
        credentials: 'include',
    }
    let groupEndpoint =  `http://localhost:8080/group?groupId=${id}`
    fetch(groupEndpoint, options)
        .then(response => response.json())
        .then(group => {
            setGroup(group)
            setLoading(false);
        })
        .catch(error => console.log(error))
}

const handleJoinRequest = () => {
    const options = {
        method: 'POST',
        credentials: 'include',
    }
    let groupJoinEndpoint =  `http://localhost:8080/group/join?groupId=${id}`
    fetch(groupJoinEndpoint, options)
    .then(response => {
        return response.json();
    })
    .then(() => {
        fetchGroup()
        setAlertClassName("alert-success");
        setAlertMessage("Join request sent!");
        setLoading(false);
    })
    .catch(error => console.log(error));
}

const removeJoinRequest = () => {
    const options = {
        method: 'DELETE',
        credentials: 'include',
    }
    let groupJoinEndpoint =  `http://localhost:8080/group/join?groupId=${id}`
    fetch(groupJoinEndpoint, options)
    .then(response => {
        return response.json();
    })
    .then(() => {
        fetchGroup()
        setAlertClassName("alert-success");
        setAlertMessage("Join request removed!");
        setLoading(false);
    })
    .catch(error => console.log(error));
}

const handleSubmit = (event) => {
    event.preventDefault();
    if (title === '') {
        setAlertClassName("alert-danger");
        setAlertMessage("Title can't be empty!");
    } else if (body === '') {
        setAlertClassName("alert-danger");
        setAlertMessage("Content can't be empty!");
    } else if (image && !image.name.match(/(gif|jpg|jpeg|png)$/gi)) {
        setAlertClassName("alert-danger");
        setAlertMessage("Only .gif .jpg .jpeg .png files are allowed");
    } else if (image && image.size > 500000) {
        setAlertClassName("alert-danger");
        setAlertMessage("Maximum 500kB");
    } else {
        const post = { groupId: group.id, title: title, body: body }
        if (image) {
            const reader = new FileReader()
            reader.onload = function () {
                post.img = reader.result
                addNewPostToDatabase(post)
            }
            reader.readAsDataURL(image)
        } else {
            addNewPostToDatabase(post)
        }
    }
}

const addNewPostToDatabase = (post) => {
    const options = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(post)
    }
    fetch('http://localhost:8080/posts/group', options)
        .then(response => {response.json()})
        .then(data => {
            setAlertClassName("alert-success");
            setAlertMessage("Post successfully created!");
            setGroupId(data)
            clearGroupPostForm()
        })
        .catch(error => {
            setAlertClassName("alert-danger");
            setAlertMessage(error);
        });
}

const viewGroupEvents = () => {
    setViewPage("events");
    setShowPostForm(false)
    setShowInviteForm(false)
    toggleFollowers("d-none");
}

const viewGroupPosts = () => {
    setViewPage("posts");
    setShowPostForm(false)
    setShowInviteForm(false)
    toggleFollowers("d-none");
}

const viewGroupMembers = () => {
    setViewPage("members");
    setShowPostForm(false)
    setShowInviteForm(false)
    toggleFollowers("d-none");
}

const viewGroupChat = () => {
    setViewPage("chat");
    setShowPostForm(false)
    setShowInviteForm(false)
    toggleFollowers("d-none");
}

const handleShowPostForm = () => {
    setShowPostForm(!showPostForm)
    setShowInviteForm(false)
    toggleFollowers("d-none");
}

const handleShowInviteForm = () => {
    setShowInviteForm(!showInviteForm)
    setShowPostForm(false)
    fetchFollowers();
    toggleFollowers("");
}

const handleCheckbox = (event, index) => {
    const newFollowers = [...filterFollowersById(followers, group.groupMembers)]
    newFollowers[index] = { ...newFollowers[index], checked: event.target.checked }
    setFollowers(newFollowers)
}

const handleInvite = (event) => {
    event.preventDefault();
    let targetIds = []
    followers.forEach(follower => {
        if (follower.checked) { 
            targetIds.push(follower.followerId) 
        }
    });
    if (targetIds.length === 0){
        setAlertClassName("alert-danger");
        setAlertMessage("No followers selected!");
        return
    }
    const invite = { userId: userId, groupId: parseInt(id), invitedUsersIds: targetIds }
    const options = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(invite)
    }
    let groupInviteEndpoint =  `http://localhost:8080/group/invite?groupId=${id}`
    fetch(groupInviteEndpoint, options)
    .then(response => {
        return response.json();
    })
    .then(() => {
        fetchGroup()
        setAlertClassName("alert-success");
        setAlertMessage("Invite sent!");
        setShowInviteForm(false)
        setLoading(false);
    })
    .catch(error => console.log(error));
}

const fetchFollowers = () => {
    const options = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
    }
    fetch('http://localhost:8080/followers', options)
        .then(response => response.json())
        .then(data => {
            if (data) {
                setFollowers(data.map((follower) => (
                    {
                        ...follower,
                        checked: false
                    }
                )))
            }
        })
        .catch(error => {
            setAlertClassName("alert-danger");
            setAlertMessage(error);
        });
}

const filterFollowersById = (followers, groupMembers) => {
    return followers.filter(follower => !groupMembers.some(member => member.id === follower.followerId));
}
// const alreadyInvited = (invitedUsersIds, followerId) => {
//     if (invitedUsersIds == null) {
//       return false;
//     }
//     return invitedUsersIds.includes(followerId);
// }


const alreadyInvited = (invitedUsersIds, followerId, invitingUsersIds, currentUserId) => {
    if (invitedUsersIds == null) {
      return false;
    }
    for (let i = 0; i < invitedUsersIds.length; i++){
        if (invitedUsersIds[i] === followerId && invitingUsersIds[i] === currentUserId){
            return true
        }
    }
    return false
}

const leaveGroup = () => {
    const options = {
        method: 'DELETE',
        credentials: 'include',
    }
    let leaveGroupEndpoint =  `http://localhost:8080/group/leave?groupId=${id}`
    fetch(leaveGroupEndpoint, options)
        .then(response => response.json())
        .then(data => {
            setAlertClassName("alert-success");
            setAlertMessage("Group left!");
            fetchGroup()
        })
        .catch(error => console.log(error))
}

const removeUser = (userId) => {
    const options = {
        method: 'DELETE',
        credentials: 'include',
        body: JSON.stringify({ id: groupId })
    }

    let removeUserEndpoint =  `http://localhost:8080/group/rmuser?groupId=${id}&userId=${userId}`
    fetch(removeUserEndpoint, options)
        .then(response => response.json())
        .then(data => {
            setAlertClassName("alert-dark");
            setAlertMessage("User removed!");
            fetchGroup()
        })
        .catch(error => console.log(error))
}

useEffect(() => {
    fetchGroup();
}, []);

console.log(group)
    return (
        <>
        {loading ? (<div>Loading...</div>) : (
        <div>
            <img src={group.img} alt="Group image" className="group-image" />
            <div className="d-flex justify-content-between align-items-center">
                <div>
                    <h1>{group.title}</h1>
                    <p className="w-75">{group.description}</p>
                </div>
                <div style={{minWidth: '133px', marginRight: "10px", float: 'right'}}>
                {group.createdBy != userId && group.inGroup && 
                                    (<>
                                        {confirmLeaveGroup[group.id] ? (
                                            <div>
                                                <button
                                                    type='button'
                                                    className='btn btn-light'
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        e.preventDefault();
                                                        toggleLeaveGroupBtn(group.id)
                                                    }}
                                                    >
                                                    Cancel
                                                </button>
                                                <button
                                                    type='button'
                                                    className='btn btn-primary'
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        e.preventDefault();
                                                        leaveGroup()
                                                    }}
                                                    >
                                                    Confirm
                                                </button>
                                            </div>
                                        ) : (
                                            <button
                                            type='button'
                                            className='btn btn-primary'
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                e.preventDefault();
                                                toggleLeaveGroupBtn(group.id)
                                            }}
                                            >
                                                Leave group
                                            </button>
                                        )}
                                    </>)}
                    {/* {group.createdBy != userId && group.inGroup &&  <button onClick={leaveGroup} className="btn btn-primary">Leave group</button>} */}
                    {group.id && !group.inGroup && !group.joinRequested && <button onClick={handleJoinRequest} className="btn btn-primary mt-2">Join Group</button>}
                    {group.id && !group.inGroup && group.joinRequested && <button onClick={removeJoinRequest} className="btn btn-primary mt-2 faded-button">Join Requested</button>}
                </div>
            </div>
            <hr />
            {group.id && group.inGroup && 
            <div>
                <div style={{ display: 'flex', justifyContent: "space-evenly" }}>
                    <h3 onClick={viewGroupPosts} className={viewPage == "posts" ?  'active group-page-links' :  'group-page-links'}>Posts</h3>
                    <h3 onClick={viewGroupEvents} className={viewPage == "events" ? 'active group-page-links' : 'group-page-links'}>Events</h3>
                    <h3 onClick={viewGroupMembers} className={viewPage == "members" ? 'active group-page-links' : 'group-page-links'}>Members</h3>
                    <h3 onClick={viewGroupChat} className={viewPage == "chat" ? 'active group-page-links' : 'group-page-links'}>Chat</h3>
                </div>
                <hr />
                {viewPage == "posts" && <div>
                    <button onClick={handleShowPostForm} className="btn btn-outline-primary mb-2 mt-2 w-100">+ Create new post</button>
                    {showPostForm && (
                    <form onSubmit={handleSubmit}>
                        <Input
                            placeholder="Title"
                            type="text"
                            className="form-control"
                            name="title"
                            autoComplete="title-new"
                            value={title}
                            onChange={(event) => setTitle(event.target.value)}
                        />
                        <Textarea
                            placeholder="Content..."
                            type="text"
                            className="form-control mt-2"
                            name="body"
                            autoComplete="body-new"
                            rows={3}
                            value={body}
                            onChange={(event) => setBody(event.target.value)}
                        />
                        <div className="d-flex mt-2 align-items-center">
                        <Input
                            type="file"
                            name="postImgUpload"
                            className="form-control hidden"
                            onChange={(event) => {
                                setImage(event.target.files[0])
                                setImageName(event.target.files[0].name)
                            }}
                        />
                        <button className="btn btn-outline-secondary ms-2 w-50 text-truncate" type='button' onClick={() => document.getElementById('postImgUpload').click()}>{imageName}</button>
                        <button type="submit" className="btn btn-primary ms-2 w-50">Post</button>
                        </div>
                    </form>
                    )}
                </div>}
                {viewPage == "events" && <Events />}
                {viewPage == "members" && <div>
                    <button onClick={handleShowInviteForm} className="btn btn-outline-primary mb-2 mt-2 w-100">+ Invite followers to group</button>          
                    {showInviteForm && (
                    <form onSubmit={handleInvite}>
                        <div className={`${followersToggle} postFollowers bg-light rounded-3 p-2 mt-2`}>
                            {filterFollowersById(followers, group.groupMembers).length > 0
                                ? filterFollowersById(followers, group.groupMembers).map((follower, index) => (
                                    <div key={index}>
                                        {!alreadyInvited(group.invitedUsersIds, follower.followerId, group.invitingUsersIds, userId) && <Checkbox
                                            name={`follower-${index}`}
                                            className="bg-dark"
                                            value={follower.followerId}
                                            title={follower.followerName}
                                            onChange={(event) => handleCheckbox(event, index)}
                                            checked={follower.checked}
                                        />}
                                        {alreadyInvited(group.invitedUsersIds, follower.followerId, group.invitingUsersIds, userId) && <Checkbox
                                            name={`follower-${index}`}
                                            className="bg-dark"
                                            value={follower.followerId}
                                            title={follower.followerName + " already invited"}
                                            disabled={true}
                                            onChange={e => {}}
                                            checked={false}
                                        />}
                                    </div>
                                ))
                                : <p>No followers to invite</p>}
                        </div>
                        <button type="submit" className="btn btn-primary mt-3 w-100">Invite</button>
                    </form>
                    )}       
                    {group.groupMembers
                        .map((user) => (
                            <div key={user.id}>
                                <Link to={`/profile/${user.id}`} className="users-users-link">
                                    <div className='d-flex users-user-container justify-content-between p-1'>
                                        <div className="d-flex">
                                            <img src={`http://localhost:8080/profile/image/${user.id}`} alt="profile image" className="users-user-image img-fluid" />
                                            {user.id == group.createdBy && <p className="users-user-name">{user.firstName} {user.lastName} (Group creator)</p>}
                                            {user.id != group.createdBy && <p className="users-user-name">{user.firstName} {user.lastName}</p>}
                                        </div>
                                        <div style={{marginRight: "10px"}}>
                                            {user.id != group.createdBy && group.createdBy == userId && 
                                            (<>
                                                {confirmRemoveUser[user.id] ? (
                                                    <div >
                                                        <button
                                                            type='button'
                                                            className='btn btn-sm btn-outline-secondary'
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                e.preventDefault();
                                                                toggleRemoveUserBtn(user.id)
                                                            }}
                                                            >
                                                            Cancel
                                                        </button>
                                                        <button
                                                            type='button'
                                                            className='btn btn-sm btn-outline-danger ms-2'
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                e.preventDefault();
                                                                removeUser(user.id)
                                                            }}
                                                            >
                                                            Confirm
                                                        </button>
                                                    </div>
                                                ) : (
                                                    <button
                                                    type='button'
                                                    className='btn btn-sm btn-outline-danger'
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        e.preventDefault();
                                                        toggleRemoveUserBtn(user.id)
                                                    }}
                                                    >
                                                        Remove
                                                    </button>
                                                )}
                                            </>)}
                                        </div>
                                    </div>
                                </Link>
                            </div>
                        ))}
                </div>}
                {viewPage == "chat" && <GroupChat members={group.groupMembers} />}
            </div>
            }
        </div>
        )}
        {group.id && group.inGroup && viewPage=="posts" && <div><PostsList groupId={group.id} /></div>}
        </>
    )
}

export default Group