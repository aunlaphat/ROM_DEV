import styled from "styled-components";

const TextStyle = styled.p<{
  color?: string;
  bold?: string;
  size?: string;
  align?: string;
}>`
  font-weight: ${(props) => (props.bold ? "bold" : "")};
  color: ${(props) => (props.color ? props.color : "black")};
  font-size: ${(props) => (props.size ? props.size : "12px")};
  text-align: ${(props) => (props.align ? props.align : "")};
  margin: 0;
`;

export const TextXSMall = ({ color, bold, text, align }: any) => {
  return (
    <TextStyle color={color} bold={bold} size={"12px"} align={align}>
      {text}
    </TextStyle>
  );
};

export const TextSmall = ({ color, bold, text, align, onClick }: any) => {
  return (
    <TextStyle
      color={color}
      bold={bold}
      size={"14px"}
      align={align}
      onClick={onClick}
    >
      {text}
    </TextStyle>
  );
};

export const TextLarge = ({ color, bold, text, align }: any) => {
  return (
    <TextStyle color={color} bold={bold} size={"18px"} align={align}>
      {text}
    </TextStyle>
  );
};

export const TextLogoLogin = ({ color, bold, text, align, size }: any) => {
  return (
    <TextStyle color={color} bold={bold} size={size || "24px"} align={align}>
      {text}
    </TextStyle>
  );
};

export const TextInputLabel = ({ color, bold, text, align, ...props }: any) => {
  return (
    <div>
      <TextStyle
        style={{ color }}
        bold={true}
        size={"11px"}
        align={align}
        {...props}
      >
        {text}
      </TextStyle>
    </div>
  );
};
