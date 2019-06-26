import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { formatCents } from 'shared/formatters';
import WeightTicketListItem from './WeightTicketListItem';
import ExpenseTicketListItem from './ExpenseTicketListItem';
import WizardHeader from '../../WizardHeader';
import './PaymentReview.css';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';
import { ppmPaymentLegal } from 'scenes/Legalese/legaleseText';
import './PaymentReview.css';
import PPMPaymentRequestActionBtns from 'scenes/Moves/Ppm/PPMPaymentRequestActionBtns';

const nextBtnLabel = 'Submit Request';

class PaymentReview extends Component {
  state = {
    acceptTerms: false,
  };

  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.moveId);
  }

  handleOnAcceptTermsChange = acceptTerms => {
    this.setState({ acceptTerms });
  };

  getExpenses(expenses) {
    return expenses.map(expense => {
      return {
        id: expense.id,
        amount: formatCents(expense.requested_amount_cents),
        type: this.formatExpenseType(expense.moving_expense_type),
        paymentMethod: expense.payment_method,
      };
    });
  }

  formatExpenseType(expenseType) {
    if (typeof expenseType !== 'string') return '';
    let type = expenseType.toLowerCase().replace('_', ' ');
    return type.charAt(0).toUpperCase() + type.slice(1);
  }

  render() {
    const { moveId, moveDocuments, submitting } = this.props;
    const expenses = this.getExpenses(moveDocuments.expenses);
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
        <div className="usa-grid">
          <div className="review-payment-request-header">
            <h3>Review Payment Request</h3>
            <p>
              Make sure <strong>all</strong> your documents are uploaded.
            </p>
          </div>

          <div className="doc-summary-container">
            <h3>Document summary - {weightTickets.length + expenses.length} total</h3>
            <h4>{weightTickets.length} sets of weight tickets</h4>
            <div className="tickets">
              {weightTickets.map((ticket, index) => <WeightTicketListItem key={ticket.id} num={index} {...ticket} />)}
            </div>
            <Link data-cy="weight-ticket-link" to={`/moves/${moveId}/ppm-weight-ticket`}>
              <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add weight ticket
            </Link>
            <hr id="doc-summary-separator" />
            <h4>
              {expenses.length} expense{expenses.length > 1 ? 's' : ''}
            </h4>
            <div className="tickets">
              {expenses.map(expense => <ExpenseTicketListItem key={expense.id} {...expense} />)}
            </div>
            <div className="add-expense-link">
              <Link data-cy="expense-link" to={`/moves/${moveId}/ppm-expenses`}>
                <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add expense
              </Link>
            </div>
          </div>

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
  };
};

const mapDispatchToProps = {
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(PaymentReview);
