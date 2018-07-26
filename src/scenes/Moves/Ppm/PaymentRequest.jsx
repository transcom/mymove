import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get } from 'lodash';

import DocumentUploader from 'scenes/Office/DocumentViewer/DocumentUploader';
import DocumentList from 'scenes/Office/DocumentViewer/DocumentList';
import {
  selectAllDocumentsForMove,
  getMoveDocumentsForMove,
} from 'shared/Entities/modules/moveDocuments';

import './PaymentRequest.css';
export class PaymentRequest extends Component {
  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }
  render() {
    const { moveDocuments } = this.props;
    const { moveId } = this.props.match.params;
    const numMoveDocs = get(moveDocuments, 'length', 'TBD');
    return (
      <div className="usa-grid payment-request">
        <div className="usa-width-two-thirds">
          <h2>Request Payment </h2>
          <div className="instructions">
            Please upload all your weight tickets, expenses, and storage fee
            documents one at a time. For expenses, you’ll need to enter
            additional details.
          </div>
          <DocumentUploader moveId={moveId} />
          <h4> Done uploading documents?</h4>
          <Link to="/" className="usa-button ">
            Submit Request
          </Link>
        </div>
        <div className="usa-width-one-third">
          <h4 className="doc-list-title">All Documents ({numMoveDocs})</h4>
          <DocumentList
            moveDocuments={moveDocuments}
            moveId={moveId}
            disableLinks={true}
          />
        </div>
      </div>
    );
  }
}
PaymentRequest.propTypes = {
  moveDocuments: PropTypes.array,
  moveId: PropTypes.string,
};

const mapStateToProps = (state, props) => ({
  moveDocuments: selectAllDocumentsForMove(state, props.match.params.moveId),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ getMoveDocumentsForMove }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequest);
