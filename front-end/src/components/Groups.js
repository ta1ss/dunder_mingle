import React, { useEffect, useState } from 'react';
import Input from "./form/Input";
import { Link, useOutletContext } from "react-router-dom";

const Groups = () => {
    const { userId } = useOutletContext();

    const [showForm, setShowForm] = useState(false);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [img, setImg] = useState(null);
    const [groups, setGroups] = useState([])
    const [confirmDeleteGroup, setConfirmDeleteGroup] = useState({});


    const {setAlertMessage} = useOutletContext();
    const {setAlertClassName} = useOutletContext();

    const toggleDeleteGroupBtn = (groupId) => {
        setConfirmDeleteGroup((prevState) => ({
            ...prevState,
            [groupId]: !prevState[groupId],
        }));
    };

    const deleteGroup = (groupId) => {
        const options = {
            method: 'DELETE',
            credentials: 'include',
            body: JSON.stringify({ id: groupId })
        }

        fetch("http://localhost:8080/groups", options)
            .then(response => response.json())
            .then(data => {
                setAlertClassName("alert-dark");
                setAlertMessage("Group successfully deleted!");
                fetchGroups()
            })
            .catch(error => console.log(error))
    }

    const createGroup = (payload) => {
        const options = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(payload)
        };
        fetch('http://localhost:8080/groups', options)
            .then(() => {
                setAlertClassName("alert-success");
                setAlertMessage("Group successfully created!");
                setShowForm(!showForm);
                fetchGroups();
            })
            .catch(error => {
                setAlertClassName("alert alert-danger");
                setAlertMessage(error);
            });
    }

    const fetchGroups = () => {
        const options = {
            method: 'GET',
            credentials: 'include',
        }
        fetch('http://localhost:8080/groups', options)
            .then(response => response.json())
            .then(groups => setGroups(groups))
            .catch(error => console.log(error))
    }

    useEffect(() => {
        fetchGroups();
        setAlertClassName('d-none');
    }, []);

    const handleShowForm = () => {
        setShowForm(!showForm)
    }

    const handleSubmit = (event) => {
        let payload = {
            title: title,
            description: description
        }
        event.preventDefault();
        if (title === "" || description === ""){
            setAlertClassName("alert alert-danger");
            setAlertMessage("Empty fields");
        } else if (img !== null && !img.name.match(/(gif|jpg|jpeg|png)$/gi) ){
            setAlertClassName("alert alert-danger");
            setAlertMessage("File type invalid");
        } else if (img && img.size > 733333) {
            setAlertClassName("alert alert-danger");
            setAlertMessage("File too large");
        } else if (img) {
            const reader = new FileReader();
            reader.onload = function () {
                payload.img = reader.result;
                createGroup(payload)
            };
            reader.readAsDataURL(img)
        } else {
            createGroup(payload)
        }

    }

    return (
        <div>
            <button onClick={handleShowForm} className="btn btn-outline-primary mt-2 w-100">+ Create new group</button>
            {showForm && (
                <form onSubmit={handleSubmit}>
                <Input 
                    title="Title"
                    type="text"
                    className="form-control"
                    name="title"
                    autoComplete="title-new"
                    onChange={(event) => setTitle(event.target.value)}
                />
                <Input 
                    title="Description"
                    type="text"
                    className="form-control"
                    name="description"
                    autoComplete="description-new"
                    onChange={(event) => setDescription(event.target.value)}
                />
                <Input 
                    title="Group picture"
                    type="file"
                    className="form-control"
                    name="img"
                    onChange={(event) => setImg(event.target.files[0])}
                />
                <hr />
                <input 
                    type="submit"
                    className="btn btn-primary"
                    value="Create"
                />
            </form>
            )}
            <div>
                {groups.map((group, index) => (
                    <Link to={`/group/${group.id}`} style={{ textDecoration: 'none', color: '#212529' }} key={index}>
                        <div className='post border rounded-3'>
                            <div className='d-flex justify-content-between'>
                            <h5 id='group' className='m-0'>{group.title}</h5>
                                {userId && userId === group.createdBy &&
                                    (<>
                                        {confirmDeleteGroup[group.id] ? (
                                            <div>
                                                <button
                                                    type='button'
                                                    className='btn btn-sm btn-outline-secondary'
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        e.preventDefault();
                                                        toggleDeleteGroupBtn(group.id)
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
                                                        deleteGroup(group.id)
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
                                                toggleDeleteGroupBtn(group.id)
                                            }}
                                            >
                                                Delete
                                            </button>
                                        )}
                                    </>)}
                            </div>
                            <input id='groupId' type='hidden' value={group.id} />
                            <p className='m-0'>{group.description}</p>
                            <p className='postInfo mb-1'>Members: {group.groupMembers.length}</p>
                        </div>
                    </Link> 
                ))}
            </div>
        </div>
    )
}

export default Groups