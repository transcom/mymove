import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import Summary from './Summary';

export default class Review extends Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  render() {
    const { pages, pageKey } = this.props;

    return (
      <WizardPage handleSubmit={no_op} pageList={pages} pageKey={pageKey} pageIsValid={true}>
        <h1>Review</h1>
        <p>You're almost done! Please review your details before we finalize the move.</p>
        <Summary />
      </WizardPage>
    );
  }
}
