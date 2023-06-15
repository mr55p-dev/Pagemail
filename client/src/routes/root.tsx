import { Outlet } from "react-router-dom";
import { Navbar } from "../components/Navbar/Navbar.component";

const Root = () => {
  return (
    <>
      <Navbar />
      <Outlet />
    </>
  );
};

export default Root;
