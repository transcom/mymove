import React from 'react';

import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';
import Agreement from 'scenes/Legalese';

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <h1>Placeholder for {props.title}</h1>
    </WizardPage>
  );
};

const stub = (key, pages) => () => (
  <Placeholder pageList={pages} pageKey={key} title={key} />
);

export default () => {
  const pages = {
    '/moves/:moveId': { render: stub },
    '/moves/:moveId/ppm-transition': { render: stub },
    '/moves/:moveId/ppm-size': { render: stub },
    '/moves/:moveId/ppm-incentive': { render: stub },
    '/moves/:moveId/agreement': {
      render: (key, pages) => () => (
        <WizardPage
          handleSubmit={() => undefined}
          pageList={pages}
          pageKey={key}
        >
          <Agreement />
        </WizardPage>
      ),
    },
  };
  const pageList = Object.keys(pages);
  return pageList.map(key => {
    const render = pages[key].render(key, pageList);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};
