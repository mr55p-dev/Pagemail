import { Outlet } from "react-router-dom";
import { useUser } from "../lib/pocketbase";

const Root = () => {
  const { user, logout } = useUser();
  return (
    <>
      <div className="root">
        <h1>Pagemail</h1>
        {user ? <button onClick={() => logout()}>Log out</button> : undefined}
      </div>
      <div className="content">
        <Outlet />
      </div>
    </>
  );
};

export default Root;
