import { Outlet, useNavigate } from "react-router-dom";
import { pb, useUser } from "../lib/pocketbase";

const Root = () => {
  const { user } = useUser();
  const nav = useNavigate()
  const handleSignout = () => {
    pb.authStore.clear();
  };
  return (
    <>
      <h1>Pagemail</h1>
      <div className="root">
        <h1>This is the root page talking</h1>
      </div>
      <div className="root-content">
        <Outlet />
      </div>
      {user ? <button onClick={() => handleSignout()} >Log out</button> : <button onClick={() => nav("/auth")}>Log in</button>}
    </>
  );
};

export default Root;
