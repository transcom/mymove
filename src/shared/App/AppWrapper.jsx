import React from 'react';
import { Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import Feedback from 'scenes/Feedback/Feedback';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header/Header';
import { history } from 'shared/store';
import Footer from 'shared/Footer/Footer';
import DD1299 from 'scenes/DD1299';

const AppWrapper = () => (
  <ConnectedRouter history={history}>
    <div className="App site">
      <Header />
      <main className="site__content">
        <Route exact path="/" component={Feedback} />
        <Route path="/submitted" component={SubmittedFeedback} />
        <Route path="/DD1299" component={DD1299} />
      </main>
      <Footer />
    </div>
  </ConnectedRouter>
);

export default AppWrapper;
