import React, { Component } from 'react';
// import { get} from 'lodash';
import './PPMPaymentRequest.css';
import AlertWithConfirmation from 'shared/AlertWithConfirmation';
import withRouter from 'utils/routing';

class PPMPaymentRequestActionBtns extends Component {
  state = {
    hasConfirmation: this.props.hasConfirmation,
    displayConfirmation: false,
  };

  showConfirmationOrFinishLater = (formValues) => {
    const {
      router: { navigate },
      hasConfirmation,
    } = this.props;

    if (!hasConfirmation) {
      return navigate('/ppm');
    }

    this.setState({ displayConfirmation: true });
    return undefined;
  };

  cancelConfirmationHandler = () => {
    this.setState({ displayConfirmation: false });
  };

  confirmFinishLater = () => {
    const {
      router: { navigate },
    } = this.props;
    navigate('/ppm');
  };

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
      <div className="grid-row">
        <div className="grid-col-12">
          {hasConfirmation && this.state.displayConfirmation && (
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
            <div className="ppm-payment-request-footer align-right">
              <button
                type="button"
                className="usa-button usa-button--secondary"
                onClick={this.showConfirmationOrFinishLater}
              >
                Finish Later
              </button>
              {displaySkip && (
                <button
                  data-testid="skip"
                  type="button"
                  className="usa-button usa-button--secondary"
                  onClick={skipHandler}
                >
                  Skip
                </button>
              )}
              <button
                type="button"
                className="usa-button"
                onClick={saveAndAddHandler}
                disabled={submitButtonsAreDisabled || submitting}
              >
                {nextBtnLabel}
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }
}

export default withRouter(PPMPaymentRequestActionBtns);
