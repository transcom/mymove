import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

class SSNField extends Component {
  constructor(props) {
    super(props);

    this.state = {
      focused: false,
    };

    this.localOnBlur = this.localOnBlur.bind(this);
    this.localOnFocus = this.localOnFocus.bind(this);
  }

  localOnBlur(value) {
    const { input } = this.props;
    this.setState({ focused: false });
    input.onBlur(value);
  }

  localOnFocus(value) {
    const { input } = this.props;
    this.setState({ focused: true });
    input.onFocus(value);
  }

  render() {
    const {
      input,
      meta: { touched, error },
      ssnOnServer,
    } = this.props;

    const { name, value } = input;
    const { focused } = this.state;

    let displayedValue = value;
    if (!focused && (value !== '' || ssnOnServer)) {
      displayedValue = '•••-••-••••';
    }
    const displayError = touched && error;

    // This is copied from JsonSchemaField to match the styling
    return (
      <div className={classNames('usa-form-group', { 'usa-form-group--error': displayError })}>
        <label className={classNames('usa-label', { 'usa-label--error': displayError })} htmlFor="ssnInput">
          Social Security number
        </label>
        {touched && error && (
          <span className="usa-error-message" id={`${name}-error`} role="alert">
            {error}
          </span>
        )}
        <input
          id="ssnInput"
          type="text"
          // eslint-disable-next-line react/jsx-props-no-spreading
          {...input}
          className="usa-input"
          onFocus={this.localOnFocus}
          onBlur={this.localOnBlur}
          value={displayedValue}
        />
      </div>
    );
  }
}

SSNField.propTypes = {
  input: PropTypes.shape({
    name: PropTypes.string.isRequired,
    value: PropTypes.string,
    onFocus: PropTypes.func.isRequired,
    onBlur: PropTypes.func.isRequired,
  }).isRequired,
  meta: PropTypes.shape({
    touched: PropTypes.bool,
    error: PropTypes.string,
  }).isRequired,
  ssnOnServer: PropTypes.bool.isRequired,
};

export default SSNField;
