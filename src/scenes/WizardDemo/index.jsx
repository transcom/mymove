import React from 'react';
import { Route } from 'react-router-dom';
import WizardPage from 'shared/WizardPage';
import { withRouter } from 'react-router-dom';
import intro from './intro.png';
import moveType from './select-move-type.png';
import dateSelection from './select-date.png';
import mover from './select-mover.png';
import review from './review-locations.png';

const C = props => (
  <WizardPage
    handleSubmit={() => undefined}
    pageList={props.pageList}
    pageKey={props.pageKey}
    history={props.history}
  >
    <img src={props.src} alt={props.path} />
  </WizardPage>
);

const ImagePage = withRouter(C);

export default () => {
  const pages = {
    '/mymove/intro': { src: intro },
    '/mymove/moveType': { src: moveType },
    '/mymove/dateSelection': { src: dateSelection },
    '/mymove/mover': { src: mover },
    '/mymove/review': { src: review },
  };
  const pageList = Object.keys(pages);

  return pageList.map(key => (
    <Route
      path={key}
      key={key}
      render={() => (
        <ImagePage pageList={pageList} pageKey={key} src={pages[key].src} />
      )}
    />
  ));
};
