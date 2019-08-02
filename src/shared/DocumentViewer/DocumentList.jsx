import React from 'react';
import PropTypes from 'prop-types';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/fontawesome-free-solid';
import { Link } from 'react-router-dom';
import { renderStatusIcon, openLinkInNewWindow } from 'shared/utils';
import styles from 'shared/DocumentViewer/DocumentList.module.scss';
import moveInfoStyles from 'scenes/Office/MoveInfo.module.scss';

const documentUploadIcon = faPlusCircle;

const DocumentList = ({ currentMoveDocumentId, moveDocuments, detailUrlPrefix, disableLinks, uploadDocumentUrl }) => (
  <div>
    {moveDocuments.map(doc => {
      const chosenDocument = currentMoveDocumentId === doc.id ? styles['chosen-document'] : 'link-blue';
      const status = renderStatusIcon(doc.status);
      const detailUrl = `${detailUrlPrefix}/${doc.id}`;
      return (
        <div className={`panel-field ${chosenDocument}`} key={doc.id}>
          <span className="status">{status}</span>
          {!disableLinks &&
            (window.name === 'docViewer' ? (
              <Link className={`${chosenDocument} ${moveInfoStyles.doctitle}`} to={detailUrl}>
                {doc.title}
              </Link>
            ) : (
              <div
                className={`${chosenDocument} ${moveInfoStyles.doctitle}`}
                onClick={openLinkInNewWindow.bind(this, detailUrl, 'docViewer', window)}
              >
                {doc.title}
              </div>
            ))}
          {disableLinks && <span>{doc.title}</span>}
        </div>
      );
    })}
    {window.name === 'docViewer' ? (
      <div className={`${styles['document-upload-link']} link-blue`}>
        <Link to={uploadDocumentUrl}>
          <FontAwesomeIcon className="icon link-blue" icon={documentUploadIcon} />
          Upload new document
        </Link>
      </div>
    ) : (
      <div
        className={`${styles['document-upload-link']} link-blue`}
        data-cy="document-upload-link"
        onClick={openLinkInNewWindow.bind(this, uploadDocumentUrl, 'docViewer', window)}
      >
        <FontAwesomeIcon className="icon link-blue" icon={documentUploadIcon} />
        Upload new document
      </div>
    )}
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
