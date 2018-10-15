import React, { Component, Fragment } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import Alert from 'shared/Alert'; // eslint-disable-line
import { get } from 'lodash';

import DocumentUploader from 'scenes/Office/DocumentViewer/DocumentUploader';

import { selectAllDocumentsForMove, getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

import { submitExpenseDocs } from './ducks.js';

import './PaymentRequest.css';

function RequestPaymentSection(props) {
  const { ppm, updatingPPM, submitDocs, disableSubmit } = props;

  if (!ppm) {
    return null;
  }

  if (ppm.status === 'APPROVED') {
    return (
      <Fragment>
        <h4>Done uploading documents?</h4>
        <button onClick={submitDocs} className="usa-button" disabled={updatingPPM || disableSubmit}>
          Submit Payment Request
        </button>
      </Fragment>
    );
  } else if (ppm.status === 'PAYMENT_REQUESTED') {
    return (
      <Fragment>
        <h4>Payment requested, awaiting approval.</h4>
      </Fragment>
    );
  } else {
    console.error('Unexpectedly got to PaymentRequest screen without PPM approval');
  }
}

export class PaymentRequest extends Component {
  constructor(props) {
    super(props);
    this.submitDocs = this.submitDocs.bind(this);
  }

  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }

  submitDocs() {
    this.props
      .submitExpenseDocs()
      .then(() => {
        this.props.history.push('/');
      })
      .catch(() => {
        window.scrollTo(0, 0);
      });
  }

  render() {
    const { moveDocuments, updateError } = this.props;
    const { moveId } = this.props.match.params;
    const numMoveDocs = get(moveDocuments, 'length', 'TBD');
    const disableSubmit = numMoveDocs === 0;
    return (
      <div className="usa-grid payment-request">
        <div className="usa-width-two-thirds">
          {updateError && (
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                There was an error requesting payment, please try again.
              </Alert>
            </div>
          )}
          <h2>Request Payment </h2>
          <div className="instructions">
            Please upload all your weight tickets, expenses, and storage fee documents one at a time. For expenses,
            you’ll need to enter additional details.
          </div>
          <DocumentUploader moveId={moveId} />
          <RequestPaymentSection
            ppm={this.props.currentPpm}
            updatingPPM={this.props.updatingPPM}
            submitDocs={this.submitDocs}
            disableSubmit={disableSubmit}
          />
        </div>
        <div className="usa-width-one-third">
          <h4 className="doc-list-title">All Documents ({numMoveDocs})</h4>
          {(moveDocuments || []).map(doc => {
            return (
              <div className="panel-field" key={doc.id}>
                <span>{doc.title}</span>
              </div>
            );
          })}
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
  currentPpm: state.ppm.currentPpm,
  updatingPPM: state.ppm.hasSubmitInProgress,
  updateError: state.ppm.hasSubmitError,
});

const mapDispatchToProps = dispatch => bindActionCreators({ getMoveDocumentsForMove, submitExpenseDocs }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequest);
