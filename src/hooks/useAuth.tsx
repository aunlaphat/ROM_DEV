import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../resources/routes-name";
import { checkAuthen, login, logout } from "../redux/authen/action";
import { getCookies } from "../store/useCookies";
import { openAlert } from "../components/alert/useAlert";

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

  function checkLoginToken() {
    // if (!getCookies('accessToken')) { // open this and comment below
    // if (!getCookies("jwt") && window.location.pathname !== "/") {
    //   openAlert({
    //     type: "error",
    //     message: "Token is unauthorized!",
    //     title: "ERROR!",
    //   });
    //   dispatch(logout());
    // } else {
    //   dispatch(checkAuthen());
    // }
  }

  function redirectToMain() {
    // if (getCookies('accessToken')) { // open this and comment below
    if (!getCookies("jwt")) {
      NavigateTo({ pathname: ROUTES_PATH.ROUTE_MAIN.PATH });
    }
  }

  return {
    toLogin,
    checkLoginToken,
    redirectToMain,
    onLogin,
    onLogout,
  };
};
