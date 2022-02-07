import React, { useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import moment from 'moment';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { FilesShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content/Content';
import Menu from './Menu/Menu';

import { formatDate } from 'shared/dates';
import { filenameFromPath } from 'shared/formatters';

/**
 * TODO
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - handle fetch doc errors
 */

const DocumentViewer = ({ files }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [menuIsOpen, setMenuOpen] = useState(false);

  const sortedFiles = files.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));

  const selectedFile = sortedFiles[parseInt(selectedFileIndex, 10)];

  if (!selectedFile) {
    return <h2>File Not Found</h2>;
  }

  let fileType = selectedFile.contentType;
  switch (selectedFile.contentType) {
    case 'application/pdf': {
      fileType = 'pdf';
      break;
    }
    case 'image/png': {
      fileType = 'png';
      break;
    }
    case 'image/jpeg': {
      fileType = 'jpg';
      break;
    }
    case 'image/gif': {
      fileType = 'gif';
      break;
    }
    default: {
      break;
    }
  }

  const openMenu = () => {
    setMenuOpen(true);
  };
  const closeMenu = () => {
    setMenuOpen(false);
  };

  const handleSelectFile = (index) => {
    selectFile(index);
    closeMenu();
  };

  const selectedFilename = filenameFromPath(selectedFile.filename);

  const selectedFileDate = formatDate(moment(selectedFile.createdAt), 'DD MMM YYYY');

  return (
    <div className={styles.DocumentViewer}>
      <div className={styles.titleBar}>
        <Button data-testid="openMenu" type="button" onClick={openMenu} aria-label="Open menu" unstyled>
          <FontAwesomeIcon icon="th-list" />
        </Button>
        <p title={selectedFilename} data-testid="documentTitle">
          <span>{selectedFilename}</span> <span>- Added on {selectedFileDate}</span>
        </p>
      </div>
      <Content fileType={fileType} filePath={selectedFile.url} />
      {menuIsOpen && <div className={styles.overlay} />}
      <Menu
        isOpen={menuIsOpen}
        files={sortedFiles}
        handleClose={closeMenu}
        selectedFileIndex={selectedFileIndex}
        handleSelectFile={handleSelectFile}
      />
    </div>
  );
};

DocumentViewer.propTypes = {
  files: FilesShape,
};

DocumentViewer.defaultProps = {
  files: [],
};

export default DocumentViewer;
