import React from 'react';
import { ppmInfoPacket } from 'shared/constants';
import ppmCar from 'scenes/Landing/images/ppm-car.svg';
import PPMStatusTimeline from 'scenes/Landing/PPMStatusTimeline';
import FindWeightScales from 'scenes/Landing/MoveSummary/FindWeightScales';
import PpmMoveDetails from 'scenes/Landing/MoveSummary/SubmittedPpmMoveDetails';

const SubmittedPpmMoveSummary = props => {
  const { ppm } = props;
  return (
    <div className="shipment_box">
      <div className="shipment_type">
        <img className="move_sm" src={ppmCar} alt="ppm-car" />
        Move your own stuff (PPM)
      </div>
      <div className="shipment_box_contents">
        <PPMStatusTimeline ppm={ppm} />
        <div className="step-contents">
          <div className="status_box usa-width-two-thirds">
            <div className="step">
              <div className="title">Next Step: Wait for approval &amp; get ready</div>
              <div className="next-step">
                You'll be notified when your move is approved (up to 5 days). To get ready to move:
                <ul>
                  <li>
                    Go to{' '}
                    <a
                      href="https://move.mil/resources/locator-maps"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="usa-link"
                    >
                      certified weight scales
                    </a>{' '}
                    to get empty &amp; full weight tickets.
                  </li>
                  <li>Save expense receipts, including for storage.</li>
                </ul>
              </div>
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
        <a className="usa-button usa-button-secondary" href={ppmInfoPacket} target="_blank" rel="noopener noreferrer">
          Read PPM Info Sheet
        </a>
        <div className="step-links">
          <FindWeightScales />
        </div>
      </div>
    </div>
  );
};

export default SubmittedPpmMoveSummary;
