import React from 'react';
import { string, arrayOf, shape } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFile } from '@fortawesome/free-solid-svg-icons';

import styles from './DocsUploaded.module.scss';

const DocsUploaded = ({ files }) => (
  <div className={styles['doc-list-container']} data-testid="doc-list-container">
    <h6 className={styles['doc-list-header']}>
      {files.length} File{files.length > 1 ? 's' : ''} uploaded
    </h6>
    {files.map((file) => (
      <div key={file.filename} className={styles['doc-list-item']}>
        <FontAwesomeIcon icon={faFile} className={styles['docs-icon']} />
        {file.filename}
      </div>
    ))}
  </div>
);

DocsUploaded.propTypes = {
  files: arrayOf(shape({ filename: string.isRequired })).isRequired,
};

export default DocsUploaded;
