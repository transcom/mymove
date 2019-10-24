import React from 'react';
import PropTypes from 'prop-types';
import './WizardHeader.css';

const WizardHeader = ({ icon, right, title }) => (
  <div className="wizard-header">
    <div className="grid-row grid-gap">
      <div className="tablet:grid-col-6 desktop:grid-col-8 wizard-left">
        {icon && <img className="icon" src={icon} alt="" />}
        <h1>{title}</h1>
      </div>
      <div className="tablet:grid-col-6 desktop:grid-col-4 wizard-right">{right}</div>
    </div>
    <div className="grid-row">
      <div className="grid-row-12">
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
