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
import { UPLOAD_DOC_STATUS, UPLOAD_SCAN_STATUS } from 'shared/constants';
import Alert from 'shared/Alert';

/**
 * TODO
 * - implement next/previous pages instead of scroll through pages
 * - implement rotate left/right
 * - handle fetch doc errors
 */

const DocumentViewer = ({ files, isFileUploading, allowDownload, paymentRequestId }) => {
  const [selectedFileIndex, selectFile] = useState(0);
  const [fileStatus, setFileStatus] = useState(null);
  const [disableSaveButton, setDisableSaveButton] = useState(false);
  const [menuIsOpen, setMenuOpen] = useState(false);
  const sortedFiles = files.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));
  const selectedFile = sortedFiles[parseInt(selectedFileIndex, 10)];

  const [rotationValue, setRotationValue] = useState(selectedFile?.rotation || 0);

  const mountedRef = useRef(true);

  const queryClient = useQueryClient();

  useEffect(() => {
    if (isFileUploading) {
      setFileStatus(UPLOAD_DOC_STATUS.UPLOADING);
    } else {
      setFileStatus(UPLOAD_DOC_STATUS.ESTABLISHING);
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
  }, [rotationValue, selectedFile]);

  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    selectFile(0);
  }, [files.length]);

  useEffect(() => {
    setRotationValue(selectedFile?.rotation || 0);

    if (isFileUploading) return undefined;

    const handleFileProcessing = async (newStatus) => {
      if (newStatus === UPLOAD_SCAN_STATUS.PROCESSING) {
        setFileStatus(UPLOAD_DOC_STATUS.SCANNING);
      } else if (newStatus === UPLOAD_SCAN_STATUS.CLEAN) {
        setFileStatus(UPLOAD_DOC_STATUS.ESTABLISHING);
      } else if (newStatus === UPLOAD_SCAN_STATUS.INFECTED) {
        setFileStatus(UPLOAD_DOC_STATUS.INFECTED);
      } else {
        setFileStatus(null);
      }
    };

    const sse = new EventSource(`/internal/uploads/${selectedFile.id}/status`, { withCredentials: true });
    sse.onmessage = (event) => {
      if (event.data === UPLOAD_SCAN_STATUS.CLEAN || event.data === UPLOAD_SCAN_STATUS.INFECTED) {
        sse.close();
      }
      handleFileProcessing(event.data);
    };
    sse.onerror = () => {
      setFileStatus(null);
    };

    return () => {
      sse.close();
    };
  }, [selectedFile, isFileUploading]);

  useEffect(() => {
    if (fileStatus === 'ESTABLISHING') {
      new Promise((resolve) => {
        setTimeout(resolve, 3000);
      }).then(() => setFileStatus(UPLOAD_DOC_STATUS.LOADED));
    }
  }, [fileStatus]);

  const fileType = useRef(selectedFile?.contentType);

  if (!selectedFile || !fileStatus || selectedFile?.status === UPLOAD_SCAN_STATUS.INFECTED) {
    return <Alert heading="File Not Found" />;
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

  if (fileStatus && fileStatus !== 'LOADED') {
    return (
      <Alert type="info" className="usa-width-one-whole" heading="Document Status">
        {fileStatus === UPLOAD_DOC_STATUS.UPLOADING && 'Uploading'}
        {fileStatus === UPLOAD_DOC_STATUS.SCANNING && 'Scanning'}
        {fileStatus === UPLOAD_DOC_STATUS.ESTABLISHING && 'Establishing Document for View'}
      </Alert>
    );
  }

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
      <Content
        fileType={fileType.current}
        filePath={selectedFile?.url}
        rotationValue={rotationValue}
        disableSaveButton={disableSaveButton}
        setRotationValue={setRotationValue}
        saveRotation={saveRotation}
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
