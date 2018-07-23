import React, { Fragment } from 'react';

import { get, includes } from 'lodash';
import moment from 'moment';

import TransportationOfficeContactInfo from 'shared/TransportationOffices/TransportationOfficeContactInfo';
import './MoveSummary.css';
import ppmCar from './images/ppm-car.svg';
import truck from 'shared/icon/truck-gray.svg';
import ppmDraft from './images/ppm-draft.png';
import ppmSubmitted from './images/ppm-submitted.png';
import ppmApproved from './images/ppm-approved.png';
import ppmInProgress from './images/ppm-in-progress.png';
import { ppmInfoPacket } from 'shared/constants';
import Alert from 'shared/Alert';
import { formatCents, formatCentsRange } from 'shared/formatters';
import { Link } from 'react-router-dom';
import { withContext } from 'shared/AppContext';

export const CanceledMoveSummary = props => {
  const { profile, reviewProfile } = props;
  const currentStation = get(profile, 'current_station');
  const stationPhone = get(
    currentStation,
    'transportation_office.phone_lines.0',
  );
  return (
    <Fragment>
      <h2>New move</h2>
      <br />
      <div className="usa-width-three-fourths">
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
                    Make sure you have a copy of your move orders before you get
                    started. Questions or need to help? Contact your local
                    Transportation Office (PPPO) at{' '}
                    {get(currentStation, 'name')}
                    {stationPhone ? ` at ${stationPhone}` : ''}.
                  </div>
                </div>
              </div>
            </div>
            <div className="step-links">
              <button onClick={reviewProfile}>Start</button>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const DraftMoveSummary = props => {
  const { orders, profile, move, entitlement, resumeMove } = props;
  return (
    <Fragment>
      <MoveInfoHeader
        orders={orders}
        profile={profile}
        move={move}
        entitlement={entitlement}
      />
      <br />
      <div className="usa-width-three-fourths">
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={truck} alt="ppm-car" />
            Move to be scheduled
          </div>

          <div className="shipment_box_contents">
            <div>
              <img className="status_icon" src={ppmDraft} alt="status" />
              <div className="step-contents">
                <div className="status_box usa-width-two-thirds">
                  <div className="step">
                    <div className="title">
                      Next Step: Finish setting up your move
                    </div>
                    <div>
                      Questions or need help? Contact your local Transportation
                      Office (PPPO) at {get(profile, 'current_station.name')}.
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
                <button onClick={resumeMove}>Continue Move Setup</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const SubmittedMoveSummary = props => {
  const { ppm, orders, profile, move, entitlement } = props;
  return (
    <Fragment>
      <MoveInfoHeader
        orders={orders}
        profile={profile}
        move={move}
        entitlement={entitlement}
      />
      <br />
      <div className="usa-width-three-fourths">
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={ppmCar} alt="ppm-car" />
            Move your own stuff (PPM)
          </div>

          <div className="shipment_box_contents">
            <img className="status_icon" src={ppmSubmitted} alt="status" />
            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                <div className="step">
                  <div className="title">Next Step: Awaiting approval</div>
                  <div>
                    Your shipment is awaiting approval. This can take up to 3
                    business days. Questions or need help? Contact your local
                    Transportation Office (PPPO) at{' '}
                    {profile.current_station.name}.
                  </div>
                </div>
              </div>
              <div className="usa-width-one-third">
                <MoveDetails ppm={ppm} />
                <div className="titled_block">
                  <div className="title">Documents</div>
                  <div className="details-links">
                    <a
                      href={ppmInfoPacket}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
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
    </Fragment>
  );
};

export const ApprovedMoveSummary = withContext(props => {
  const { ppm, orders, profile, move, entitlement } = props;
  const canRequestPayment = props.context.flags.paymentRequest;
  const moveInProgress = moment(
    ppm.planned_move_date,
    'YYYY-MM-DD',
  ).isSameOrBefore();
  return (
    <Fragment>
      <MoveInfoHeader
        orders={orders}
        profile={profile}
        move={move}
        entitlement={entitlement}
      />
      <br />
      <div className="usa-width-three-fourths">
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={ppmCar} alt="ppm-car" />
            Move your own stuff (PPM)
          </div>

          <div className="shipment_box_contents">
            {moveInProgress ? (
              <img className="status_icon" src={ppmInProgress} alt="status" />
            ) : (
              <img className="status_icon" src={ppmApproved} alt="status" />
            )}

            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                {!moveInProgress && (
                  <div className="step">
                    <div className="title">Next Step: Get ready to move</div>
                    <div>
                      Remember to save your weight tickets and expense receipts.
                      For more information, read the PPM info packet.
                    </div>
                    <a
                      href={ppmInfoPacket}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <button className="usa-button-secondary">
                        Read PPM Info Packet
                      </button>
                    </a>
                  </div>
                )}
                <div className="step">
                  <div className="title">Next Step: Request payment</div>
                  <div>
                    Request a PPM payment, a storage payment, or an advance
                    against your PPM payment before your move is done.
                  </div>
                  {canRequestPayment && (
                    <Link
                      to={`moves/${move.id}/request-payment`}
                      className="usa-button usa-button-secondary"
                    >
                      Request Payment
                    </Link>
                  )}
                  {!canRequestPayment && (
                    <button className="usa-button-secondary" disabled={true}>
                      Request Payment - Coming Soon!
                    </button>
                  )}
                </div>
              </div>
              <div className="usa-width-one-third">
                <MoveDetails ppm={ppm} />
                <div className="titled_block">
                  <div className="title">Documents</div>
                  <div className="details-links">
                    <a
                      href={ppmInfoPacket}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
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
    </Fragment>
  );
});

const MoveDetails = props => {
  const { ppm } = props;
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString = ppm.has_requested_advance
    ? `Advance Requested: $${formatCents(ppm.advance.requested_amount)}`
    : '';
  const hasSitString = `Temp. Storage: ${
    ppm.days_in_storage
  } days ${privateStorageString}`;

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
      <div>
        Incentive (est.):{' '}
        {formatCentsRange(
          ppm.incentive_estimate_min,
          ppm.incentive_estimate_max,
        )}
      </div>
      {ppm.has_sit && <div>{hasSitString}</div>}
      {ppm.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const FindWeightScales = () => (
  <span>
    <a
      href="https://www.move.mil/resources/locator-maps"
      target="_blank"
      rel="noopener noreferrer"
    >
      Find Weight Scales
    </a>
  </span>
);

const MoveInfoHeader = props => {
  const { orders, profile, move, entitlement } = props;
  return (
    <Fragment>
      <h2>
        {get(orders, 'new_duty_station.name', 'New move')} from{' '}
        {get(profile, 'current_station.name', '')}
      </h2>
      {get(move, 'locator') && <div>Move Locator: {get(move, 'locator')}</div>}
      {entitlement && (
        <div>
          Weight Entitlement:{' '}
          <span>{entitlement.sum.toLocaleString()} lbs</span>
        </div>
      )}
    </Fragment>
  );
};

const moveSummaryStatusComponents = {
  DRAFT: DraftMoveSummary,
  SUBMITTED: SubmittedMoveSummary,
  APPROVED: ApprovedMoveSummary,
  CANCELED: CanceledMoveSummary,
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
    reviewProfile,
  } = props;
  const move_status = get(move, 'status', 'DRAFT');
  const ppm_status = get(ppm, 'status', 'DRAFT');
  const status =
    move_status === 'APPROVED' && ppm_status === 'SUBMITTED'
      ? ppm_status
      : move_status;
  const StatusComponent = moveSummaryStatusComponents[status]; // eslint-disable-line security/detect-object-injection
  return (
    <Fragment>
      {status === 'CANCELED' && (
        <Alert type="info" heading="Your move was canceled">
          Your move from {get(profile, 'current_station.name')} to{' '}
          {get(orders, 'new_duty_station.name')} with the move locator ID{' '}
          {get(move, 'locator')} was canceled.
        </Alert>
      )}

      <div className="whole_box">
        <StatusComponent
          className="status-component"
          ppm={ppm}
          orders={orders}
          profile={profile}
          move={move}
          entitlement={entitlement}
          resumeMove={resumeMove}
          reviewProfile={reviewProfile}
        />

        <div className="sidebar usa-width-one-fourth">
          <button
            className="usa-button-secondary"
            onClick={() => editMove(move)}
            disabled={includes(['DRAFT', 'CANCELED'], status)}
          >
            Edit Move
          </button>

          <div className="contact_block">
            <div className="title">Contacts</div>
            <TransportationOfficeContactInfo
              dutyStation={profile.current_station}
              isOrigin={true}
            />
            {status !== 'CANCELED' && (
              <TransportationOfficeContactInfo
                dutyStation={get(orders, 'new_duty_station')}
              />
            )}
          </div>
        </div>
      </div>
    </Fragment>
  );
};
