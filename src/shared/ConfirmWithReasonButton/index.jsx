import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import Alert from 'shared/Alert'; // eslint-disable-line

export default class ConfirmWithReasonButton extends Component {
  state = { displayState: 'Button', cancelReason: '' };

  setConfirmState = () => {
    this.setState({ displayState: 'Confirm' });
  };

  setCancelState = () => {
    if (this.state.cancelReason !== '') {
      this.setState({ displayState: 'Cancel' });
    }
  };

  setButtonState = () => {
    this.setState({ displayState: 'Button' });
  };

  handleChange = event => {
    this.setState({ cancelReason: event.target.value });
  };

  cancel = event => {
    event.preventDefault();
    this.props.onConfirm(this.state.cancelReason);
    this.setState({ displayState: 'Redirect' });
  };

  render() {
    const { buttonTitle, reasonPrompt, warningPrompt } = this.props;

    const reasonPresent = this.state.cancelReason !== '';

    if (this.state.displayState === 'Cancel') {
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
    } else if (this.state.displayState === 'Confirm') {
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
    } else if (this.state.displayState === 'Button') {
      return (
        <button className="usa-button-secondary" onClick={this.setConfirmState}>
          {buttonTitle}
        </button>
      );
    } else if (this.state.displayState === 'Redirect') {
      return <Redirect to="/" />;
    }
  }
}
