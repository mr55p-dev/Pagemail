import { IosShare } from "@mui/icons-material";
import { Box, Link, Typography } from "@mui/joy";
import ImageURL1 from "../assets/screenshot-ios-1.png";
import ImageURL2 from "../assets/screenshot-ios-2.png";
import ImageURL3 from "../assets/screenshot-ios-3.png";
import {
  Carousel,
  CarouselItem,
} from "../components/Carousel/Carousel.component";

export const Tutorial = () => {
  return (
    <Box mt={1} mx="auto">
      <Typography level="h2">iOS setup instructions</Typography>
      <Typography>To install Pagemail on your device:</Typography>
      <Carousel>
        <CarouselItem idx={1} image={<img src={ImageURL1} height="500px" />}>
          <Typography level="body1">
            Open <Link href="https://pagemail.io">Pagemail in safari</Link>, and
            press the <IosShare /> icon at the bottom of the screen
          </Typography>
        </CarouselItem>
        <CarouselItem idx={2} image={<img src={ImageURL2} height="500px" />}>
          <Typography>Find "Add to homescreen" and save it</Typography>
        </CarouselItem>
        <CarouselItem idx={3} image={<img src={ImageURL3} height="500px" />}>
          <Typography>
            Launch the app, and continue to create an account!
          </Typography>
        </CarouselItem>
      </Carousel>
    </Box>
  );
};
