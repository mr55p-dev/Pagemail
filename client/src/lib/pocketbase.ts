import PocketBase, { Record } from "pocketbase";
import { UserRecord } from "./datamodels";

// export const pb = new PocketBase('http://127.0.0.1:8090')
export const pb = new PocketBase("http://192.168.68.60:8090");
pb.autoCancellation(false);

export const getCurrentUser = (): UserRecord | null => {
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
