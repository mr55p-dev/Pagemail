import { Outlet, useNavigate } from "react-router-dom";
import { pb, useUser } from "../lib/pocketbase";

const Root = () => {
  const { user, logout } = useUser();
  const nav = useNavigate()
  return (
    <>
      <h1>Pagemail</h1>
      <div className="root">
        <h1>This is the root page talking</h1>
      </div>
      <div className="root-content">
        <Outlet />
      </div>
      {user ? <button onClick={() => logout()} >Log out</button> : undefined}
    </>
  );
};

export default Root;
