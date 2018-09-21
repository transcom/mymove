import React, { Component } from 'react';

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
    const { schema, onSubmit, formComponent } = this.props;

    const formProps = {
      onCancel: this.setButtonState,
      schema: schema,
      onSubmit: onSubmit,
    };
    return React.createElement(formComponent, formProps, null);
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
  formComponent: PropTypes.element.isRequired,
  schema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
  buttonTitle: PropTypes.string.isRequired,
};
