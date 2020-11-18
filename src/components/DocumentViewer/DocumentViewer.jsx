import React, { useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faThList } from '@fortawesome/free-solid-svg-icons';

import { FilesShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content/Content';
import Menu from './Menu/Menu';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';
import { filenameFromPath } from 'shared/formatters';

/**
 * TODO
 * - implement open in a new window
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - handle fetch doc errors
 */

const DocumentViewer = ({ files }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [menuIsOpen, setMenuOpen] = useState(false);

  const selectedFile = files[parseInt(selectedFileIndex, 10)];

  if (!selectedFile) {
    return (
      <>
        <h2>File Not Found</h2>
      </>
    );
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
    // eslint-disable-next-line no-empty
    default: {
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
  const openInNewWindow = () => {
    // TODO - do we need to stream the file or can we just open the URL?
  };

  const selectedFilename = filenameFromPath(selectedFile.filename);

  return (
    <div className={styles.DocumentViewer}>
      <div className={styles.titleBar}>
        <Button data-testid="openMenu" type="button" onClick={openMenu} unstyled>
          <FontAwesomeIcon icon={faThList} />
        </Button>

        <p title={selectedFilename}>{selectedFilename}</p>
        {/* TODO */}
        <Button type="button" unstyled onClick={openInNewWindow}>
          <span>Open in a new window</span>
          <ExternalLink />
        </Button>
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
};

DocumentViewer.defaultProps = {
  files: [],
};

export default DocumentViewer;
