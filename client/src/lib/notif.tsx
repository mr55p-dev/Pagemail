/* eslint-disable @typescript-eslint/no-empty-function */
import React from "react";

export enum NotifState {
  OK,
  INFO,
  ERR,
}


export interface NotificationCtxAttrs {
  notification?: Notification;
  style?: NotifState;
  notifOk: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifInfo: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifErr: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifClear: () => void;
}

export interface Notification {
  title: string;
  text?: string;
  icon?: React.ReactNode;
}

export const NotificationCtx = React.createContext<NotificationCtxAttrs>({
  notification: {
    title: "",
    text: "",
  },
  notifOk: () => {},
  notifInfo: () => {},
  notifErr: () => {},
  notifClear: () => {},
});

export const NotificationProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [notification, setNotif] = React.useState<Notification>();
  const [style, setStyle] = React.useState<NotifState>();

  const ok = (title: string, text?: string, icon?: React.ReactNode) => {
    setNotif({
      title,
      text,
      icon,
    });
    setStyle(NotifState.OK);
  };

  const info = (title: string, text?: string, icon?: React.ReactNode) => {
    setNotif({
      title,
      text,
      icon,
    });
    setStyle(NotifState.INFO);
  };

  const err = (title: string, text?: string, icon?: React.ReactNode) => {
    setNotif({
      title,
      text,
      icon,
    });
    setStyle(NotifState.ERR);
  };

  const clear = () => {
    setNotif(undefined);
  };

  return (
    <NotificationCtx.Provider
      value={{
        notification,
        style,
        notifOk: ok,
        notifInfo: info,
        notifErr: err,
        notifClear: clear,
      }}
    >
      {children}
    </NotificationCtx.Provider>
  );
};
