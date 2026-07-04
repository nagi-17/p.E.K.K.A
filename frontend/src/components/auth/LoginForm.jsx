import { useState } from "react";
import { loginUser } from "../../api/auth";
import { useAuthStore } from "../../store/authStore";
import { useNavigate } from 'react-router-dom';
import styles from './AuthForm.module.css';

export default function LoginForm() {
    const[username, setUsername]=useState('');
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
            const data=await loginUser(username, pass);
            login(data.token, data.player_id);
            navigate('/');
        }
        catch(err) {
            setError(err.message);
        }
    };

    function handleUsernameChange(event) {setUsername(event.target.value);}
    function handlePassChnage(event) {setPass(event.target.value);}

    let errMsg=null;
    if (error!==null) {
        errMsg=<p className={styles.error}>{error}</p>;
    }

    return (
        <form onSubmit={handleSubmit} className={styles.form}>
            <h2>Login</h2>
            {errMsg}
            <input type="text" placeholder="Username" value={username} onChange={handleUsernameChange} className={styles.input} required></input>
            <input type="password" placeholder="Password" value={pass} onChange={handlePassChnage} className={styles.input} required></input>
            <button type="submit" className={styles.button}>Enter Village</button>
        </form>
    );
};
