import { Box, Card } from "@mui/joy";

export const CarouselItem = ({
  children,
  image,
  idx,
}: {
  children: React.ReactNode;
  image?: React.ReactNode;
  idx: string | number;
}) => {
  return (
    <>
      <div id={"card-" + idx.toString}>
        <Card variant="outlined">
          {children}

          {!!image ? (
              image
          ) : undefined}
        </Card>
      </div>
    </>
  );
};

export const Carousel = ({ children }: { children: React.ReactNode }) => {
  return (
    <>
      <Box
        sx={{
          display: "flex",
          alignItems: "stretch",
          gap: 1,
          py: 1,
		  overflowX: "auto",
		  maxWidth: "100%",
          scrollSnapType: "x mandatory",
          "& > *": {
            scrollSnapAlign: "center",
          },
        }}
      >
        {children}
      </Box>
    </>
  );
};
