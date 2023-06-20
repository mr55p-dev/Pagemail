import { Link, useNavigate } from "react-router-dom";
import { AuthState } from "../../lib/data";
import { useUser } from "../../lib/pocketbase";
import NavBrandDark from "../../assets/default-monochrome-a.svg";
import NavBrandLight from "../../assets/default-monochrome-light-a.svg";

import {
  Box,
  CircularProgress,
  Divider,
  Grid,
  IconButton,
  useColorScheme,
} from "@mui/joy";
import {
  AccountBoxOutlined,
  ArticleOutlined,
  DarkModeOutlined,
  LightModeOutlined,
  LogoutOutlined,
} from "@mui/icons-material";

export const Navbar = () => {
  const { authState, logout } = useUser();
  const { mode, setMode } = useColorScheme();
  const nav = useNavigate();
  return (
    <Box>
      <Grid
        container
        direction="row"
        justifyContent="space-between"
        alignItems="center"
        spacing={1}
        sx={{ flexGrow: 1, px: 1 }}
      >
        <Grid
          xs={4}
          display="flex"
          flexDirection="row"
          alignItems="center"
          gap={1}
        >
          <Link
            style={{
              textDecoration: "none",
              color: "unset",
              display: "grid",
              placeItems: "end center",
            }}
            to="/"
          >
            <img
              height="40px"
              src={mode === "light" ? NavBrandDark : NavBrandLight}
            />
          </Link>
        </Grid>
        <Grid xs="auto" display="flex" direction="row" gap={1} my={1}>
          {authState === AuthState.AUTH ? (
            <>
              <IconButton size="md" onClick={() => nav("/pages")}>
                <ArticleOutlined />
              </IconButton>

              <IconButton size="md" onClick={logout}>
                <LogoutOutlined />
              </IconButton>
            </>
          ) : authState === AuthState.PENDING ? (
            <>
              <IconButton size="md">
                <CircularProgress thickness={2} />
              </IconButton>
            </>
          ) : (
            <>
              <IconButton size="md" onClick={() => nav("/pages")}>
                <AccountBoxOutlined />
              </IconButton>
            </>
          )}

          <IconButton
            onClick={() => setMode(mode === "light" ? "dark" : "light")}
            sx={{ justifySelf: "end" }}
            size="md"
          >
            {mode === "light" ? <DarkModeOutlined /> : <LightModeOutlined />}
          </IconButton>
        </Grid>
      </Grid>
      <Divider />
    </Box>
  );
};
