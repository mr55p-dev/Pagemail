import React from "react";
import ReactDOM from "react-dom/client";
import "normalize.css";
import "./index.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Root from "./routes/root.tsx";
import ErrorPage from "./error-page.tsx";
import AuthPage from "./routes/auth.tsx";
import PagesPage from "./routes/pages.tsx";
import Protected from "./components/Protected/Protected.component.tsx";
import { Index } from "./routes/index.tsx";

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
      },
      {
        path: "pages",
        element: (
          <Protected>
            <PagesPage />,
          </Protected>
        ),
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
