import { useRouteError } from "react-router";

const ErrorPage = () => {
  const error = useRouteError();
  return (
    <>
      <div className="error-wrapper">
        <h1>Error!</h1>
        <p>This page is not found</p>
        <p>{error.statusText || error.message}</p>
      </div>
    </>
  );
};

export default ErrorPage;
