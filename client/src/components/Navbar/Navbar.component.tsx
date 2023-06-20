import { Link, useNavigate } from "react-router-dom";
import { AuthState } from "../../lib/data";
import { useUser } from "../../lib/pocketbase";
import {
  Box,
  CircularProgress,
  Divider,
  Grid,
  IconButton,
  Typography,
  useColorScheme,
} from "@mui/joy";
import IconDefault from "../../assets/default-icon.svg";
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
          sx={{
            ["& > *"]: {
              my: 1,
            },
          }}
        >
          <Link to="/">
            <img height="36px" src={IconDefault} />
          </Link>
          <Typography level="h4" sx={{ py: 1 }}>
            <Link style={{ textDecoration: "none", color: "unset" }} to="/">
              Pagemail
            </Link>
          </Typography>
        </Grid>
        <Grid xs="auto" display="flex" direction="row" gap={1}>
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
