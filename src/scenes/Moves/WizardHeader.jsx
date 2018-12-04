import React from 'react';
import PropTypes from 'prop-types';
import './WizardHeader.css';

const WizardHeader = ({ icon, right, title }) => (
  <div className="wizard-header">
    <div className="usa-grid">
      <div className="usa-width-one-half">
        <img className="icon" src={icon} alt="" />
        <p>{title}</p>
      </div>
      <div className="usa-width-one-half" style={{ textAlign: 'right' }}>
        {right}
      </div>
    </div>
    <div className="usa-grid">
      <div className="usa-width-one-whole">
        <hr />
      </div>
    </div>
  </div>
);

WizardHeader.defaultProps = {
  title: <span>&nbsp;</span>,
};

WizardHeader.propTypes = {
  icon: PropTypes.string,
  title: PropTypes.oneOfType([PropTypes.string, PropTypes.element]),
  right: PropTypes.element,
};

export default WizardHeader;
