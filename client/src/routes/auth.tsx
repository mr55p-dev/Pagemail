import React from "react";
import { Login, SignUp } from "../components/AuthForm/AuthForm.component";

enum AuthMethod {
  LOGIN,
  SIGNUP,
}
const AuthPage = () => {
  const [method, setMethod] = React.useState<AuthMethod>(AuthMethod.LOGIN);

  return (
    <>
      <div className="auth-wrapper">
        <h1>Authenticate</h1>
        <div className="button-container">
          <button onClick={() => setMethod(AuthMethod.LOGIN)}>Log in</button>
          <button onClick={() => setMethod(AuthMethod.SIGNUP)}>Sign up</button>
        </div>
        <div className="form">
          {method === AuthMethod.LOGIN ? <Login /> : <SignUp />}
        </div>
      </div>
    </>
  );
};

export default AuthPage;
