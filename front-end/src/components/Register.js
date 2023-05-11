import { useState } from "react"
import Input from "./form/Input";
import Select from "./form/Select";
import { useNavigate, useOutletContext } from "react-router-dom";

const Register = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [passwordConfirm, setPasswordConfirm] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [dateOfBirth, setDateOfBirth] = useState("");
    const [image, setImage] = useState(null)
    const [nickname, setNickname] = useState("");
    const [about, setAbout] = useState("");
    const [profileP, setProfileP] = useState(parseInt("1"))

    const mapProfilePOptions = [
        {id: "Public", value: 1},
        {id: "Private", value: 0}
    ]

    const {setLoggedIn} = useOutletContext();
    const {setAlertMessage} = useOutletContext();
    const {setAlertClassName} = useOutletContext();

    const navigate = useNavigate();

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
              setAlertClassName("alert alert-danger");
              setAlertMessage(data.message);
            } else {
              setLoggedIn(true);
              setAlertClassName("alert-success");
                setAlertMessage("Successfully registered!");
              navigate("/");
            }
          })
          .catch(error => {
            setAlertClassName("alert alert-danger");
            setAlertMessage(error);
          });
      };

    const handleSubmit = (event) => {
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
        if (email === "" || firstName === "" || lastName === "" || password === "" || passwordConfirm === "" || dateOfBirth === ""){
            setAlertClassName("alert alert-danger");
            setAlertMessage("Empty fields");
        }else if (!email.match(/^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$/g)){
            setAlertClassName("alert alert-danger");
            setAlertMessage("Not valid email");
        }else if (password.length < 5){
            setAlertClassName("alert alert-danger");
            setAlertMessage("Passwords too short");
        } else if (password !== passwordConfirm) {
            setAlertClassName("alert alert-danger");
            setAlertMessage("Passwords dont match");
        } else if (image && !image.name.match(/(gif|jpg|jpeg|png)$/gi) ){
            setAlertClassName("alert alert-danger");
            setAlertMessage("file type invalid");
        }else if (image && image.size > 733333){
            setAlertClassName("alert alert-danger");
            setAlertMessage("File too large");
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
        <div className="col-md-6 offset-md-3">
            <h1>Register</h1>
            <hr />
            <form onSubmit={handleSubmit}>
                <Input 
                    title="Email (required)"
                    type="email"
                    className="form-control"
                    name="email"
                    autoComplete="email-new"
                    onChange={(event) => setEmail(event.target.value)}
                />
                <Input 
                    title="First name (required)"
                    type="text"
                    className="form-control"
                    name="firstName"
                    autoComplete="firstName-new"
                    onChange={(event) => setFirstName(event.target.value)}
                />
                <Input 
                    title="Last name (required)"
                    type="text"
                    className="form-control"
                    name="lastName"
                    autoComplete="lastName-new"
                    onChange={(event) => setLastName(event.target.value)}
                />
                <Input 
                    title="Date of birth (required)"
                    type="date"
                    className="form-control"
                    name="dateOfBirth"
                    autoComplete="dateOfBirth-new"
                    onChange={(event) => setDateOfBirth(new Date(event.target.value))}
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
                <Select 
                    title="Profile status"
                    type="select"
                    className="form-control"
                    name="profileP"
                    autoComplete="profilePublic-new"
                    options={mapProfilePOptions}
                    onChange={(event) => setProfileP(parseInt(event.target.value))}
                />
                <Input 
                    title="Profile picture"
                    type="file"
                    className="form-control"
                    name="image"
                    onChange={(event) => setImage(event.target.files[0])}
                />
                <Input 
                    title="Password (required)"
                    type="password"
                    className="form-control"
                    name="password"
                    autoComplete="password-new"
                    onChange={(event) => setPassword(event.target.value)}
                />
                <Input 
                    title="Confirm password (required)"
                    type="password"
                    className="form-control"
                    name="passwordConfirm"
                    autoComplete="passwordConfirm-new"
                    onChange={(event) => setPasswordConfirm(event.target.value)}
                />
                <hr />

                <input 
                    type="submit"
                    className="btn btn-primary"
                    value="Register"
                />

            </form>
        </div>
    )
}

export default Register