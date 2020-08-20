import React from 'react';
import { string, arrayOf, shape } from 'prop-types';

import styles from './Home.module.scss';

import { ReactComponent as DocsIcon } from 'shared/icon/documents.svg';

const DocsUploaded = ({ files }) => (
  <div className={`${styles['doc-list-container']} padding-left-2 padding-right-2`}>
    <h6 className="margin-top-2 margin-bottom-2">
      {files.length} FILE{files.length > 1 ? 'S' : ''} UPLOADED
    </h6>
    {files.map((file) => (
      <div key={file.filename} className={`margin-bottom-2 ${styles['doc-list-item']}`}>
        <DocsIcon className={styles['docs-icon']} />
        {file.filename}
      </div>
    ))}
  </div>
);

DocsUploaded.propTypes = {
  files: arrayOf(shape({ filename: string.isRequired })).isRequired,
};

export default DocsUploaded;
