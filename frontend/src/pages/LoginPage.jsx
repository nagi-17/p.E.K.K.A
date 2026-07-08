import LoginForm from '../components/auth/LoginForm'
import { useState } from 'react';
import RegisterForm from '../components/auth/RegisterForm';
import styles from './LoginPage.module.css';

function LoginPage() {
    const[isLogin, setIsLogin]=useState(true);
    function tglReg() {
        setIsLogin(!isLogin);
    }
    return (
        <div className={styles.body}>
            <div className={styles.contain}>
                {isLogin?<LoginForm />:<RegisterForm onSuccess={() => setIsLogin(true)} />}
                <div onClick={tglReg} className={styles.tgl}>
                    {isLogin?<p>Register here</p>:<p>Login here</p>}
                </div>
            </div>
        </div>
    )
}



export default LoginPage;