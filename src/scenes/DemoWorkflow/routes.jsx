import React from 'react';
import { Route } from 'react-router-dom';
import DemoWorkflow from '.';
//from https://github.com/ReactTraining/react-router/issues/4105
const renderMergedProps = (component, ...rest) => {
  const finalProps = Object.assign({}, ...rest);
  return React.createElement(component, finalProps);
};

const PropsRoute = ({ component, ...rest }) => {
  return (
    <Route
      {...rest}
      render={routeProps => {
        return renderMergedProps(component, routeProps, rest);
      }}
    />
  );
};
export default () => {
  const pages = {
    '/demoWorkflow1': { subsetOfUiSchema: ['service_member_information'] },
    '/demoWorkflow2': { subsetOfUiSchema: ['orders_information'] },
  };
  const pageList = Object.keys(pages);

  return pageList.map(key => (
    <PropsRoute
      key={key}
      path={key}
      component={DemoWorkflow}
      subsetOfUiSchema={pages[key].subsetOfUiSchema}
      pageList={pageList}
    />
  ));
};
