import React from 'react';

import './MoveSummary.css';
import ppmCar from './images/ppm-car.svg';
import ppmStatus from './images/progress.png';

const DutyStationContactInfo = props => {
  const { dutyStation } = props;
  return (
    <div className="titled_block">
      <a>{dutyStation.name}</a>
      <div className="Todo">Origin Transportation Office</div>
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
        <div>
          Move Locator: <span className="Todo">{move.id}</span>
        </div>
        <div>
          Weight Entitlement: <span className="Todo">10,500 lbs</span>
        </div>
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="sm" src={ppmCar} alt="ppm-car" />
            Move your own stuff (PPM)
          </div>
          <div className="shipment_box_contents">
            <img src={ppmStatus} alt="status" />
            <div className="status_box usa-width-two-thirds">
              <div className="title">STATUS TEXT GOES HERE</div>
              <div>
                Your shipment is awaiting approval from the transporation
                office. This process can take up to 3 business days. If you have
                questions or need expedited processing contact contact your
                local Transportation Office (PPPO) at Lackland AFB at (210)
                671-2821.
              </div>
              <p>
                <a>Find Weight Scales</a> | <a>Report a Problem</a> |{' '}
                <a>Cancel Shipment</a>
              </p>
            </div>
            <div className="usa-width-one-third">
              <div className="titled_block">
                <div className="title">Details</div>
                <div>Incentive (est.): {ppm.estimated_incentive}</div>
              </div>
              <div className="titled_block">
                <div className="title">Documents</div>
                <div>
                  <a>PPM Info Packet</a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="sidebar usa-width-one-fourth">
        <div>
          <button onClick={() => editMove(move)}>Edit Move Details</button>
        </div>
        <a>✚ Add Amended Orders</a>
        <hr />
        <a>✚ Add Shipment</a>
        <hr />
        <a>𝗫 Cancel Move</a>

        <div className="contact_block">
          <div className="title">Contacts</div>
          <DutyStationContactInfo dutyStation={orders.new_duty_station} />
          <DutyStationContactInfo dutyStation={profile.current_station} />
        </div>
      </div>
    </div>
  );
};
