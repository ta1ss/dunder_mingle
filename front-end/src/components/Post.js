import { Link, useOutletContext } from "react-router-dom";
import React, { useState, useEffect } from 'react'

const Post = ({ post, type, groupId, onDelete }) => {

    const { userId } = useOutletContext();
    const [confirmDeletePost, setConfirmDeletePost] = useState({});
    const [commentsURL, setCommentsURL] = useState("");

    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"
    const postImagesEndpoint = "http://localhost:8080/media/post_images/"

    const { setAlertMessage } = useOutletContext();
    const { setAlertClassName } = useOutletContext();
    
    const toggleDeletePostBtn = (postId) => {
        setConfirmDeletePost((prevState) => ({
            ...prevState,
            [postId]: !prevState[postId],
        }));
    };

    const deletePost = (postId) => {
        const options = {
            method: 'DELETE',
            credentials: 'include',
            body: JSON.stringify({ id: postId, userId: userId })
        }

        let deletePostEndpoint = "http://localhost:8080/posts/"
        if (type === "detail" && !groupId) {
            deletePostEndpoint += "user"
        } else if (type === "detail" && groupId) {
            deletePostEndpoint += "group"
        } else {
            deletePostEndpoint += type
        }

        fetch(deletePostEndpoint, options)
            .then(response => response.json())
            .then(data => {
                setAlertClassName("alert-dark");
                setAlertMessage("Post successfully deleted!");
                onDelete()
            })
            .catch(error => console.log(error))
    }

    useEffect(() => {
        if (type === "user") {
            setCommentsURL(`/post/${post.id}`)
        } else if (type === "group" && groupId !== 0) {
            setCommentsURL(`/group/${groupId}/post/${post.id}`)
        }
    }, [post])

    return (
        <div className='post border rounded-3'>
            <div className="d-flex">
                <div className='col-md-9 d-flex'>
                    <Link to={`/profile/${post.userId}`}>
                        <img src={`${profileImagesEndpoint}${post.userImg}`} className="postUserImg" />
                    </Link>
                    <div>
                        <h5 className='m-0'>{post.title}</h5>
                        <input id='postId' type='hidden' value={post.id} />
                        <input id='userId' type='hidden' value={post.userId} />
                        <p className='postInfo mb-1'>
                            <Link to={`/profile/${post.userId}`}>{post.createdBy}</Link> {new Date(post.createdAt).toLocaleString('en-GB', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })}
                        </p>
                    </div>
                </div>
                <div className='col-md-3 text-end'>
                    {userId && userId === post.userId &&
                        (<>
                            {confirmDeletePost[post.id] ? (
                                <div className="d-flex justify-content-end">
                                    <button
                                        type='button'
                                        className='btn btn-sm btn-outline-secondary'
                                        onClick={() => toggleDeletePostBtn(post.id)}
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type='button'
                                        className='btn btn-sm btn-outline-danger ms-2'
                                        onClick={() => deletePost(post.id)}
                                    >
                                        Confirm
                                    </button>
                                </div>
                            ) : (
                                <button
                                    type='button'
                                    className='btn btn-sm btn-outline-danger'
                                    onClick={() => toggleDeletePostBtn(post.id)}
                                >
                                    Delete
                                </button>
                            )}
                        </>)}
                </div>
            </div>

            <p className='m-0'>{post.body}</p>
            {post.img && <img src={`${postImagesEndpoint}${post.img}`} className="postImg rounded-1 mt-1" alt="Post image" />}
            {type !== "detail"
                ? <div className="text-end">
                    <Link className="commentsLink" to={commentsURL}>
                        {post.comments === 0 ? `Comment`
                            : post.comments === 1 ? `1 comment`
                                : `${post.comments} comments`}
                    </Link>
                </div>
                : ""
            }

        </div>
    )

}
export default Post