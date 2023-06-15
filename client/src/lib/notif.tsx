import React from "react";

export const useNotification = () => {
  const [show, setShow] = React.useState(false);
  const [title, setTitle] = React.useState<string>("");

  React.useEffect(() => {
    if (!show) return;

    const timer = setTimeout(() => setShow(false), 1000);

    return () => {
      clearTimeout(timer);
      setShow(false);
    };
  }, [show, setShow]);

  const trigger = (ttl: string) => {
    setTitle(ttl);
    setShow(true);
  };

  return {
    trigger,
    component: (
      <div style={{ visibility: show ? "visible" : "hidden" }}>
        <p>{title}</p>
      </div>
    ),
  };
};
