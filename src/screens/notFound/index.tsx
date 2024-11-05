import notFoundImg from "../../assets/images/404.png";
import { Link } from "react-router-dom";

const NotfoundScene = () => {
  return (
    <div
      style={{
        textAlignLast: "center",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <img
        alt={notFoundImg}
        src={notFoundImg}
        style={{ minHeight: "100%", width: "50vw" }}
      />
      <Link to="/home">Click to Back Home</Link>
    </div>
  );
};

export default NotfoundScene;
