import { Box, Typography } from "@mui/joy";

export function Verify() {
  return (
    <>
      <Box>
        <Typography level="h1">Please check your inbox</Typography>
        <Typography>
          You should have received an email, click the link inside to verify
          your account.
        </Typography>
      </Box>
    </>
  );
}
