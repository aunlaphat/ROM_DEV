import { useSelector } from 'react-redux';
import { RootState } from '../../redux/store';
import { AuthState, User } from '../../types/auth.types';

export const useAuthState = () => {
  const auth = useSelector<RootState, AuthState>((state) => state.auth);
  
  return {
    user: auth.user,
    isAuthenticated: auth.isAuthenticated,
    loading: auth.loading,
    error: auth.error,
    token: auth.token,
  };
};