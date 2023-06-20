import { Container, Typography } from "@mui/joy";
import { PageAdd } from "../components/PageAdd/PageAdd.component";
import { PageView } from "../components/PageView/PageView.component";

const PagesPage = () => {
  return (
    <>
      <Container maxWidth="md">
        <PageAdd />
        <PageView />
      </Container>
    </>
  );
};

export default PagesPage;
