import React from 'react';
import PropTypes from 'prop-types';
import FileViewer from 'react-file-viewer';
import { Button } from '@trussworks/react-uswds';

import styles from './DocumentViewer.module.scss';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';

/**
 * TODO
 * - implement open in a new window
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - fix styling of controls bar (need to modify react-file-viewer)
 * - support images in addition to PDFs
 * - menu bar for browsing multiple documents
 * - handle fetch doc errors
 */

const DocumentViewer = ({ filename, fileType, filePath }) => {
  const onError = () => {
    // console.log('file viewer error', e);
  };

  return (
    <div className={styles.DocumentViewer}>
      <div className={styles.titleBar}>
        <p>{filename}</p>
        <Button unstyled>
          Open in a new window
          <ExternalLink />
        </Button>
      </div>
      <FileViewer fileType={fileType} filePath={filePath} onError={onError} />
    </div>
  );
};

DocumentViewer.propTypes = {
  filename: PropTypes.node.isRequired,
  filePath: PropTypes.string.isRequired,
  fileType: PropTypes.string.isRequired,
};

export default DocumentViewer;
