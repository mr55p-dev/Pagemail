import { Outlet } from "react-router-dom";
// import { useUser } from "../lib/pocketbase";

const Root = () => {
  // const { user, logout } = useUser();
  return (
    <>
      <div className="content">
        <Outlet />
      </div>
    </>
  );
};

export default Root;
