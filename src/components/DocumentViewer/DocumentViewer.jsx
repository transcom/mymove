import React, { useEffect, useState } from 'react';
import { bool } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import moment from 'moment';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { FilesShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content/Content';
import Menu from './Menu/Menu';

import { formatDate } from 'shared/dates';
import { filenameFromPath } from 'utils/formatters';

/**
 * TODO
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - handle fetch doc errors
 */

const DocumentViewer = ({ files, allowDownload }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [menuIsOpen, setMenuOpen] = useState(false);
  const selectedFile = files[parseInt(selectedFileIndex, 10)];

  useEffect(() => {
    selectFile(0);
  }, [files]);

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
        <p title={selectedFilename} className={styles.documentTitle} data-testid="documentTitle">
          <span>{selectedFilename}</span> <span>- Added on {selectedFileDate}</span>
        </p>
        {allowDownload && (
          <p className={styles.downloadLink}>
            <a href={selectedFile.url} download tabIndex={menuIsOpen ? '-1' : '0'}>
              <span>Download file</span> <FontAwesomeIcon icon="download" />
            </a>
          </p>
        )}
      </div>
      <Content fileType={fileType} filePath={selectedFile.url} />
      {menuIsOpen && <div className={styles.overlay} />}
      <Menu
        isOpen={menuIsOpen}
        files={files}
        handleClose={closeMenu}
        selectedFileIndex={selectedFileIndex}
        handleSelectFile={handleSelectFile}
      />
    </div>
  );
};

DocumentViewer.propTypes = {
  files: FilesShape,
  allowDownload: bool,
};

DocumentViewer.defaultProps = {
  files: [],
  allowDownload: false,
};

export default DocumentViewer;
