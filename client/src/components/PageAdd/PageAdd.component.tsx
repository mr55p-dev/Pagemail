import React from "react";
import { pb, useUser } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import {
  Button,
  CircularProgress,
  IconButton,
  Input,
  Stack,
  Typography,
} from "@mui/joy";
import { ContentPaste } from "@mui/icons-material";
import { NotificationCtx } from "../../lib/notif";

export const PageAdd = () => {
  const { user } = useUser();
  const { notifOk, notifErr } = React.useContext(NotificationCtx);
  const [clipboardEnabled, setClipboardEnabled] = React.useState<boolean>(true);
  const [url, setUrl] = React.useState<string>("");
  const [dataState, setDataState] = React.useState<DataState>(
    DataState.UNKNOWN
  );

  const handlePaste = () => {
    navigator.clipboard
      .readText()
      .then((txt) => setUrl(txt))
      .catch(() => setClipboardEnabled(false));
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setDataState(DataState.PENDING);
    const data = {
      url: url,
      user_id: user?.id,
    };
    pb.collection("pages")
      .create(data)
      .then(() => {
        setDataState(DataState.SUCCESS);
        notifOk("Success!", );
      })
      .then(() => setUrl(""))
      .catch(() => {
        notifErr("Failed");
        setDataState(DataState.FAILED);
      });
  };

  const isPending = dataState === DataState.PENDING;
  return (
    <>
      <Typography level="h4" my={1}>
		Save your pages
      </Typography>
      <form onSubmit={handleSubmit}>
        <Stack direction="row" gap={1}>
          <Input
            type="url"
            id="url-input"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            autoComplete="off"
            placeholder="URL"
            sx={{ width: "100%" }}
            disabled={isPending}
			size="lg"
            endDecorator={
              <IconButton
                onClick={handlePaste}
                variant="plain"
                disabled={!clipboardEnabled || isPending}
              >
                <ContentPaste />
              </IconButton>
            }
          />
          <Button
            type="submit"
            startDecorator={
              isPending && <CircularProgress variant="outlined" />
            }
            disabled={isPending}
          >
            Submit
          </Button>
        </Stack>
      </form>
    </>
  );
};
