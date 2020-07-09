import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { FilesShape } from '../types';

import styles from './Menu.module.scss';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';

const DocViewerMenu = ({ isOpen, files, handleClose, selectedFileIndex, handleSelectFile }) => (
  <div data-testid="DocViewerMenu" className={classnames(styles.docViewerMenu, { [styles.collapsed]: !isOpen })}>
    <div className={styles.menuHeader}>
      <h3>Documents</h3>
      <div className={styles.menuControls}>
        <Button data-testid="closeMenu" type="button" onClick={handleClose} unstyled className={styles.menuClose}>
          <XLightIcon />
        </Button>
      </div>
    </div>
    <ul className={styles.menuList}>
      {files.map((file, i) => {
        const itemClasses = classnames(styles.menuItemBtn, {
          [styles.active]: i === selectedFileIndex,
        });
        return (
          // eslint-disable-next-line react/no-array-index-key
          <li key={`menu_file_${i}`}>
            <Button unstyled className={itemClasses} type="button" onClick={() => handleSelectFile(i)}>
              <p>{file.filename}</p>
            </Button>
          </li>
        );
      })}
    </ul>
  </div>
);

DocViewerMenu.propTypes = {
  isOpen: PropTypes.bool,
  files: FilesShape.isRequired,
  handleClose: PropTypes.func.isRequired,
  selectedFileIndex: PropTypes.number,
  handleSelectFile: PropTypes.func.isRequired,
};

DocViewerMenu.defaultProps = {
  isOpen: false,
  selectedFileIndex: null,
};

export default DocViewerMenu;
