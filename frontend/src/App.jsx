import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import LoginPage from './pages/LoginPage';
import HomePage from './pages/HomePage';
import BattlePage from './pages/BattlePage';

export default function App() {
  const isAuthDone = useAuthStore(function(state) {
      return state.isAuthDone;
  });

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={!isAuthDone ? <LoginPage /> : <Navigate to="/" />}/>
        <Route path="/" element={isAuthDone ? <HomePage />: <Navigate to="/login" />}/> 
        <Route path="/battle" element={isAuthDone ? <BattlePage /> : <Navigate to="/login" />}/> 
      </Routes>
    </BrowserRouter>
  );
}

const styles = {
    h1_text: {fontFamily: '"Luckiest Guy", cursive', fontSize: '1rem', color: 'white'}
}