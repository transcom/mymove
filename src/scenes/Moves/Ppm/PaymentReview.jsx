import React, { Component } from 'react';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import carImg from 'shared/images/car_mobile.png';
import boxTruckImg from 'shared/images/box_truck_mobile.png';
import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import './PaymentReview.css';
import WizardHeader from '../WizardHeader';
const WEIGHT_TICKET_IMAGES = {
  CAR: carImg,
  BOX_TRUCK: boxTruckImg,
  CAR_TRAILER: carTrailerImg,
};

const weightTickets = [
  { nickname: 'Moving truck', empty_weight: 2000, full_weight: 3000, type: 'BOX_TRUCK' },
  { nickname: 'My Car', empty_weight: 2000, full_weight: 3000, type: 'CAR' },
  { nickname: 'My Trailer', empty_weight: 2000, full_weight: 3000, type: 'CAR_TRAILER' },
];
const expenses = [{}, {}, {}];
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
            <h4>{weightTickets.length + expenses.length} sets of weight tickets</h4>
            {weightTickets.map((ticket, index) => (
              <div style={{ display: 'flex' }}>
                <div>
                  <img src={WEIGHT_TICKET_IMAGES[ticket.type]} alt={ticket.type} />
                </div>
                <div>
                  <div style={{ display: 'flex' }}>
                    <h4>
                      {ticket.nickname} ({index} set)
                    </h4>
                    <input alt="delete document button" type="image" src={deleteButtonImg} />
                  </div>
                  <p>Empty weight ticket {ticket.empty_weight}</p>
                  <p>Full weight ticket {ticket.full_weight}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </>
    );
  }
}

export default PaymentReview;
