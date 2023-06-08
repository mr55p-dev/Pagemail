import PocketBase, { Record } from "pocketbase";
import { UserRecord } from "./datamodels";
import React from "react";

const pb_url = process.env.NODE_ENV === "development" ? 'http://127.0.0.1:8090' : 'https://v2.pagemail.io'
export const pb = new PocketBase(pb_url)
pb.autoCancellation(false);

export const getCurrentUser = (): UserRecord | null => {
  console.log("Checking current user")
  if (pb.authStore.model instanceof Record) {
    const mdl = pb.authStore.model;
    return {
      id: mdl.id,
      email: mdl.email,
      created: mdl.created,
      updated: mdl.updated,
      verified: mdl.verified ? true : false,
      subscribed: mdl.subscribed,
    };
  }
  return null;
};

export const useUser = () => {
  const [user, setUser] = React.useState<UserRecord | null>(getCurrentUser());
  // const [authStatus, setAuthStatus] = React.useState<DataState>( user ? DataState.SUCCESS : DataState.UNKNOWN);
  React.useEffect(() => {
	console.log(user)
  }, [user])

  pb.authStore.onChange(() => {
  console.log("Firing on auth store change")
    setUser(getCurrentUser());
  });

  return { user, setUser }
};


