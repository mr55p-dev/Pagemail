import React from "react";
import { pb } from "../../lib/pocketbase";
import { DataState } from "../../lib/data";
import UserContext from "../../lib/context";

export const PageAdd = () => {
  const [clipboardEnabled, setClipboardEnabled] = React.useState<boolean>(true);
  const [url, setUrl] = React.useState<string>("");
  const [showSuccess, setShowSuccess] = React.useState<boolean>(false);
  const [dataState, setDataState] = React.useState<DataState>(
    DataState.UNKNOWN
  );

  // Get our user and do not render anything if not auth'd
  const { user } = React.useContext(UserContext);
  if (!user) {
    return undefined;
  }

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
      user_id: user.id,
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
          <button type="button" onClick={handlePaste} disabled={!clipboardEnabled}>Paste</button>
          <button type="reset" onClick={() => setUrl("")}>
            Clear
          </button>
        </form>
      )}
    </div>
  );
};
