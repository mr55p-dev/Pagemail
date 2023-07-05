import { IosShare } from "@mui/icons-material";
import { Box, Link, Sheet, Stack, Typography } from "@mui/joy";
import ImageURL1 from "../assets/screenshot-ios-1.png";
import ImageURL2 from "../assets/screenshot-ios-2.png";
import ImageURL3 from "../assets/screenshot-ios-3.png";

export const Tutorial = () => {
  return (
    <Box mt={1} mx="auto">
      <Typography level="h2">iOS setup instructions</Typography>
      <Typography>
        To install Pagemail on your device:
        <ol>
          <li>
            Open <Link href="https://pagemail.io">Pagemail in safari</Link>
          </li>
          <li>
            Press the <IosShare /> icon at the bottom of the screen
          </li>
          <li>Find "Add to homescreen" and save it</li>
        </ol>
      </Typography>
      <Stack direction="row" flexWrap="wrap" justifyContent="center">
        <img src={ImageURL1} height="500px" />
        <img src={ImageURL2} height="500px" />
        <img src={ImageURL3} height="500px" />
      </Stack>
    </Box>
  );
};
