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

enum AccViews {
  TOKEN,
  SETTINGS,
}

export const AccountPage = () => {
  const [tkn, setTkn] = React.useState<string | undefined>();
  const { notifOk, notifErr } = React.useContext(NotificationCtx);

  const handleCopyToken = () => {
    tkn &&
      navigator.clipboard
        .writeText(tkn)
        .then(() => notifOk("Copied"))
        .catch(() => notifErr("Could not copy to clipboard"));
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
        <TabList sx={{ m: 1 }}>
          <Tab value={AccViews.SETTINGS}>Settings</Tab>
          <Tab value={AccViews.TOKEN}>Token</Tab>
        </TabList>
        <TabPanel></TabPanel>
        <TabPanel value={AccViews.TOKEN} sx={{ ["& > *"]: { my: 1 } }}>
          <Typography level="body1">
            You can generate a new token for use with the iOS shortcut here. The
            token will revoke any past ones, so you'll need to update it
            everywhere. The token will be displayed here for 30 seconds.
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
