import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES, ROUTE_LOGIN } from "../resources/routes";
import { login, logout } from "../redux/auth/action";
import { getCookies } from "../store/useCookies";
import { logger } from "../utils/logger";

export const useAuthLogin = () => {
  const dispatch = useDispatch();

  /**
   * ไปยังหน้า Login
   */
  const toLogin = () => {
    logger.auth("info", "Navigating to Login page");
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  };

  /**
   * ดำเนินการ Login
   */
  const onLogin = async (values: any) => {
    try {
      logger.auth("info", "Dispatching login request", { username: values.username });
      dispatch(login(values));
    } catch (error) {
      logger.auth("error", "Login failed", { error });
      throw error;
    }
  };

  /**
   * ดำเนินการ Logout
   */
  const onLogout = async () => {
    try {
      logger.auth("info", "Processing logout request");
      dispatch(logout());

      // ล้างข้อมูลการ Authentication
      localStorage.removeItem("access_token"); // Clear Local Storage
      document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"; // Clear Cookie

      logger.auth("info", "User logged out successfully");
    } catch (error) {
      logger.auth("error", "Logout failed", { error });
    }
  };

  /**
   * Redirect ไปยังหน้า Home ถ้าผู้ใช้มีสิทธิ์เข้าใช้งาน
   */
  const redirectToMain = () => {
    if (!getCookies("jwt")) {
      logger.auth("info", "Redirecting to Home page");
      NavigateTo({ pathname: ROUTES.ROUTE_MAIN.PATH });
    }
  };

  return {
    toLogin,
    redirectToMain,
    onLogin,
    onLogout,
  };
};
