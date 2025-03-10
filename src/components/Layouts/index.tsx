import React, { useEffect, useState, useCallback } from "react";
import { useLocation, useNavigate, Outlet } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { checkAuthen } from "../../redux/auth/action";
import Loading from "../loading";
import { Layout } from "antd";
import { logger } from "../../utils/logger";
import { ROUTES_NO_AUTH } from "../../resources/routes";
import SiderLayout from "../Layouts/siderLayout";
import ContentLayout from "../Layouts/contentLayout";
import HeaderBar from "../Layouts/headerLayout";
import { AppDispatch } from "../../redux/store";
import { RootState } from "../../redux/types";

const LayoutPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const authState = useSelector((state: RootState) => state.auth);

  const [collapsed, setCollapsed] = useState(window.innerWidth <= 767);
  const [collapsedWidth, setCollapsedWidth] = useState<number | undefined>(
    window.innerWidth <= 767 ? 0 : undefined
  );

  // ✅ โหลด token แค่ครั้งเดียว
  const token = localStorage.getItem("access_token");

  /**
   * ✅ ปรับปรุงการจัดการ Responsive Layout
   */
  useEffect(() => {
    const handleResize = () => {
      const isMobile = window.innerWidth <= 767;
      setCollapsed(isMobile);
      setCollapsedWidth(isMobile ? 0 : undefined);
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  /**
   * ✅ ใช้ useCallback เพื่อลดการ re-run ของ useEffect
   */
  const validateAuth = useCallback(async () => {
    logger.log("debug", "Validating auth state", {
      isAuthenticated: authState?.isAuthenticated,
      hasUser: !!authState?.user,
      path: location.pathname,
      hasToken: !!token,
    });

    if (token && !authState?.user) {
      logger.log("info", "Has token but no user data, fetching user info");
      await dispatch(checkAuthen());
      return;
    }

    if (!authState?.isAuthenticated && location.pathname !== ROUTES_NO_AUTH.ROUTE_LOGIN.PATH) {
      logger.log("warn", "Unauthorized access, redirecting to login");
      navigate(ROUTES_NO_AUTH.ROUTE_LOGIN.PATH);
    }
  }, [authState?.isAuthenticated, authState?.user, location.pathname, dispatch, navigate, token]);

  useEffect(() => {
    validateAuth();
  }, [validateAuth]);

  /**
   * ✅ ปรับปรุงการโหลดข้อมูล Auth
   */
  if (authState?.loading) {
    logger.log("debug", "Loading auth state...");
    return <Loading />;
  }

  if (!authState?.isAuthenticated && location.pathname !== "/") {
    logger.log("warn", "User not authenticated, blocking access");
    return null;
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <SiderLayout
        collapsed={collapsed}
        collapsedWidth={collapsedWidth}
        toggle={() => setCollapsed(!collapsed)}
      />
      <Layout className="site-layout">
        <HeaderBar collapsed={collapsed} toggle={() => setCollapsed(!collapsed)} />
        <ContentLayout>
          <Outlet />
        </ContentLayout>
      </Layout>
    </Layout>
  );
};

export default LayoutPage;
