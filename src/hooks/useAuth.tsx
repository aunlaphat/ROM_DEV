import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../resources/routes-name";
import { login, logout, login_lark } from "../redux/authen/action";
import { getCookies } from "../store/useCookies";

export const useAuthLogin = () => {
  const dispatch = useDispatch();

  function toLogin() {
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  }

  const onLogin = async (values: any) => {
    dispatch(login(values));
  };

  const onLogout = async (values?: any) => {
    dispatch(logout());
  };

  const onLarkLogin = async (values: any) => {
    dispatch(login_lark(values));
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
    onLarkLogin,
  };
};
