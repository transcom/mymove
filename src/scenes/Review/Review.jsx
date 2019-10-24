import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import Summary from './Summary';
import { connect } from 'react-redux';
import scrollToTop from 'shared/scrollToTop';

import './Review.css';

class Review extends Component {
  componentDidMount() {
    scrollToTop();
  }
  render() {
    const { pages, pageKey } = this.props;

    return (
      <div>
        <WizardPage handleSubmit={no_op} pageList={pages} pageKey={pageKey} pageIsValid={true}>
          <div className="grid-row">
            <div className="grid-col-12 edit-title">
              <h2>Review Move Details</h2>
              <p>You're almost done! Please review your details before we finalize the move.</p>
            </div>
          </div>
          <Summary />
        </WizardPage>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => ({
  ...ownProps,
});

export default connect(mapStateToProps)(Review);
