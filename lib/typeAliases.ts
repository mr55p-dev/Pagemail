import { User } from "firebase/auth";

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
  username: string;
}
