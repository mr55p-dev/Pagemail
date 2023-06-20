import { Typography } from "@mui/joy";
import { useNavigate } from "react-router";

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
            <Typography level="display1">Never forget a link again</Typography>
          </div>
          <div className="cta">
            <button onClick={handleCta}>
              <p>Get started!</p>
            </button>
          </div>
        </div>
      </div>
    </>
  );
};
