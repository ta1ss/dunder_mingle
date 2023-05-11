import { Link, useOutletContext, useParams } from "react-router-dom";
import React, { useEffect, useState } from 'react'
import Post from './Post';
import CommentForm from "./form/CommentForm";

const PostDetail = () => {
    let { postId } = useParams();
    let { groupId } = useParams();

    const [post, setPost] = useState(null)
    const [comments, setComments] = useState([])
    const [confirmDelete, setConfirmDelete] = useState({});
    const { userId } = useOutletContext();

    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"
    const commentImagesEndpoint = "http://localhost:8080/media/comment_images/"

    const { setAlertMessage } = useOutletContext();
    const { setAlertClassName } = useOutletContext();

    const newComment = (newComment) => {
        setComments([...comments, newComment])
    };

    const handleDeleteToggle = (postId) => {
        setConfirmDelete((prevState) => ({
            ...prevState,
            [postId]: !prevState[postId],
        }));
    };

    const handlePostDelete = () => {
        setPost(null)
        setComments([])
    }

    useEffect(() => {
        const options = {
            method: 'GET',
            credentials: 'include',
        }

        let endpoint = "http://localhost:8080/"

        if (groupId) {
            endpoint += `group/post/${postId}`
        } else {
            endpoint += `user/post/${postId}`
        }

        fetch(endpoint, options)
            .then(response => response.json())
            .then(data => {
                if (data.post) {
                    setPost(data.post)
                }
                if (data.comments) {
                    setComments(data.comments)
                }
            })
            .catch(error => console.log(error))
    }, [])

    const handleCommentDelete = (commentId) => {
        const options = {
            method: 'DELETE',
            credentials: 'include',
            body: JSON.stringify({ id: commentId, userId: userId })
        }

        let endpoint = "http://localhost:8080/comments/"

        if (groupId) {
            endpoint += "group"
        } else {
            endpoint += "user"
        }

        fetch(endpoint, options)
            .then(response => response.json())
            .then(data => {
                const newData = comments.filter(comment => comment.id !== commentId)
                setComments(newData)
                setAlertClassName("alert-dark");
                setAlertMessage("Comment successfully deleted!");
                setConfirmDelete((prevState) => ({
                    ...prevState,
                    [commentId]: false,
                }));
            })
            .catch(error => console.log(error))
    }

    return (
        <>
            {post
                ?
                <Post
                    post={post}
                    type="detail"
                    groupId={groupId}
                    onDelete={() => handlePostDelete()}
                />
                : ""
            }
            {comments
                ?
                comments.map((comment, index) => (
                    <div key={index} className='ms-4 comment border rounded-3'>
                        <div className="d-flex">
                            <div className='col-md-3 d-flex'>
                                <Link to={`/profile/${comment.userId}`}>
                                    <img src={`${profileImagesEndpoint}${comment.userImg}`} className="commentUserImg" />
                                </Link>
                                <div>
                                    <input id='commentId' type='hidden' value={comment.id} />
                                    <p className='postInfo m-0'>
                                        <Link to={`/profile/${comment.userId}`}>{comment.createdBy}</Link>
                                    </p>
                                    <p className='commentInfo mb-0'>{new Date(comment.createdAt).toLocaleString('en-GB', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })}</p>
                                </div>

                            </div>
                            <div className='col-md-6'>
                                <p className='mb-0'>{comment.body}</p>
                                {comment.img && <img src={`${commentImagesEndpoint}${comment.img}`} className="commentImg rounded-1 mt-1" alt="Comment image" />}
                            </div>
                            <div className="col-md-3 text-end">
                                {userId && userId === comment.userId &&
                                    (<>
                                        {confirmDelete[comment.id] ? (
                                            <div className="d-flex justify-content-end">
                                                <button
                                                    type='button'
                                                    className='btn btn-sm btn-outline-secondary'
                                                    onClick={() => handleDeleteToggle(comment.id)}
                                                >
                                                    Cancel
                                                </button>
                                                <button
                                                    type='button'
                                                    className='btn btn-sm btn-outline-danger ms-2'
                                                    onClick={() => handleCommentDelete(comment.id)}
                                                >
                                                    Confirm
                                                </button>
                                            </div>
                                        ) : (
                                            <button
                                                type='button'
                                                className='btn btn-sm btn-outline-danger'
                                                onClick={() => handleDeleteToggle(comment.id)}
                                            >
                                                Delete
                                            </button>
                                        )}
                                    </>)}
                            </div>
                        </div>
                    </div>
                ))
                : ""
            }
            {post ? <CommentForm postId={post.id} groupId={groupId} newComment={newComment} /> : ""}
        </>
    )

}
export default PostDetail