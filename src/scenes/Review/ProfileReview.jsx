import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import Summary from './Summary';

export default class ProfileReview extends Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  render() {
    const { pages, pageKey } = this.props;
    console.log(pages, pageKey);
    console.log(this.props);

    return (
      <WizardPage
        handleSubmit={no_op}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={true}
      >
        <h1>Profile Review</h1>
        <p>
          Has anything changed since your last move? Please check your info
          below, especially your Rank.
        </p>
        <Summary />
      </WizardPage>
    );
  }
}
