import React from 'react';
import PropTypes from 'prop-types';

const CustomerAgreement = ({ onChange, link, checked, agreementText, className }) => {
  const handleOnChange = e => {
    onChange(e.target.checked);
  };

  const handleClick = e => {
    // Prevent this from checking the box after opening the alert.
    e.preventDefault();
    alert(agreementText);
  };

  return (
    <div className={className || 'customer-agreement'}>
      <p>
        <strong>Customer Agreement</strong>
      </p>
      <input id="agree-checkbox" type="checkbox" checked={checked} onChange={handleOnChange} />
      <label htmlFor="agree-checkbox">
        I agree to the{' '}
        <a href={link} onClick={link ? null : handleClick}>
          {' '}
          Legal Agreement / Privacy Act
        </a>
      </label>
    </div>
  );
};

CustomerAgreement.propTypes = {
  onChange: PropTypes.func,
  checked: PropTypes.bool.isRequired,
  agreementText: PropTypes.string.isRequired,
  link: PropTypes.string,
  className: PropTypes.string,
};

export default CustomerAgreement;
