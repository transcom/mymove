import React from 'react';
import { get } from 'lodash';
import truck from 'shared/icon/truck-gray.svg';

const CanceledMoveSummary = (props) => {
  const { profile, reviewProfile } = props;
  const currentStation = get(profile, 'current_location');
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');
  return (
    <div>
      <h1>New move</h1>
      <br />
      <div>
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={truck} alt="ppm-car" />
            Start here
          </div>

          <div className="shipment_box_contents">
            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                <div className="step">
                  <div>
                    Make sure you have a copy of your move orders before you get started. Questions or need to help?
                    Contact your local Transportation Office (PPPO) at {get(currentStation, 'name')}
                    {stationPhone ? ` at ${stationPhone}` : ''}.
                  </div>
                </div>
              </div>
            </div>
            <div className="step-links">
              <button className="usa-button" onClick={reviewProfile}>
                Start
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CanceledMoveSummary;
