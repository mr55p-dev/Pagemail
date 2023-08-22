import React from "react";
import {
  AudioFileOutlined,
  ContentCopy,
  DeleteOutline,
  OpenInNew,
} from "@mui/icons-material";
import LinesEllipsis from "react-lines-ellipsis";
import { pb } from "../../lib/pocketbase";
import { PageRecord } from "../../lib/datamodels";
import {
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

export const Page = (pageProps: PageRecord) => {
  const { notifOk, notifErr } = React.useContext(NotificationCtx);

  const dt = new Date(pageProps.created);
  const dest = new URL(pageProps.url);

  const handleDelete = () => {
    pb.collection("pages")
      .delete(pageProps.id)
      .then(() => {
        notifOk("Deleted", `Page at ${dest.hostname} removed.`);
      })
      .catch(() => {
        notifErr("Failed");
      });
  };

  function requestReadability() {
    pb.send("/api/page/readability", {
      method: "GET",
      params: { page_id: pageProps.id },
      cache: "no-cache",
    }).then((res) => console.log(res));
  }

  const title = pageProps.title || dest.host + dest.pathname;

  return (
    <Grid xs={12} sm={6} md={4} maxHeight="400px">
      <Card variant="outlined" sx={{ height: "100%", boxShadow: "md" }}>
        <CardContent>
          <Link href={pageProps.url} target="_blank" maxWidth="100%">
            <Typography
              level="h5"
              sx={{
                maxWidth: "100%",
                wordBreak: pageProps.title ? "break-word" : "break-all",
              }}
            >
              <LinesEllipsis
                maxLine={pageProps.title ? "3" : "2"}
                text={title}
              />
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
          <Typography level="body1" mt={1}>
			{pageProps.description && <LinesEllipsis maxLine="4" text={pageProps.description} />}
          </Typography>
        </CardContent>
        <ButtonGroup
          variant="outlined"
          color="neutral"
          sx={{ mx: "auto", width: 1, ["& > *"]: { flexGrow: 1 } }}
        >
          <IconButton
            size="sm"
            color="primary"
            onClick={() => window.open(pageProps.url)}
          >
            <OpenInNew />
          </IconButton>
          {pageProps.is_readable && (
            <IconButton size="sm" onClick={requestReadability}>
              <AudioFileOutlined />
            </IconButton>
          )}
          <IconButton size="sm" onClick={handleDelete} color="danger">
            <DeleteOutline />
          </IconButton>
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
          <Page {...e} key={e.id} />
        ))}
      </Grid>
    </>
  );
}
