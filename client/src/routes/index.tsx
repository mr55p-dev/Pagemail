import { Box, Button, Link, Typography } from "@mui/joy";
import { useNavigate } from "react-router";

export const Index = () => {
  const nav = useNavigate();
  const handleCta = () => {
    nav("/auth");
  };
  return (
    <Box
      width="100%"
    >
      <Typography level="display1" textAlign="center">Never forget a link again</Typography>
      <Button
        size="lg"
        variant="solid"
        fullWidth
        sx={{ display: "block", mx: "auto", my: 2, maxWidth: "sm" }}
        color="primary"
        onClick={handleCta}
      >
        Get started!
      </Button>
      <Typography
        level="body1"
		textAlign="center"
		justifyContent="center"
        endDecorator={
          <Link href="https://www.icloud.com/shortcuts/6da395d20b9542d8aa5ee56e884f0c4b">
            Get the shortcut here!
          </Link>
        }
      >
        On iOS? 
      </Typography>
    </Box>
  );
};
