import React from 'react';
import { Route } from 'react-router-dom';
import Feedback from 'scenes/Feedback/Feedback';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header/Header';
import Footer from 'shared/Footer/Footer';

const AppWrapper = () => (
  <div className="App site">
    <Header />
    <main className="site__content">
      <Route exact path="/" component={Feedback} />
      <Route path="/submitted" component={SubmittedFeedback} />
    </main>
    <Footer />
  </div>
);

export default AppWrapper;
