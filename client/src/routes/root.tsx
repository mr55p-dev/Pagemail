import { Outlet } from "react-router-dom";
import { Navbar } from "../components/Navbar/Navbar.component";
import { NotificationBanner } from "../lib/notif";
import { Box, Container } from "@mui/joy";

const Root = () => {
  return (
    <Box maxWidth="md" mx="auto">
      <Navbar />
      <NotificationBanner />
      <Container maxWidth="md">
        <Outlet />
      </Container>
    </Box>
  );
};

export default Root;
