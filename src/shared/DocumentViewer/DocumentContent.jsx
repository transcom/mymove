import React from 'react';
import PropTypes from 'prop-types';

import './index.css';

const DocumentContent = ({ url, filename, contentType, uploadId, rotate, orientation }) => {
  const imgHeight = document.querySelector('.page img').getBoundingClientRect().height;

  return (
    <div
      className="page"
      style={{ minHeight: imgHeight + 50, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}
    >
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
        <div style={{ marginTop: imgHeight / 5, marginBottom: imgHeight / 5 }}>
          <img src={url} style={{ transform: `rotate(${orientation}deg)` }} alt="document upload" />
        </div>
      )}

      <div>
        <button onClick={rotate.bind(this, uploadId, 'left')} data-direction="left">
          rotate left
        </button>
        <button onClick={rotate.bind(this, uploadId, 'right')} data-direction="right">
          rotate right
        </button>
      </div>
    </div>
  );
};

DocumentContent.propTypes = {
  contentType: PropTypes.string.isRequired,
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  uploadId: PropTypes.string.isRequired,
  orientation: PropTypes.string,
  rotate: PropTypes.func.isRequired,
};

export default DocumentContent;
