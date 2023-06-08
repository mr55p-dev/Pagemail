import { PageAdd } from "../components/PageAdd/PageAdd.component";
import { PageView } from "../components/PageView/PageView.component";

const PagesPage = () => {
  return (
    <>
      <div className="pages-wrapper">
        <h1>Pages</h1>
        <p>This is the view which will render user pages</p>
		<PageAdd />
		<PageView />
      </div>
    </>
  );
};

export default PagesPage;
