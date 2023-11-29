import React from 'react';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import moment from 'moment';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { FilesShape } from '../types';

import styles from './Menu.module.scss';

import { filenameFromPath } from 'utils/formatters';
import { formatDate } from 'shared/dates';

const DocViewerMenu = ({ isOpen, files, handleClose, selectedFileIndex, handleSelectFile }) => (
  <div data-testid="DocViewerMenu" className={classnames(styles.docViewerMenu, { [styles.collapsed]: !isOpen })}>
    <div className={styles.menuHeader}>
      <h3>Documents</h3>
      <div className={styles.menuControls}>
        <Button data-testid="closeMenu" type="button" onClick={handleClose} unstyled className={styles.menuClose}>
          <FontAwesomeIcon icon="times" title="Close menu" aria-label="Close menu" />
        </Button>
      </div>
    </div>
    <ul className={styles.menuList}>
      {files.map((file, i) => {
        const fileNameClass = file.isWeightTicket ? styles.fileNameWithTag : styles.fileNameFullWidth;
        const itemClasses = classnames(styles.menuItemBtn, {
          [styles.active]: i === selectedFileIndex,
        });
        const fileName = filenameFromPath(file.filename);
        const fileDate = formatDate(moment(file.createdAt), 'DD-MMM-YYYY');

        return (
          <li key={fileName + i}>
            <div title={fileName}>
              <Button unstyled className={itemClasses} type="button" onClick={() => handleSelectFile(i)}>
                <p>
                  <span className={fileNameClass}>{fileName}</span> {file.isWeightTicket && <Tag>Weight Ticket</Tag>}
                </p>
                <p>
                  <span>Uploaded on {fileDate}</span>
                </p>
              </Button>
            </div>
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
