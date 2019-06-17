import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import carImg from 'shared/images/car_mobile.png';
import boxTruckImg from 'shared/images/box_truck_mobile.png';
import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import './PaymentReview.css';
import WizardHeader from '../../WizardHeader';
const WEIGHT_TICKET_IMAGES = {
  CAR: carImg,
  BOX_TRUCK: boxTruckImg,
  CAR_TRAILER: carTrailerImg,
};

const weightTickets = [
  { id: 1, nickname: 'Moving truck', empty_weight: 2000, full_weight: 3000, type: 'BOX_TRUCK' },
  { id: 2, nickname: 'My Car', empty_weight: 2000, full_weight: 3000, type: 'CAR' },
  { id: 3, nickname: 'My Trailer', empty_weight: 2000, full_weight: 3000, type: 'CAR_TRAILER' },
];
const expenses = [
  { id: 1, title: 'Storage expense 1', amount: 336.18, type: 'Storage', paymentMethod: 'GTC' },
  { id: 2, title: 'Uhaul truck rental', amount: 632.24, type: 'Rental equipment', paymentMethod: 'GTC' },
  { id: 3, title: 'Texaco gas', amount: 106.35, type: 'Gas', paymentMethod: 'GTC' },
];

const WeightTicketListItem = ({ id, type, nickname, num, empty_weight, full_weight }) => (
  <div style={{ display: 'flex' }}>
    <div style={{ minWidth: 95 }}>
      <img src={WEIGHT_TICKET_IMAGES[type]} alt={type} />
    </div>
    <div style={{ flex: 1 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', maxWidth: 800 }}>
        <h4>
          {nickname} ({num + 1} set)
        </h4>
        <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />
      </div>
      <p>Empty weight ticket {empty_weight} lbs</p>
      <p>Full weight ticket {full_weight} lbs</p>
    </div>
  </div>
);

const ExpenseTicketListItem = ({ title, amount, type, paymentMethod }) => (
  <div>
    <div style={{ display: 'flex', justifyContent: 'space-between', maxWidth: 895 }}>
      <h4>
        {type} - ${amount}
      </h4>
      <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />
    </div>
    <div>
      {type} {paymentMethod}
    </div>
  </div>
);

class PaymentReview extends Component {
  render() {
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
          <h3>Review Payment Request</h3>
          <p>
            Make sure <strong>all</strong> your documents are uploaded.
          </p>

          <div className="doc-summary-container">
            <h3>Document summary - {weightTickets.length + expenses.length} total</h3>
            <h4>{weightTickets.length} sets of weight tickets</h4>
            {weightTickets.map((ticket, index) => <WeightTicketListItem key={ticket.id} num={index} {...ticket} />)}
            <Link to="">
              <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add weight ticket
            </Link>
            <hr />
            <h4>{expenses.length} expenses</h4>
            {expenses.map(expense => <ExpenseTicketListItem key={expense.id} {...expense} />)}
            <Link to="">
              <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add expense
            </Link>
          </div>

          <div>
            <h4>You're requesting a payment of $11,982.23</h4>
            <p>
              Finance will determine your final reimbursement after reviewing the information you’ve submitted. That
              final total will reflect the weight of your completed move (including any household goods move, if
              applicable); any advances you requested and were given; and withheld taxes.
            </p>
          </div>
        </div>
      </>
    );
  }
}

export default PaymentReview;
