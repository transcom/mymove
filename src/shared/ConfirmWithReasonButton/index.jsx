import React, { Component } from 'react';
import Alert from 'shared/Alert'; // eslint-disable-line

export default class ConfirmWithReasonButton extends Component {
  state = { displayState: 'BUTTON', cancelReason: '' };

  setConfirmState = () => {
    this.setState({ displayState: 'CONFIRM' });
  };

  setCancelState = () => {
    if (this.state.cancelReason !== '') {
      this.setState({ displayState: 'CANCEL' });
    }
  };

  setButtonState = () => {
    this.setState({ displayState: 'BUTTON', cancelReason: '' });
  };

  handleChange = event => {
    this.setState({ cancelReason: event.target.value });
  };

  cancel = event => {
    event.preventDefault();
    this.props.onConfirm(this.state.cancelReason);
  };

  render() {
    const {
      buttonTitle,
      reasonPrompt,
      warningPrompt,
      buttonDisabled,
    } = this.props;

    const reasonPresent = this.state.cancelReason !== '';

    if (this.state.displayState === 'CANCEL') {
      return (
        <div className="cancel-panel">
          <h2 className="extras usa-heading">{buttonTitle}</h2>
          <div className="extras content">
            <Alert type="warning" heading="Cancelation Warning">
              {warningPrompt}
            </Alert>
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>No, never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.cancel}>Yes, {buttonTitle}</button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'CONFIRM') {
      return (
        <div className="cancel-panel">
          <h2 className="extras usa-heading">{buttonTitle}</h2>
          <div className="extras content">
            {reasonPrompt}
            <textarea required onChange={this.handleChange} />
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>Never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.setCancelState} disabled={!reasonPresent}>
                  {buttonTitle}
                </button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'BUTTON') {
      return (
        <button
          className="usa-button-secondary"
          onClick={this.setConfirmState}
          disabled={buttonDisabled}
        >
          {buttonTitle}
        </button>
      );
    } else {
      console.error('Unknown State: ', this.state.displayState);
    }
  }
}
