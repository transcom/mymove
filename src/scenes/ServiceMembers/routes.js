import React from 'react';
import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';

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

const stub = (key, pages, component) => ({ match }) => {
  if (component) {
    const pageComponent = React.createElement(component, { match }, null);
    return (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        {pageComponent}
      </WizardPage>
    );
  } else {
    return <Placeholder pageList={pages} pageKey={key} title={key} />;
  }
};

export default () => {
  const pages = {
    '/service-member/:id/create': { render: stub },
    '/service-member/:id/name': { render: stub },
    '/service-member/:id/contact-info': { render: stub },
    '/service-member/:id/duty-station': { render: stub },
    '/service-member/:id/residence-address': { render: stub },
    '/service-member/:id/backup-mailing-address': { render: stub },
    '/service-member/:id/backup-contacts': { render: stub },
    '/service-member/:id/transition': { render: stub },
  };
  const pageList = Object.keys(pages);
  const componentMap = {};
  return pageList.map(key => {
    const step = key.split('/').pop();
    var component = componentMap[step];
    const render = pages[key].render(key, pageList, component);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};
