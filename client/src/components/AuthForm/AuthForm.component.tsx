import React from "react";
import UserContext from "../../lib/context";
import { DataState } from "../../lib/data";
import { pb } from "../../lib/pocketbase";

export const AuthForm = () => {
  const { checkUser, authStatus, setAuthStatus } =
    React.useContext(UserContext);

  const [email, setEmail] = React.useState<string>("");
  const [username, setUsername] = React.useState<string>("");
  const [password, setpassword] = React.useState<string>("");
  const [errMsg, setErrMsg] = React.useState<string>("");

  const handleEmailChange = (e: React.FormEvent<HTMLInputElement>) => {
    setEmail(e.currentTarget.value);
  };

  const handleUsernameChange = (e: React.FormEvent<HTMLInputElement>) => {
    setUsername(e.currentTarget.value);
  };

  const handlepasswordChange = (e: React.FormEvent<HTMLInputElement>) => {
    setpassword(e.currentTarget.value);
  };

  const handleSignin = () => {
    setAuthStatus(DataState.PENDING);
    pb.collection("users")
      .authWithPassword(email, password)
      .then(() => {
        setAuthStatus(DataState.SUCCESS);
        checkUser();
      })
      .catch((err) => {
        setAuthStatus(DataState.FAILED);
        setErrMsg(`${err.status}: ${err.data.message}`);
      });
  };

  const handleSignup = () => {
    setAuthStatus(DataState.PENDING);
    // example create data
    const data = {
      username: username || undefined,
      email: email,
      emailVisibility: true,
      password: password,
      passwordConfirm: password,
      // name: "test",
    };
    pb.collection("users")
      .create(data)
      .then(() => setAuthStatus(DataState.SUCCESS))
      .catch((err) => {
        setAuthStatus(DataState.FAILED);
        setErrMsg(`${err.status}: ${err.data.message}`);
      });
  };

  switch (authStatus) {
    case DataState.PENDING:
      return <h3>Submitting...</h3>;
    case DataState.SUCCESS:
      return <p>Authenticated! This does not need to be displayed</p>;
    case DataState.UNKNOWN:
    case DataState.FAILED:
    default:
      return (
        <div>
          <div>
            <h3>Login</h3>
            {authStatus === DataState.FAILED ? <p>{errMsg}</p> : undefined}
            <input
              type="email"
              onChange={handleEmailChange}
              value={email}
              id="email-field"
            />
            <label htmlFor="email-field">Email</label>
            <input
              type="text"
              onChange={handleUsernameChange}
              value={username}
              id="username-field"
            />
            <label htmlFor="username-field">Username</label>
            <input
              type="password"
              onChange={handlepasswordChange}
              value={password}
              id="password-field"
            />
            <label htmlFor="password-field">password</label>
            <button onClick={handleSignin}>Sign in</button>
            <button onClick={handleSignup}>Sign up</button>
          </div>
        </div>
      );
  }
};
