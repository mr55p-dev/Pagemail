import { Link, useNavigate } from "react-router-dom";
import { AuthState } from "../../lib/data";
import { useUser } from "../../lib/pocketbase";
import brandUrl from "../../assets/default-monochrome-white.svg";

export const Navbar = () => {
  const { authState, logout } = useUser();
  const nav = useNavigate();
  const authed_controls = (
    <>
      <Link to="/pages">Pages</Link>
      <button onClick={logout}> Log out </button>
    </>
  );
  const unauthed_controls = <Link to="/auth">Log in</Link>;
  return (
    <nav>
      <img onClick={() => nav("/")} className="brand-img" src={brandUrl} />
      <div className="buttons">
        {authState === AuthState.AUTH ? authed_controls : unauthed_controls}
      </div>
    </nav>
  );
};
