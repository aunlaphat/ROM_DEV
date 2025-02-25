import { Avatar, Button, Dropdown, MenuProps, Tag, Modal } from "antd";
import React from "react";
import { MenuUnfoldOutlined, MenuFoldOutlined } from "@ant-design/icons";
import { HeaderBarStyle, TopBarDropDown, TopBarUser } from "./style";
import { Icon } from "../../../resources";
import { useAuthLogin } from "../../../hooks/useAuth";
import { useSelector } from "react-redux";
import { TextSmall } from "../../text";
import { logger } from '../../../utils/logger';

const HeaderBar = ({ collapsed, toggle }: any) => {
  const { onLogout } = useAuthLogin();
  const user = useSelector((state: any) => state.auth.user);

  const handleLogout = () => {
    Modal.confirm({
      title: 'ยืนยันการออกจากระบบ',
      content: 'คุณต้องการออกจากระบบใช่หรือไม่?',
      okText: 'ออกจากระบบ',
      cancelText: 'ยกเลิก',
      onOk: () => {
        logger.auth('info', 'User logging out', {
          userId: user?.userID,
          userName: user?.userName
        });
        onLogout();
      }
    });
  };

  // Add null check for user
  const userId = user?.userID || 'N/A';
  const userName = user?.userName || 'N/A';
  const roleName = user?.roleName || 'N/A';

  const items: MenuProps["items"] = [
    {
      key: "logout",
      label: (
        <Button
          type="text"
          icon={Icon.logout()}
          onClick={handleLogout}
          danger
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
                {userId}
              </Tag>
              <Tag className="item-right-topbar" color="volcano">
                {userName}
              </Tag>
              <Tag className="item-right-topbar" color="green">
                {roleName}
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
