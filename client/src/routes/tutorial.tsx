import { IosShare } from "@mui/icons-material";
import { Box, Link, Typography } from "@mui/joy";
import ImageURL1 from "../assets/screenshot-ios-1.png";
import ImageURL2 from "../assets/screenshot-ios-2.png";
import ImageURL3 from "../assets/screenshot-ios-3.png";
import { iosShortcutLink } from "../lib/const";
import ShortcutShot1 from "../assets/screenshot-ios-shortcut-1.png";
import ShortcutShot2 from "../assets/screenshot-ios-shortcut-2.png";
import ShortcutShot3 from "../assets/screenshot-ios-shortcut-3.png";
import ShortcutShot4 from "../assets/screenshot-ios-shortcut-4.png";
import ShortcutShot5 from "../assets/screenshot-ios-shortcut-5.png";
import ShortcutShot6 from "../assets/screenshot-ios-shortcut-6.png";
import ShortcutShot7 from "../assets/screenshot-ios-shortcut-7.png";

import {
  Carousel,
  CarouselItem,
} from "../components/Carousel/Carousel.component";

export const Tutorial = () => {
  return (
    <Box mt={1} mx="auto">
      <Typography level="h2">iOS setup instructions</Typography>
      <Typography level="h6" mt={2}>To install Pagemail on your device:</Typography>
      <Carousel ids={["i1", "i2", "i3"]}>
        <CarouselItem id="i1" image={<img src={ImageURL1} height="500px" />}>
          <Typography level="body1">
            Open <Link href="https://pagemail.io">Pagemail in safari</Link>, and
            press the <IosShare /> icon at the bottom of the screen.
          </Typography>
        </CarouselItem>
        <CarouselItem id="i2" image={<img src={ImageURL2} height="500px" />}>
          <Typography>Find "Add to homescreen" and save it.</Typography>
        </CarouselItem>
        <CarouselItem id="i3" image={<img src={ImageURL3} height="500px" />}>
          <Typography>
            Launch the app, and continue to create an account!
          </Typography>
        </CarouselItem>
      </Carousel>

      <Typography level="h6" mt={2}>To install the shortcut</Typography>
      <Carousel ids={["s1", "s2", "s3", "s4", "s5", "s6", "s7"]}>
        <CarouselItem id="s1" image={<img src={ShortcutShot1} height="500px" />}>
          <Typography>
            Get your token from the account settings screen. Keep this token to
            yourself.
          </Typography>
        </CarouselItem>
        <CarouselItem id="s2" image={<img src={ShortcutShot2} height="500px" />}>
          <Typography>
            Copy it to your clipboard with the copy button.
          </Typography>
        </CarouselItem>
        <CarouselItem id="s3" image={<img src={ShortcutShot3} height="500px" />}>
          <Typography>
            Click <a href={iosShortcutLink}>here</a> to start installing the
            shortcut. You will need the token you just generated when setting it
            up.
          </Typography>
        </CarouselItem>
        <CarouselItem  id="s4" image={<img src={ShortcutShot4} height="500px" />}>
          <Typography>Paste your token in the field</Typography>
        </CarouselItem>
        <CarouselItem id="s5" image={<img src={ShortcutShot5} height="500px" />}>
          <Typography>
            Now the shortcut is installed, you need to enable it. Open up Safari
            to any webpage, and hit the <IosShare /> share icon. Scroll to the
            end and hit "edit actions" to bring up all actions installed on your
            device.
          </Typography>
        </CarouselItem>
        <CarouselItem id="s6" image={<img src={ShortcutShot6} height="500px" />}>
          <Typography>
            Click the green add button for Pagemail to enable it in any share
            sheet from any app!
          </Typography>
        </CarouselItem>
        <CarouselItem id="s7" image={<img src={ShortcutShot7} height="500px" />}>
          <Typography>
            Now that's enabled, you should see the "Share to Pagemail" action in
            the share sheet. Wherever you can share a URL, you can save it to
            Pagemail!
          </Typography>
        </CarouselItem>
      </Carousel>
    </Box>
  );
};
