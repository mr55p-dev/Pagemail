import { Login, SignUp } from "../components/AuthForm/AuthForm.component";
import { Sheet, Typography, Tab, TabList, TabPanel, Tabs } from "@mui/joy";

enum AuthMethod {
  LOGIN,
  SIGNUP,
}
const AuthPage = () => {
  return (
    <>
      <Sheet
        variant="outlined"
        color="primary"
        sx={{
          maxWidth: "320px",
          mx: "auto", // margin left & right
          my: 4, // margin top & bottom
          py: 3, // padding top & bottom
          px: 2, // padding left & right
          display: "flex",
          flexDirection: "column",
          gap: 2,
          borderRadius: "sm",
          boxShadow: "md",
        }}
      >
        <Typography level="h2" sx={{ textAlign: "center" }}>
          Authenticate
        </Typography>
        <Tabs
          aria-label="Authentication tabs"
          sx={{ borderRadius: "lg" }}
          size="md"
          defaultValue={AuthMethod.LOGIN}
        >
          <TabList sx={{ mx: 2, mt: 2 }}>
            <Tab value={AuthMethod.LOGIN}>Log in</Tab>
            <Tab value={AuthMethod.SIGNUP}>Sign up</Tab>
          </TabList>
          <TabPanel sx={{ p: 2 }} value={AuthMethod.LOGIN}>
            <Login />
          </TabPanel>
          <TabPanel sx={{ p: 2 }} value={AuthMethod.SIGNUP}>
            <SignUp />
          </TabPanel>
        </Tabs>
      </Sheet>
    </>
  );
};

export default AuthPage;
