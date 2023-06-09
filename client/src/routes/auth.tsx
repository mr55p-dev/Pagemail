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
      <div className="auth-wrapper" style={{ margin: "16px"}}>
        <h1>Authenticate</h1>
        {method === AuthMethod.LOGIN ? (
          <div>
            <h4>Sign in</h4>
            <Login />
            <button onClick={() => setMethod(AuthMethod.SIGNUP)}>
              Not signed up?
            </button>
          </div>
        ) : (
          <div>
            <h4>Sign up</h4>
            <SignUp />
            <button onClick={() => setMethod(AuthMethod.LOGIN)}>
              Already signed up?
            </button>
          </div>
        )}
      </div>
    </>
  );
};

export default AuthPage;
