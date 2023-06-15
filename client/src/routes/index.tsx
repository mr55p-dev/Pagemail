import { useNavigate } from "react-router";
import { useUser } from "../lib/pocketbase";
import { AuthState } from "../lib/data";
import { Navbar } from "../components/Navbar/Navbar.component";

export const Index = () => {
  const nav = useNavigate();
  const handleCta = () => {
    nav("/auth");
  };
  const { authState } = useUser();
  return (
    <>
      <Navbar />
      <div className="index-content">
        <div className="content">
          <div className="title-box">
            <h1 className="title">Never forget a link again</h1>
          </div>
          <div className="cta">
            <button onClick={handleCta}>
              <p>Get started!</p>
            </button>
            <p>{authState === AuthState.AUTH ? "AUTH" : "Not AUTH"}</p>
          </div>
        </div>
      </div>
    </>
  );
};
