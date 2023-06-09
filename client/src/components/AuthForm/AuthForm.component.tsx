import React, { ChangeEventHandler } from "react";
import { AuthState, DataState } from "../../lib/data";
import { pb, useUser } from "../../lib/pocketbase";
import signinUrl from "../../assets/google-auth/2x/btn_google_signin_light_normal_web@2x.png";

function useFormComponent(
  init: boolean
): [boolean, ChangeEventHandler<HTMLInputElement>];
function useFormComponent(
  init: string
): [string, ChangeEventHandler<HTMLInputElement>];
function useFormComponent(
  init: string | boolean
): [string | boolean, ChangeEventHandler<HTMLInputElement>] {
  const [val, setVal] = React.useState<string | boolean>(init);

  const handleValChange = (e: React.FormEvent<HTMLInputElement>) => {
    setVal(e.currentTarget.value);
  };

  return [val, handleValChange];
}

export const Login = () => {
  const { login, authState, authErr } = useUser();

  const [email, handleEmail] = useFormComponent("");
  const [password, handlePassword] = useFormComponent("");

  const handleSignin = () => {
    login(async () => {
      await pb.collection("users").authWithPassword(email, password);
    });
  };

  return (
    <div>
      {authErr ? <div>{authErr.message}</div> : undefined}
      <label htmlFor="email-field">Email</label>
      <input
        type="email"
        onChange={handleEmail}
        value={email}
        id="email-field"
      />
      <label htmlFor="password-field">password</label>
      <input
        type="password"
        onChange={handlePassword}
        value={password}
        id="password-field"
      />
      <button onClick={handleSignin} disabled={authState === AuthState.PENDING}>
        Sign in
      </button>
    </div>
  );
};

export const SignUp = () => {
  const { login, authState, authErr } = useUser();

  const [email, handleEmail] = useFormComponent("");
  const [password, handlePassword] = useFormComponent("");
  const [passwordCheck, handlePasswordCheck] = useFormComponent("");
  const [username, handleUsername] = useFormComponent("");
  const [subscribed, handleSubscribed] = useFormComponent(true);

  const handleSignup = () => {
    login(async () => {
      const data = {
        email: email,
        emailVisibility: true,
        password: password,
        passwordConfirm: passwordCheck,
        subscribed: subscribed,
        name: username,
      };
      await pb.collection("users").create(data);
      await pb.collection("users").authWithPassword(email, password);
    });
  };

  const [valid, setValid] = React.useState(true);
  React.useEffect(() => {
    setValid(password === passwordCheck);
  }, [password, passwordCheck]);

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "300px",
        margin: "16px auto",
      }}
    >
      {authErr ? <div>{authErr.message}</div> : undefined}
      <label htmlFor="username-field">Name</label>
      <input
        type="text"
        onChange={handleUsername}
        value={username}
        id="username-field"
      />
      <label htmlFor="email-field">Email</label>
      <input
        type="email"
        onChange={handleEmail}
        value={email}
        id="email-field"
      />
      <label htmlFor="password-field">Password</label>
      <input
        type="password"
        onChange={handlePassword}
        value={password}
        id="password-field"
      />
      <label htmlFor="password-check-field">Repeat password</label>
      <input
        type="password"
        onChange={handlePasswordCheck}
        value={passwordCheck}
        id="password-check-field"
      />
      <label htmlFor="subscribe-field">Subscribe?</label>
      <input
        type="checkbox"
        onChange={handleSubscribed}
        checked={subscribed}
        id="subscribed-check-field"
      />
      <button
        onClick={handleSignup}
        disabled={!valid || authState === AuthState.PENDING}
      >
        Sign Up
      </button>
    </div>
  );
};

export const AuthForm = () => {
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
