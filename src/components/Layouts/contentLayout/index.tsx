import { Content } from "antd/es/layout/layout";
import { Route, Routes } from "react-router-dom";
import {
  ROUTES_PATH,
  ROUTES_PATH_NOPERMISSION,
  ROUTES_PATH_WORKER,
} from "../../../resources/routes";
import { useSelector } from "react-redux";

const ContentLayout = ({ children }: any) => {
  const user = useSelector((state: any) => state.authen);

  const { userRoleID } = user.users || 2;
  const renderRoute = () => {
    const renderRoutes = (routes: any) => {
      return Object.values(routes).map((item: any) => (
        <Route path={item.PATH} key={item.KEY} Component={item.COMPONENT} />
      ));
    };

    // if (userRoleID !== undefined) {
      return <Routes>{renderRoutes(ROUTES_PATH)}</Routes>;
    // }
  };

  return (
    <Content
      className="site-layout-background"
      style={{
        // margin: "24px 16px",
        padding: 24,
        minHeight: 280,
      }}
    >
      {/* {renderRoute()} */}
      <Routes>
        {Object.values(ROUTES_PATH).map((item: any) => (
          <Route path={item.PATH} key={item.KEY} Component={item.COMPONENT} />
        ))}
      </Routes>
      {/* {children} */}
    </Content>
  );
};

export default ContentLayout;
