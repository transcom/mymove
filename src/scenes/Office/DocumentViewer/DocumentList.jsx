import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import '../office.css';

import { renderStatusIcon } from 'shared/utils';
import '../office.css';

export class DocumentList extends Component {
  render() {
    const { moveDocuments, moveId, disableLinks } = this.props;
    return (
      <div>
        {(moveDocuments || []).map(doc => {
          const status = renderStatusIcon(doc.status);
          const detailUrl = `/moves/${moveId}/documents/${doc.id}`;
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
  }
}

DocumentList.propTypes = {
  moveId: PropTypes.string.isRequired,
  moveDocuments: PropTypes.array.isRequired,
  disableLinks: PropTypes.bool,
};

export default DocumentList;
