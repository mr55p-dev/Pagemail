import React from "react";
import { pb } from "../lib/pocketbase";
import { Record } from "pocketbase";
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

const Page = ({ url, id }: PageProps) => {
  const [dataState, setDataState] = React.useState<DataState>(
    DataState.UNKNOWN
  );

  const handleDelete = () => {
    pb.collection("pages")
      .delete(id)
      .then(() => setDataState(DataState.SUCCESS))
      .catch(() => setDataState(DataState.FAILED));
  };
  return (
    <div>
      {dataState === DataState.PENDING ? (
        <p>Deleting...</p>
      ) : dataState === DataState.FAILED ? (
        <>
          <p>Failed to delete!</p>
          <button onClick={handleDelete}>X</button>
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
