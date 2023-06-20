import { Box, Button, Typography } from "@mui/joy";
import { useNavigate } from "react-router";

export const Index = () => {
  const nav = useNavigate();
  const handleCta = () => {
    nav("/auth");
  };
  return (
    <div className="index-content">
      <Typography level="display1">Never forget a link again</Typography>
      <Box className="cta">
        <Button
          size="lg"
          variant="solid"
          sx={{ mx: "auto" }}
          color="success"
          onClick={handleCta}
        >
          Get started!
        </Button>
      </Box>
    </div>
  );
};
