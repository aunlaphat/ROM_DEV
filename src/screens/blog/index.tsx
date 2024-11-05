import { Link } from "react-router-dom";

const BlogApp = () => {
  return (
    <div>
      <h1>Blog Index</h1>
      <Link to="/home">Home</Link> section
      <Link to="/">Login</Link> section
    </div>
  );
};

export default BlogApp;
