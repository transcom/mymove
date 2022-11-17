import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom-old';

const CustomerAgreement = ({ onChange, link, checked, agreementText, className }) => {
  const handleOnChange = (e) => {
    onChange(e.target.checked);
  };

  const handleClick = (e) => {
    // Prevent this from checking the box after opening the alert.
    e.preventDefault();
    alert(agreementText);
  };

  return (
    <div className={className || 'customer-agreement'}>
      <p>
        <strong>Customer Agreement</strong>
      </p>
      <div className="usa-checkbox">
        <input
          id="agree-checkbox"
          type="checkbox"
          checked={checked}
          onChange={handleOnChange}
          className="usa-checkbox__input"
        />
        <label htmlFor="agree-checkbox" className="usa-checkbox__label">
          I agree to the{' '}
          {link ? (
            <Link to={link} className="usa-link">
              {' '}
              Legal Agreement / Privacy Act
            </Link>
          ) : (
            <a onClick={handleClick} className="usa-link">
              {' '}
              Legal Agreement / Privacy Act
            </a>
          )}
        </label>
      </div>
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
