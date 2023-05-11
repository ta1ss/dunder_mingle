import { Link, useOutletContext, useParams } from "react-router-dom";
import React, { useEffect, useState } from 'react'

const Comments = () => {
    let { postId } = useParams();
    let { groupId } = useParams();

    console.log({postId})
    console.log({groupId})

    const [post, setPost] = useState({})
    const [comments, setComments] = useState([])
    const [confirmDelete, setConfirmDelete] = useState({});
    const { userId } = useOutletContext();

    const postAndCommentsEndpoint = "http://localhost:8080/user/post/" + postId
    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"
    const postImagesEndpoint = "http://localhost:8080/media/post_images/"
    const commentImagesEndpoint = "http://localhost:8080/media/comment_images/"

    const { setAlertMessage } = useOutletContext();
    const { setAlertClassName } = useOutletContext();

    const handleDeleteToggle = (postId) => {
        setConfirmDelete((prevState) => ({
            ...prevState,
            [postId]: !prevState[postId],
        }));
    };

    useEffect(() => {
        const options = {
            method: 'GET',
            credentials: 'include',
        }
        console.log(postAndCommentsEndpoint)
        fetch(postAndCommentsEndpoint, options)
            .then(response => response.json())
            .then(data => {
                console.log(data)
            })
            .catch(error => console.log(error))
    }, [])

    const handleConfirmDelete = (postId) => {
        const options = {
            method: 'DELETE',
            credentials: 'include',
            body: JSON.stringify({ id: postId, userId: userId })
        }

        fetch('http://localhost:8080/post/user', options)
            .then(response => response.json())
            .then(data => {
                const newData = post.filter(post => post.id !== postId)
                setPost(newData)
                setAlertClassName("alert-dark");
                setAlertMessage("Post successfully deleted!");
                setConfirmDelete((prevState) => ({
                    ...prevState,
                    [postId]: false,
                }));
            })
            .catch(error => console.log(error))
    }

    if (post && post.error) {
        return (
            <div>
                <p>{post.message}</p>
            </div>
        )
    } else if (post) {
        return (
            <div>
                {comments.map((comment, index) => (
                    <div key={index} className='comment border rounded-3'>
                        <div className="d-flex">
                            <div className='col-md-9 d-flex'>
                                <img src={`${profileImagesEndpoint}${comment.userImg}`} className="postUserImg" />
                                <div>
                                    <h5 className='m-0'>{comment.title}</h5>
                                    <input id='postId' type='hidden' value={comment.id} />
                                    <input id='userId' type='hidden' value={comment.userId} />
                                    <p className='postInfo mb-1'>
                                        <Link to={`/profile/${comment.userId}`}>{comment.createdBy}</Link> | {new Date(comment.createdAt).toLocaleString('en-GB', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })} | Privacy: {comment.privacy}
                                    </p>
                                </div>
                            </div>
                            <div className='col-md-3 text-end'>
                                {userId && userId === comment.userId &&
                                    (<>
                                        {confirmDelete[comment.id] ? (
                                            <div className="d-flex justify-content-end">
                                                <button
                                                    type='button'
                                                    className='btn btn-sm btn-secondary'
                                                    onClick={() => handleDeleteToggle(comment.id)}
                                                >
                                                    Cancel
                                                </button>
                                                <button
                                                    type='button'
                                                    className='btn btn-sm btn-danger ms-2'
                                                    onClick={() => handleConfirmDelete(comment.id)}
                                                >
                                                    Confirm
                                                </button>
                                            </div>
                                        ) : (
                                            <button
                                                type='button'
                                                className='btn btn-sm btn-danger'
                                                onClick={() => handleDeleteToggle(comment.id)}
                                            >
                                                Delete
                                            </button>
                                        )}
                                    </>)}
                            </div>
                        </div>

                        <p className='m-0'>{comment.body}</p>
                        {comment.img && <img src={`${postImagesEndpoint}${comment.img}`} className="postImg rounded-1" alt="Post image" />}
                    </div>
                ))}
            </div>
        )
    } else {
        return (
            <div>
                <p>No post yet</p>
            </div>
        )
    }
}
export default Comments