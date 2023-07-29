import React from "react";
import { ContentCopy, DeleteOutline, OpenInNew } from "@mui/icons-material";
import LinesEllipsis from "react-lines-ellipsis";
import { pb } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import { PageRecord } from "../../lib/datamodels";
import {
  Button,
  ButtonGroup,
  Card,
  CardContent,
  CardOverflow,
  Grid,
  IconButton,
  Link,
  Stack,
  Typography,
} from "@mui/joy";
import { NotificationCtx } from "../../lib/notif";

export interface PageProps {
  url: string;
  id: string;
  created: string;
}

export interface PageMetadataResponse {
  title?: string;
  description?: string;
}

export const Page = ({ url, id, created }: PageProps) => {
  const [previewState, setPreviewState] = React.useState<DataState>(
    DataState.UNKNOWN
  );
  const [previewData, setPreviewData] = React.useState<
    PageMetadataResponse | undefined
  >(undefined);
  const { notifOk, notifErr } = React.useContext(NotificationCtx);

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
      .then(() => {
        notifOk("Deleted", `Page at ${dest.hostname} removed.`);
      })
      .catch(() => {
        notifErr("Failed");
      });
  };

  const previewTitle = previewData?.title;
  const ttl = previewData?.title ?? dest.host + dest.pathname;

  return (
    <Grid xs={12} sm={6} md={4} maxHeight="400px">
      <Card variant="outlined" sx={{ height: "100%", boxShadow: "md" }}>
        <CardContent>
          <Link href={url} target="_blank" maxWidth="100%">
            <Typography
              level="h5"
              sx={{
                maxWidth: "100%",
                wordBreak: previewTitle ? "break-word" : "break-all",
              }}
            >
              <LinesEllipsis maxLine={previewTitle ? "3" : "2"} text={ttl} />
            </Typography>
          </Link>
          <Typography
            level="body3"
            sx={{
              overflow: "hidden",
              whiteSpace: "nowrap",
              textOverflow: "ellipsis",
            }}
            startDecorator={
              <IconButton
                onClick={() =>
                  navigator.clipboard
                    .writeText(dest.toString())
                    .then(() => notifOk("Copied"))
                    .catch(() => notifErr("Could not save to the clipboard"))
                }
                variant="plain"
                size="sm"
              >
                <ContentCopy />
              </IconButton>
            }
          >
            {dest.toString()}
          </Typography>
          <Typography level="body2" mt={1}>
            {previewState === DataState.PENDING ? (
              "Loading preview..."
            ) : (
              <LinesEllipsis maxLine="4" text={previewData?.description} />
            )}
          </Typography>
        </CardContent>
        <ButtonGroup
          variant="outlined"
          color="neutral"
          sx={{ mx: "auto", width: 1, ["& > *"]: { flexGrow: 1 } }}
        >
          <Button
            startDecorator={<OpenInNew />}
            color="primary"
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
          <Typography level="body3" sx={{ py: 1 }}>
            {dt.getHours().toString().padStart(2, "0") +
              ":" +
              dt.getMinutes().toString().padStart(2, "0")}
          </Typography>
        </CardOverflow>
      </Card>
    </Grid>
  );
};

interface PageGroup {
  date: string;
  pages: PageRecord[];
}

function groupPages(pages: PageRecord[]): PageGroup[] {
  const groups = {} as Record<string, PageRecord[]>;
  pages.forEach((page) => {
    const dt = new Date(page.created).toLocaleDateString("en-gb") ?? "unknown";
    const grp = groups[dt];
    if (!grp) {
      groups[dt] = [];
    }
    groups[dt].push(page);
  });
  const today = new Date().toLocaleDateString("en-gb");
  const yesterday_dt = new Date();
  yesterday_dt.setDate(yesterday_dt.getDate() - 1);
  const yesterday = yesterday_dt.toLocaleDateString("en-gb");
  return Object.keys(groups)
    .sort((l, r) => new Date(r).getTime() - new Date(l).getTime())
    .map((k) => ({
      date: k === today ? "Today" : k === yesterday ? "Yesterday" : k,
      pages: groups[k],
    }));
}

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
    <Stack spacing={1} mt={2}>
      {groupPages(pages).map((g) => (
        <PageGroup {...g} />
      ))}
    </Stack>
  );
};

function PageGroup({ pages, date }: PageGroup) {
  return (
    <>
      <Stack direction="row" px={0.5}>
        <Typography level="body3" sx={{ pr: 1 }}>
          {date}
        </Typography>
        <div style={{ width: "100%" }}>
          <hr />
        </div>
      </Stack>
      <Grid container spacing={1} sx={{ flexGroup: 1, mt: 1 }}>
        {pages.map((e) => (
          <Page url={e.url} id={e.id} created={e.created} key={e.id} />
        ))}
      </Grid>
    </>
  );
}
