import { Layout, Menu, Divider } from "antd";
import { Link, useLocation } from "react-router-dom";
import { ScrollMenu } from "./style";
import { MenuItemProps } from "../../../resources/interface";
import logo from "../../../assets/images/logo_org.png";
import { useAuthLogin } from "../../../hooks/useAuth";
import { useEffect } from "react";
import { ROUTES_PATH } from "../../../resources/routes-name";
import { Icon } from "../../../resources";
import { FormOutlined} from '@ant-design/icons';


const { Sider } = Layout;

const SiderLayout = ({ collapsed, collapsedWidth }: any) => {
  const location = useLocation();
  // const { checkLoginToken } = useAuthLogin();

  // useEffect(() => {
  //   checkLoginToken();
  // }, [location]);

  const ROUTES_MENU: MenuItemProps[] = [
    // {
    //   title: ROUTES_PATH.ROUTE_MAIN.LABEL,
    //   key: ROUTES_PATH.ROUTE_MAIN.PATH,
    //   icon: Icon.dashboard(),
    //   link: ROUTES_PATH.ROUTE_MAIN.PATH,
    //   role: [1, 2],
    // },
    // {
    //   title: ROUTES_PATH.ROUTE_ADJUST.LABEL,
    //   key: ROUTES_PATH.ROUTE_ADJUST.PATH,
    //   icon: Icon.dashboard(),
    //   link: ROUTES_PATH.ROUTE_ADJUST.PATH,
    //   role: [1, 2],
    // },
    // {
    //   title: ROUTES_PATH.ROUTE_PLATFORM.LABEL,
    //   key: ROUTES_PATH.ROUTE_PLATFORM.PATH,
    //   icon: Icon.dashboard(),
    //   link: ROUTES_PATH.ROUTE_PLATFORM.PATH,
    //   role: [1, 2],
    // },
    // {
    //   title: ROUTES_PATH.ROUTE_MANAGEMENT.LABEL,
    //   key: ROUTES_PATH.ROUTE_MANAGEMENT.PATH,
    //   icon: Icon.dashboard(),
    //   link: ROUTES_PATH.ROUTE_MANAGEMENT.PATH,
    //   role: [1, 2],
    // },
    // {
    //   title: "Sync MKP",
    //   key: "",
    //   icon: Icon.dashboard(),
    //   link:"",
    //   role: [1],
    //   subItems: [
    //     {
    //       title: ROUTES_PATH.ROUTE_SYNC.LABEL,
    //       key: ROUTES_PATH.ROUTE_SYNC.PATH,
    //       icon: Icon.dashboard(),
    //       link: ROUTES_PATH.ROUTE_SYNC.PATH,
    //       role: [1, 2],
    //     },
    //     {
    //       title: "Uplift",
    //       key: "/sync/subitem2",
    //       icon: Icon.dashboard(), // Replace with your icon
    //       link: "/sync/subitem2",
    //     },
    //   ],
    // },
    {
      title: ROUTES_PATH.   ROUTE_RETURNORDER.LABEL,
      key: ROUTES_PATH.  ROUTE_RETURNORDER.PATH,
      icon: Icon.Home(),
      link: ROUTES_PATH.   ROUTE_RETURNORDER.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.   ROUTE_CREATERETURN.LABEL,
      key: ROUTES_PATH.  ROUTE_CREATERETURN.PATH,
      icon: Icon.Edit1(),
      link: ROUTES_PATH.   ROUTE_CREATERETURN.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.   ROUTE_CREATETRADERETURN.LABEL,
      key: ROUTES_PATH.  ROUTE_CREATETRADERETURN.PATH,
      icon: Icon.Edit1(),
      link: ROUTES_PATH.   ROUTE_CREATETRADERETURN.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.   ROUTE_DRAFTANDCONFIRM.LABEL,
      key: ROUTES_PATH.  ROUTE_DRAFTANDCONFIRM.PATH,
      icon: Icon.Check(),
      link: ROUTES_PATH.   ROUTE_DRAFTANDCONFIRM.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.     ROUTE_CREATEBLIND.LABEL,
      key: ROUTES_PATH.    ROUTE_CREATEBLIND.PATH,
      icon: Icon.Edit1(),
      link: ROUTES_PATH.     ROUTE_CREATEBLIND.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.     ROUTE_CONFIRMRETURNTRADE.LABEL,
      key: ROUTES_PATH.    ROUTE_CONFIRMRETURNTRADE.PATH,
      icon: Icon.Check(),
      link: ROUTES_PATH.     ROUTE_CONFIRMRETURNTRADE.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH. ROUTE_IMPORTORDER.LABEL,
      key: ROUTES_PATH. ROUTE_IMPORTORDER.PATH,
      icon: Icon.BoxArrow(),
      link: ROUTES_PATH. ROUTE_IMPORTORDER.PATH,
      role: [1, 2],
    },
    
    {
      title: ROUTES_PATH.     ROUTE_SALERETURN.LABEL,
      key: ROUTES_PATH.    ROUTE_SALERETURN.PATH,
      icon: Icon.Return(),
      link: ROUTES_PATH.     ROUTE_SALERETURN.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.     ROUTE_OTHER.LABEL,
      key: ROUTES_PATH.    ROUTE_OTHER.PATH,
      icon: Icon.Return(),
      link: ROUTES_PATH.     ROUTE_OTHER.PATH,
      role: [1, 2],
    },
   
    {
      title: ROUTES_PATH.     ROUTE_REPORT.LABEL,
      key: ROUTES_PATH.    ROUTE_REPORT.PATH,
      icon: Icon.Report(),
      link: ROUTES_PATH.     ROUTE_REPORT.PATH,
      role: [1, 2],
    },
    {
      title: ROUTES_PATH.     ROUTE_MANAGEUSER.LABEL,
      key: ROUTES_PATH.    ROUTE_MANAGEUSER.PATH,
      icon: Icon.manageUser(),
      link: ROUTES_PATH.     ROUTE_MANAGEUSER.PATH,
      role: [1, 2],
    },
    // {
    //   title: ROUTES_PATH.   ROUTE_CREATEEXPENSE.LABEL,
    //   key: ROUTES_PATH.  ROUTE_CREATEEXPENSE.PATH,
    //   icon: Icon.dashboard(),
    //   link: ROUTES_PATH.   ROUTE_CREATEEXPENSE.PATH,
    //   role: [1, 2],
    // },
   
    
   
   
    
    // {
    //   title: "จัดการผู้ใช้งาน",
    //   key: "/management_user",
    //   icon: Icon.users(),
    //   link: "",
    //   role: [1],
    //   subItems: [],
    // },
  ];

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      collapsedWidth={collapsedWidth}
      theme="light"
      style={{
        position: "sticky",
        top: 0,
        height: "100vh",
        overflow: "auto",
      }}
    >
      <div className="logo">
        <img
          src={logo}
          alt="Avatar"
          className="avatar"
          style={{ width: "90%" }}
        />
      </div>
      <Divider />
      <ScrollMenu>
        <Menu
          theme="light"
          mode="inline"
          defaultSelectedKeys={[location.pathname]}
        >
          {ROUTES_MENU.map((item) => {
            // if (item.role?.includes(1)) {
              if (item.subItems) {
                return (
                  <Menu.SubMenu
                    key={item.key}
                    icon={item.icon}
                    title={item.title}
                  >
                    {item.subItems.map((subItem) => (
                      <Menu.Item key={subItem.key} icon={subItem.icon}>
                        <Link to={subItem.link}>{subItem.title}</Link>
                      </Menu.Item>
                    ))}
                  </Menu.SubMenu>
                );
              }
              return (
                <Menu.Item key={item.key} icon={item.icon}>
                  <Link to={item.link}>{item.title}</Link>
                </Menu.Item>
              );
            // } else {
            //   return null;
            // }
          })}
        </Menu>
      </ScrollMenu>
    </Sider>
  );
};

export default SiderLayout;
