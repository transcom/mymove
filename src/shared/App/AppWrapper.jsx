import React from 'react';
import { Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import Feedback from 'scenes/Feedback/Feedback';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header/Header';
import history from 'shared/store';
import Footer from 'shared/Footer/Footer';

const AppWrapper = () => (
  <ConnectedRouter history={history}>
    <div className="App site">
      <Header />
      <main className="site__content">
        <Route exact path="/" component={Feedback} />
        <Route path="/submitted" component={SubmittedFeedback} />
      </main>
      <Footer />
    </div>
  </ConnectedRouter>
);

export default AppWrapper;
