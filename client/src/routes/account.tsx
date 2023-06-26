import {
  Button,
  IconButton,
  Input,
  LinearProgress,
  Sheet,
  Stack,
  Typography,
} from "@mui/joy";
import { pb } from "../lib/pocketbase";
import React from "react";
import { NotificationCtx } from "../lib/notif";
import { useTimeoutProgress } from "../lib/timeout";
import { Cancel, ContentCopy } from "@mui/icons-material";

interface NewTokenRes {
  data: {
    token: string;
  };
}

export const AccountPage = () => {
  const [tkn, setTkn] = React.useState<string | undefined>();
  const { notifOk, notifErr } = React.useContext(NotificationCtx);

  const handleCopyToken = () => {
	tkn && navigator.clipboard.writeText(tkn)
	.then(() => notifOk("Copied"))
	.catch(() => notifErr("Could not copy to clipboard"))
  }

  const generateToken = async () => {
    try {
      const res = await pb.send<NewTokenRes>("/api/user/token/new", {});
      if (!res.data.token) {
        throw new Error("token not found in response");
      }
      setTkn(res.data.token);
    } catch (err) {
      notifErr(
        "Something went wrong fetching your new token.",
        (err as Error).message
      );
    }
  };

  const { progress, cancel } = useTimeoutProgress(120, !!tkn, () =>
    setTkn(undefined)
  );

  return (
    <Sheet
      variant="outlined"
      color="primary"
      sx={{
        my: 2,
        p: 2,
        borderRadius: "sm",
        boxShadow: "md",
        ["& > *"]: { my: 1 },
      }}
    >
      <Typography level="h4" textAlign="center">
        Account
      </Typography>
      <Typography>
        You can generate a new token for use with the iOS shortcut here. The
        token will revoke any past ones, so you'll need to update it everywhere.
        The token will be displayed here for 30 seconds.
      </Typography>
      {tkn ? (
        <>
          <Typography level="body1">Your token is:</Typography>
          <Input
            defaultValue={tkn}
			onClick={(e) => e.currentTarget.select()}
            endDecorator={
              <IconButton onClick={handleCopyToken} variant="plain">
                <ContentCopy />
              </IconButton>
            }
          />
          <Stack direction="row" alignItems="center" gap={2}>
            <LinearProgress determinate value={progress} />
            <IconButton onClick={cancel}>
              <Cancel />
            </IconButton>
          </Stack>
        </>
      ) : (
        <Button onClick={generateToken}>Generate a new token</Button>
      )}
    </Sheet>
  );
};
