import React from "react";
import { pb, useUser } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import {
  Button,
  FormControl,
  FormLabel,
  IconButton,
  Input,
  Stack,
  Typography,
} from "@mui/joy";
import { AddCircleOutlined, ContentPaste } from "@mui/icons-material";

export const PageAdd = () => {
  const { user } = useUser();
  const [clipboardEnabled, setClipboardEnabled] = React.useState<boolean>(true);
  const [url, setUrl] = React.useState<string>("");
  const [showSuccess, setShowSuccess] = React.useState<boolean>(false);
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
        setShowSuccess(true);
        setTimeout(() => setShowSuccess(false), 1000);
      })
      .then(() => setUrl(""))
      .catch(() => setDataState(DataState.FAILED));
  };
  return (
    <>
      <Typography level="h4" my={1}>
        Add a page
      </Typography>
      {showSuccess ? <p>Success!</p> : undefined}
      {dataState === DataState.PENDING ? (
        <p>Loading...</p>
      ) : (
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
              endDecorator={
                <IconButton
                  onClick={handlePaste}
                  variant="plain"
                  disabled={!clipboardEnabled}
                >
                  <ContentPaste />
                </IconButton>
              }
            />
            <Button type="submit">Submit</Button>
          </Stack>
        </form>
      )}
    </>
  );
};
