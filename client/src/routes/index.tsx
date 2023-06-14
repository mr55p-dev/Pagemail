import { useNavigate } from "react-router";
import "../styles/index.css";

export const Index = () => {
  const nav = useNavigate();
  const handleCta = () => {
    nav("/auth");
  };
  return (
    <>
      <div className="index-content">
        <div className="content">
          <div className="title-box">
            <h1 className="title">Pagemail</h1>
            <p className="subtitle sans">Super simple Read-it-later!</p>
          </div>
          <div className="cta sans">
            <button onClick={handleCta}>
              <p>Get started!</p>
            </button>
          </div>
        </div>
      </div>
    </>
  );
};
