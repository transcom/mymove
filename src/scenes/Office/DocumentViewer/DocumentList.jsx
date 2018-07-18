import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import { selectAllDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { renderStatusIcon } from 'shared/utils';

import '../office.css';

export class DocumentList extends Component {
  render() {
    const { moveDocuments, moveId } = this.props;
    return (
      <div>
        {moveDocuments.map(doc => {
          const status = renderStatusIcon(doc.status);
          const detailUrl = `/moves/${moveId}/documents/${doc.id}`;
          return (
            <div className="panel-field" key={doc.id}>
              <span className="status">{status}</span>
              <Link to={detailUrl}>{doc.title}</Link>
            </div>
          );
        })}
      </div>
    );
  }
}

DocumentList.propTypes = {
  moveDocuments: PropTypes.array,
  moveId: PropTypes.string,
};

const mapStateToProps = (state, props) => ({
  moveDocuments: selectAllDocumentsForMove(state, props.moveId),
});

export default connect(mapStateToProps)(DocumentList);
