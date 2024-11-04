import React, { useEffect, useState, useRef } from 'react';
import { bool, PropTypes } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import moment from 'moment';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import { FileShape } from './types';
import styles from './DocumentViewer.module.scss';
import Content from './Content/Content';
import Menu from './Menu/Menu';

import { milmoveLogger } from 'utils/milmoveLog';
import { UPLOADS } from 'constants/queryKeys';
import { bulkDownloadPaymentRequest, updateUpload } from 'services/ghcApi';
import { formatDate } from 'shared/dates';
import { filenameFromPath } from 'utils/formatters';
import AsyncPacketDownloadLink from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';

/**
 * TODO
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 */

const DocumentViewer = ({ files, allowDownload, paymentRequestId }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [disableSaveButton, setDisableSaveButton] = useState(false);
  const [menuIsOpen, setMenuOpen] = useState(false);
  const [showContentError, setShowContentError] = useState(false);
  const sortedFiles = files.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));
  const selectedFile = sortedFiles[parseInt(selectedFileIndex, 10)];

  const [rotationValue, setRotationValue] = useState(selectedFile?.rotation || 0);

  const mountedRef = useRef(true);

  const queryClient = useQueryClient();

  const { mutate: mutateUploads } = useMutation(updateUpload, {
    onSuccess: async (data, variables) => {
      if (mountedRef.current) {
        await queryClient.setQueryData([UPLOADS, variables.uploadID], data);
        await queryClient.invalidateQueries(UPLOADS);
      }
    },
    onError: (error) => {
      const errorMsg = error;
      milmoveLogger.error(errorMsg);
    },
  });

  useEffect(() => {
    const selectedFileHasRotation = selectedFile?.rotation !== undefined;
    if (
      (selectedFileHasRotation && selectedFile?.rotation !== rotationValue) ||
      (!selectedFileHasRotation && rotationValue !== 0)
    ) {
      setDisableSaveButton(false);
    } else {
      setDisableSaveButton(true);
    }
  }, [rotationValue, selectedFile, selectFile]);

  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    selectFile(0);
  }, [files.length]);

  useEffect(() => {
    setShowContentError(false);
    setRotationValue(selectedFile?.rotation || 0);
  }, [selectedFile]);

  const fileType = useRef(selectedFile?.contentType);

  if (!selectedFile) {
    return <h2>File Not Found</h2>;
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

  const fileTypeMap = {
    'application/pdf': 'pdf',
    'image/png': 'png',
    'image/jpeg': 'jpg',
    'image/jpg': 'jpg',
    'image/gif': 'gif',
  };

  fileType.current = fileTypeMap[selectedFile?.contentType] || '';

  const selectedFilename = filenameFromPath(selectedFile?.filename);

  const selectedFileDate = formatDate(moment(selectedFile?.createdAt), 'DD MMM YYYY');

  const onContentError = (errorObject) => {
    setShowContentError(true);
    milmoveLogger.error(errorObject);
  };

  const saveRotation = () => {
    if (fileType.current !== 'pdf' && mountedRef.current === true) {
      const uploadBody = {
        rotation: rotationValue,
      };
      mutateUploads({ uploadID: selectedFile?.id, body: uploadBody });
      setDisableSaveButton(true);
    }
  };

  const paymentPacketDownload = (
    <div>
      <dd data-testid="bulkPacketDownload">
        <p className={styles.bulkDownload}>
          <AsyncPacketDownloadLink
            id={paymentRequestId}
            label="Download All Files (PDF)"
            asyncRetrieval={bulkDownloadPaymentRequest}
          />
        </p>
      </dd>
    </div>
  );

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
            <a href={selectedFile?.url} download tabIndex={menuIsOpen ? '-1' : '0'}>
              <span>Download file</span> <FontAwesomeIcon icon="download" />
            </a>
          </p>
        )}
        {paymentRequestId !== undefined ? paymentPacketDownload : null}
      </div>
      {showContentError && (
        <div className={styles.errorMessage}>
          <h2>File Not Found</h2>The document is not yet available for viewing. Please try again in a moment.
        </div>
      )}
      <Content
        fileType={fileType.current}
        filePath={selectedFile?.url}
        rotationValue={rotationValue}
        disableSaveButton={disableSaveButton}
        setRotationValue={setRotationValue}
        saveRotation={saveRotation}
        onError={onContentError}
      />
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
  files: PropTypes.arrayOf(FileShape),
  allowDownload: bool,
};

DocumentViewer.defaultProps = {
  files: [],
  allowDownload: false,
};

export default DocumentViewer;
