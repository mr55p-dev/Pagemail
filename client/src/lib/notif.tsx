/* eslint-disable @typescript-eslint/no-empty-function */
import { CheckCircle, Info, Warning } from "@mui/icons-material";
import {
  Alert,
  Button,
  Typography,
  ColorPaletteProp,
  CircularProgress,
} from "@mui/joy";
import React from "react";

enum NotifState {
  OK,
  INFO,
  ERR,
}

const colors: Record<NotifState, ColorPaletteProp> = {
  [NotifState.OK]: "success",
  [NotifState.INFO]: "info",
  [NotifState.ERR]: "danger",
};

const icons: Record<NotifState, React.ReactNode> = {
  [NotifState.OK]: <CheckCircle />,
  [NotifState.INFO]: <Info />,
  [NotifState.ERR]: <Warning />,
};

interface NotificationCtxAttrs {
  notification?: Notification;
  style?: NotifState;
  notifOk: (title: string, body?: string, icon?: React.ReactNode) => void;
  notifInfo: (title: string, body?: string, icon?: React.ReactNode) => void;
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
const DURATION = 400;

export const NotificationBanner = () => {
  const notif = React.useContext(NotificationCtx);
  const [progress, setProgress] = React.useState(0);
  const interval = React.useRef<NodeJS.Timeout>();

  React.useEffect(() => {
    if (notif.notification) {
      const increment = 100 / DURATION; // Calculate the increment value per millisecond
      interval.current = setInterval(() => {
        setProgress((prevProgress) => {
          const newProgress = prevProgress + increment;
          return newProgress >= 100 ? 100 : newProgress;
        });
      }, 1); // Increase progress every 1 millisecond

      return () => clearInterval(interval.current);
    }
  }, [notif]);

  React.useEffect(() => {
    if (progress >= 100) {
      notif.notifClear();
      setProgress(0);
    }
  }, [notif, progress]);

  if (!notif.notification || notif.style == null) return null;

  const handleClear = () => {
    notif.notifClear();
    setProgress(0);
  };
  return (
    <Alert
      variant="soft"
      color={colors[notif.style]}
      startDecorator={
        <CircularProgress
          determinate
          value={progress}
          variant="soft"
          color={colors[notif.style]}
        >
          {notif.notification.icon || icons[notif.style]}
        </CircularProgress>
      }
      endDecorator={
        <Button
          onClick={handleClear}
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
