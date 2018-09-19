import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { renderStatusIcon } from 'shared/utils';

const DocumentList = ({ moveDocuments, detailUrlPrefix, disableLinks }) => (
  <div>
    {moveDocuments.map(doc => {
      const status = renderStatusIcon(doc.status);
      const detailUrl = `${detailUrlPrefix}/${doc.id}`;
      return (
        <div className="panel-field" key={doc.id}>
          <span className="status">{status}</span>
          {!disableLinks && <Link to={detailUrl}>{doc.title}</Link>}
          {disableLinks && <span>{doc.title}</span>}
        </div>
      );
    })}
  </div>
);

DocumentList.propTypes = {
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
