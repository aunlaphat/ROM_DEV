import { Avatar, Button, Dropdown, MenuProps, Tag } from "antd";
import React from "react";
import { MenuUnfoldOutlined, MenuFoldOutlined } from "@ant-design/icons";
import { HeaderBarStyle, TopBarDropDown, TopBarUser } from "./style";
import { Icon } from "../../../resources";
import { useAuthLogin } from "../../../hooks/useAuth";
import { useSelector } from "react-redux";
import { TextSmall } from "../../text";

const HeaderBar = ({ collapsed, toggle }: any) => {
  const { onLogout } = useAuthLogin();
  const user = useSelector((state: any) => state.authen);

  const items: MenuProps["items"] = [
    {
      key: "logout",
      label: (
        <Button
          type="text"
          icon={Icon.logout()}
          onClick={() => {
            onLogout();
          }}
        >
          ออกจากระบบ
        </Button>
      ),
    },
  ];

  return (
    <HeaderBarStyle
      className="site-layout-background"
      style={{
        padding: 0,
        backgroundColor: "white",
      }}
    >
      {React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
        className: "trigger",
        onClick: toggle,
        style: { margin: "0 20px", color: "black" },
      })}
      <TopBarUser>
        <TextSmall
          className="item-right-topbar"
          text={
            <>
              <Tag className="item-right-topbar" color="blue">
                {user.users.userID}
              </Tag>
              <Tag className="item-right-topbar" color="volcano">
                {user.users.userFullName}
              </Tag>
            </>
          }
        />
      </TopBarUser>
      <TopBarDropDown>
        <Dropdown menu={{ items }}>
          <Avatar
            src={`https://api.dicebear.com/7.x/miniavs/svg?seed=1`}
            size="large"
          />
        </Dropdown>
      </TopBarDropDown>
    </HeaderBarStyle>
  );
};

export default HeaderBar;
