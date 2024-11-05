import * as routes from "../../resources/routes-name";
import { useNavigate } from "react-router-dom";

export const NavigateTo = ({ pathname, state, search }: any) => {
  const navigate = useNavigate();

  if (!navigate) {
    return;
  }

  if (!pathname) {
    navigate(routes.ROUTES_PATH.ROUTE_MAIN.PATH);
  }

  const fullPath = search ? `${pathname}?${search}` : pathname;

  navigate(fullPath);
};

export const NavigateReplaceTo = ({ pathname, state = {}, search }: any) => {
  const navigate = useNavigate();
  if (!pathname) {
    navigate(routes.ROUTES_PATH.ROUTE_MAIN.PATH);
  } else {
    const fullPathname = search ? `${pathname}?${search}` : pathname;
    navigate(fullPathname, { state: { ...state } });
  }
};

export const windowNavigateReplaceTo = ({ pathname }: any) => {
  if (!pathname) {
    return;
  }
  window.location.replace(pathname);
};
