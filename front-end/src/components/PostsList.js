import { useOutletContext, } from "react-router-dom";
import React, { useEffect, useState } from 'react'
import Post from './Post';

const Posts = (props) => {

    const [posts, setPosts] = useState([])
    const [type, setType] = useState("user")
    const [groupId, setGroupId] = useState(0)

    const { setAlertMessage } = useOutletContext();
    const { setAlertClassName } = useOutletContext();

    const handlePostDelete = (postId) => {
        setPosts(posts.filter(post => post.id !== postId));
      }

    useEffect(() => {
        if (props.newPost) {
            setPosts([props.newPost, ...posts])
        }
    }, [props.newPost])

    useEffect(() => {
        const options = {
            method: 'GET',
            credentials: 'include',
        }

        let postsEndpoint = 'http://localhost:8080/posts/user'
        if (props) {
            if (props.userId) {
                postsEndpoint += `?id=${props.userId}`
            } else if (props.groupId) {
                setType("group")
                setGroupId(parseInt(props.groupId))
                postsEndpoint = `http://localhost:8080/posts/group?id=${props.groupId}`
            }
        }

        fetch(postsEndpoint, options)
            .then(response => response.json())
            .then(data => {
                if (data && data.error) {
                    setAlertClassName("alert-danger");
                    setAlertMessage(data.message);
                } else if (data) {
                    setPosts(data)
                }
            })
            .catch(error => { console.log(error) })
    }, [props])

    return (
        <div>
            {posts
                ? posts.map((post, index) => (
                    <Post
                        key={index}
                        post={post}
                        type={type}
                        groupId={groupId}
                        onDelete={() => handlePostDelete(post.id)}
                    />))
                : <p>No posts yet</p>
            }
        </div>
    )
}

export default Posts