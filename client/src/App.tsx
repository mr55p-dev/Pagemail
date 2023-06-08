import "./App.css";
import React from "react";
import { getCurrentUser, pb, useUser } from "./lib/pocketbase";
import { DataState } from "./lib/data";
import { PageAdd } from "./components/PageAdd/PageAdd.component";
import { PageView } from "./components/PageView/PageView.component";
import UserContext from "./lib/context";
import { AuthForm } from "./components/AuthForm/AuthForm.component";
import { UserRecord } from "./lib/datamodels";

function App() {
  const { user, setUser } = useUser();
  function handleSignout(
    event: MouseEvent<HTMLButtonElement, MouseEvent>
  ): void {
	pb.authStore.clear()
  }

  return (
    <>
      <h1>Pagemail</h1>
      <UserContext.Provider
        value={{
          user: user,
          checkUser: () => setUser(getCurrentUser()),
          authStatus: DataState.UNKNOWN,
          setAuthStatus: () => {},
        }}
      >
        {user ? (
          <button onClick={handleSignout}>Sign out</button>
        ) : (
          <AuthForm />
        )}
        <div>
          {authStatus === DataState.SUCCESS && user ? (
            <>
              <h3>Welcome, {user.username || user.email || "user"}</h3>
              <PageAdd />
              <PageView />
            </>
          ) : undefined}
        </div>
      </UserContext.Provider>
    </>
  );
}

export default App;
