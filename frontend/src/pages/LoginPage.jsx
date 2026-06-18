import LoginForm from '../components/auth/LoginForm'
import { useState } from 'react';
import RegisterForm from '../components/auth/RegisterForm';

function LoginPage() {
    const[isLogin, setIsLogin]=useState(true);
    function tglReg() {
        setIsLogin(!isLogin);
    }
    return (
        <div style={styles.body}>
            <div style={styles.contain}>
                {isLogin?<LoginForm />:<RegisterForm />}
                <div onClick={tglReg} style={styles.tgl}>
                    {isLogin?<p>Register here</p>:<p>Login here</p>}
                </div>
            </div>
        </div>
    )
}

const styles={
    body: {
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100dvh', 
        margin: 0,
        fontFamily: '"Luckiest Guy", cursive',
        fontSize: '1.5rem',
        backgroundImage: 'url(login.jpg)',
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'center',
        backgroundSize: 'cover',
        backgroundAttachment: 'fixed'
    },
    contain: {
        display: 'flex', flexDirection: 'column', alignItems: 'center'
    },
    tgl: {
        textAlign: 'center', marginTop: '0.75rem', fontSize: '2rem', cursor: 'pointer', color: '#000000', textShadow: '0.125rem 0.125rem 0.25rem rgba(0, 0, 0, 0.8)'
    }
};

export default LoginPage;