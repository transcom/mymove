import { Component } from 'react';
import React from 'react';
import PropTypes from 'prop-types';
import PopUp from 'shared/PopUp';
import CheckBox from 'shared/CheckBox';

class CustomerAgreement extends Component {
  handleAcceptTermsChange = acceptTerms => {
    this.props.onAcceptTermsChange(acceptTerms);
  };

  render() {
    return (
      <div className="customer-agreement">
        <p>
          <strong>Customer Agreement</strong>
        </p>
        <CheckBox onChangeHandler={this.handleAcceptTermsChange} checked={this.props.checked}>
          I agree to the
          <PopUp alertMessage={this.props.agreementText}> Legal Agreement / Privacy Act</PopUp>
        </CheckBox>
      </div>
    );
  }
}

CustomerAgreement.propTypes = {
  onAcceptTermsChange: PropTypes.func.isRequired,
  checked: PropTypes.bool.isRequired,
  agreementText: PropTypes.string.isRequired,
};

export default CustomerAgreement;
