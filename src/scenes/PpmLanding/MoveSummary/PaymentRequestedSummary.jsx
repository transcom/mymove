import React from 'react';
import { ppmInfoPacket } from 'shared/constants';
import moment from 'moment';
import ppmCar from 'scenes/PpmLanding/images/ppm-car.svg';
import PPMStatusTimeline from 'scenes/PpmLanding/PPMStatusTimeline';
import FindWeightScales from 'scenes/PpmLanding/MoveSummary/FindWeightScales';
import PpmMoveDetails from 'scenes/PpmLanding/MoveSummary/SubmittedPpmMoveDetails';

const PaymentRequestedSummary = (props) => {
  const { ppm } = props;
  const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
  return (
    <div>
      <div className="shipment_box">
        <div className="shipment_type">
          <img className="move_sm" src={ppmCar} alt="ppm-car" />
          Handle your own move (PPM)
        </div>

        <div className="shipment_box_contents">
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
                  <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer" className="usa-link">
                    <button className="usa-button usa-button--secondary">Read PPM Info Packet</button>
                  </a>
                </div>
              )}
              <div className="step">
                <div className="title">Your payment is in review</div>
                <div>You will receive a notification from your destination PPPO office when it has been reviewed.</div>
              </div>
            </div>
            <div className="usa-width-one-third">
              <PpmMoveDetails ppm={ppm} />
              <div className="titled_block">
                <div className="title">Documents</div>
                <div className="details-links">
                  <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer" className="usa-link">
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
