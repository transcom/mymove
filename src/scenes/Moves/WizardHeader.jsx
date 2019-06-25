import React from 'react';
import PropTypes from 'prop-types';
import './WizardHeader.css';

const WizardHeader = ({ icon, right, title }) => (
  <div className="wizard-header">
    <div className="usa-grid">
      <div className="wizard-left">
        <img className="icon" src={icon} alt="" />
        <h3>{title}</h3>
      </div>
      <div className="wizard-right">{right}</div>
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
