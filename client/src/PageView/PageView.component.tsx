import React from "react";
import { pb } from "../lib/pocketbase";
import { DataState } from "../lib/data";

interface PageRecord {
  id: string;
  url: string;
  user_id: string;
}

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
      console.log("Making request");
      try {
        const res = await pb.send<PageMetadataResponse>("/api/preview", {
          method: "GET",
          params: { target: url },
        });
        setPreviewData(res);
        setPreviewState(DataState.SUCCESS);
      } catch (e) {
        console.log(e);
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

  return (
    <div>
      {deleteState === DataState.PENDING ? (
        <p>Deleting...</p>
      ) : deleteState === DataState.FAILED ? (
        <>
          <p>Failed to delete!</p>
          <button onClick={handleDelete}>X</button>
        </>
      ) : previewState === DataState.SUCCESS && previewData ? (
        <>
          <h4>{previewData.title}</h4>
          <p>{previewData.description}</p>
        </>
      ) : (
        <>
          <p>{url}</p>
          <button onClick={handleDelete}>X</button>
        </>
      )}
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
