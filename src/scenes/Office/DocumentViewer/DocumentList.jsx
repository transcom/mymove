import React, { Component, Fragment } from 'react';
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
          className="icon approval-waiting"
          icon={faExclamationCircle}
        />
      );
    }
  }

  render() {
    const { moveDocuments } = this.props;
    return (
      <Fragment>
        {moveDocuments.map(doc => {
          const status = this.renderDocStatus(doc.status);
          return (
            <div key={doc.id}>
              <span className="status">{status}</span>
              <Link to="/" target="_blank">
                {doc.document.name}
              </Link>
            </div>
          );
        })}
      </Fragment>
    );
  }
}

DocumentList.propTypes = {
  moveDocuments: PropTypes.array,
};

const mapStateToProps = state => ({
  moveDocuments: get(state, 'moveDocuments.moveDocuments', {}),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ indexMoveDocuments }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(DocumentList);
