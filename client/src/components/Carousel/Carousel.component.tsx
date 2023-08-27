import {
  ArrowBackIosRounded,
  ArrowForwardIosRounded,
} from "@mui/icons-material";
import { Box, Card, IconButton, Stack } from "@mui/joy";
import React from "react";

export const CarouselItem = ({
  children,
  image,
  id,
}: {
  children: React.ReactNode;
  image?: React.ReactNode;
  id?: string;
}) => {
  return (
    <>
      <Box>
        <Card
          id={id}
          variant="outlined"
          sx={{
            height: "100%",
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
          }}
        >
          {children}
          {image ? image : undefined}
        </Card>
      </Box>
    </>
  );
};

export const Carousel = ({
  children,
  ids,
}: {
  children: React.ReactNode;
  ids: string[];
}) => {
  const [current, setCurrent] = React.useState<number>(0);
  const handleNext = () => {
    window.location.assign(
      window.location.protocol +
        "//" +
        window.location.host +
        window.location.pathname +
        "#" +
        ids[current + 1]
    );
    setCurrent((p) => p + 1);
  };
  const handlePrev = () => {
    window.location.assign(
      window.location.protocol +
        "//" +
        window.location.host +
        window.location.pathname +
        "#" +
        ids[current - 1]
    );
    setCurrent((p) => p - 1);
  };
  return (
    <Box maxWidth="100%">
      <Box
        sx={{
          display: "flex",
          alignItems: "stretch",
          gap: 1,
          py: 1,
          overflowX: "auto",
          maxWidth: "100%",
          scrollSnapType: "x mandatory",
          ["::-webkit-scrollbar"]: { display: "none" },
          "& > *": {
            scrollSnapAlign: "center",
          },
        }}
      >
        {children}
      </Box>
      <Stack justifyContent="space-between" direction="row">
        <IconButton onClick={handlePrev} disabled={current <= 0}>
          <ArrowBackIosRounded />
        </IconButton>
        <IconButton onClick={handleNext} disabled={current >= ids.length - 1}>
          <ArrowForwardIosRounded />
        </IconButton>
      </Stack>
    </Box>
  );
};
