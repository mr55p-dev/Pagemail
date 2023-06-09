import { PageAdd } from "../components/PageAdd/PageAdd.component";
import { PageView } from "../components/PageView/PageView.component";

const PagesPage = () => {
  return (
    <>
      <div className="pages-wrapper">
        <h1>Pages</h1>
		<PageAdd />
		<PageView />
      </div>
    </>
  );
};

export default PagesPage;
