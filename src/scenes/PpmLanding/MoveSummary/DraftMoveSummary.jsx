import React from 'react';
import { ProfileStatusTimeline } from 'scenes/PpmLanding/StatusTimeline';
import truck from 'shared/icon/truck-gray.svg';

const DraftMoveSummary = (props) => {
  const { profile, resumeMove } = props;
  return (
    <div>
      <div className="shipment_box">
        <div className="shipment_type">
          <img className="move_sm" src={truck} alt="ppm-car" />
          Move to be scheduled
        </div>

        <div className="shipment_box_contents">
          <div>
            <ProfileStatusTimeline profile={profile} />
            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                <div className="step">
                  <div className="title">Next Step: Finish setting up your move</div>
                  <div>
                    Questions or need help? Contact your local Transportation Office (PPPO) at{' '}
                    {profile.current_location.name}.
                  </div>
                </div>
              </div>
              <div className="usa-width-one-third">
                <div className="titled_block">
                  <div className="title">Details</div>
                  <div>No details</div>
                </div>
                <div className="titled_block">
                  <div className="title">Documents</div>
                  <div className="details-links">No documents</div>
                </div>
              </div>
            </div>
            <div className="step-links">
              <button className="usa-button" onClick={resumeMove}>
                Continue Move Setup
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DraftMoveSummary;
