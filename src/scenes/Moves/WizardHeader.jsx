import React from 'react';
import PropTypes from 'prop-types';

const WizardHeader = ({ right, title }) => (
  <div>
    <div className="usa-grid">
      <div className="usa-width-one-half">
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

WizardHeader.propTypes = {
  title: PropTypes.string,
  right: PropTypes.element,
};

export default WizardHeader;
