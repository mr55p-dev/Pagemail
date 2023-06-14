import { Link, Outlet } from "react-router-dom";
import "normalize.css";
import { useUser } from "../lib/pocketbase";

const Root = () => {
  const { logout } = useUser();
  return (
    <>
      <div className="content">
        <div className="backdrop">
          <div className="backdrop-fade" />
        </div>
        <nav>
            <Link to="/pages">Pages</Link>
            <Link to="/auth">Log in</Link>
            <button onClick={logout}>Log out</button>
        </nav>
        <Outlet />
      </div>
    </>
  );
};

export default Root;
