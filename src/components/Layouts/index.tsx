import { Layout } from "antd";
import { useEffect, useState } from "react";
import SiderLayout from "./siderLayout";
import HeaderBar from "./headerLayout";
import ContentLayout from "./contentLayout";
import { useAuthLogin } from "../../hooks/useAuth";
import { useLocation } from "react-router-dom";
import { getCookies } from "../../store/useCookies";
import Loading from "../loading";

const LayoutPage = ({ children }: any) => {
  const location = useLocation();
  // const { checkLoginToken } = useAuthLogin();
  const [collapsed, setCollapsed] = useState(
    window.innerWidth <= 767 ? true : false
  );
  const [collapsedWidth, setCollapsedWidth] = useState(
    window.innerWidth <= 767 ? 0 : undefined
  );

  useEffect(() => {
    const handleResize = () => {
      setCollapsed(window.innerWidth <= 767 ? true : false);
      setCollapsedWidth(window.innerWidth <= 767 ? 0 : undefined);
    };

    window.addEventListener("resize", handleResize);

    handleResize();
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  const toggle = () => {
    setCollapsed(!collapsed);
  };

  // useEffect(() => {
  //   const handleClick = () => {
  //     if (!getCookies("jwt")) {
  //       // checkLoginToken();
  //     }
  //   };
  //   document.addEventListener("click", handleClick);
  //   return () => {
  //     document.removeEventListener("click", handleClick);
  //   };
  // }, [location]);

  // if (!getCookies("jwt")) {
  //   checkLoginToken();
  //   return <Loading />;
  // }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <SiderLayout
        collapsed={collapsed}
        collapsedWidth={collapsedWidth}
        toggle={toggle}
      />
      <Layout className="site-layout">
        <HeaderBar collapsed={collapsed} toggle={toggle} />
        <ContentLayout>{children}</ContentLayout>
      </Layout>
    </Layout>
  );
};

export default LayoutPage;
