import { Box, Link, Typography } from "@mui/joy";

export function Legacy() {
  return (
    <Box my={2}>
      <Typography mb={2} level="h2">Legacy users</Typography>
      <Typography>
        For any users who have data on the previous version of this site, it can
        still be reached <Link href="https://legacy.pagemail.io">here</Link>.
      </Typography>
      <Typography>
        Accounts have not, and will not be migrated, but I do encourage everyone
        to sign up for the V2 version of the platform (it's much better!). To
        get started, click <Link href="/tutorial">here</Link>.
      </Typography>
    </Box>
  );
}
