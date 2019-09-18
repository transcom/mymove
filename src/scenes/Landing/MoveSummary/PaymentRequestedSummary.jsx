import React from 'react';
import { Link } from 'react-router-dom';
import { ppmInfoPacket } from 'shared/constants';
import Alert from 'shared/Alert';
import moment from 'moment';
import ppmCar from 'scenes/Landing/images/ppm-car.svg';
import PPMStatusTimeline from 'scenes/Landing/PPMStatusTimeline';
import FindWeightScales from 'scenes/Landing/MoveSummary/FindWeightScales';
import PpmMoveDetails from 'scenes/Landing/MoveSummary/SubmittedPpmMoveDetails';

const PaymentRequestedSummary = props => {
  const { ppm, move, requestPaymentSuccess } = props;
  const paymentRequested = ppm.status === 'PAYMENT_REQUESTED';
  const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
  const ppmPaymentRequestIntroRoute = `moves/${move.id}/request-payment`;
  return (
    <div>
      <div className="shipment_box">
        <div className="shipment_type">
          <img className="move_sm" src={ppmCar} alt="ppm-car" />
          Move your own stuff (PPM)
        </div>

        <div className="shipment_box_contents">
          {requestPaymentSuccess && (
            <Alert type="success" heading="">
              Payment request submitted
            </Alert>
          )}

          <PPMStatusTimeline ppm={ppm} />
          <div className="step-contents">
            <div className="status_box usa-width-two-thirds">
              {!moveInProgress && (
                <div className="step">
                  <div className="title">Next Step: Get ready to move</div>
                  <div>
                    Remember to save your weight tickets and expense receipts. For more information, read the PPM info
                    packet.
                  </div>
                  <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer">
                    <button className="usa-button-secondary">Read PPM Info Packet</button>
                  </a>
                </div>
              )}
              {paymentRequested ? (
                <div className="step">
                  <div className="title">Your payment is in review</div>
                  <div>
                    You will receive a notification from your destination PPPO office when it has been reviewed.
                  </div>
                </div>
              ) : (
                <div className="step">
                  <div className="title">Next Step: Request payment</div>
                  <div>
                    Request a PPM payment, a storage payment, or an advance against your PPM payment before your move is
                    done.
                  </div>
                  <Link to={ppmPaymentRequestIntroRoute} className="usa-button usa-button-secondary">
                    Request Payment
                  </Link>
                </div>
              )}
            </div>
            <div className="usa-width-one-third">
              <PpmMoveDetails ppm={ppm} />
              <div className="titled_block">
                <div className="title">Documents</div>
                <div className="details-links">
                  <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer">
                    PPM Info Packet
                  </a>
                </div>
              </div>
            </div>
          </div>
          <div className="step-links">
            <FindWeightScales />
          </div>
        </div>
      </div>
    </div>
  );
};

export default PaymentRequestedSummary;
