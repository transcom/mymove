import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { renderStatusIcon } from 'shared/utils';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/fontawesome-free-solid';
import styles from './DocumentList.module.scss';

const documentUploadIcon = faPlusCircle;

const DocumentList = ({ currentMoveDocumentId, moveDocuments, detailUrlPrefix, disableLinks, uploadDocumentUrl }) => (
  <div>
    {moveDocuments.map(doc => {
      const chosenDocument = currentMoveDocumentId === doc.id ? 'chosen-document' : null;
      const status = renderStatusIcon(doc.status);
      const detailUrl = `${detailUrlPrefix}/${doc.id}`;
      return (
        <div className={`panel-field ${chosenDocument}`} key={doc.id}>
          <span className="status">{status}</span>
          {!disableLinks && (
            <Link className={chosenDocument} to={detailUrl}>
              {doc.title}
            </Link>
          )}
          {disableLinks && <span>{doc.title}</span>}
        </div>
      );
    })}
    <div className={styles['document-upload-link']} data-cy="document-upload-link">
      <FontAwesomeIcon className="icon link-blue" icon={documentUploadIcon} />
      <Link to={uploadDocumentUrl} target="_blank">
        Upload new document
      </Link>
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
