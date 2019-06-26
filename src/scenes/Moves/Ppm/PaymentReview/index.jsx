import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import carImg from 'shared/images/car_mobile.png';
import boxTruckImg from 'shared/images/box_truck_mobile.png';
import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { connect } from 'react-redux';
import {
  selectAllDocumentsForMove,
  selectExpenseTicketsFromDocuments,
  selectWeightTicketsFromDocuments,
} from 'shared/Entities/modules/moveDocuments';
import { formatCents } from 'shared/formatters';
import { intToOrdinal } from '../utility';
import PPMPaymentRequestActionBtns from '../PPMPaymentRequestActionBtns';
import WizardHeader from '../../WizardHeader';
import './PaymentReview.css';

const WEIGHT_TICKET_IMAGES = {
  CAR: carImg,
  BOX_TRUCK: boxTruckImg,
  CAR_TRAILER: carTrailerImg,
};

const MissingLabel = ({ children }) => (
  <p className="missing-doc">
    <em>{children}</em> <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
  </p>
);
const WeightTicketListItem = ({
  vehicle_options,
  vehicle_nickname,
  num,
  empty_weight,
  full_weight,
  empty_weight_ticket_missing,
  full_weight_ticket_missing,
  trailer_ownership_missing,
}) => (
  <div className="ticket-item" style={{ display: 'flex' }}>
    <div style={{ minWidth: 95 }}>
      {/*eslint-disable security/detect-object-injection*/}
      <img className="weight-ticket-image" src={WEIGHT_TICKET_IMAGES[vehicle_options]} alt={vehicle_options} />
    </div>
    <div style={{ flex: 1 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', maxWidth: 820 }}>
        <h4>
          {vehicle_nickname} ({intToOrdinal(num + 1)} set)
        </h4>
        <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />
      </div>
      {empty_weight_ticket_missing ? (
        <MissingLabel>Missing empty weight ticket</MissingLabel>
      ) : (
        <p>Empty weight ticket {empty_weight} lbs</p>
      )}
      {full_weight_ticket_missing ? (
        <MissingLabel>Missing full weight ticket</MissingLabel>
      ) : (
        <p>Full weight ticket {full_weight} lbs</p>
      )}
      {vehicle_options === 'CAR_TRAILER' &&
        trailer_ownership_missing && <MissingLabel>Missing ownership documentation</MissingLabel>}
      {vehicle_options === 'CAR_TRAILER' && !trailer_ownership_missing && <p>Ownership documentation</p>}
    </div>
  </div>
);

const ExpenseTicketListItem = ({ amount, type, paymentMethod }) => (
  <div className="ticket-item">
    <div style={{ display: 'flex', justifyContent: 'space-between', maxWidth: 916 }}>
      <h4>
        {type} - ${amount}
      </h4>
      <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />
    </div>
    <div>
      {type} ({paymentMethod})
    </div>
  </div>
);

class PaymentReview extends Component {
  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.moveId);
  }

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
    const expenses = this.getExpenses(this.props.moveDocuments.expenses);
    const weightTickets = this.props.moveDocuments.weightTickets;

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
            <Link to="">
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
              <Link to="">
                <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add expense
              </Link>
            </div>
          </div>

          <div className="doc-review">
            <h4>You're requesting a payment of $11,982.23</h4>
            <p>
              Finance will determine your final reimbursement after reviewing the information you’ve submitted. That
              final total will reflect the weight of your completed move (including any household goods move, if
              applicable); any advances you requested and were given; and withheld taxes.
            </p>
          </div>
          <PPMPaymentRequestActionBtns
            nextBtnLabel="Submit Request"
            submitButtonsAreDisabled
            submitting
            saveForLaterHandler={() => {}}
            saveAndAddHandler={() => {}}
            displaySaveForLater
          />
        </div>
      </>
    );
  }
}

const mapStateToProps = (state, props) => {
  const { moveId } = props.match.params;
  const documents = selectAllDocumentsForMove(state, moveId);

  return {
    moveDocuments: {
      expenses: selectExpenseTicketsFromDocuments(documents),
      weightTickets: selectWeightTicketsFromDocuments(documents),
    },
    moveId,
  };
};

const mapDispatchToProps = {
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(PaymentReview);
