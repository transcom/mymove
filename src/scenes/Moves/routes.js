import React from 'react';
import { Route } from 'react-router-dom';
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
const stub = (key, pages) => {
  return () => <Placeholder pageList={pages} pageKey={key} title={key} />;
};
export default () => {
  const pages = {
    '/moves/:moveId': { render: stub },
    '/moves/:moveId/ppm-transition': { render: stub },
    '/moves/:moveId/ppm-size': { render: stub },
    '/moves/:moveId/ppm-incentive': { render: stub },
    '/moves/:moveId/agreement': {
      render: (key, pages) => {
        return () => (
          <WizardPage
            handleSubmit={() => undefined}
            pageList={pages}
            pageKey={key}
          >
            <Agreement />
          </WizardPage>
        );
      },
    },
  };
  const pageList = Object.keys(pages);
  const val = pageList.map(key => {
    const render = pages[key].render(key, pageList);
    return <Route exact path={key} key={key} render={render} />;
  });
  return val;
};
