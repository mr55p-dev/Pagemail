import React from "react";
import UserContext from "../../lib/context";
import { DataState } from "../../lib/data";
import { pb } from "../../lib/pocketbase";
import signinUrl from "../../assets/web/1x/btn_google_signin_light_normal_web.png";

export const AuthForm = () => {
  const { authStatus, setAuthStatus } = React.useContext(UserContext);

  const [email, setEmail] = React.useState<string>("");
  // const [username, setUsername] = React.useState<string>("");
  const [password, setpassword] = React.useState<string>("");
  const [subscribe, setSubscribe] = React.useState<boolean>(true);
  const [errMsg, setErrMsg] = React.useState<string>("");

  const handleEmailChange = (e: React.FormEvent<HTMLInputElement>) => {
    setEmail(e.currentTarget.value);
  };

  // const handleUsernameChange = (e: React.FormEvent<HTMLInputElement>) => {
  //   setUsername(e.currentTarget.value);
  // };

  const handlePasswordChange = (e: React.FormEvent<HTMLInputElement>) => {
    setpassword(e.currentTarget.value);
  };

  const handleGoogle = () => {
    setAuthStatus(DataState.PENDING);
    const handler = async () => {
      try {
        await pb.collection("users").authWithOAuth2({ provider: "google" });
        setAuthStatus(DataState.SUCCESS);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        setAuthStatus(DataState.FAILED);
        console.error(err);
        setErrMsg(`${err.status}: ${err.data.message}`);
      }
    };
    handler();
  };
  const handleSignin = () => {
    setAuthStatus(DataState.PENDING);
    const handler = async () => {
      try {
        await pb.collection("users").authWithPassword(email, password);
        setAuthStatus(DataState.SUCCESS);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        setAuthStatus(DataState.FAILED);
        setErrMsg(`${err.status}: ${err.data.message}`);
      }
    };
    handler();
  };

  const handleSignup = () => {
    setAuthStatus(DataState.PENDING);
    const handler = async () => {
      // example create data
      const data = {
        email: email,
        emailVisibility: true,
        password: password,
        passwordConfirm: password,
        subscribed: subscribe,
        // name: "test",
      };
      try {
        await pb.collection("users").create(data);
        await pb.collection("users").authWithPassword(email, password);
        setAuthStatus(DataState.SUCCESS);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        setAuthStatus(DataState.FAILED);
        setErrMsg(`${err.status}: ${err.data.message}`);
      }
    };
    handler();
  };

  switch (authStatus) {
    case DataState.PENDING:
      return <h3>Submitting...</h3>;
    case DataState.SUCCESS:
      return null;
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
            {/*
            <input
              type="text"
              onChange={handleUsernameChange}
              value={username}
              id="username-field"
            />
            <label htmlFor="username-field">Username</label>
			*/}
            <input
              type="password"
              onChange={handlePasswordChange}
              value={password}
              id="password-field"
            />
            <label htmlFor="password-field">password</label>
            <input
              type="checkbox"
              onChange={() => setSubscribe((prev) => !prev)}
              checked={subscribe}
              id="subscribe-field"
            />
            <label htmlFor="subscribe-field">Subscribe</label>
            <button onClick={handleSignin}>Sign in</button>
            <button onClick={handleSignup}>Sign up</button>
            <button onClick={handleGoogle}>
              <img src={signinUrl} width="200px" />
            </button>
          </div>
        </div>
      );
  }
};
