import { Box, Button, Typography } from "@mui/joy";
import { pb, useUser } from "../lib/pocketbase";
import { NotificationCtx } from "../lib/notif";
import React from "react";

export function Verify() {
  const { user } = useUser();
  const { notifOk, notifErr } = React.useContext(NotificationCtx);

  function requestVerification() {
    if (user && user.email) {
      pb.collection("users")
        .requestVerification(user?.email)
        .then(() => {
          notifOk("Sent verification email");
        })
        .catch(e => {
          notifErr("Failed to send verification email", e.status);
        });
    }
  }
  return (
    <>
      <Box>
        <Typography level="h1">Please check your inbox</Typography>
        <Typography>
          You will receive an email, click the link inside to verify your
          account.
        </Typography>
        <Typography>
          No email? <Button onClick={requestVerification}>Click here</Button>
        </Typography>
      </Box>
    </>
  );
}
