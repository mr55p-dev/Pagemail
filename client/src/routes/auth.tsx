import React from "react";
import { useNavigate } from "react-router-dom";
import { Form } from "react-router-dom";
import { pb } from "../lib/pocketbase";

const AuthPage = () => {
  const nav = useNavigate();
  const [email, setEmail] = React.useState<string>("");
  // const [username, setUsername] = React.useState<string>("");
  const [password, setpassword] = React.useState<string>("");
  const [subscribe, setSubscribe] = React.useState<boolean>(true);

  const handleEmailChange = (e: React.FormEvent<HTMLInputElement>) => {
    setEmail(e.currentTarget.value);
  };

  const handlePasswordChange = (e: React.FormEvent<HTMLInputElement>) => {
    setpassword(e.currentTarget.value);
  };

  const handleSignin = () => {
    const handler = async () => {
      try {
        await pb.collection("users").authWithPassword(email, password);
        nav("/pages");
      } catch (err) {
        console.error(err);
      }
    };
    handler();
  };

  return (
    <>
      <div className="auth-wrapper">
        <h1>Authenticate</h1>
        <div>
          <h4>Sign in</h4>
          <Form></Form>
        </div>
        <div>
          <input
            type="email"
            onChange={handleEmailChange}
            value={email}
            id="email-field"
          />
          <label htmlFor="email-field">Email</label>
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
        </div>
      </div>
    </>
  );
};

export default AuthPage;
