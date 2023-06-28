import {
  Button,
  IconButton,
  Input,
  LinearProgress,
  Sheet,
  Stack,
  Tab,
  TabList,
  TabPanel,
  Tabs,
  Typography,
} from "@mui/joy";
import { pb, useUser } from "../lib/pocketbase";
import React from "react";
import { NotificationCtx } from "../lib/notif";
import { useTimeoutProgress } from "../lib/timeout";
import { Cancel, ContentCopy, Delete } from "@mui/icons-material";
import { UserRecord } from "../lib/datamodels";

interface NewTokenRes {
  data: {
    token: string;
  };
}

enum AccViews {
  TOKEN,
  SETTINGS,
}

export const AccountPage = () => {
  const { notifOk, notifErr } = React.useContext(NotificationCtx);
  const { user } = useUser();

  const [tkn, setTkn] = React.useState<string | undefined>();
  const [subscribed, setSubscribed] = React.useState<boolean>(
    user?.subscribed || false
  );

  const handleCopyToken = () => {
    tkn &&
      navigator.clipboard
        .writeText(tkn)
        .then(() => notifOk("Copied"))
        .catch(() => notifErr("Could not copy to clipboard"));
  };

  const handleSubscribeToggle = () => {
    if (user?.id) {
      pb.collection("users")
        .update<UserRecord>(user.id, {
          subscribed: !subscribed,
        })
        .then((data) => setSubscribed(data.subscribed))
        .catch((e) => notifErr("Failed to chanege subscription", e.message));
    }
  };

  const handleAccountDelete = () => {
    if (user?.id) {
      pb.collection("users")
        .delete(user.id)
        .then(() => notifOk("Account deleted"))
        .catch(() =>
          notifErr(
            "Something went wrong",
            "Your account has not been deleted. Please contact support."
          )
        );
    }
  };

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
        maxWidth: "sm",
        mx: "auto", // margin left & right
        my: 4, // margin top & bottom
        py: 3, // padding top & bottom
        px: 2, // padding left & right
        display: "flex",
        flexDirection: "column",
        gap: 2,
        borderRadius: "sm",
        boxShadow: "md",
      }}
    >
      <Typography level="h4" textAlign="center">
        Account
      </Typography>

      <Tabs
        aria-label="Account view tabs"
        sx={{ borderRadius: "lg" }}
        size="md"
        defaultValue={AccViews.SETTINGS}
      >
        <TabList sx={{ my: 1 }}>
          <Tab value={AccViews.SETTINGS}>Settings</Tab>
          <Tab value={AccViews.TOKEN}>Token</Tab>
        </TabList>
        <TabPanel value={AccViews.SETTINGS} sx={{ ["& > *"]: { my: 1 } }}>
          <Stack direction="row" useFlexGap justifyContent="space-evenly">
            <Button onClick={handleSubscribeToggle}>
              {subscribed ? "Unsubscribe" : "Subscribe"}
            </Button>
            <Button
              endDecorator={<Delete />}
              color="danger"
              onClick={handleAccountDelete}
            >
              Delete Account
            </Button>
          </Stack>
        </TabPanel>
        <TabPanel value={AccViews.TOKEN} sx={{ ["& > *"]: { mb: 2 } }}>
          <Typography level="body1">
            You can generate a new token for use with the iOS shortcut here. The
            token will revoke any past ones, so you'll need to update it
            everywhere.
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
            <Button onClick={generateToken} fullWidth>
              Generate a new token
            </Button>
          )}
        </TabPanel>
      </Tabs>
    </Sheet>
  );
};
