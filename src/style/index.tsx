import { Col, Row } from "antd";
import styled from "styled-components";
import { color } from "../resources";

export const RowComponent = styled(Row)`
  display: flex;
  align-items: center;
`;

export const Required = styled.div`
  color: ${color.red};
  position: relative;
  top: 3px;
  left: 3px;
`;

export const ContainerButton = styled(Col)<{ align: string }>`
  text-align-last: ${(props) =>
    props.align === "left"
      ? "start"
      : props.align === "center"
      ? "center"
      : props.align === "right"
      ? "end"
      : ""};
  width: 100%;
`;

export const CenterContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 75vh;
`;

export const A5Container = styled.div`
  width: 148mm;
  height: 210mm;
  background-color: white;
  border: 1px solid black;
  padding: 10px;
`;

export const CenterInput = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;
