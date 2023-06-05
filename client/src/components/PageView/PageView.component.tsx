import React from "react";
import { pb } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import { PageRecord } from "../../lib/datamodels";

interface PageProps {
  url: string;
  id: string;
}

interface PageMetadataResponse {
  title?: string;
  description?: string;
}

const Page = ({ url, id }: PageProps) => {
  const [deleteState, setDeleteState] = React.useState<DataState>(
    DataState.UNKNOWN
  );
  const [previewState, setPreviewState] = React.useState<DataState>(
    DataState.UNKNOWN
  );
  const [previewData, setPreviewData] = React.useState<
    PageMetadataResponse | undefined
  >(undefined);

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
        <h4>{previewData.title}</h4>
        <p>{previewData.description}</p>
        <span>{url}</span>
      </>
    );
  } else {
    body = (
      <>
        <p>{url}</p>
        {previewState === DataState.PENDING ? (
          <p>Loading preview...</p>
        ) : undefined}
      </>
    );
  }
  return (
    <div>
      {body}
      <button onClick={handleDelete}>X</button>
    </div>
  );
};

export const PageView = () => {
  const [pages, setPages] = React.useState<PageRecord[]>([]);

  React.useEffect(() => {
    pb.collection("pages")
      .getList<PageRecord>()
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
      pb.collection("pages").unsubscribe("*");
    };
  }, []);
  return (
    <div>
      {pages.map((e) => (
        <Page url={e.url} id={e.id} key={e.id} />
      ))}
    </div>
  );
};
