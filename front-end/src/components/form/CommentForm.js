import Input from "./Input";
import Textarea from "./Textarea";
import React, { useState } from 'react'
import { useOutletContext} from "react-router-dom";

const CommentForm = ({ postId, groupId, newComment }) => {

    const [body, setBody] = useState('');
    const [image, setImage] = useState(null)
    const [imageName, setImageName] = useState("Add Image")

    const { setAlertMessage } = useOutletContext();
    const { setAlertClassName } = useOutletContext();

    const postComment = (event) => {
        event.preventDefault();
        if (body === '') {
            setAlertClassName("alert-danger");
            setAlertMessage("Content can't be empty!");
        } else if (image && !image.name.match(/(gif|jpg|jpeg|png)$/gi)) {
            setAlertClassName("alert-danger");
            setAlertMessage("Only .gif .jpg .jpeg .png allowed");
        } else if (image && image.size > 500000) {
            setAlertClassName("alert-danger");
            setAlertMessage("Maximum image size 500kB");
        } else {
            const comment = { postId: postId, body: body, groupId: parseInt(groupId) }
            if (image) {
                const reader = new FileReader()
                reader.onload = function () {
                    comment.img = reader.result
                    addCommentToDatabase(comment)
                }
                reader.readAsDataURL(image)
            } else {
                addCommentToDatabase(comment)
            }
        }
    }

    const addCommentToDatabase = (comment) => {
        const options = {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(comment)
        }

        let endpoint = "http://localhost:8080/comments/"

        if (!groupId) {
            endpoint += "user"
        } else {
            endpoint += "group"
        }

        fetch(endpoint, options)
            .then(response => response.json())
            .then(data => {
                newComment(data)
                setBody("")
                setImage(null)
                setImageName("Add Image")
                setAlertClassName("alert-success");
                setAlertMessage("Comment successfully created!");
            })
            .catch(error => console.log(error))
    }

    return (
        <form className="commentForm ms-4" onSubmit={postComment}>
            <Textarea
                placeholder="Comment..."
                type="text"
                className="form-control"
                name="body"
                autoComplete="body-new"
                rows={3}
                value={body}
                onChange={(event) => setBody(event.target.value)}
            />
            <div className="d-flex justify-content-end" style={{marginTop: "12px"}}>
                <Input
                    type="file"
                    name="commentImgUpload"
                    className="form-control hidden"
                    onChange={(event) => {
                        setImage(event.target.files[0])
                        setImageName(event.target.files[0].name)
                    }}
                />
                <button className="btn btn-outline-secondary w-25 text-truncate" type='button' onClick={() => document.getElementById('commentImgUpload').click()}>{imageName}</button>
                <button type="submit" className="btn btn-primary w-25 ms-2" onClick={postComment}>Comment</button>
            </div>

        </form>
    )
}

export default CommentForm