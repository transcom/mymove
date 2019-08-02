import React from 'react';
import PropTypes from 'prop-types';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/fontawesome-free-solid';

import { renderStatusIcon, openLinkInNewWindow } from 'shared/utils';
import styles from 'shared/DocumentViewer/DocumentList.module.scss';

const documentUploadIcon = faPlusCircle;

const DocumentList = ({ currentMoveDocumentId, moveDocuments, detailUrlPrefix, disableLinks, uploadDocumentUrl }) => (
  <div>
    {moveDocuments.map(doc => {
      const chosenDocument = currentMoveDocumentId === doc.id ? styles.chosenDocument : 'link-blue';
      const status = renderStatusIcon(doc.status);
      const detailUrl = `${detailUrlPrefix}/${doc.id}`;
      return (
        <div className={`panel-field ${chosenDocument}`} key={doc.id}>
          <span className="status">{status}</span>
          {!disableLinks && (
            <div
              className={(chosenDocument, styles.doctitle)}
              onClick={openLinkInNewWindow.bind(this, detailUrl, '_blank', window)}
            >
              {doc.title}
            </div>
          )}
          {disableLinks && <span>{doc.title}</span>}
        </div>
      );
    })}
    <div
      className={`${styles['document-upload-link']} link-blue`}
      data-cy="document-upload-link"
      onClick={openLinkInNewWindow.bind(this, uploadDocumentUrl, '_blank', window)}
    >
      <FontAwesomeIcon className="icon link-blue" icon={documentUploadIcon} />
      Upload new document
    </div>
  </div>
);

DocumentList.propTypes = {
  currentMoveDocumentId: PropTypes.string,
  detailUrlPrefix: PropTypes.string.isRequired,
  disableLinks: PropTypes.bool,
  moveDocuments: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      status: PropTypes.string.isRequired,
      title: PropTypes.string.isRequired,
    }),
  ).isRequired,
  uploadDocumentUrl: PropTypes.string.isRequired,
};

export default DocumentList;
