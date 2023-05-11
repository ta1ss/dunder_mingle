import { useNavigate } from "react-router-dom"
import { useState } from "react";
import Input from "./form/Input";
import Select from "./form/Select";
import Alert from "./form/Alert";

const Welcome = (props) => {
    const [showLoginForm, setShowLoginForm] = useState(true);
    const navigate = useNavigate();

    // Login
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    // Register
    const [passwordConfirm, setPasswordConfirm] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [dateOfBirth, setDateOfBirth] = useState("");
    const [image, setImage] = useState(null)
    const [nickname, setNickname] = useState("");
    const [about, setAbout] = useState("");
    const [profileP, setProfileP] = useState(parseInt("1"))

    const mapProfilePOptions = [
        { id: "Public", value: 1 },
        { id: "Private", value: 0 }
    ]


    const handleLoginButtonClick = () => {
        setShowLoginForm(true);
    };

    const handleRegisterButtonClick = () => {
        setShowLoginForm(false);
    };




    // Login
    const handleLoginSubmit = (event) => {
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
                    props.setAlertMessage(data.message);
                    props.setAlertClassName("alert alert-danger");
                } else {
                    console.log("user logged in")
                    props.setLoggedIn(true);
                    props.setUserId(data.id);
                    props.setAlertMessage("");
                    props.setAlertClassName("d-none");
                    navigate('/');
                }
            })
            .catch(error => {
                props.setAlertClassName("alert alert-danger");
                props.setAlertMessage(error);
            })
    }


    // Register
    const register = (payload) => {
        const requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(payload)
        };
        fetch("http://localhost:8080/register", requestOptions)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    props.setAlertClassName("alert alert-danger");
                    props.setAlertMessage(data.message);
                } else {
                    props.setLoggedIn(true);
                    props.setAlertClassName("alert-success");
                    props.setAlertMessage("Successfully registered!");
                    navigate("/");
                }
            })
            .catch(error => {
                props.setAlertClassName("alert alert-danger");
                props.setAlertMessage(error);
            });
    };

    const handleRegisterSubmit = (event) => {
        let payload = {
            email,
            password,
            passwordConfirm,
            firstName,
            lastName,
            dateOfBirth,
            nickname,
            about,
            profileP
        }
        event.preventDefault();
        if (email === "" || firstName === "" || lastName === "" || password === "" || passwordConfirm === "" || dateOfBirth === "") {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("Empty fields");
        } else if (!email.match(/^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$/g)) {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("Not valid email");
        } else if (password.length < 5) {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("Passwords too short");
        } else if (password !== passwordConfirm) {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("Passwords dont match");
        } else if (image && !image.name.match(/(gif|jpg|jpeg|png)$/gi)) {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("file type invalid");
        } else if (image && image.size > 733333) {
            props.setAlertClassName("alert alert-danger");
            props.setAlertMessage("File too large");
        } else if (image) {
            const reader = new FileReader();
            reader.onload = function () {
                payload.image = reader.result;
                register(payload);
            };
            reader.readAsDataURL(image);
        } else {
            register(payload);
        }
    }


    return (
        <div className="row text-center welcome-container">
            <h2 className="welcome-title">Welcome to the dunder mingle</h2>
            <p className="welcome-subtitle">Please login or register to start mingling..</p>
            <div className="col-md-4" style={{ borderRight: "1px solid black" }}>
                <div className=" welcome-logo-div" >
                    <img src={"http://localhost:8080/media/various/logo.png"} alt="Social Network Logo" className="welcome-logo" />
                </div>
            </div>

            <div className="col-md-8">
            <Alert
                message={props.alertMessage}
                className={props.alertClassName}
              />
                <div className=" welcome-forms " >
                    <div className="welcome-option-buttons text-center">
                        <button
                            className={`welcome-login-option-button btn ${showLoginForm ? 'btn-dark btn-border' : 'btn-light btn-border'}`}
                            onClick={handleLoginButtonClick}
                        >
                            Login
                        </button>
                        <button
                            className={`welcome-register-option-button btn ${!showLoginForm ? 'btn-dark btn-border' : 'btn-light btn-border'}`}
                            onClick={handleRegisterButtonClick}
                        >
                            Register
                        </button>
                    </div>

                    {showLoginForm ? (
                        <div className="welcome-login-form">
                            <form onSubmit={handleLoginSubmit}>
                                <Input
                                    title="Email"
                                    type="email"
                                    className="form-control"
                                    name="email"
                                    autoComplete="email-new"
                                    onChange={(event) => setEmail(event.target.value)}
                                />
                                <Input
                                    title="Password"
                                    type="password"
                                    className="form-control"
                                    name="password"
                                    autoComplete="password-new"
                                    onChange={(event) => setPassword(event.target.value)}
                                />

                                <input
                                    type="submit"
                                    className="btn btn-dark welcome-submit-button"
                                    value="Login"
                                />
                            </form>
                        </div>
                    ) : (
                        <div className="welcome-register-form">
                            <form onSubmit={handleRegisterSubmit}>
                                <div className="row">
                                    <div className="col-md-6">
                                        <Input
                                            title="First name *"
                                            type="text"
                                            className="form-control"
                                            name="firstName"
                                            autoComplete="firstName-new"
                                            onChange={(event) => setFirstName(event.target.value)}
                                        />
                                        <Input
                                            title="Last name *"
                                            type="text"
                                            className="form-control"
                                            name="lastName"
                                            autoComplete="lastName-new"
                                            onChange={(event) => setLastName(event.target.value)}
                                        />
                                        <Input
                                            title="Email *"
                                            type="email"
                                            className="form-control"
                                            name="email"
                                            autoComplete="email-new"
                                            onChange={(event) => setEmail(event.target.value)}
                                        />

                                        <Input
                                            title="Date of birth *"
                                            type="date"
                                            className="form-control"
                                            name="dateOfBirth"
                                            autoComplete="dateOfBirth-new"
                                            onChange={(event) => setDateOfBirth(new Date(event.target.value))}
                                        />

                                        <Select
                                            title="Profile status"
                                            type="select"
                                            className="form-control"
                                            name="profileP"
                                            autoComplete="profilePublic-new"
                                            options={mapProfilePOptions}
                                            onChange={(event) => setProfileP(parseInt(event.target.value))}
                                        />

                                    </div>
                                    <div className="col-md-6">
                                        <Input
                                            title="Profile picture"
                                            type="file"
                                            className="form-control"
                                            name="image"
                                            onChange={(event) => setImage(event.target.files[0])}
                                        />
                                        <Input
                                            title="Nickname"
                                            type="text"
                                            className="form-control"
                                            name="nickname"
                                            autoComplete="nickname-new"
                                            onChange={(event) => setNickname(event.target.value)}
                                        />
                                        <Input
                                            title="About"
                                            type="text"
                                            className="form-control"
                                            name="about"
                                            autoComplete="about-new"
                                            onChange={(event) => setAbout(event.target.value)}
                                        />


                                        <Input
                                            title="Password *"
                                            type="password"
                                            className="form-control"
                                            name="password"
                                            autoComplete="password-new"
                                            onChange={(event) => setPassword(event.target.value)}
                                        />
                                        <Input
                                            title="Confirm password *"
                                            type="password"
                                            className="form-control"
                                            name="passwordConfirm"
                                            autoComplete="passwordConfirm-new"
                                            onChange={(event) => setPasswordConfirm(event.target.value)}
                                        />
                                    </div>
                                </div>

                                <input
                                    type="submit"
                                    className="btn btn-dark welcome-submit-button"
                                    value="Register"
                                />

                            </form>
                        </div>

                    )}
                </div>
            </div>
        </div>

    );
}
export default Welcome