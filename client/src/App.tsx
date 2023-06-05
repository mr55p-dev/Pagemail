import "almond.css/dist/almond.min.css";
import "./App.css";
import React from "react";
import { getCurrentUser, pb } from "./lib/pocketbase";
import { DataState } from "./lib/data";
import { PageAdd } from "./components/PageAdd/PageAdd.component";
import { PageView } from "./components/PageView/PageView.component";
import UserContext from "./lib/context";
import { AuthForm } from "./components/AuthForm/AuthForm.component";
import { UserRecord } from "./lib/datamodels";

function App() {
  const [user, setUser] = React.useState<UserRecord | null>(getCurrentUser());
  const [authStatus, setAuthStatus] = React.useState<DataState>(
    user ? DataState.SUCCESS : DataState.UNKNOWN
  );
  React.useEffect(() => {
    console.log(user);
  }, [user]);
  React.useEffect(() => {
    console.log(authStatus);
  }, [authStatus]);

  pb.authStore.onChange(() => {
    setUser(getCurrentUser());
  });

  const handleSignout = () => {
    pb.authStore.clear();
    setAuthStatus(DataState.UNKNOWN);
  };

  return (
    <>
      <h1>Pagemail</h1>
      <UserContext.Provider
        value={{
          user: user,
          checkUser: () => setUser(getCurrentUser()),
          authStatus: authStatus,
          setAuthStatus: setAuthStatus,
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
