import React, { Component, Fragment } from 'react';
import { string, arrayOf, object, shape, bool } from 'prop-types';
import { connect } from 'react-redux';
import Alert from 'shared/Alert'; // eslint-disable-line
import { get } from 'lodash';
import { includes } from 'lodash';
import qs from 'query-string';

import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';
import { convertDollarsToCents } from 'shared/utils';
import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import { createMovingExpenseDocument } from 'shared/Entities/modules/movingExpenseDocuments';

import { selectAllDocumentsForMove, getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

import { submitExpenseDocs } from './ducks.js';
import scrollToTop from 'shared/scrollToTop';

import './PaymentRequest.css';

function openLegalese(event) {
  // Prevent this from checking the box after opening the alert.
  event.preventDefault();

  const legalese = `\
LEGAL AGREEMENT / PRIVACY ACT

FINANCIAL LIABILITY:
If this shipment(s) incurs costs above the allowance I am entitled to, I will pay the difference to the government, \
or consent to the collection from my pay as necessary to cover all excess costs associated by this shipment(s).

ADVANCE OBLIGATION:
I understand that the maximum advance allowed is based on the estimated weight and scheduled departure date of my \
shipment(s). In the event, less weight is moved or my move occurs on a different scheduled departure date, I may have \
to remit the difference with the balance of my incentive disbursement and/or from the collection of my pay as may be \
necessary.

I understand that the maximum advance allowed is based on the estimated weight and scheduled departure date of my \
shipment(s). In the event, less weight is moved or my move occurs on a different scheduled departure date, I may \
have to remit the difference with the balance of my incentive disbursement and/or from the collection of my pay as may \
be necessary. If I receive an advance for my PPM shipment, I agree to furnish weight tickets within 45 days of final \
delivery to my destination. I understand that failure to furnish weight tickets within this timeframe may lead to the \
collection of my pay as necessary to cover the cost of the advance.`;

  alert(legalese);
}

function RequestPaymentSection(props) {
  const { ppm, updatingPPM, submitDocs, disableSubmit } = props;

  if (!ppm) {
    return null;
  }

  if (ppm.status === 'APPROVED') {
    return (
      <Fragment>
        <h4>Done uploading documents?</h4>
        <div className="customer-agreement">
          <p>
            <strong>Customer Agreement</strong>
          </p>
          <input id="agree-checkbox" type="checkbox" />
          <label for="agree-checkbox">
            I agree to the
            <a onClick={openLegalese}> Legal Agreement / Privacy Act</a>
          </label>
        </div>
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
  static propTypes = {
    currentPpm: shape({ id: string.isRequired }).isRequired,
    docTypes: arrayOf(string),
    moveDocuments: arrayOf(object).isRequired,
    genericMoveDocSchema: object.isRequired,
    moveDocSchema: object.isRequired,
    updatingPPM: bool.isRequired,
    updateError: bool.isRequired,
  };

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
        scrollToTop();
      });
  }

  handleSubmit = (uploadIds, formValues) => {
    const {
      match: {
        params: { moveId },
      },
      currentPpm,
    } = this.props;
    const {
      title,
      moving_expense_type: movingExpenseType,
      move_document_type: moveDocumentType,
      requested_amount_cents: requestedAmountCents,
      payment_method: paymentMethod,
      notes,
    } = formValues;
    const personallyProcuredMoveId = currentPpm ? currentPpm.id : null;

    if (get(formValues, 'move_document_type', false) === 'EXPENSE') {
      return this.props.createMovingExpenseDocument({
        moveId,
        personallyProcuredMoveId,
        uploadIds,
        title,
        movingExpenseType,
        moveDocumentType,
        requestedAmountCents: convertDollarsToCents(requestedAmountCents),
        paymentMethod,
        notes,
      });
    }
    return this.props.createMoveDocument({
      moveId,
      personallyProcuredMoveId,
      uploadIds,
      title,
      moveDocumentType,
      notes,
    });
  };

  render() {
    const { location, moveDocuments, updateError, docTypes } = this.props;
    const numMoveDocs = get(moveDocuments, 'length', 'TBD');
    const disableSubmit = numMoveDocs === 0;
    const moveDocumentType = qs.parse(location.search).moveDocumentType;
    const initialValues = {};

    // Verify the provided doc type against the schema
    if (includes(docTypes, moveDocumentType)) {
      initialValues.move_document_type = moveDocumentType;
    }

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
          <DocumentUploader
            form="payment-docs"
            genericMoveDocSchema={this.props.genericMoveDocSchema}
            initialValues={initialValues}
            isPublic={false}
            location={location}
            moveDocSchema={this.props.moveDocSchema}
            onSubmit={this.handleSubmit}
          />
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

const mapStateToProps = (state, props) => ({
  moveDocuments: selectAllDocumentsForMove(state, props.match.params.moveId),
  currentPpm: state.ppm.currentPpm,
  updatingPPM: state.ppm.hasSubmitInProgress,
  updateError: state.ppm.hasSubmitError,
  docTypes: get(state, 'swaggerInternal.spec.definitions.MoveDocumentType.enum', []),
  genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
  moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
});

const mapDispatchToProps = {
  getMoveDocumentsForMove,
  submitExpenseDocs,
  createMoveDocument,
  createMovingExpenseDocument,
};
export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequest);
