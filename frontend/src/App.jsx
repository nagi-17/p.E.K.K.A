import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import LoginPage from './pages/LoginPage';

function App() {
  const isAuthDone = useAuthStore(function(state) {
      return state.isAuthDone;
  });

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={isAuthDone ? <h1 style={styles.h1_text}>Welcome to the Village!</h1> : <Navigate to="/login" />}/> 
        <Route path="/login" element={!isAuthDone ? <LoginPage /> : <Navigate to="/" />}/>
      </Routes>
    </BrowserRouter>
  );
}

const styles = {
    h1_text: {fontFamily: '"Luckiest Guy", cursive', fontSize: '1rem', color: 'white'}
}

export default App;