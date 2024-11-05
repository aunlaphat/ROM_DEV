import styled from "styled-components";
import { Card } from "antd";

const CardStyle = styled(Card)`
  border-radius: 15px;
  border-width: 0px;
`;

export const CardTheme = ({ title, content, style, ...props }: any) => {
  return (
    <CardStyle title={title || ""} style={style} {...props}>
      {content}
    </CardStyle>
  );
};
