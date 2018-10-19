import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { renderStatusIcon } from 'shared/utils';

const DocumentList = ({ currentMoveDocumentId, moveDocuments, detailUrlPrefix, disableLinks }) => (
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
};

export default DocumentList;
