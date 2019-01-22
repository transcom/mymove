import React from 'react';
import PropTypes from 'prop-types';

import './index.css';

const DocumentContent = ({ contentType, filename, url }) => (
  <div className="page">
    {contentType === 'application/pdf' ? (
      <div className="pdf-placeholder">
        {filename && <span className="filename">{filename}</span>}
        This PDF can be{' '}
        <a target="_blank" href={url}>
          viewed here
        </a>
        .
      </div>
    ) : (
      <img src={url} width="100%" height="100%" alt="document upload" />
    )}
  </div>
);

DocumentContent.propTypes = {
  contentType: PropTypes.string.isRequired,
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
};

export default DocumentContent;
