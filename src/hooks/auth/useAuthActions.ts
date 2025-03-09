import { useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { login, logout, checkAuth, loginLark } from '../../redux/auth/action';
import { LoginPayload } from '../../types/auth.types';
import { logger } from '../../utils/logger';

export const useAuthActions = () => {
  const dispatch = useDispatch();

  const handleLogin = useCallback(
    (credentials: LoginPayload) => {
      logger.perf.start('Auth: Login Process');
      logger.log('info', 'Login Initiated', {
        username: credentials.username,
        timestamp: new Date().toISOString(),
      });

      dispatch(login(credentials));

      logger.perf.end('Auth: Login Process');
    },
    [dispatch]
  );

  const handleLogout = useCallback(() => {
    logger.perf.start('Auth: Logout Process');
    logger.log('info', 'Logout Initiated', {
      timestamp: new Date().toISOString(),
    });

    dispatch(logout());

    logger.perf.end('Auth: Logout Process');
  }, [dispatch]);

  const handleCheckAuth = useCallback(() => {
    logger.log('info', 'Auth Check Initiated', {
      timestamp: new Date().toISOString(),
    });

    dispatch(checkAuth());
  }, [dispatch]);

  const handleLarkLogin = useCallback(
    (payload: any) => {
      logger.log('info', 'Lark Login Initiated', {
        timestamp: new Date().toISOString(),
      });

      dispatch(loginLark(payload));
    },
    [dispatch]
  );

  return {
    login: handleLogin,
    logout: handleLogout,
    checkAuth: handleCheckAuth,
    loginLark: handleLarkLogin,
  };
};