import { Outlet } from "react-router-dom";
import { Navbar } from "../components/Navbar/Navbar.component";
import { NotificationBanner } from "../lib/notif";
import { Container } from "@mui/joy";

const Root = () => {
  return (
    <Container maxWidth="md">
      <Navbar />
      <NotificationBanner />
      <Outlet />
    </Container>
  );
};

export default Root;
