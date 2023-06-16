import React, { ChangeEventHandler } from "react";
import { AuthState } from "../../lib/data";
import { pb, useUser } from "../../lib/pocketbase";
import signinUrl from "../../assets/google-auth/2x/btn_google_signin_light_normal_web@2x.png";
import { useNotification } from "../../lib/notif";
import FormControl from "@mui/joy/FormControl";
import FormLabel from "@mui/joy/FormLabel";
import Input from "@mui/joy/Input";
import Button from "@mui/joy/Button";

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
    <>
      {authErr ? <div>{authErr.message}</div> : undefined}
      <FormControl>
        <FormLabel>Email</FormLabel>
        <Input
          type="email"
          onChange={handleEmail}
          value={email}
          id="email-field"
        />
      </FormControl>
      <FormControl>
        <FormLabel>Password</FormLabel>
        <Input
          type="password"
          onChange={handlePassword}
          value={password}
          id="password-field"
        />
      </FormControl>
      <Button
        onClick={handleSignin}
        disabled={authState === AuthState.PENDING}
        sx={{ m: 1 }}
      >
        Sign in
      </Button>
      <GoogleAuth />
    </>
  );
};

export const SignUp = () => {
  const { login, authState, authErr } = useUser();
  const { trigger, component } = useNotification();

  const [email, handleEmail] = useFormComponent("");
  const [password, handlePassword] = useFormComponent("");
  const [passwordCheck, handlePasswordCheck] = useFormComponent("");
  const [username, handleUsername] = useFormComponent("");
  const [subscribed, setSubscribed] = React.useState(true);

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
    <>
      {authErr ? <div>{authErr.message}</div> : undefined}
      {component}
      <button onClick={() => trigger("hello")}>Click meee</button>
      <div className="form-input">
        <label htmlFor="username-field">Name</label>
        <input
          type="text"
          onChange={handleUsername}
          value={username}
          id="username-field"
        />
      </div>
      <div className="form-input">
        <label htmlFor="email-field">Email</label>
        <input
          type="email"
          onChange={handleEmail}
          value={email}
          id="email-field"
        />
      </div>
      <div className="form-input">
        <label htmlFor="password-field">Password</label>
        <input
          type="password"
          onChange={handlePassword}
          value={password}
          id="password-field"
        />
      </div>
      <div className="form-input">
        <label htmlFor="password-check-field">Repeat password</label>
        <input
          type="password"
          onChange={handlePasswordCheck}
          value={passwordCheck}
          id="password-check-field"
        />
      </div>
      <div className="form-input">
        <label htmlFor="subscribe-field">Subscribe?</label>
        <input
          type="checkbox"
          onChange={() => setSubscribed((prev) => !prev)}
          checked={subscribed}
          id="subscribed-field"
        />
      </div>
      <button
        onClick={handleSignup}
        disabled={!valid || authState === AuthState.PENDING}
      >
        Sign Up
      </button>
      <GoogleAuth />
    </>
  );
};

const GoogleAuth = () => {
  const { login } = useUser();
  const handleGoogle = () => {
    login(async () => {
      await pb.collection("users").authWithOAuth2({ provider: "google" });
    });
  };

  return (
    <Button
      sx={{ mx: 2 }}
      onClick={handleGoogle}
    >
      <img src={signinUrl} width="200px" />
    </Button>
  );
};
