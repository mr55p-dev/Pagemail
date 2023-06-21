import React from "react";

export const useTimeoutProgress = (
  duration: number,
  shouldShow: boolean,
  onTimeout: () => void
) => {
  const [progress, setProgress] = React.useState(0);
  const interval = React.useRef<NodeJS.Timeout>();

  React.useEffect(() => {
    if (!shouldShow) {
      return;
    }
    const increment = 1 / duration; // Calculate the increment value per millisecond
    interval.current = setInterval(() => {
      setProgress((prevProgress) => {
        const newProgress = prevProgress + increment;
        return newProgress >= 100 ? 100 : newProgress;
      });
    }, 1); // Increase progress every 1 millisecond

    return () => clearInterval(interval.current);
  }, [duration, shouldShow]);

  const cancel = React.useCallback(() => {
    onTimeout();
    setProgress(0);
  }, [onTimeout]);

  React.useEffect(() => {
    if (progress >= 100) {
      cancel();
    }
  }, [cancel, onTimeout, progress]);

  return {
    cancel,
    progress,
  };
};

