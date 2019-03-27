import React from 'react';
import PropTypes from 'prop-types';

const CustomerAgreement = ({ onChange, checked, agreementText }) => {
  const handleOnChange = e => {
    onChange(e.target.checked);
  };

  const handleClick = e => {
    // Prevent this from checking the box after opening the alert.
    e.preventDefault();
    alert(agreementText);
  };

  return (
    <div className="customer-agreement">
      <p>
        <strong>Customer Agreement</strong>
      </p>
      <input id="agree-checkbox" type="checkbox" checked={checked} onChange={handleOnChange} />
      <label htmlFor="agree-checkbox">
        I agree to the <a onClick={handleClick}> Legal Agreement / Privacy Act</a>
      </label>
    </div>
  );
};

CustomerAgreement.propTypes = {
  onChange: PropTypes.func,
  checked: PropTypes.bool.isRequired,
  agreementText: PropTypes.string.isRequired,
};

export default CustomerAgreement;
