import React, { useEffect, useState } from "react";
import { useLocation, useNavigate, Outlet } from "react-router-dom";
import { Layout } from "antd";
import { logger } from "../../utils/logger";
import { ROUTES_NO_AUTH } from "../../resources/routes";
import SiderLayout from "./siderLayout";
import ContentLayout from "./contentLayout";
import HeaderBar from "./headerLayout";
import Loading from "../loading";
import { useAuth } from "../../hooks/auth"; // นำเข้า hook ที่ refactor แล้ว

const LayoutPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  
  // ใช้ hook ที่เรา refactor แล้ว
  const { 
    isAuthenticated, 
    loading, 
    checkAuth, 
    hasInitialized,
    userID
  } = useAuth();

  const [collapsed, setCollapsed] = useState(window.innerWidth <= 767);
  const [collapsedWidth, setCollapsedWidth] = useState<number | undefined>(
    window.innerWidth <= 767 ? 0 : undefined
  );

  // จัดการ responsive layout
  useEffect(() => {
    const handleResize = () => {
      const isMobile = window.innerWidth <= 767;
      setCollapsed(isMobile);
      setCollapsedWidth(isMobile ? 0 : undefined);
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  // ตรวจสอบและจัดการการ authentication
  useEffect(() => {
    const validateAuth = async () => {
      // ตรวจสอบว่าหน้าปัจจุบันเป็นหน้า login หรือไม่
      const isLoginPage = location.pathname === ROUTES_NO_AUTH.ROUTE_LOGIN.PATH;
      // ตรวจสอบว่ามี token หรือไม่
      const hasToken = localStorage.getItem("access_token");

      logger.log("debug", "Validating auth state", {
        isAuthenticated,
        userID,
        path: location.pathname,
        hasToken: !!hasToken,
        hasInitialized
      });

      // ถ้ามี token แต่ยังไม่ได้ authenticate และไม่ได้กำลังโหลด
      if (hasToken && !isAuthenticated && !loading) {
        logger.log("info", "Has token but not authenticated, checking auth status");
        checkAuth();
        return;
      }

      // ถ้าไม่ได้ authenticate และไม่ได้อยู่ที่หน้า login ให้ redirect ไปหน้า login
      if (hasInitialized && !isAuthenticated && !isLoginPage) {
        logger.log("warn", "Unauthorized access, redirecting to login");
        navigate(ROUTES_NO_AUTH.ROUTE_LOGIN.PATH);
      }
      
      // ถ้า authenticated แล้วและอยู่ที่หน้า login ให้ redirect ไปหน้าหลัก
      if (hasInitialized && isAuthenticated && isLoginPage) {
        logger.log("info", "Already authenticated, redirecting to main page");
        navigate("/");
      }
    };

    validateAuth();
  }, [isAuthenticated, location.pathname, navigate, checkAuth, loading, userID, hasInitialized]);

  // แสดง loading ระหว่างตรวจสอบ auth
  if (loading || !hasInitialized) {
    logger.log("debug", "Loading auth state...");
    return <Loading />;
  }

  // หน้าที่ต้องการ authentication แต่ยังไม่ได้ authenticate
  if (!isAuthenticated && location.pathname !== ROUTES_NO_AUTH.ROUTE_LOGIN.PATH) {
    return null; // ไม่แสดงอะไรเพราะ useEffect จะทำการ redirect
  }

  // Main layout สำหรับผู้ใช้ที่ authenticate แล้ว
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