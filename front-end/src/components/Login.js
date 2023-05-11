import { useContext, useEffect, useState } from "react"
import Input from "./form/Input";
import { useNavigate, useOutletContext } from "react-router-dom";

const Login = () => {

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const {setLoggedIn} = useOutletContext();
    const {setAlertMessage} = useOutletContext();
    const {setAlertClassName} = useOutletContext();
    const {setUserId} = useOutletContext();

    const navigate = useNavigate();

    const handleSubmit = (event) => {
        event.preventDefault();
        
        let payload = {
            email: email,
            password: password,
        }

        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(payload)
        }

        fetch('http://localhost:8080/login', requestOptions)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    console.log("invalid login")
                    setAlertMessage(data.message);
                    setAlertClassName("alert alert-danger");
                } else {
                    console.log("user logged in")
                    setLoggedIn(true);
                    setUserId(data.id);
                    setAlertMessage("");
                    setAlertClassName("d-none");
                    navigate('/');
                }
            })
            .catch(error => {
                setAlertClassName("alert alert-danger");
                setAlertMessage(error);
            })
    }

    return (
        <div className="col-md-6 offset-md-3">
            <h1>Login</h1>
            <hr />

            <form onSubmit={handleSubmit}>
                <Input 
                    title="Email"
                    type="email"
                    className="form-control"
                    name="email"
                    autoComplete="email-new"
                    onChange={(event) => setEmail(event.target.value)}
                />
                 <Input 
                    title="password"
                    type="password"
                    className="form-control"
                    name="password"
                    autoComplete="password-new"
                    onChange={(event) => setPassword(event.target.value)}
                />
                <hr />
                <input 
                    type="submit"
                    className="btn btn-primary"
                    value="Login"
                />
            </form>
        </div>
    )
}

export default Login