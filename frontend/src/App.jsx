import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
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
      <Toaster position="top-center" reverseOrder={false} />
      <Routes>
        <Route path="/login" element={!isAuthDone ? <LoginPage /> : <Navigate to="/" />}/>
        <Route path="/" element={isAuthDone ? <HomePage />: <Navigate to="/login" />}/> 
        <Route path="/battle" element={isAuthDone ? <BattlePage /> : <Navigate to="/login" />}/> 
      </Routes>
    </BrowserRouter>
  );
}