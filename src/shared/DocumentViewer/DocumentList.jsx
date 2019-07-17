import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { renderStatusIcon } from 'shared/utils';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/fontawesome-free-solid';

const documentUploadIcon = faPlusCircle;

const DocumentList = ({ currentMoveDocumentId, moveDocuments, detailUrlPrefix, disableLinks, uploadUrlPrefix }) => (
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
    {uploadUrlPrefix && (
      <div className="document-upload-link">
        <FontAwesomeIcon className="icon link-blue" icon={documentUploadIcon} />
        <Link to={uploadUrlPrefix}>Upload new document</Link>
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
  uploadUrlPrefix: PropTypes.string.isRequired,
};

export default DocumentList;
