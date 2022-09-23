import React from 'react';
import PropTypes from 'prop-types';

const NotFound = ({ handleOnClick }) => {
  return (
    <div className="usa-grid">
      <div className="grid-container usa-prose">
        <h1>Page not found</h1>
        <p>Looks like you&apos;ve followed a broken link or entered a URL that doesn&apos;t exist on this site.</p>
        <button type="button" className="usa-button" onClick={handleOnClick}>
          Go Back
        </button>
      </div>
    </div>
  );
};

NotFound.propTypes = {
  handleOnClick: PropTypes.func.isRequired,
};

export default NotFound;
