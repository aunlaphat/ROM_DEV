import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../resources/routes";
import { login, logout } from "../redux/auth/action";
import { getCookies } from "../store/useCookies";
import { logger } from '../utils/logger';

export const useAuthLogin = () => {
  const dispatch = useDispatch();

  function toLogin() {
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  }

  const onLogin = async (values: any) => {
    try {
      console.log("Dispatching login with values:", values);
      // Just dispatch and let the saga handle the rest
      dispatch(login(values));
    } catch (error) {
      console.error("Login Error:", error);
      throw error;
    }
  };

  const onLogout = async () => {
    try {
      logger.auth('info', 'Processing logout request');
      dispatch(logout());
      localStorage.removeItem("access_token");  // Clear local storage
      document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";  // Clear cookie
    } catch (error) {
      logger.auth('error', 'Logout failed', { error });
    }
  };

  function redirectToMain() {
    // if (getCookies('accessToken')) { // open this and comment below
    if (!getCookies("jwt")) {
      NavigateTo({ pathname: ROUTES_PATH.ROUTE_MAIN.PATH });
    }
  }

  return {
    toLogin,
    redirectToMain,
    onLogin,
    onLogout,
    //onLarkLogin,
  };
};
