import axios from "axios";
import "./register.css";
import { useRef } from "react";
import { useHistory, Link } from "react-router-dom";
import { makeRequest } from "../../axios";
import { registerCall } from "../../apiCalls";

export default function Register() {
    const username = useRef();
    const email = useRef();
    const password = useRef();
    const paswordAgain = useRef();
    const history = useHistory();

    const handleClick = async (e) => {
        e.preventDefault();
        if (paswordAgain.current.value !== password.current.value) {
            paswordAgain.current.setCustomValidity("Passwords do not matched");
        } else {
            const user = {
                user_name: username.current.value,
                email: email.current.value,
                password: password.current.value,
            };
            registerCall(user);
            history.push("/login");
        }
    };

    return (
        <div className="register">
            <div className="registerWrapper">
                <div className="registerLeft">
                    <h3 className="registerLogo">Lamasocial</h3>
                    <span className="registerDesc">
                        Connect with friends and the world around you on
                        Lamasocial.
                    </span>
                </div>
                <div className="registerRight">
                    <form className="registerBox" onSubmit={handleClick}>
                        <input
                            placeholder="Username"
                            required
                            ref={username}
                            className="registerInput"
                        />
                        <input
                            placeholder="Email"
                            required
                            type="email"
                            ref={email}
                            className="registerInput"
                        />
                        <input
                            placeholder="Password"
                            required
                            type="password"
                            ref={password}
                            className="registerInput"
                        />
                        <input
                            placeholder="Password Again"
                            required
                            type="password"
                            ref={paswordAgain}
                            className="registerInput"
                        />
                        <button className="registerButton" type="submit">
                            Sign Up
                        </button>
                    </form>
                    <Link to="/login">
                        <button className="registerLoginButton">
                            Log into Account
                        </button>
                    </Link>
                </div>
            </div>
        </div>
    );
}
