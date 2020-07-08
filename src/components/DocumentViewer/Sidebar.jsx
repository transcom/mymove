import React from 'react';
import PropTypes from 'prop-types';

const DocViewerSidebar = ({ header, body, footer }) => (
  <div data-testid="DocViewerSidebar">
    <div>
      <h3>{header}</h3>
    </div>
    <div>{body}</div>
    <div>{footer}</div>
  </div>
);

DocViewerSidebar.propTypes = {
  header: PropTypes.string.isRequired,
  body: PropTypes.element.isRequired,
  footer: PropTypes.element.isRequired,
};

export default DocViewerSidebar;
