import React from "react";
import { pb } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";

export const PageAdd = ({ user_id }: { user_id: string }) => {
  const [url, setUrl] = React.useState<string>("");
  const [showSuccess, setShowSuccess] = React.useState<boolean>(false);
  const [dataState, setDataState] = React.useState<DataState>(
    DataState.UNKNOWN
  );

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setDataState(DataState.PENDING);
    const data = {
      url: url,
      user_id: user_id,
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
    <div>
      <h3>Add a page</h3>
      {showSuccess ? <p>Success!</p> : undefined}
      {dataState === DataState.PENDING ? (
        <p>Loading...</p>
      ) : (
        <form onSubmit={handleSubmit}>
          <input
            type="url"
            id="url-input"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            autoComplete="off"
          />
          <label htmlFor="url-input">URL</label>
          <button type="submit">Submit</button>
          <button type="reset" onClick={() => setUrl("")}>
            Clear
          </button>
        </form>
      )}
    </div>
  );
};
