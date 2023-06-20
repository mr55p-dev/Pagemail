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

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/",
        element: <Index />,
        errorElement: <ErrorPage />,
      },
      {
        path: "auth",
        element: <AuthPage />,
        errorElement: <ErrorPage />,
      },
      {
        path: "pages",
        element: (
          <Protected>
            <PagesPage />
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
