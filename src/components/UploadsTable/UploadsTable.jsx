import React from 'react';
import classNames from 'classnames';
import PropTypes from 'prop-types';
import bytes from 'bytes';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './UploadsTable.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';

const UploadsTable = ({ uploads, onDelete }) => {
  const getIcon = (fileType) => {
    if (fileType === 'application/pdf') {
      return 'file-pdf';
    }
    return 'file-image';
  };

  return (
    <SectionWrapper className={classNames(styles.wrapper)}>
      <h6>{uploads.length} Files Uploaded</h6>
      <ul>
        {uploads.map((upload) => (
          <li className={classNames(styles.uploadListItem)} key={upload.id}>
            <div className={classNames(styles.fileInfoContainer)}>
              <FontAwesomeIcon size="lg" icon={getIcon(upload.content_type)} className={classNames(styles.faIcon)} />
              <div className={classNames(styles.fileInfo)}>
                <p>{upload.filename}</p>
                <p className={classNames(styles.fileSizeAndTime)}>
                  <span className={classNames(styles.uploadFileSize)}>{bytes(upload.bytes)}</span>
                  <span>Uploaded {moment(upload.created_at).format('DD MMM YYYY h:mm A')}</span>
                </p>
              </div>
            </div>
            <Button type="button" unstyled onClick={() => onDelete(upload.id)}>
              Delete
            </Button>
          </li>
        ))}
      </ul>
    </SectionWrapper>
  );
};

UploadsTable.propTypes = {
  uploads: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      created_at: PropTypes.string.isRequired,
      bytes: PropTypes.number.isRequired,
      url: PropTypes.string.isRequired,
      filename: PropTypes.string.isRequired,
    }),
  ).isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default UploadsTable;
