import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import '../office.css';

import { indexMoveDocuments } from './ducks.js';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';

export class DocumentList extends Component {
  renderDocStatus(status) {
    if (status === 'AWAITING_REVIEW') {
      return (
        <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />
      );
    }
    if (status === 'OK') {
      return <FontAwesomeIcon className="icon approval-ready" icon={faCheck} />;
    }
    if (status === 'HAS_ISSUE') {
      return (
        <FontAwesomeIcon
          className="icon approval-problem"
          icon={faExclamationCircle}
        />
      );
    }
  }

  render() {
    const { moveDocuments, moveId } = this.props;
    return (
      <div className="doclist">
        {moveDocuments.map(doc => {
          const status = this.renderDocStatus(doc.status);
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

const mapStateToProps = state => ({
  moveDocuments: get(state, 'moveDocuments.moveDocuments', {}),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ indexMoveDocuments }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(DocumentList);
