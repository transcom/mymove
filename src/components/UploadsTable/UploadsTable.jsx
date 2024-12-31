import React, { useState, useEffect } from 'react';
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
  const [fileAvailable, setFileAvailable] = useState({});
  const pollingInterval = 5000; // Poll every 5 seconds

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

  const checkFileAvailability = async (url, fileId) => {
    try {
      const response = await fetch(url, { method: 'HEAD' }); // Send a HEAD request to check availability
      if (response.ok) {
        setFileAvailable((prev) => ({ ...prev, [fileId]: true })); // Mark as available
      } else {
        setFileAvailable((prev) => ({ ...prev, [fileId]: false })); // Mark as unavailable
      }
    } catch (error) {
      setFileAvailable((prev) => ({ ...prev, [fileId]: false })); // Mark as unavailable if error occurs
    }
  };

  useEffect(() => {
    const intervalIds = {}; // Store interval IDs for each file to clear later

    uploads.forEach((upload) => {
      if (upload.url && !Object.hasOwn(fileAvailable, upload.id)) {
        // Start polling if the file URL is available and polling isn't already in progress
        intervalIds[upload.id] = setInterval(() => {
          checkFileAvailability(upload.url, upload.id);
        }, pollingInterval);
      }
    });

    // Cleanup polling on component unmount
    return () => {
      Object.values(intervalIds).forEach(clearInterval); // Clear all intervals
    };
  }, [uploads, fileAvailable]);

  const renderFileContent = (upload) => {
    if (showDownloadLink && upload.url) {
      // If the file is available, show a link; otherwise, just the filename as plain text
      return fileAvailable[upload.id] ? (
        <a href={upload.url} download>
          {upload.filename}
        </a>
      ) : (
        upload.filename // Plain text filename if file is not available
      );
    }

    return upload.filename; // Plain text filename if download link is not shown
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
                  <p>{renderFileContent(upload)}</p>
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
