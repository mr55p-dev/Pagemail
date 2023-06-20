import React from "react";
import { pb } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import { PageRecord } from "../../lib/datamodels";
import {
  Box,
  Button,
  ButtonGroup,
  Card,
  CardContent,
  CardOverflow,
  Grid,
  Link,
  Typography,
} from "@mui/joy";
import { DeleteOutline, OpenInNew } from "@mui/icons-material";

interface PageProps {
  url: string;
  id: string;
  created: string;
}

interface PageMetadataResponse {
  title?: string;
  description?: string;
}

const Page = ({ url, id, created }: PageProps) => {
  const [deleteState, setDeleteState] = React.useState<DataState>(
    DataState.UNKNOWN
  );
  const [previewState, setPreviewState] = React.useState<DataState>(
    DataState.UNKNOWN
  );
  const [previewData, setPreviewData] = React.useState<
    PageMetadataResponse | undefined
  >(undefined);

  const dt = new Date(created);
  const dest = new URL(url);

  React.useEffect(() => {
    setPreviewState(DataState.PENDING);
    const fetchLocal = async () => {
      try {
        const res = await pb.send<PageMetadataResponse>("/api/preview", {
          method: "GET",
          params: { target: url },
          // mode: "same-origin",
          cache: "default",
        });
        if (!res.title && !res.description) {
          throw new Error("Service returned no title or description");
        }
        setPreviewData(res);
        setPreviewState(DataState.SUCCESS);
      } catch (e) {
        console.error(e);
        setPreviewState(DataState.FAILED);
      }
    };
    fetchLocal();
  }, [url]);

  const handleDelete = () => {
    pb.collection("pages")
      .delete(id)
      .then(() => setDeleteState(DataState.SUCCESS))
      .catch(() => setDeleteState(DataState.FAILED));
  };

  switch (deleteState) {
    case DataState.PENDING:
      return <p>Deleting...</p>;
    case DataState.FAILED:
      return <p>Something went wrong deleting this!</p>;
  }
  let body;
  if (previewData) {
    body = (
      <>
        <Link href={url} target="_blank">
          <Typography level="h4">{previewData.title || url}</Typography>
        </Link>
        <Typography level="body2">{url}</Typography>
        <Typography level="body1" mt={1}>
          {previewData.description}
        </Typography>
      </>
    );
  } else {
    body = (
      <>
        <Link href={url} target="_blank">
          <Typography level="h4">{dest.hostname}</Typography>
        </Link>
		<Typography level="body2">{dest.origin}</Typography>
        {previewState === DataState.PENDING ? (
          <p>Loading preview...</p>
        ) : undefined}
      </>
    );
  }
  return (
    <Grid xs={12} sm={6} md={4}>
      <Card variant="outlined" sx={{ height: "100%", boxShadow: "md" }}>
        <CardContent>{body}</CardContent>
        <ButtonGroup variant="outlined" color="neutral" sx={{ mx: "auto" }}>
          <Button
            startDecorator={<OpenInNew />}
            color="success"
            onClick={() => window.open(url)}
          >
            Open
          </Button>
          <Button
            startDecorator={<DeleteOutline />}
            onClick={handleDelete}
            color="danger"
          >
            Delete
          </Button>
        </ButtonGroup>

        <CardOverflow sx={{ w: 1, bgcolor: "background.level1" }}>
          <Typography level="body1" sx={{ py: 1 }}>
            {dt.toLocaleDateString()} @ {dt.toLocaleTimeString()}
          </Typography>
        </CardOverflow>
      </Card>
    </Grid>
  );
};

export const PageView = () => {
  const [pages, setPages] = React.useState<PageRecord[]>([]);

  React.useEffect(() => {
    pb.collection("pages")
      .getList<PageRecord>(1, 50, {
        sort: "-created",
      })
      .then((records) => setPages(records.items));

    pb.collection("pages").subscribe<PageRecord>("*", function (e) {
      setPages((prev) => {
        switch (e.action) {
          case "create":
            return [e.record, ...prev];
          case "delete":
            return [...prev.filter((i) => i.id !== e.record.id)];
          default:
            return [...prev];
        }
      });
    });
    return () => {
      try {
        pb.collection("pages").unsubscribe("*");
      } catch (e) {
        console.error(e);
      }
    };
  }, []);
  return (
    <Grid container spacing={1} sx={{ flexGrow: 1, mt: 1 }}>
      {pages.map((e) => (
        <Page url={e.url} id={e.id} created={e.created} key={e.id} />
      ))}
    </Grid>
  );
};
