import React from 'react';
import PropTypes from 'prop-types';
import bytes from 'bytes';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './UploadsTable.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { ExistingUploadsShape } from 'types/uploads';

const UploadsTable = ({ className, uploads, onDelete, showDeleteButton, showDownloadLink = false }) => {
  const getIcon = (fileType) => {
    switch (fileType) {
      case 'application/pdf':
        return 'file-pdf';
      case 'application/vnd.ms-excel':
      case 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet':
        return 'file-excel';
      case 'image/png':
      case 'image/jpeg':
        return 'file-image';
      default:
        return 'file';
    }
  };

  return (
    uploads?.length > 0 && (
      <SectionWrapper className={classnames(styles.uploadsTableContainer, className)} data-testid="uploads-table">
        <h6>{uploads.length} Files Uploaded</h6>
        <ul>
          {uploads.map((upload) => (
            <li className={styles.uploadListItem} key={upload.id}>
              <div className={styles.fileInfoContainer}>
                <FontAwesomeIcon size="lg" icon={getIcon(upload.contentType)} className={styles.faIcon} />
                <div className={styles.fileInfo}>
                  <p>
                    {showDownloadLink ? (
                      <a href={upload.url} download>
                        {upload.filename}
                      </a>
                    ) : (
                      upload.filename
                    )}
                  </p>
                  <p className={styles.fileSizeAndTime}>
                    <span className={styles.uploadFileSize}>{bytes(upload.bytes)}</span>
                    <span>Uploaded {moment(upload.createdAt).format('DD MMM YYYY h:mm A')}</span>
                  </p>
                </div>
              </div>
              {showDeleteButton && (
                <Button type="button" unstyled onClick={() => onDelete(upload.id)}>
                  Delete
                </Button>
              )}
            </li>
          ))}
        </ul>
      </SectionWrapper>
    )
  );
};

UploadsTable.propTypes = {
  className: PropTypes.string,
  uploads: ExistingUploadsShape.isRequired,
  onDelete: PropTypes.func.isRequired,
  showDeleteButton: PropTypes.bool,
};

UploadsTable.defaultProps = {
  className: '',
  showDeleteButton: true,
};

export default UploadsTable;
