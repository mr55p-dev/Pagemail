/* eslint-disable @typescript-eslint/no-empty-function */
import { CheckCircle, Warning } from "@mui/icons-material";
import { Alert, Button, Typography, ColorPaletteProp } from "@mui/joy";
import React from "react";

enum NotifState {
  OK,
  ERR,
}

const colors: Record<NotifState, ColorPaletteProp> = {
  [NotifState.OK]: "success",
  [NotifState.ERR]: "danger",
};

const icons: Record<NotifState, React.ReactNode> = {
  [NotifState.OK]: <CheckCircle />,
  [NotifState.ERR]: <Warning />,
};

interface NotificationCtxAttrs {
  notification?: Notification;
  style?: NotifState;
  notifOk: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifErr: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifClear: () => void;
}

interface Notification {
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
        notifErr: err,
        notifClear: clear,
      }}
    >
      {children}
    </NotificationCtx.Provider>
  );
};

export const NotificationBanner = () => {
  const notif = React.useContext(NotificationCtx);
  if (!notif.notification || notif.style == null) return;
  return (
    <Alert
      variant="soft"
      sx={{
        my: 1,
      }}
      color={colors[notif.style]}
      startDecorator={notif.notification.icon || icons[notif.style]}
      endDecorator={
        <Button
          onClick={notif.notifClear}
          variant="outlined"
          color={colors[notif.style]}
        >
          Close
        </Button>
      }
    >
      <Typography level="body1">{notif.notification.title}</Typography>
      <Typography level="body2">{notif.notification.text}</Typography>
    </Alert>
  );
};
