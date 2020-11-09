import React, { Component } from 'react';
import { connect } from 'react-redux';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { get, isEmpty } from 'lodash';
import moment from 'moment';
import Alert from 'shared/Alert';
import { formatCents } from 'shared/formatters';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { createSignedCertification } from 'shared/Entities/modules/signed_certifications';
import scrollToTop from 'shared/scrollToTop';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';
import { ppmPaymentLegal } from 'scenes/Legalese/legaleseText';
import PPMPaymentRequestActionBtns from 'scenes/Moves/Ppm/PPMPaymentRequestActionBtns';
import { loadEntitlementsFromState } from 'shared/entitlements';

import { submitExpenseDocs } from '../ducks';
import DocumentsUploaded from './DocumentsUploaded';
import { calcNetWeight } from '../utility';
import WizardHeader from '../../WizardHeader';
import './PaymentReview.css';
import {
  selectActivePPMForMove,
  loadPPMs,
  selectPPMEstimateRange,
  getPpmWeightEstimate,
  updatePPM,
} from 'shared/Entities/modules/ppms';

const nextBtnLabel = 'Submit Request';

class PaymentReview extends Component {
  state = {
    acceptTerms: false,
    moveSubmissionError: false,
  };

  componentDidMount() {
    const { originDutyStationZip, currentPPM, moveId } = this.props;
    const { original_move_date, pickup_postal_code } = currentPPM;

    this.props.loadPPMs(moveId).then(() => {
      if (!isEmpty(currentPPM)) {
        this.props.getMoveDocumentsForMove(moveId).then(({ obj: documents }) => {
          const weightTicketNetWeight = calcNetWeight(documents);
          const netWeight =
            weightTicketNetWeight > this.props.entitlement.sum ? this.props.entitlement.sum : weightTicketNetWeight;
          // TODO: make not async, make sure this happens

          this.props.getPpmWeightEstimate(
            original_move_date,
            pickup_postal_code,
            originDutyStationZip,
            this.props.orders.id,
            netWeight,
          );
        });
      }
    });
  }

  componentDidUpdate(prevProps) {
    const { originDutyStationZip, currentPPM, moveDocuments } = this.props;
    const { original_move_date, pickup_postal_code } = currentPPM;
    if (moveDocuments.weightTickets.length !== prevProps.moveDocuments.weightTickets.length) {
      if (!isEmpty(currentPPM)) {
        this.props.getMoveDocumentsForMove(this.props.moveId).then(({ obj: documents }) => {
          const weightTicketNetWeight = calcNetWeight(documents);
          const netWeight =
            weightTicketNetWeight > this.props.entitlement.sum ? this.props.entitlement.sum : weightTicketNetWeight;
          this.props.getPpmWeightEstimate(
            original_move_date,
            pickup_postal_code,
            originDutyStationZip,
            this.props.orders.id,
            netWeight,
          );
        });
      }
    }
  }

  handleOnAcceptTermsChange = (acceptTerms) => {
    this.setState({ acceptTerms });
  };

  submitCertificate = () => {
    const signatureTime = moment().format();
    const { currentPPM, moveId } = this.props;
    const certificate = {
      certification_text: ppmPaymentLegal,
      date: signatureTime,
      signature: 'CHECKBOX',
      personally_procured_move_id: currentPPM.id,
      certification_type: SIGNED_CERT_OPTIONS.PPM_PAYMENT,
    };
    return this.props.createSignedCertification(moveId, certificate);
  };

  applyClickHandlers = () => {
    this.setState({ moveSubmissionError: false }, () =>
      Promise.all([this.submitCertificate(), this.props.submitExpenseDocs()])
        .then(() => {
          // TODO: path may change to home after ppm integration with new home page
          this.props.history.push('/ppm');
        })
        .catch(() => {
          this.setState({ moveSubmissionError: true });
          scrollToTop();
        }),
    );
  };

  render() {
    const { moveId, moveDocuments, submitting, history, incentiveEstimateMin } = this.props;
    const weightTickets = moveDocuments.weightTickets;
    const missingSomeWeightTicket = weightTickets.some(
      ({ empty_weight_ticket_missing, full_weight_ticket_missing }) =>
        empty_weight_ticket_missing || full_weight_ticket_missing,
    );

    return (
      <div className="grid-container usa-prose">
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
            {missingSomeWeightTicket ? (
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
            ) : (
              <>
                <h4>You're requesting a payment of ${formatCents(incentiveEstimateMin)}</h4>
                <p>
                  Finance will determine your final reimbursement after reviewing the information you’ve submitted. That
                  final total will reflect the weight of your completed move (including any household goods move, if
                  applicable); any advances you requested and were given; and withheld taxes.
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
            finishLaterHandler={() => history.push('/')}
            submitButtonsAreDisabled={!this.state.acceptTerms}
            saveAndAddHandler={this.applyClickHandlers}
            submitting={submitting}
          />
        </div>
      </div>
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
    currentPPM: selectActivePPMForMove(state, moveId),
    incentiveEstimateMin: selectPPMEstimateRange(state).range_min,
    incentiveEstimateMax: selectPPMEstimateRange(state).range_max,
    originDutyStationZip: get(state, 'serviceMember.currentServiceMember.current_station.address.postal_code'),
    entitlement: loadEntitlementsFromState(state),
    orders: get(state, 'orders.currentOrders', {}),
  };
};

const mapDispatchToProps = {
  submitExpenseDocs,
  createSignedCertification,
  getMoveDocumentsForMove,
  getPpmWeightEstimate,
  loadPPMs,
  updatePPM,
};

export default connect(mapStateToProps, mapDispatchToProps)(PaymentReview);
