import React, { useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import { useQueryClient } from '@tanstack/react-query';
import { Button, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './DocumentViewerFileManager.module.scss';

import {
  createUploadForDocument,
  createUploadForAmdendedOrders,
  createUploadForSupportingDocuments,
  deleteUploadForDocument,
  getOrder,
} from 'services/ghcApi';
import { ORDERS_DOCUMENTS, MOVES, ORDERS } from 'constants/queryKeys';
import FileUpload from 'components/FileUpload/FileUpload';
import Hint from 'components/Hint';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import DeleteDocumentFileConfirmationModal from 'components/ConfirmationModals/DeleteDocumentFileConfirmationModal';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';

const DocumentViewerFileManager = ({
  className,
  move,
  orderId,
  documentId,
  files,
  documentType,
  updateAmendedDocument,
  fileUploadrequired,
}) => {
  const queryClient = useQueryClient();
  const filePondEl = useRef();
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isFileProcessing, setIsFileProcessing] = useState(false);
  const [currentFile, setCurrentFile] = useState(null);
  const [serverError, setServerError] = useState('');
  const [showUpload, setShowUpload] = useState(false);
  const [isExpandedView, setIsExpandedView] = useState(false);
  const [buttonHeaderText, setButtonHeaderText] = useState('Manage Documents');

  const moveId = move?.id;
  const moveCode = move?.locator;

  function appendTimestampToFilename(file) {
    // Create a date-time stamp in the format "yyyymmddhh24miss"
    const now = new Date();
    const timestamp =
      now.getFullYear().toString() +
      (now.getMonth() + 1).toString().padStart(2, '0') +
      now.getDate().toString().padStart(2, '0') +
      now.getHours().toString().padStart(2, '0') +
      now.getMinutes().toString().padStart(2, '0') +
      now.getSeconds().toString().padStart(2, '0');

    // Create a new filename with the timestamp prepended
    const newFileName = `${file.name}-${timestamp}`;

    // Create and return a new File object with the new filename
    return new File([file], newFileName, { type: file.type });
  }

  useEffect(() => {
    if (documentType === MOVE_DOCUMENT_TYPE.ORDERS) {
      setButtonHeaderText('Manage Orders');
    } else if (documentType === MOVE_DOCUMENT_TYPE.AMENDMENTS) {
      setButtonHeaderText('Manage Amended Orders');
    } else if (documentType === MOVE_DOCUMENT_TYPE.SUPPORTING) {
      setShowUpload(true);
      setIsExpandedView(true);
    }
    if (fileUploadrequired) {
      setShowUpload(true);
    }
  }, [documentType, fileUploadrequired]);

  const closeDeleteFileModal = () => {
    setCurrentFile(null);
    setIsDeleteModalOpen(false);
  };

  const toggleUploadVisibility = (e) => {
    e.preventDefault();
    setShowUpload((show) => !show);
    setServerError('');
  };

  const openDeleteFileModal = (uploadId) => {
    const selectedFile = files?.find((file) => file.id === uploadId);
    setCurrentFile(selectedFile);
    if (selectedFile) {
      setIsDeleteModalOpen(true);
      setServerError('');
    }
  };

  const uploadOrders = async (file) => {
    const newFile = appendTimestampToFilename(file);
    return createUploadForDocument(newFile, documentId)
      .catch((e) => {
        const { response } = e;
        const error = `Failed to upload due to server error: ${response?.body?.detail}`;
        setServerError(error);
      })
      .finally(() => {
        queryClient.invalidateQueries([ORDERS_DOCUMENTS, documentId]);
        setIsFileProcessing(false);
      });
  };

  const uploadAmdendedOrders = async (file) => {
    const newFile = appendTimestampToFilename(file);
    return createUploadForAmdendedOrders(newFile, orderId)
      .then(async () => {
        return getOrder(null, orderId)
          .then((res) => {
            const updatedOrder = res.orders[orderId];
            const amendedOrderDocumentId = updatedOrder?.uploadedAmendedOrderID;
            const newOrderEtag = updatedOrder?.eTag;

            updateAmendedDocument(amendedOrderDocumentId);
            queryClient.invalidateQueries([ORDERS_DOCUMENTS, amendedOrderDocumentId]);

            queryClient.setQueryData([ORDERS, orderId], (oldData) => {
              if (!oldData) return oldData;
              return {
                ...oldData,
                orders: {
                  ...oldData.orders,
                  [orderId]: {
                    ...oldData.orders[orderId],
                    eTag: newOrderEtag,
                  },
                },
              };
            });
          })
          .catch((e) => {
            const { response } = e;
            const error = `Failed to upload due to server error: ${response?.body?.detail}`;
            setServerError(error);
          });
      })
      .catch((e) => {
        const { response } = e;
        const error = `Failed to upload due to server error: ${response?.body?.detail}`;
        setServerError(error);
      })
      .finally(() => {
        setIsFileProcessing(false);
      });
  };

  const uploadSupportingDocuments = async (file) => {
    const newFile = appendTimestampToFilename(file);
    return createUploadForSupportingDocuments(newFile, moveId)
      .catch((e) => {
        const { response } = e;
        const error = `Failed to upload due to server error: ${response?.body?.detail}`;
        setServerError(error);
      })
      .finally(() => {
        queryClient.invalidateQueries([MOVES, moveCode]);
        setIsFileProcessing(false);
      });
  };

  const deleteDocuments = async () => {
    return deleteUploadForDocument(currentFile.id, orderId)
      .then(() => {
        if (documentType === MOVE_DOCUMENT_TYPE.SUPPORTING) {
          queryClient.invalidateQueries([MOVES, moveCode]);
        } else {
          queryClient.invalidateQueries([ORDERS_DOCUMENTS, documentId]);
        }
        closeDeleteFileModal();
      })
      .catch((e) => {
        const { response } = e;
        const error = `Failed to delete due to server error: ${response?.body?.detail}`;
        setServerError(error);
      })
      .finally(() => {
        setIsFileProcessing(false);
      });
  };

  const handleUpload = async (file) => {
    setIsFileProcessing(true);
    if (documentType === MOVE_DOCUMENT_TYPE.ORDERS) {
      uploadOrders(file);
    } else if (documentType === MOVE_DOCUMENT_TYPE.AMENDMENTS) {
      uploadAmdendedOrders(file);
    } else if (documentType === MOVE_DOCUMENT_TYPE.SUPPORTING) {
      uploadSupportingDocuments(file);
    }
  };

  const handleChange = () => {
    filePondEl.current?.removeFiles();
    queryClient.invalidateQueries([ORDERS_DOCUMENTS, documentId]);
    setServerError('');
  };

  const handleDeleteSubmit = () => {
    if (!isFileProcessing) {
      setIsFileProcessing(true);
      deleteDocuments();
    }
  };

  return (
    <div className={classnames(styles.documentViewerFileManager, className)}>
      {currentFile && (
        <DeleteDocumentFileConfirmationModal
          isOpen={isDeleteModalOpen}
          closeModal={closeDeleteFileModal}
          submitModal={handleDeleteSubmit}
          fileInfo={currentFile}
        />
      )}
      {!isExpandedView && (
        <Button disabled={isFileProcessing || fileUploadrequired} onClick={toggleUploadVisibility}>
          {buttonHeaderText}
        </Button>
      )}
      <div>
        {showUpload && (
          <>
            <br />
            {serverError && (
              <Alert type="error" headingLevel="h4" heading="An error occurred">
                {serverError}
              </Alert>
            )}
            <UploadsTable className={styles.sectionWrapper} uploads={files} onDelete={openDeleteFileModal} />
            <div className={classnames(styles.upload, className)}>
              {fileUploadrequired && (
                <Alert type="error" id="fileRequiredAlert" data-testid="fileRequiredAlert">
                  File upload is required
                </Alert>
              )}
              <FileUpload
                required={fileUploadrequired}
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={handleChange}
                labelIdle={'Drag files here or <span class="filepond--label-action">click</span> to upload'}
              />
              <Hint>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</Hint>
              {!isExpandedView && (
                <Button disabled={isFileProcessing || fileUploadrequired} onClick={toggleUploadVisibility}>
                  Done
                </Button>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

DocumentViewerFileManager.propTypes = {
  orderId: PropTypes.string.isRequired,
  documentId: PropTypes.string.isRequired,
  files: PropTypes.array.isRequired,
  documentType: PropTypes.string.isRequired,
};

export default DocumentViewerFileManager;
