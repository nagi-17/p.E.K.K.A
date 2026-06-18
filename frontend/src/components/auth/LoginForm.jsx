import { useState } from "react";
import { loginUser } from "../../api/auth";
import { useAuthStore } from "../../store/authStore";
import {useNavigate} from 'react-router-dom';

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
        errMsg=<p style={styles.error}>{error}</p>;
    }

    return (
        <form onSubmit={handleSubmit} style={styles.form}>
            <h2>Login to p.E.K.K.A</h2>
            {errMsg}
            <input type="text" placeholder="Username" value={username} onChange={handleUsernameChange} style={styles.input} required></input>
            <input type="password" placeholder="Password" value={pass} onChange={handlePassChnage} style={styles.input} required></input>
            <button type="submit" style={styles.button}>Enter Village</button>
        </form>
    );
};

const styles = {
    form: { display:'flex', flexDirection:'column', gap:'1rem', width:'18.75rem', margin:'0 auto', padding:'1.25rem', backgroundColor:'transparent', borderRadius:'0.5rem', color:'black', alignItems: 'center'},
    input: { width: '100%', padding:'0.875rem' ,borderRadius:'0.5rem', border:'none'},
    button: { fontFamily: '"Luckiest Guy"', width:'100%', padding:'0.75rem', backgroundColor:'#00aaff', color:'white', border:'none', borderRadius:'0.25rem', cursor:'pointer', fontWeight:'light', fontSize: '1.5rem', alignItems: 'center'},
    error: { color: '#e74c3c', fontSize: '0.875rem', margin: 0}
};