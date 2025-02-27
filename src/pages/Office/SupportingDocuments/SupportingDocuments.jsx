import React, { useState } from 'react';
import moment from 'moment';
import classNames from 'classnames';

import DocumentViewerFileManager from 'components/DocumentViewerFileManager/DocumentViewerFileManager';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import Restricted from 'components/Restricted/Restricted';
import styles from 'styles/documentViewerWithSidebar.module.scss';
import { permissionTypes } from 'constants/permissions';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';

const SupportingDocuments = ({ move, uploads }) => {
  const [isFileUploading, setFileUploading] = useState(false);
  const filteredAndSortedUploads = Object.values(uploads || {})
    ?.filter((file) => {
      return !file.deletedAt;
    })
    ?.sort((a, b) => moment(b.createdAt) - moment(a.createdAt));
  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        {!filteredAndSortedUploads ||
        filteredAndSortedUploads.constructor !== Array ||
        filteredAndSortedUploads?.length <= 0 ? (
          <h2>No supporting documents have been uploaded.</h2>
        ) : (
          <DocumentViewer files={filteredAndSortedUploads} allowDownload isFileUploading={isFileUploading} />
        )}
      </div>
      <Restricted to={permissionTypes.createSupportingDocuments}>
        <div className={styles.sidebar}>
          <div className={styles.content}>
            <div className={classNames(styles.top, styles.noBottomBorder)}>
              <DocumentViewerFileManager
                move={move}
                orderId={move.orderId}
                documentId={move.additionalDocuments?.id}
                files={filteredAndSortedUploads}
                documentType={MOVE_DOCUMENT_TYPE.SUPPORTING}
                onAddFile={() => {
                  setFileUploading(true);
                }}
              />
            </div>
          </div>
        </div>
      </Restricted>
    </div>
  );
};

export default SupportingDocuments;
