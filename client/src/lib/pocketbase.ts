import PocketBase, { Record } from "pocketbase";
import { UserRecord } from "./datamodels";

const pb_url = process.env.NODE_ENV === "development" ? 'http://127.0.0.1:8090' : 'https://v2.pagemail.io'
export const pb = new PocketBase(pb_url)
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
