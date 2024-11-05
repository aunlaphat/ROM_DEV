import { color } from "../../../resources";
import styled from "styled-components";

export const ScrollMenu = styled.div`
  overflow-y: auto;
  height: calc(100vh - 200px);
  padding-bottom: 15px;

  &::-webkit-scrollbar {
    background-color: transparent;
    width: 10px;
  }
  &::-webkit-scrollbar-track {
    background-color: transparent;
  }
  &::-webkit-scrollbar-thumb {
    background-color: ${color.theme};
    box-shadow: inset 2px 2px 5px 0 rgba(#fff, 0.5);
    border-radius: 100px;
  }
  &::-webkit-scrollbar-button {
    display: none;
  }
`;
