import "almond.css/dist/almond.min.css";
import "./App.css";
import React from "react";
import { pb } from "./lib/pocketbase";
import { Record, RecordAuthResponse } from "pocketbase";

enum DataState {
  UNKNOWN,
  SUCCESS,
  FAILED,
  PENDING,
}

function App() {
  const [email, setEmail] = React.useState<string>("");
  const [password, setpassword] = React.useState<string>("");
  const [errMsg, setErrMsg] = React.useState<string>("");
  const [authState, setAuthState] = React.useState<DataState>(
    DataState.UNKNOWN
  );

  const handleEmailChange = (e: React.FormEvent<HTMLInputElement>) => {
    setEmail(e.currentTarget.value);
  };
  const handlepasswordChange = (e: React.FormEvent<HTMLInputElement>) => {
    setpassword(e.currentTarget.value);
  };

  const handleSignin = () => {
    setAuthState(DataState.PENDING);
    pb.collection("users")
      .authWithPassword(email, password)
      .then(() => setAuthState(DataState.SUCCESS))
      .catch(() => setAuthState(DataState.FAILED));
  };

  const handleSignup = () => {
    setAuthState(DataState.PENDING);
    // example create data
    const data = {
      // username: "test_username",
      email: email,
      emailVisibility: true,
      password: password,
      passwordConfirm: password,
      // name: "test",
    };
    pb.collection("users")
      .create(data)
      .then(() => setAuthState(DataState.SUCCESS))
      .catch((err) => {
        setAuthState(DataState.FAILED);
        setErrMsg(err.data.message);
      });
  };

  const handleSignout = () => {
    pb.authStore.clear();
    setAuthState(DataState.UNKNOWN);
  };

  return (
    <>
      <h1>Pagemail</h1>
      {authState !== DataState.SUCCESS ? (
        <>
          <div>
            <h3>Login</h3>
            <input
              type="email"
              onChange={handleEmailChange}
              value={email}
              id="email-field"
            />
            <label htmlFor="email-field">Email</label>
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
        </>
      ) : undefined}
      <div>
        {authState === DataState.SUCCESS ? (
          <>
            <h3>Welcome, {pb.authStore.model?.id}</h3>
            <button onClick={handleSignout}>Sign out</button>
          </>
        ) : authState === DataState.PENDING ? (
          <h3>Loading...</h3>
        ) : authState === DataState.FAILED ? (
          <>
            <h3>Failed to login</h3>
            <p>{errMsg}</p>
          </>
        ) : authState === DataState.UNKNOWN ? (
          <h3>Please sign in or sign up!</h3>
        ) : undefined}
      </div>
    </>
  );
}

export default App;
