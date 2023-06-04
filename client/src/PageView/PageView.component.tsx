import React from "react";
import { pb } from "../lib/pocketbase";
import { Record } from "pocketbase";

interface PageRecord {
  id: string;
  url: string;
  user_id: string;
}

const Page = ({ url }: { url: string }) => {
  return <p>{url}</p>;
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
        <Page url={e.url} key={e.id} />
      ))}
    </div>
  );
};
