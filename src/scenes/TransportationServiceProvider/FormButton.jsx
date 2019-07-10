import React, { Component } from 'react';
import PropTypes from 'prop-types';

export default class FormButton extends Component {
  state = {
    displayState: 'BUTTON',
  };

  setButtonState = () => {
    this.setState({ displayState: 'BUTTON' });
  };

  setFormState = () => {
    this.setState({ displayState: 'FORM' });
  };

  buttonView = () => {
    const { buttonTitle } = this.props;
    return <button onClick={this.setFormState}>{buttonTitle}</button>;
  };

  enterFormView = () => {
    const { FormComponent } = this.props;

    const formProps = {
      onCancel: this.setButtonState,
      ...this.props,
    };
    return <FormComponent {...formProps} />;
  };

  render() {
    const viewStates = {
      BUTTON: this.buttonView(),
      FORM: this.enterFormView(),
    };
    return viewStates[this.state.displayState];
  }
}

FormButton.propTypes = {
  FormComponent: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
  buttonTitle: PropTypes.string.isRequired,
};
