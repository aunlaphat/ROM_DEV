import { Spin } from "antd";
import { useSelector } from "react-redux";

const Loading = (props: any) => {
  const alert = useSelector((state: any) => state.alert);
  return (
    <div className="div-loading-center">
      <Spin className="loading-center" spinning={alert.loading} size={"large"}>
        {props.children}
      </Spin>
    </div>
  );
};

export default Loading;
