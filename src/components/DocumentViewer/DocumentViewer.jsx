import React, { useEffect, useState, useRef, useMemo } from 'react';
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
import { UPLOAD_DOC_STATUS, UPLOAD_SCAN_STATUS, UPLOAD_DOC_STATUS_DISPLAY_MESSAGE } from 'shared/constants';
import Alert from 'shared/Alert';
import { hasRotationChanged, toRotatedDegrees, toRotatedPosition } from 'shared/utils';

const DocumentViewer = ({ files, allowDownload, paymentRequestId, isFileUploading }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [disableSaveButton, setDisableSaveButton] = useState(false);
  const [menuIsOpen, setMenuOpen] = useState(false);
  const [showContentError, setShowContentError] = useState(false);
  const sortedFiles = files.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));
  const selectedFile = sortedFiles[parseInt(selectedFileIndex, 10)];
  const [isJustUploadedFile, setIsJustUploadedFile] = useState(false);
  const [fileStatus, setFileStatus] = useState(null);

  const [rotationValue, setRotationValue] = useState(selectedFile?.rotation || 0);

  const mountedRef = useRef(true);

  const queryClient = useQueryClient();

  useEffect(() => {
    if (isFileUploading) {
      setIsJustUploadedFile(true);
      setFileStatus(UPLOAD_DOC_STATUS.UPLOADING);
    } else {
      setIsJustUploadedFile(false);
    }
  }, [isFileUploading]);

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

  const fileTypeMap = useMemo(
    () => ({
      'application/pdf': 'pdf',
      'image/png': 'png',
      'image/jpeg': 'jpg',
      'image/jpg': 'jpg',
      'image/gif': 'gif',
    }),
    [],
  );

  const fileType = useRef(selectedFile?.contentType);

  useEffect(() => {
    const savedRotation = selectedFile?.rotation;
    const rotationChanged = hasRotationChanged(rotationValue, savedRotation, fileType.current);
    setDisableSaveButton(!rotationChanged);
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

    fileType.current = fileTypeMap[selectedFile?.contentType] || '';
    const initialRotation = toRotatedDegrees(selectedFile?.rotation, fileType.current);
    setRotationValue(initialRotation);
    const handleFileProcessing = async (status) => {
      switch (status) {
        case UPLOAD_SCAN_STATUS.PROCESSING:
          setFileStatus(UPLOAD_DOC_STATUS.SCANNING);
          break;
        case UPLOAD_SCAN_STATUS.CLEAN:
          setFileStatus(UPLOAD_DOC_STATUS.ESTABLISHING);
          break;
        case UPLOAD_SCAN_STATUS.INFECTED:
          setFileStatus(UPLOAD_DOC_STATUS.INFECTED);
          break;
        default:
          throw new Error(`unrecognized file status`);
      }
    };
    if (!isFileUploading && isJustUploadedFile) {
      setFileStatus(UPLOAD_DOC_STATUS.UPLOADING);
    }

    let sse;
    if (selectedFile) {
      sse = new EventSource(`/ghc/v1/uploads/${selectedFile.id}/status`, { withCredentials: true });
      sse.onmessage = (event) => {
        handleFileProcessing(event.data);
        if (
          event.data === UPLOAD_SCAN_STATUS.CLEAN ||
          event.data === UPLOAD_SCAN_STATUS.INFECTED ||
          event.data === 'Connection closed'
        ) {
          sse.close();
        }
      };
      sse.onerror = () => {
        sse.close();
        setFileStatus(null);
      };
    }

    return () => {
      sse?.close();
    };
  }, [selectedFile, isFileUploading, isJustUploadedFile, fileTypeMap]);
  useEffect(() => {
    if (fileStatus === UPLOAD_DOC_STATUS.ESTABLISHING) {
      setTimeout(() => {
        setFileStatus(UPLOAD_DOC_STATUS.LOADED);
        setShowContentError(false);
      }, 2000);
    }
  }, [fileStatus]);

  const getStatusMessage = (currentFileStatus, currentSelectedFile) => {
    switch (currentFileStatus) {
      case UPLOAD_DOC_STATUS.UPLOADING:
        return UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.UPLOADING;
      case UPLOAD_DOC_STATUS.SCANNING:
        return UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.SCANNING;
      case UPLOAD_DOC_STATUS.ESTABLISHING:
        return UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.ESTABLISHING_DOCUMENT_FOR_VIEWING;
      case UPLOAD_DOC_STATUS.INFECTED:
        return UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.INFECTED_FILE_MESSAGE;
      default:
        if (!currentSelectedFile) {
          return UPLOAD_DOC_STATUS_DISPLAY_MESSAGE.FILE_NOT_FOUND;
        }
        return null;
    }
  };

  const alertMessage = getStatusMessage(fileStatus, selectedFile);
  const alertType = fileStatus === UPLOAD_SCAN_STATUS.INFECTED ? 'error' : 'info';
  const alertHeading = fileStatus === UPLOAD_SCAN_STATUS.INFECTED ? 'Ask for a new file' : 'Document Status';
  if (alertMessage) {
    return (
      <Alert type={alertType} className="usa-width-one-whole" heading={alertHeading} data-testid="documentAlertHeading">
        <span data-testid="documentAlertMessage">{alertMessage}</span>
      </Alert>
    );
  }

  const openMenu = () => {
    setMenuOpen(true);
  };
  const closeMenu = () => {
    setMenuOpen(false);
  };

  const handleSelectFile = (index) => {
    selectFile(index);
    setFileStatus(UPLOAD_DOC_STATUS.ESTABLISHING);
    closeMenu();
  };

  const selectedFilename = filenameFromPath(selectedFile?.filename);

  const selectedFileDate = formatDate(moment(selectedFile?.createdAt), 'DD MMM YYYY');

  const onContentError = (errorObject) => {
    setShowContentError(true);
    milmoveLogger.error(errorObject);
  };

  const saveRotation = () => {
    if (mountedRef.current === true) {
      const rotationPosition = toRotatedPosition(rotationValue, fileType.current);

      const uploadBody = {
        rotation: rotationPosition,
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
        <Button
          data-testid="openMenu"
          type="button"
          onClick={openMenu}
          aria-label="Open menu"
          unstyled
          style={{ maxWidth: 'fit-content' }}
        >
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
          <Alert type="error" className="usa-width-one-whole" heading={alertHeading} data-testid="documentAlertHeading">
            <span data-testid="documentAlertMessage">
              MilMove encountered an issue during the scanning phase of this document. Contact the service member. Ask
              them to upload a photo of the original document instead.
            </span>
          </Alert>
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
