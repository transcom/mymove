import React, { useState } from 'react';
import { Button } from '@trussworks/react-uswds';

import { FilesShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content/Content';
import Menu from './Menu/Menu';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';

/**
 * TODO
 * - implement open in a new window
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - handle fetch doc errors
 */

const DocumentViewer = ({ files }) => {
  // TODO - show msg if there are no files
  const [selectedFileIndex, selectFile] = useState(0);
  const [menuIsOpen, setMenuOpen] = useState(false);

  const selectedFile = files[parseInt(selectedFileIndex, 10)];

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

  return (
    <div className={styles.DocumentViewer}>
      <div className={styles.titleBar}>
        <Button data-testid="openMenu" type="button" onClick={openMenu} unstyled>
          <DocMenu />
        </Button>

        <p>{selectedFile.filename}</p>
        {/* TODO */}
        <Button type="button" unstyled onClick={openInNewWindow}>
          <span>Open in a new window</span>
          <ExternalLink />
        </Button>
      </div>
      <Content fileType={selectedFile.fileType} filePath={selectedFile.filePath} />
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
