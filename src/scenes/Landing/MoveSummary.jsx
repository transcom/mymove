import React from 'react';

import { get } from 'lodash';
import './MoveSummary.css';
import ppmCar from './images/ppm-car.svg';
import truck from 'shared/icon/truck-gray.svg';
import ppmDraft from './images/ppm-draft.png';
import ppmSubmitted from './images/ppm-submitted.png';
import ppmApproved from './images/ppm-approved.png';
import ppmInProgress from './images/ppm-in-progress.png';

const DutyStationContactInfo = props => {
  const { dutyStation, origin } = props;
  const stationName = get(dutyStation, 'name');
  if (!stationName) return <div />;
  return (
    <div className="titled_block">
      <a>{stationName}</a>
      <div className="Todo">
        {origin ? 'Origin' : 'Destination'} Transportation Office
      </div>
      <div className="Todo">PPPO</div>
      <div className="Todo">(210) 671-2821</div>
    </div>
  );
};

export const MoveSummary = props => {
  const {
    profile,
    move,
    orders,
    ppm,
    editMove,
    entitlement,
    resumeMove,
  } = props;
  const status = get(move, 'status', 'DRAFT');
  return (
    <div className="whole_box">
      <h2>
        {get(orders, 'new_duty_station.name', 'New move')} from{' '}
        {get(profile, 'current_station.name', '')}
      </h2>
      <div className="usa-width-three-fourths">
        {move && <div>Move Locator: {get(move, 'locator')}</div>}
        {entitlement && (
          <div>
            Weight Entitlement:{' '}
            <span>{entitlement.sum.toLocaleString()} lbs</span>
          </div>
        )}
        <div className="shipment_box">
          {status === 'DRAFT' && (
            <div className="shipment_type">
              <img className="move_sm" src={truck} alt="ppm-car" />
              Move to be scheduled
            </div>
          )}
          {status !== 'DRAFT' && (
            <div className="shipment_type">
              <img className="move_sm" src={ppmCar} alt="ppm-car" />
              Move your own stuff (PPM)
            </div>
          )}

          <div className="shipment_box_contents">
            {status === 'DRAFT' && (
              <div>
                <img src={ppmDraft} alt="status" />
                <div className="step-contents">
                  <div className="status_box usa-width-two-thirds">
                    <div className="step">
                      <div className="title">
                        Next Step: Finish setting up your move
                      </div>
                      <div>
                        Questions or need help? Contact your local
                        Transportation Office (PPPO) at{' '}
                        {get(profile, 'current_station.name')}.
                      </div>
                    </div>
                  </div>
                  <div className="usa-width-one-third">
                    <div className="titled_block">
                      <div className="title">Details</div>
                      <div>No detail</div>
                    </div>
                    <div className="titled_block">
                      <div className="title">Documents</div>
                      <div className="details-links">No documents</div>
                    </div>
                  </div>
                </div>
                <div className="step-links">
                  {status !== 'DRAFT' && (
                    <span>
                      <a
                        href="https://www.move.mil/resources/locator-maps"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Find Weight Scales
                      </a>
                    </span>
                  )}
                  {status === 'DRAFT' && (
                    <button onClick={resumeMove}>Continue Move Setup</button>
                  )}
                </div>
              </div>
            )}
            {/* Submitted Move */}
            {status === 'SUBMITTED' && (
              <div>
                <img src={ppmSubmitted} alt="status" />
                <div className="step-contents">
                  <div className="status_box usa-width-two-thirds">
                    <div className="step">
                      <div className="title">Next Step: Awaiting approval</div>
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
                    <a
                      href="https://www.move.mil/resources/locator-maps"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      Find Weight Scales
                    </a>
                  </span>
                </div>
              </div>
            )}
            {/* Approved Move */}
            {status === 'APPROVED' && (
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
            disabled={status === 'DRAFT'}
          >
            Edit Move
          </button>
        </div>

        <div className="contact_block">
          <div className="title">Contacts</div>
          <DutyStationContactInfo
            dutyStation={profile.current_station}
            origin
          />
          <DutyStationContactInfo
            dutyStation={get(orders, 'new_duty_station')}
          />
        </div>
      </div>
    </div>
  );
};
