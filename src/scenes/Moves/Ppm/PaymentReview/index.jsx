import React, { Component } from 'react';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';

import WizardHeader from 'scenes/Moves/WizardHeader';
import { Link } from 'react-router-dom';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';
import { ppmPaymentLegal } from 'scenes/Legalese/legaleseText';
import './PaymentReview.css';
import PPMPaymentRequestActionBtns from 'scenes/Moves/Ppm/PPMPaymentRequestActionBtns';

const nextBtnLabel = 'Submit Request';

class PaymentReview extends Component {
  state = {
    acceptTerms: false,
  };
  handleOnAcceptTermsChange = acceptTerms => {
    this.setState({ acceptTerms });
  };
  render() {
    const { submitting } = this.props;
    const moveId = this.props.match.params.moveId;
    const weightTicketsPage = `/moves/${moveId}/ppm-weight-ticket`;
    const expensePage = `/moves/${moveId}/ppm-expenses`;
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
        <div className="usa-grid expenses-container">
          <h3 className="expenses-header">Review Payment Request</h3>
          <p>
            {' '}
            Make sure <strong>all</strong> your documents are uploaded.
          </p>
          <div>
            <h3 style={{ color: 'red' }}>DOCUMENT SUMMARY TBU</h3>
            <ul style={{ marginBottom: '30em' }}>
              <li>
                <Link to={weightTicketsPage} data-cy="weight-ticket-link">
                  Weight Ticket
                </Link>
              </li>
              <li>
                <Link to={expensePage} data-cy="expense-link">
                  Expenses
                </Link>
              </li>
            </ul>
          </div>
          <div className="usa-grid" style={{ marginBottom: '4.5em' }}>
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
            submitting={submitting}
            displaySaveForLater={true}
          />
        </div>
      </>
    );
  }
}

export default PaymentReview;
