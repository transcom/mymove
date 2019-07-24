import React from 'react';
import PropTypes from 'prop-types';

import { detectFirefox } from 'shared/utils';

import './index.css';
import styles from './DocumentContent.module.scss';

const downloadOnlyView = (filename, url) => (
  <div className="pdf-placeholder">
    {filename && <span className="filename">{filename}</span>}
    This PDF can be{' '}
    <a target="_blank" rel="noopener noreferrer" href={url}>
      viewed here
    </a>
    .
  </div>
);

const DocumentContent = ({ contentType, filename, url }) => (
  <div className="page">
    {contentType === 'application/pdf' ? (
      detectFirefox() ? (
        downloadOnlyView(filename, url)
      ) : (
        <object className={styles.pdf} data={url} type="application/pdf" alt="document upload">
          {downloadOnlyView(filename, url)}
        </object>
      )
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
