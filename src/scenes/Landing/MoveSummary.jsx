import React from 'react';

import './MoveSummary.css';
import ppmCar from './images/ppm-car.svg';
import ppmSubmitted from './images/ppm-submitted.png';
import ppmApproved from './images/ppm-approved.png';
import ppmInProgress from './images/ppm-in-progress.png';

const DutyStationContactInfo = props => {
  const { dutyStation, origin } = props;
  return (
    <div className="titled_block">
      <a>{dutyStation.name}</a>
      <div className="Todo">
        {origin ? 'Origin' : 'Destination'} Transportation Office
      </div>
      <div className="Todo">PPPO</div>
      <div className="Todo">(210) 671-2821</div>
    </div>
  );
};

export const MoveSummary = props => {
  const { profile, move, orders, ppm, editMove } = props;
  return (
    <div className="whole_box">
      <h2>
        {orders.new_duty_station.name} (from {profile.current_station.name})
      </h2>
      <div className="usa-width-three-fourths">
        <div>Move Locator: {move.locator}</div>
        <div>
          Weight Entitlement: <span className="Todo">10,500 lbs</span>
        </div>
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={ppmCar} alt="ppm-car" />
            Move your own stuff (PPM)
          </div>
          <div className="shipment_box_contents">
            {/* Submitted Move */}
            {move.status === 'SUBMITTED' && (
              <div>
                <img src={ppmSubmitted} alt="status" />
                <div className="step-contents">
                  <div className="status_box usa-width-two-thirds">
                    <div className="step">
                      <div className="title">Awaiting Approval</div>
                      <div>
                        Your shipment is awaiting approval. This can take up to
                        3 business days. Questions or need help? Contact your
                        local Transportation Office (PPPO) at{' '}
                        {profile.current_station.name}.
                      </div>
                    </div>
                  </div>
                  <div className="usa-width-one-third">
                    <div className="titled_block">
                      <div className="title">Details</div>
                      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
                      <div>Incentive (est.): {ppm.estimated_incentive}</div>
                    </div>
                    <div className="titled_block">
                      <div className="title">Documents</div>
                      <div className="details-links">
                        <a>PPM Info Packet</a>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="step-links">
                  <span>
                    <a>Request Storage</a> | <a>Find Weight Scales</a> |{' '}
                    <a>Report a Problem</a> | <a>Cancel Shipment</a>
                  </span>
                </div>
              </div>
            )}

            {/* Approved Move */}
            {move.status === 'APPROVED' && (
              <div>
                <img src={ppmApproved} alt="status" />
                <div className="step-contents">
                  <div className="status_box usa-width-two-thirds">
                    <div className="step">
                      <div className="title">Next step: Get ready to move</div>
                      <div>
                        Remember to save your weight tickets and expense
                        receipts. For more information, read the PPM info
                        packet.
                      </div>
                      <button className="usa-button-secondary">
                        Read PPM Info Packet
                      </button>
                    </div>
                    <div className="step">
                      <div className="title">Next step: Request Payment</div>
                      <div>
                        Request a PPM payment, a storage payment, or an advance
                        against your PPM payment before your move is done.
                      </div>
                      <button className="usa-button-secondary">
                        Request Payment
                      </button>
                    </div>
                  </div>
                  <div className="usa-width-one-third">
                    <div className="titled_block">
                      <div className="title">Details</div>
                      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
                      <div>Incentive (est.): {ppm.estimated_incentive}</div>
                    </div>
                    <div className="titled_block">
                      <div className="title">Documents</div>
                      <div className="details-links">
                        <a>PPM Info Packet</a>
                        <a>Advance paperwork</a>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="step-links">
                  <span>
                    <a>Request Storage</a> | <a>Find Weight Scales</a> |{' '}
                    <a>Report a Problem</a> | <a>Cancel Shipment</a>
                  </span>
                </div>
              </div>
            )}

            {/* In Progress Move */}
            {/* NOTE: The above blocks rely on move.status. This in progress block
                is unviewable until we start editing PPM statuses. */}
            {ppm.status === 'IN_PROGRESS' && (
              <div>
                <img src={ppmInProgress} alt="status" />
                <div className="step-contents">
                  <div className="status_box usa-width-two-thirds">
                    <div className="step">
                      <div className="title">Next step: Request payment</div>
                      <div>
                        Request a PPM payment, a storage payment, or an advance
                        against your PPM payment before your move is done.
                      </div>
                      <button className="usa-button-secondary">
                        Request Payment
                      </button>
                    </div>
                  </div>
                  <div className="usa-width-one-third">
                    <div className="titled_block">
                      <div className="title">Details</div>
                      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
                      <div>Incentive (est.): {ppm.estimated_incentive}</div>
                    </div>
                    <div className="titled_block">
                      <div className="title">Documents</div>
                      <div className="details-links">
                        <a>PPM Info Packet</a>
                        <a>Advance paperwork</a>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="step-links">
                  <span>
                    <a>Request Storage</a> | <a>Find Weight Scales</a> |{' '}
                    <a>Report a Problem</a> | <a>Cancel Shipment</a>
                  </span>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="sidebar usa-width-one-fourth">
        <div>
          <button
            className="usa-button-secondary"
            onClick={() => editMove(move)}
          >
            Edit Move Details
          </button>
        </div>
        <a>✚ Add Amended Orders</a>
        <hr />
        <a>✚ Add Shipment</a>
        <hr />
        <a>𝗫 Cancel Move</a>

        <div className="contact_block">
          <div className="title">Contacts</div>
          <DutyStationContactInfo
            dutyStation={profile.current_station}
            origin
          />
          <DutyStationContactInfo dutyStation={orders.new_duty_station} />
        </div>
      </div>
    </div>
  );
};
