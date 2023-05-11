import React, { useContext, useEffect, useState } from 'react'
import Posts from "./PostsList";
import Input from "./form/Input";
import Select from "./form/Select";
import Checkbox from "./form/Checkbox";
import Textarea from "./form/Textarea";
import { useOutletContext } from "react-router-dom";

const Home = () => {
    const [title, setTitle] = useState('');
    const [body, setBody] = useState('');
    const [privacy, setPostPrivacy] = useState("Public");
    const [image, setImage] = useState(null)
    const [imageName, setImageName] = useState("Add Image")
    const [followers, setFollowers] = useState([]);
    const [followersToggle, toggleFollowers] = useState("d-none");

    const [newPost, setNewPost] = useState(null);

    const { setAlertMessage, setAlertClassName } = useOutletContext();
    

    const clearForm = () => {
        setTitle('')
        setBody('')
        setImage(null)
        setImageName("Add Image")
        setPostPrivacy("Public")
        toggleFollowers("d-none")
        document.getElementById("postImgUpload").value = ""
    }

    const handleCheckbox = (event, index) => {
        const newFollowers = [...followers]
        newFollowers[index] = { ...newFollowers[index], checked: event.target.checked }
        setFollowers(newFollowers)
    }

    const handlePrivacy = (event) => {
        setPostPrivacy(event.target.value)
        if (event.target.value === "Custom") {
            fetchFollowers();
            toggleFollowers("");
        } else {
            toggleFollowers("d-none");
        }
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

    const addNewPostToDatabase = (post) => {
        const options = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(post)
        }
        fetch('http://localhost:8080/posts/user', options)
            .then(response => response.json())
            .then(data => {
                setAlertClassName("alert-success");
                setAlertMessage("Post successfully created!");
                setNewPost(data)
                clearForm()
            })
            .catch(error => {
                console.log(error)
                setAlertClassName("alert-danger");
                setAlertMessage(error);
            });
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
            setAlertMessage("Only .gif .jpg .jpeg .png allowed");
        } else if (image && image.size > 500000) {
            setAlertClassName("alert-danger");
            setAlertMessage("Maximum image size 500kB");
        } else {
            let targetIds = []
            if (privacy === "Custom") {
                followers.forEach(follower => {
                    if (follower.checked) { targetIds.push(follower.followerId) }
                });
            }
            const post = { title: title, body: body, privacy: privacy, customPrivacy: targetIds }
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

    useEffect(() => {
        setAlertClassName('d-none');
    }, []);

    return (
        <div>
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
                <div className="d-flex mt-2">
                    <Input
                        type="file"
                        name="postImgUpload"
                        className="form-control hidden"
                        onChange={(event) => {
                            setImage(event.target.files[0])
                            setImageName(event.target.files[0].name)
                        }}
                    />
                    <Select
                        id="privacy"
                        name="privacy"
                        className="form-select w-25"
                        options={[
                            { id: "Public", value: "Public" },
                            { id: "Private", value: "Private" },
                            { id: "Custom", value: "Custom" }
                        ]}
                        value={privacy}
                        onChange={handlePrivacy}
                    />
                    <button className="btn btn-outline-secondary ms-2 w-25 text-truncate" type='button' onClick={() => document.getElementById('postImgUpload').click()}>{imageName}</button>
                    <button type="submit" className="btn btn-primary ms-2 w-50">Post</button>
                </div>
                <div className={`${followersToggle} postFollowers bg-light rounded-3 p-2 mt-2`}>
                    {followers.length > 0
                        ? followers.map((follower, index) => (
                            <div key={index}>
                                <Checkbox
                                    name={`follower-${index}`}
                                    className="bg-dark"
                                    value={follower.followerId}
                                    title={follower.followerName}
                                    onChange={(event) => handleCheckbox(event, index)}
                                    checked={follower.checked}
                                />
                            </div>
                        ))
                        : <p>You have no followers.</p>}
                </div>
            </form>
            <Posts type="user" newPost={newPost} />
        </div>
    )
}
export default Home