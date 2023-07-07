import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  Link,
  Stack,
  Typography,
} from "@mui/joy";
import { useNavigate } from "react-router";
import { iosShortcutLink } from "../lib/const";

const InfoCard = ({ children }: { children: React.ReactNode }) => {
  return (
    <Grid xs={12} md={4}>
      <Card variant="outlined" sx={{ height: "100%", boxShadow: "md" }}>
        <CardContent>{children}</CardContent>
      </Card>
    </Grid>
  );
};

export const Index = () => {
  const nav = useNavigate();
  const handleCta = () => {
    nav("/auth");
  };
  return (
    <Box>
      <Typography level="display1" textAlign="center" mt={4} mb={2}>
        Never forget a link again
      </Typography>
      <Box
        sx={{
          width: "auto",
          mx: "auto",
        }}
      >
        <Stack
          direction={{ xs: "column", md: "row" }}
          spacing={1}
          maxWidth="sm"
          mx="auto"
		  my={2}
        >
          <Button
            fullWidth
            variant="solid"
            color="primary"
            sx={{ p: 1 }}
            onClick={handleCta}
          >
            <Typography fontSize="xl" sx={{ color: "white" }}>
              Get started!
            </Typography>
          </Button>
          <Button
            fullWidth
            variant="solid"
            color="neutral"
            sx={{ p: 1 }}
            onClick={() => nav("/tutorial")}
          >
            <Typography fontSize="xl" sx={{ color: "white" }}>
              Setup
            </Typography>
          </Button>
        </Stack>
      </Box>

      <Box>
        <Typography level="body1" fontSize="lg">
          <Typography color="primary">
            <b>
              Save, organize, and enjoy your favorite web pages with Pagemail{" "}
            </b>
          </Typography>
          - the ultimate read-it-later application. Say goodbye to cluttered
          bookmarks and never miss out on interesting articles, blog posts, or
          web pages again. Pagemail empowers you to curate your own digital
          library, and delivers a curated collection of your saved content right
          to your inbox every morning.
        </Typography>
      </Box>

      <Grid container my={2} spacing={2}>
        <InfoCard>
          <Typography level="h4" mb={1}>
            Save and Access Your Web Pages with Ease
          </Typography>
          <Typography level="body1">
            <Typography level="body1" color="primary" fontSize="md">
              <b>Effortlessly save and access web pages whenever you want </b>
            </Typography>
            - Whether it's a captivating article, a handy tutorial, or an
            inspiring blog post, simply save it to Pagemail and keep it at your
            fingertips for future reading. No more endless searching or
            forgetting about those hidden gems.
          </Typography>
        </InfoCard>

        <InfoCard>
          <Typography level="h4" mb={1}>
            Daily Email Digests
          </Typography>
          <Typography level="body1">
            <Typography level="body1" color="primary">
              <b>Stay up to date with a personalized email digest</b>
            </Typography>{" "}
            - Each morning, Pagemail compiles all the web pages you've saved
            over the past 24 hours and sends you a beautifully designed email.
            Start your day by diving into a your own collection of articles,
            organized and ready to explore.
          </Typography>
        </InfoCard>

        <InfoCard>
          <Typography level="h4" mb={1}>
            Integration with Your Mobile Device
          </Typography>
          <Typography level="body1">
            <Typography level="body1" color="primary">
              <b>Take Pagemail with you, wherever you go</b>
            </Typography>{" "}
            - Access your saved web pages on the go with our easy-to-install web
            app, compatible with both iOS and Android devices. For iOS, open
            this page in safari, click the share button and select "Add to
            homescreen". We also offer a shortcut for iOS devices, to let you
            save to Pagemail from anywhere.{" "}
            <Link href={iosShortcutLink}>
              Check it out here!
            </Link>
          </Typography>
        </InfoCard>
      </Grid>

      <Grid container spacing={4} my={2}></Grid>
    </Box>
  );
};
