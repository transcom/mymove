import React, { useState, useEffect, useRef, useCallback } from 'react';
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
  const intervalIds = useRef({}); // Use a ref to persist interval IDs

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

  const testLinkWithoutDownload = (fileUrl) => {
    return new Promise((resolve) => {
      const iframe = document.createElement('iframe');
      iframe.style.display = 'none';
      iframe.sandbox = 'allow-scripts'; // Restrict iframe to prevent download

      iframe.onload = () => {
        resolve(true);
        document.body.removeChild(iframe);
      };

      iframe.onerror = () => {
        resolve(false);
        document.body.removeChild(iframe);
      };

      iframe.src = fileUrl;
      document.body.appendChild(iframe);
    });
  };

  const pollForValidLink = useCallback(
    async (fileUrl, fileId) => {
      const checkLink = async () => {
        const isValid = await testLinkWithoutDownload(fileUrl);
        if (isValid) {
          setFileAvailable((prev) => ({ ...prev, [fileId]: true }));
          clearInterval(intervalIds.current[fileId]); // Stop polling
        } else {
          setFileAvailable((prev) => ({ ...prev, [fileId]: false }));
        }
      };

      if (!intervalIds.current[fileId]) {
        intervalIds.current[fileId] = setInterval(checkLink, pollingInterval);
        await checkLink(); // Run immediately on start
      }
    },
    [pollingInterval],
  );

  useEffect(() => {
    const localIntervalIds = intervalIds.current; // Capture the current ref value

    uploads.forEach((upload) => {
      if (upload.url && !Object.hasOwn(fileAvailable, upload.id)) {
        pollForValidLink(upload.url, upload.id);
      }
    });

    return () => {
      Object.values(localIntervalIds).forEach(clearInterval);
    };
  }, [uploads, fileAvailable, pollForValidLink]);

  const renderFileContent = (upload) => {
    if (showDownloadLink && upload.url) {
      return fileAvailable[upload.id] ? (
        <a href={upload.url} download>
          {upload.filename}
        </a>
      ) : (
        upload.filename
      );
    }

    return upload.filename;
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
