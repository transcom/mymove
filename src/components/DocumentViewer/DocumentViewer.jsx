import React, { useState } from 'react';
import { Button } from '@trussworks/react-uswds';

import { FilesShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content';
import Menu from './Menu';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';

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

const DocumentViewer = ({ files }) => {
  // TODO - show msg if there are no files
  const [selectedFileIndex, selectFile] = useState(0);
  const [menuIsOpen, setMenuOpen] = useState(false);

  const selectedFile = files[parseInt(selectedFileIndex, 10)];

  const openMenu = () => setMenuOpen(true);
  const closeMenu = () => setMenuOpen(false);
  const handleSelectFile = (index) => {
    selectFile(index);
    closeMenu();
  };

  return (
    <div className={styles.DocumentViewer}>
      <div className={styles.titleBar}>
        <Button type="button" onClick={openMenu} unstyled>
          <DocMenu />
        </Button>

        <p>{selectedFile.filename}</p>
        {/* TODO */}
        <Button unstyled>
          <span>Open in a new window</span>
          <ExternalLink />
        </Button>
      </div>
      <Content fileType={selectedFile.fileType} filePath={selectedFile.filePath} />
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
