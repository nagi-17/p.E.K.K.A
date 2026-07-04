import { useState } from "react";
import { regUser } from "../../api/auth";
import { useAuthStore } from "../../store/authStore";
import { useNavigate } from 'react-router-dom';
import styles from './AuthForm.module.css';

export default function RegisterForm() {
    const[username, setUsername]=useState('');
    const[email, setEmail]=useState('');
    const[pass, setPass]=useState('');
    const[error, setError]=useState(null);

    const navigate=useNavigate();
    const login=useAuthStore(function(state) {
        return state.login;
    });

    async function handleSubmit(event) {
        event.preventDefault();
        setError(null);

        try {
            const data=await regUser(username, email, pass);
            if (data.token!==undefined) {
                login(data.token, data.player_id);
                navigate('/');
            }
            else {
                alert("Acc. has been created, please login")
            }
        }
        catch(err) {
            setError(err.message);
        }
    };

    function handleUsernameChange(event) {setUsername(event.target.value);}
    function handlePassChange(event) {setPass(event.target.value);}
    function handleEmailChange(event) {setEmail(event.target.value);}

    let errMsg=null;
    if (error!==null) {
        errMsg=<p className={styles.error}>{error}</p>;
    }

    return (
        <form onSubmit={handleSubmit} className={styles.form}>
            <h2>Register</h2>
            {errMsg}
            <input type="text" placeholder="Username" value={username} onChange={handleUsernameChange} className={styles.input} required></input>
            <input type="email" placeholder="email" value={email} onChange={handleEmailChange} className={styles.input} required></input>
            <input type="password" placeholder="Password" value={pass} onChange={handlePassChange} className={styles.input} required></input>
            <button type="submit" className={styles.button}>Create Account</button>
        </form>
    );
};
