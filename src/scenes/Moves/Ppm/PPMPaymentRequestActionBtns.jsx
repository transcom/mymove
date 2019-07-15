import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
// import { get} from 'lodash';

import './PPMPaymentRequest.css';
import AlertWithConfirmation from 'shared/AlertWithConfirmation';

class PPMPaymentRequestActionBtns extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasConfirmation: props.hasConfirmation,
      displayConfirmation: false,
    };

    this.showConfirmationOrFinishLater = formValues => {
      const { history, hasConfirmation } = this.props;

      if (!hasConfirmation) {
        return history.push('/');
      }

      this.setState({ displayConfirmation: true });
      return;
    };

    this.cancelConfirmationHandler = () => {
      this.setState({ displayConfirmation: false });
      return;
    };

    this.confirmFinishLater = () => {
      return this.props.history.push('/');
    };
  }

  render() {
    const {
      nextBtnLabel,
      displaySkip,
      skipHandler,
      saveAndAddHandler,
      hasConfirmation,
      submitButtonsAreDisabled,
      submitting,
    } = this.props;
    return (
      <div className="usa-width-one-whole">
        {hasConfirmation &&
          this.state.displayConfirmation && (
            <div className="ppm-payment-request-footer">
              <AlertWithConfirmation
                hasConfirmation={hasConfirmation}
                type="warning"
                cancelActionHandler={this.cancelConfirmationHandler}
                okActionHandler={this.confirmFinishLater}
                message="Go back to the home screen without saving current screen."
              />
            </div>
          )}

        {!this.state.displayConfirmation && (
          <div className="ppm-payment-request-footer">
            <button type="button" className="usa-button-secondary" onClick={this.showConfirmationOrFinishLater}>
              Finish Later
            </button>
            <div className="usa-width-one-thirds">
              {displaySkip && (
                <button data-cy="skip" type="button" className="usa-button-secondary" onClick={skipHandler}>
                  Skip
                </button>
              )}
              <button type="button" onClick={saveAndAddHandler} disabled={submitButtonsAreDisabled || submitting}>
                {nextBtnLabel}
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }
}

export default withRouter(PPMPaymentRequestActionBtns);
