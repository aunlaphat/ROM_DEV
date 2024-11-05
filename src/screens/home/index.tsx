import { Link } from "react-router-dom";


const Home = () => {
  return (
    
    <div>
      <h1>Welcome!</h1>
      <p>
        Check out the <Link to="/blog">blog</Link> or the{" "}
        <Link to="/users">users</Link> section
      </p>
    </div>
  );
};




export default Home;
