// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare const google: any;

import React from "react";
import UserContext from "../../lib/context";
import { DataState } from "../../lib/data";
import { pb } from "../../lib/pocketbase";

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
        setErrMsg(`${err.status}: ${err.data.message}`);
      }
    };
    handler();
  };

  const handleGoogle2 = (token: string) => {
    setAuthStatus(DataState.SUCCESS);
	
    const handler = async () => {
      try {
        const record = await pb.collection("users").authWithOAuth2Code({ provider: "google" });
		pb.authStore.save(token, )
        setAuthStatus(DataState.SUCCESS);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        setAuthStatus(DataState.FAILED);
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

  const btn = React.useRef(null);

  React.useEffect(() => {
    google.accounts.id.initialize({
      client_id:
        "556909502728-jjj3tpkoat64e0mot8vqlfjjqrd9l0ip.apps.googleusercontent.com",
      callback: handleGoogle,
      // login_uri: "https://v2.pagemail.io/api/oauth2-redirect",
    });

    google.accounts.id.renderButton(btn.current, {
      theme: "outline",
      size: "large",
    });
    console.log("google done");
  });

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
            <button onClick={handleGoogle}>Sign in with Google</button>

            <div ref={btn}></div>

            <div
              id="g_id_onload"
              data-client_id="556909502728-jjj3tpkoat64e0mot8vqlfjjqrd9l0ip.apps.googleusercontent.com"
              data-context="use"
              data-ux_mode="popup"
              data-callback="handleGoogle"
              data-auto_prompt="false"
            ></div>

            <div
              className="g_id_signin"
              data-type="standard"
              data-shape="rectangular"
              data-theme="outline"
              data-text="signin_with"
              data-size="large"
              data-logo_alignment="left"
            ></div>

            <script
              src="https://accounts.google.com/gsi/client"
              async
              defer
            ></script>
          </div>
        </div>
      );
  }
};
