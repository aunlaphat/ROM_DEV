import styled from "styled-components";
import { Layout } from "antd";

const { Header } = Layout;

export const HeaderBarStyle = styled(Header)`
  padding: 0 !important;
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
`;

export const TopBarDropDown = styled.div`
  padding-top: 5px;
  padding-right: 5%;
`;

export const TopBarUser = styled.div`
  text-align-last: right;
  justify-content: flex-end;
  align-items: center;
  flex: 1;
  padding-right: 5px;
`;
