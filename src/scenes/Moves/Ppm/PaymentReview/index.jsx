import React, { Component } from 'react';
import { connect } from 'react-redux';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { get } from 'lodash';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';
import { ppmPaymentLegal } from 'scenes/Legalese/legaleseText';
import PPMPaymentRequestActionBtns from 'scenes/Moves/Ppm/PPMPaymentRequestActionBtns';
import moment from 'moment';
import Alert from 'shared/Alert';
import { createSignedCertification } from 'shared/Entities/modules/signed_certifications';
import scrollToTop from 'shared/scrollToTop';
import { submitExpenseDocs } from '../ducks';
import DocumentsUploaded from './DocumentsUploaded';
import WizardHeader from '../../WizardHeader';
import './PaymentReview.css';

const nextBtnLabel = 'Submit Request';

class PaymentReview extends Component {
  state = {
    acceptTerms: false,
    moveSubmissionError: false,
  };

  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.moveId);
  }

  handleOnAcceptTermsChange = acceptTerms => {
    this.setState({ acceptTerms });
  };

  submitCertificate = () => {
    const signatureTime = moment().format();
    const { currentPpm, moveId } = this.props;
    const certificate = {
      certification_text: ppmPaymentLegal,
      date: signatureTime,
      signature: 'CHECKBOX',
      personally_procured_move_id: currentPpm.id,
      certification_type: 'PPM_PAYMENT',
    };
    return this.props.createSignedCertification(moveId, certificate);
  };

  applyClickHandlers = () => {
    this.setState({ moveSubmissionError: false });
    Promise.all([this.submitCertificate(), this.props.submitExpenseDocs()])
      .then(() => {
        this.props.history.push('/');
      })
      .catch(() => {
        this.setState({ moveSubmissionError: true });
        scrollToTop();
      });
  };

  render() {
    const { moveId, moveDocuments, submitting } = this.props;
    const weightTickets = moveDocuments.weightTickets;
    const missingSomeWeightTicket = weightTickets.some(
      ({ empty_weight_ticket_missing, full_weight_ticket_missing }) =>
        empty_weight_ticket_missing || full_weight_ticket_missing,
    );

    return (
      <>
        <WizardHeader
          title="Review"
          right={
            <ProgressTimeline>
              <ProgressTimelineStep name="Weight" completed />
              <ProgressTimelineStep name="Expenses" completed />
              <ProgressTimelineStep name="Review" current />
            </ProgressTimeline>
          }
        />
        <div className="payment-review-container usa-grid">
          <div className="review-payment-request-header">
            {this.state.moveSubmissionError && (
              <div className="usa-width-one-whole error-message">
                <Alert type="error" heading="An error occurred">
                  There was an error requesting payment, please try again.
                </Alert>
              </div>
            )}
            <h3>Review Payment Request</h3>
            <p>
              Make sure <strong>all</strong> your documents are uploaded.
            </p>
          </div>

          <DocumentsUploaded inReviewPage showLinks moveId={moveId} />

          <div className="doc-review">
            {missingSomeWeightTicket && (
              <>
                <h4 className="missing-label">
                  <FontAwesomeIcon
                    style={{ marginLeft: 0, color: 'red' }}
                    className="icon"
                    icon={faExclamationCircle}
                  />{' '}
                  Your estimated payment is unknown
                </h4>
                <p>
                  We cannot give you estimated payment because of missing weight tickets. Submit your payment request,
                  then go to the PPPO office to receive help in resolving this issue.
                </p>
              </>
            )}
          </div>
          <div className="review-customer-agreement-container">
            <CustomerAgreement
              className="review-customer-agreement"
              onChange={this.handleOnAcceptTermsChange}
              link="/ppm-customer-agreement"
              checked={this.state.acceptTerms}
              agreementText={ppmPaymentLegal}
            />
          </div>
          <PPMPaymentRequestActionBtns
            nextBtnLabel={nextBtnLabel}
            submitButtonsAreDisabled={!this.state.acceptTerms}
            saveAndAddHandler={this.applyClickHandlers}
            submitting={submitting}
            displaySaveForLater
          />
        </div>
      </>
    );
  }
}

const mapStateToProps = (state, props) => {
  const { moveId } = props.match.params;
  return {
    moveDocuments: {
      expenses: selectPPMCloseoutDocumentsForMove(state, moveId, ['EXPENSE']),
      weightTickets: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
    },
    moveId,
    currentPpm: get(state, 'ppm.currentPpm'),
  };
};

const mapDispatchToProps = {
  submitExpenseDocs,
  createSignedCertification,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(PaymentReview);
