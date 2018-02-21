import React from 'react';
import { Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import Feedback from 'scenes/Feedback';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header';
import { history } from 'shared/store';
import Shipments from 'scenes/Shipments';
import Footer from 'shared/Footer';
import DD1299 from 'scenes/DD1299';
import Landing from 'scenes/Landing';
import WizardDemo from 'scenes/WizardDemo';
import DemoWorkflow from 'scenes/DemoWorkflow';

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
const wizardPages = ['/demoWorkflow1', '/demoWorkflow2'];
const AppWrapper = () => (
  <ConnectedRouter history={history}>
    <div className="App site">
      <Header />
      <main className="site__content">
        <Route exact path="/" component={Feedback} />
        <Route path="/submitted" component={SubmittedFeedback} />
        <Route path="/shipments/:shipmentsStatus" component={Shipments} />
        <Route path="/DD1299" component={DD1299} />
        <Route path="/landing" component={Landing} />
        <Route path="/wizardDemo" component={WizardDemo} />
        <PropsRoute
          path="/demoWorkflow1"
          component={DemoWorkflow}
          subsetOfUiSchema={['service_member_information']}
          pageList={wizardPages}
        />
        <PropsRoute
          path="/demoWorkflow2"
          component={DemoWorkflow}
          subsetOfUiSchema={['orders_information']}
          pageList={wizardPages}
        />
      </main>
      <Footer />
    </div>
  </ConnectedRouter>
);

export default AppWrapper;
