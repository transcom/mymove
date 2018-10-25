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
import StatusTimelineContainer from './StatusTimeline';

export const CanceledMoveSummary = props => {
  const { profile, reviewProfile } = props;
  const currentStation = get(profile, 'current_station');
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');
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
                    Make sure you have a copy of your move orders before you get started. Questions or need to help?
                    Contact your local Transportation Office (PPPO) at {get(currentStation, 'name')}
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
      <MoveInfoHeader orders={orders} profile={profile} move={move} entitlement={entitlement} />
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
                    <div className="title">Next Step: Finish setting up your move</div>
                    <div>
                      Questions or need help? Contact your local Transportation Office (PPPO) at{' '}
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
                <button onClick={resumeMove}>Continue Move Setup</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const SubmittedPpmMoveSummary = props => {
  const { ppm, orders, profile, move, entitlement } = props;
  return (
    <Fragment>
      <MoveInfoHeader orders={orders} profile={profile} move={move} entitlement={entitlement} />
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
                  <div
                  >{`Your shipment is awaiting approval. This can take up to 3 business days. Questions or need help? Contact your local Transportation Office (PPPO) at ${
                    profile.current_station.name
                  }.`}</div>
                </div>
              </div>
              <div className="usa-width-one-third">
                <PPMMoveDetails ppm={ppm} />
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
    </Fragment>
  );
};

const getTenDaysBookedDate = bookDate => {
  let businessDays = 0;
  let newDate;
  const bookDateMoment = moment(bookDate);
  while (businessDays < 10) {
    newDate = bookDateMoment.add(1, 'day');
    // Saturday => 6, Sunday => 7
    if (newDate.isoWeek() !== 6 && newDate.isoWeek() !== 7) {
      businessDays += 1;
    }
    console.log('new Date', businessDays, newDate.toString());
  }
  return newDate;
};

const showLandingPageText = shipment => {
  const today = moment();
  if (shipment.status === 'DELIVERED' || shipment.status === 'COMPLETED') {
    return (
      <div className="step">
        <div className="title">Next Step: Complete your customer satisfaction survey</div>
        <div>
          Tell us about your move experience. You have up to one year to complete your satisfaction survey. We use this
          information to decide which movers we allow to work with you.
        </div>
      </div>
    );
  } else if (today.isBefore(getTenDaysBookedDate(shipment.book_date), 'day')) {
    return (
      <div className="step">
        <div className="title">Next Step: Prepare for move</div>
        <div>
          Your mover will contact you within ten days to schedule a pre-move survey, where they will provide you with a
          detailed weight estimate and more accurate packing and delivery dates.
        </div>
      </div>
    );
  } else {
    return (
      <div className="step">
        <div className="title">Next step: Read pre-move tips</div>
        <div>
          Read the{' '}
          <a href="hhg-premove-tips" target="_blank">
            pre-move tips
          </a>{' '}
          documents, so you know what to expect and are prepared for your move.
        </div>
      </div>
    );
  }
};

export const SubmittedHhgMoveSummary = props => {
  const { shipment, orders, profile, move, entitlement } = props;
  let today = moment();
  return (
    <Fragment>
      <MoveInfoHeader orders={orders} profile={profile} move={move} entitlement={entitlement} />
      <br />
      <div className="usa-width-three-fourths">
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={truck} alt="hhg-truck" />
            Government Movers and Packers (HHG)
          </div>

          <div className="shipment_box_contents">
            <StatusTimelineContainer
              bookDate={shipment.book_date}
              moveDates={shipment.move_dates_summary}
              shipment={shipment}
            />
            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                {showLandingPageText(shipment)}
                {(shipment.actual_pack_date || today.isSameOrAfter(shipment.pm_survey_planned_pack_date)) && (
                  <div className="step">
                    <div className="title">File a Claim</div>
                    <div>
                      If you have household goods damaged or lost during the move, contact Able Movers Claims to file a
                      claim: (567) 980-4321. If, after attempting to work with them, you do not feel that you are
                      receiving adequate compensation, contact the Military Claims Office for help.
                    </div>
                  </div>
                )}
              </div>
              <div className="usa-width-one-third">
                <HhgMoveDetails hhg={shipment} />
                <div className="titled_block">
                  <div className="title">Documents</div>
                  <div className="details-links">
                    <a href="placeholder" target="_blank" rel="noopener noreferrer">
                      Pre-move tips
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const ApprovedMoveSummary = withContext(props => {
  const { ppm, orders, profile, move, entitlement, requestPaymentSuccess } = props;
  const paymentRequested = ppm.status === 'PAYMENT_REQUESTED';
  const moveInProgress = moment(ppm.planned_move_date, 'YYYY-MM-DD').isSameOrBefore();
  return (
    <Fragment>
      <MoveInfoHeader orders={orders} profile={profile} move={move} entitlement={entitlement} />
      <br />
      <div className="usa-width-three-fourths">
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
                      Request a PPM payment, a storage payment, or an advance against your PPM payment before your move
                      is done.
                    </div>
                    <Link to={`moves/${move.id}/request-payment`} className="usa-button usa-button-secondary">
                      Request Payment
                    </Link>
                  </div>
                )}
              </div>
              <div className="usa-width-one-third">
                <PPMMoveDetails ppm={ppm} />
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
    </Fragment>
  );
});

const PPMMoveDetails = props => {
  const { ppm } = props;
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString = ppm.has_requested_advance
    ? `Advance Requested: $${formatCents(ppm.advance.requested_amount)}`
    : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
      <div>Incentive (est.): {formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max)}</div>
      {ppm.has_sit && <div>{hasSitString}</div>}
      {ppm.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const HhgMoveDetails = props => {
  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {props.hhg.weight_estimate} lbs</div>
    </div>
  );
};

const FindWeightScales = () => (
  <span>
    <a href="https://www.move.mil/resources/locator-maps" target="_blank" rel="noopener noreferrer">
      Find Weight Scales
    </a>
  </span>
);

const MoveInfoHeader = props => {
  const { orders, profile, move, entitlement } = props;
  return (
    <Fragment>
      <h2>
        {get(orders, 'new_duty_station.name', 'New move')} (from {get(profile, 'current_station.name', '')})
      </h2>
      {get(move, 'locator') && <div>Move Locator: {get(move, 'locator')}</div>}
      {entitlement && (
        <div>
          Weight Entitlement: <span>{entitlement.sum.toLocaleString()} lbs</span>
        </div>
      )}
    </Fragment>
  );
};

const ppmSummaryStatusComponents = {
  DRAFT: DraftMoveSummary,
  SUBMITTED: SubmittedPpmMoveSummary,
  APPROVED: ApprovedMoveSummary,
  CANCELED: CanceledMoveSummary,
};

const hhgSummaryStatusComponents = {
  DRAFT: DraftMoveSummary,
  SUBMITTED: SubmittedHhgMoveSummary,
  APPROVED: ApprovedMoveSummary,
  CANCELED: CanceledMoveSummary,
  AWARDED: SubmittedHhgMoveSummary,
  ACCEPTED: SubmittedHhgMoveSummary,
  COMPLETED: SubmittedHhgMoveSummary,
};

const getStatus = (moveStatus, moveType, ppm, shipment) => {
  let status = 'DRAFT';
  if (moveType === 'PPM') {
    // assign the status
    const ppmStatus = get(ppm, 'status', 'DRAFT');
    status =
      moveStatus === 'APPROVED' && (ppmStatus === 'SUBMITTED' || ppmStatus === 'DRAFT') ? 'SUBMITTED' : moveStatus;
  } else if (moveType === 'HHG') {
    // assign the status
    const shipmentStatus = get(shipment, 'status', 'DRAFT');
    status = ['SUBMITTED', 'AWARDED', 'ACCEPTED', 'APPROVED'].includes(shipmentStatus) ? shipmentStatus : 'DRAFT';
  }
  return status;
};

export const MoveSummary = props => {
  const {
    profile,
    move,
    orders,
    ppm,
    shipment,
    editMove,
    entitlement,
    resumeMove,
    reviewProfile,
    requestPaymentSuccess,
  } = props;
  const moveStatus = get(move, 'status', 'DRAFT');
  const status = getStatus(moveStatus, move.selected_move_type, ppm, shipment);
  const StatusComponent =
    move.selected_move_type === 'PPM' ? ppmSummaryStatusComponents[status] : hhgSummaryStatusComponents[status]; // eslint-disable-line security/detect-object-injection
  return (
    <Fragment>
      {status === 'CANCELED' && (
        <Alert type="info" heading="Your move was canceled">
          Your move from {get(profile, 'current_station.name')} to {get(orders, 'new_duty_station.name')} with the move
          locator ID {get(move, 'locator')} was canceled.
        </Alert>
      )}

      <div className="whole_box">
        <StatusComponent
          className="status-component"
          ppm={ppm}
          shipment={shipment}
          orders={orders}
          profile={profile}
          move={move}
          entitlement={entitlement}
          resumeMove={resumeMove}
          reviewProfile={reviewProfile}
          requestPaymentSuccess={requestPaymentSuccess}
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
            <TransportationOfficeContactInfo dutyStation={profile.current_station} isOrigin={true} />
            {status !== 'CANCELED' && <TransportationOfficeContactInfo dutyStation={get(orders, 'new_duty_station')} />}
            {['ACCEPTED', 'APPROVED', 'IN_TRANSIT', 'COMPLETED'].includes(status) && (
              <div className="titled_block">
                <strong>TSP name</strong>
                <div>phone #</div>
              </div>
            )}
          </div>
        </div>
      </div>
    </Fragment>
  );
};
