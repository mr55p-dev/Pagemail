import { User } from "firebase/auth";
import { CollectionReference, DocumentData, FieldValue } from "firebase/firestore";

export interface IPageMetadata {
    title: string;
    description: string;
    author: string;
    image: string;
  }

export interface INotifState {
    title: string;
    text: string;
    style: "default" | "error" | "success";
}

export interface NotifCallback<T> {
  (setStateValue: T): void
}

export interface INotifContext {
  setNotifShow: NotifCallback<boolean>;
  setNotifState: NotifCallback<INotifState>;
}

export interface INotifProp {
  show: boolean;
  state: INotifState;
}

export interface IUserContext {
  user: User;
}

export interface IUserDoc {
  username: string
  email: string;
  photoURL: string;
  anonymous: boolean;
  newsletter: boolean;
}

export interface IUserData extends IUserDoc {
  user: User
  pages: CollectionReference<IPage>;
}

export interface IPage extends DocumentData {
  url: string;
  timeAdded: FieldValue;
}
