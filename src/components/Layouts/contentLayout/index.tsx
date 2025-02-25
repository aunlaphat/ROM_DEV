import { Content } from "antd/es/layout/layout";
import { Route, Routes } from "react-router-dom";
import { ROUTES_PATH, ROUTES_PATH_NOPERMISSION } from "../../../resources/routes";
import { useSelector } from "react-redux";
import { RoleID } from '../../../constants/roles';

const ContentLayout = ({ children }: any) => {
  const auth = useSelector((state: any) => state.auth);
  // เข้าถึง roleID โดยตรง
  const roleID = auth?.user?.roleID;

  const getRoutesByRole = () => {
    console.log('🔍 Auth State:', {
      isAuthenticated: auth?.isAuthenticated,
      user: auth?.user,
      roleID: roleID
    });
    
    if (!roleID) {
      console.warn('⚠️ No roleID found in user data');
      return ROUTES_PATH_NOPERMISSION;
    }

    switch (roleID) {
      case RoleID.ADMIN:
        console.log('👑 Admin routes loaded');
        return ROUTES_PATH;
        
      case RoleID.TRADE_CONSIGN:
        console.log('💼 Trade consign routes loaded');
        return {
          ROUTE_MAIN: ROUTES_PATH.ROUTE_MAIN,
          ROUTE_IMPORTORDER: ROUTES_PATH.ROUTE_IMPORTORDER,
          ROUTE_RETURNORDER: ROUTES_PATH.ROUTE_RETURNORDER,
        };
        
      case RoleID.ACCOUNTING:
        // Accounting gets their routes
        return {
          ROUTE_MAIN: ROUTES_PATH.ROUTE_MAIN,
          ROUTE_IMPORTORDER: ROUTES_PATH.ROUTE_IMPORTORDER,
          // Add other allowed routes
        };
        
      case RoleID.WAREHOUSE:
        // Warehouse gets their routes
        return {
          ROUTE_MAIN: ROUTES_PATH.ROUTE_MAIN,
          ROUTE_RETURNORDER: ROUTES_PATH.ROUTE_RETURNORDER,
          // Add other allowed routes
        };
        
      case RoleID.VIEWER:
        // Viewer gets view-only routes
        return {
          ROUTE_MAIN: ROUTES_PATH.ROUTE_MAIN,
          // Add other view-only routes
        };
        
      default:
        console.warn('⚠️ Unknown role ID:', roleID);
        return ROUTES_PATH_NOPERMISSION;
    }
  };

  return (
    <Content className="site-layout-background" style={{ padding: 24, minHeight: 280 }}>
      <Routes>
        {Object.values(getRoutesByRole()).map((item: any) => (
          <Route path={item.PATH} key={item.KEY} Component={item.COMPONENT} />
        ))}
      </Routes>
    </Content>
  );
};

export default ContentLayout;
