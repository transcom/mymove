import React, { Component } from 'react';
import { connect } from 'react-redux';
import Alert from 'shared/Alert';
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
import PropTypes from 'prop-types';
import { createSignedCertification } from 'shared/Entities/modules/signed_certifications';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';
import { ppmPaymentLegal } from 'scenes/Legalese/legaleseText';

export class PaymentRequest extends Component {
  state = {
    acceptTerms: false,
  };

  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }

  submitDocs = () => {
    this.props
      .submitExpenseDocs()
      .then(() => {
        this.props.history.push('/');
      })
      .catch(() => {
        scrollToTop();
      });
  };

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

  submitCertificate = () => {
    const certificate = {
      certification_text: ppmPaymentLegal,
      date: PaymentRequest.getUserDate(),
      signature: 'CHECKBOX',
      personally_procured_move_id: this.props.currentPpm.id,
      certification_type: 'PPM_PAYMENT',
    };
    this.props.createSignedCertification(this.props.match.params.moveId, certificate);
  };

  handleOnAcceptTermsChange = acceptTerms => {
    this.setState({ acceptTerms: acceptTerms });
  };

  applyClickHandlers = () => {
    this.submitDocs();
    this.submitCertificate();
  };

  static getUserDate() {
    return new Date().toISOString().split('T')[0];
  }

  renderCustomerAgreement = ppmStatus => {
    switch (ppmStatus) {
      case null:
        //ppm hasn't loaded yet
        return;
      case 'APPROVED':
        return (
          <>
            <h4>Done uploading documents?</h4>
            <CustomerAgreement
              onAcceptTermsChange={this.handleOnAcceptTermsChange}
              checked={this.state.acceptTerms}
              agreementText={ppmPaymentLegal}
            />
          </>
        );
      case 'PAYMENT_REQUESTED':
        return <h4>Payment requested, awaiting approval.</h4>;
      default:
        console.error('Unexpectedly got to PaymentRequest screen without PPM approval');
    }
  };

  render() {
    const { location, moveDocuments, updatingPPM, updateError, docTypes, currentPpm } = this.props;
    const numMoveDocs = get(moveDocuments, 'length', 'TBD');
    const atLeastOneMoveDoc = numMoveDocs > 0;
    // TODO don't think this part does anything possibly delete
    const moveDocumentType = location ? qs.parse(location.search).moveDocumentType : null;
    const currentPpmStatus = currentPpm ? currentPpm.status : null;
    const initialValues = {};
    const canSubmitPayment = !updatingPPM && atLeastOneMoveDoc && this.state.acceptTerms;

    // TODO don't think this part does anything possibly delete
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
          <h2>Request Payment</h2>
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
          <div>{this.renderCustomerAgreement(currentPpmStatus)}</div>
          <button onClick={this.applyClickHandlers} disabled={!canSubmitPayment} className="usa-button">
            Submit Payment Request
          </button>
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
  currentPpm: PropTypes.shape({ id: PropTypes.string.isRequired }),
  docTypes: PropTypes.arrayOf(PropTypes.string),
  moveDocuments: PropTypes.arrayOf(PropTypes.object).isRequired,
  genericMoveDocSchema: PropTypes.object.isRequired,
  moveDocSchema: PropTypes.object.isRequired,
  updatingPPM: PropTypes.bool,
  updateError: PropTypes.bool.isRequired,
};

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
  createSignedCertification,
  getMoveDocumentsForMove,
  submitExpenseDocs,
  createMoveDocument,
  createMovingExpenseDocument,
};

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequest);
