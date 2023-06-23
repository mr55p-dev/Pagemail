import {
  Alert,
  Button,
  Typography,
  CircularProgress,
  Box,
  ColorPaletteProp,
} from "@mui/joy";
import React from "react";
import { useTimeoutProgress } from "../../lib/timeout";
import { NotifState, NotificationCtx } from "../../lib/notif";
import { CheckCircle, Info, Warning } from "@mui/icons-material";


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

export const NotificationBanner = () => {
    const notif = React.useContext(NotificationCtx);
    const { cancel, progress } = useTimeoutProgress(
      4,
      !!notif.notification,
      notif.notifClear
    );
  
    if (!notif.notification || notif.style == null) return null;
  
    return (
      <Box
      zIndex={10}
        sx={{
          position: "absolute",
          width: "100%",
          top: 0,
          left: 0,
        }}
      >
        <Alert
          variant="soft"
          color={colors[notif.style]}
          sx={{
            maxWidth: "sm",
            marginX: "auto",
            borderTopLeftRadius: 0,
            borderTopRightRadius: 0,
          boxShadow: "md",
          }}
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
              onClick={cancel}
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
      </Box>
    );
  };
