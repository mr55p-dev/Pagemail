import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Root from "./routes/root.tsx";
import ErrorPage from "./error-page.tsx";
import AuthPage from "./routes/auth.tsx";
import PagesPage from "./routes/pages.tsx";
import Protected from "./components/Protected/Protected.component.tsx";
import { Index } from "./routes/index.tsx";
import { CssVarsProvider, extendTheme } from "@mui/joy/styles";
import CssBaseline from "@mui/joy/CssBaseline";
import { NotificationProvider } from "./lib/notif.tsx";
import { AccountPage } from "./routes/account.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/",
        element: <Index />,
      },
      {
        path: "auth",
        element: <AuthPage />,
      },
      {
        path: "pages",
        element: (
          <Protected>
            <PagesPage />
          </Protected>
        ),
      },
      {
        path: "account",
        element: (
          <Protected>
            <AccountPage />
          </Protected>
        ),
      },
    ],
  },
]);

const thm = extendTheme({
  typography: {
    display1: {
      // `--joy` is the default CSS variable prefix.
      // If you have a custom prefix, you have to use it instead.
      // For more details about the custom prefix, go to https://mui.com/joy-ui/customization/using-css-variables/#custom-prefix
      background:
        "linear-gradient(135deg, rgba(218,89,39,1) 0%, rgba(232,172,60,1) 100%)",
      // `Webkit*` properties must come later.
      WebkitBackgroundClip: "text",
      WebkitTextFillColor: "transparent",
    },
  },
  colorSchemes: {
    light: {
      palette: {
        primary: {
          "50": "#fbf2e9",
          "100": "#f4d9be",
          "200": "#ecbf92",
          "300": "#e5a566",
          "400": "#de8c3b",
          "500": "#c47221",
          "600": "#99591a",
          "700": "#6d3f13",
          "800": "#41200b",
          "900": "#160d04",
        },
      },
    },
    dark: {
      palette: {
        primary: {
          "50": "#fbf0e9",
          "100": "#f4d3be",
          "200": "#ecb592",
          "300": "#e59767",
          "400": "#dd7a3b",
          "500": "#c46022",
          "600": "#984b1a",
          "700": "#6d3613",
          "800": "#41200b",
          "900": "#160b04",
        },
      },
    },
  },
});

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <CssVarsProvider theme={thm} defaultMode="system">
      <CssBaseline />
      <NotificationProvider>
        <RouterProvider router={router} />
      </NotificationProvider>
    </CssVarsProvider>
  </React.StrictMode>
);
