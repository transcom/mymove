import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { renderStatusIcon } from 'shared/utils';
import 'scenes/Office/office.css';

export class DocumentList extends Component {
  render() {
    const { detailUrlPrefix, moveDocuments, disableLinks } = this.props;
    return (
      <div>
        {(moveDocuments || []).map(doc => {
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
  }
}

DocumentList.propTypes = {
  detailUrlPrefix: PropTypes.string.isRequired,
  disableLinks: PropTypes.bool,
  moveDocuments: PropTypes.array.isRequired,
};

export default DocumentList;
