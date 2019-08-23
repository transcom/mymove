import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get, includes, isEmpty } from 'lodash';
import moment from 'moment';

import { ppmInfoPacket } from 'shared/constants';
import Alert from 'shared/Alert';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents, formatCentsRange } from 'shared/formatters';
import TransportationOfficeContactInfo from 'shared/TransportationOffices/TransportationOfficeContactInfo';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import truck from 'shared/icon/truck-gray.svg';
import { selectReimbursement } from 'shared/Entities/modules/ppms';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { calcNetWeight } from 'scenes/Moves/Ppm/utility';
import { getPpmWeightEstimate } from 'scenes/Moves/Ppm/ducks';
import ppmCar from './images/ppm-car.svg';
import { ProfileStatusTimeline } from './StatusTimeline';
import PPMStatusTimeline from './PPMStatusTimeline';

import './MoveSummary.css';

export const CanceledMoveSummary = props => {
  const { profile, reviewProfile } = props;
  const currentStation = get(profile, 'current_station');
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');
  return (
    <Fragment>
      <h2>New move</h2>
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
              <button onClick={reviewProfile}>Start</button>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const DraftMoveSummary = props => {
  const { profile, resumeMove } = props;
  return (
    <Fragment>
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
                      {get(profile, 'current_station.name')}.
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
                <button onClick={resumeMove}>Continue Move Setup</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export const PPMAlert = props => {
  return (
    <Alert type="success" heading={props.heading}>
      Next, wait for approval. Once approved:
      <br />
      <ul>
        <li>
          Get certified <strong>weight tickets</strong>, both empty &amp; full
        </li>
        <li>
          Save <strong>expense receipts</strong>, including for storage
        </li>
        <li>
          Read the{' '}
          <strong>
            <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer">
              PPM info sheet
            </a>
          </strong>{' '}
          for more info
        </li>
      </ul>
    </Alert>
  );
};

export const SubmittedPpmMoveSummary = props => {
  const { ppm } = props;
  return (
    <Fragment>
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
                      <a href="https://move.mil/resources/locator-maps" target="_blank" rel="noopener noreferrer">
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
          <a className="usa-button usa-button-secondary" href={ppmInfoPacket} target="_blank" rel="noopener noreferrer">
            Read PPM Info Sheet
          </a>
          <div className="step-links">
            <FindWeightScales />
          </div>
        </div>
      </div>
    </Fragment>
  );
};

//TODO remove redundant ApprovedMoveSummary component w/ ppmPaymentRequest flag
const NewApprovedMoveSummaryComponent = ({
  ppm,
  move,
  weightTicketSets,
  isMissingWeightTicketDocuments,
  incentiveEstimate,
}) => {
  const paymentRequested = ppm.status === 'PAYMENT_REQUESTED';
  const ppmPaymentRequestIntroRoute = `moves/${move.id}/ppm-payment-request-intro`;
  const ppmPaymentRequestReviewRoute = `moves/${move.id}/ppm-payment-review`;
  return (
    <Fragment>
      <div>
        <div className="shipment_box">
          <div className="shipment_type">
            <img className="move_sm" src={ppmCar} alt="ppm-car" />
            Move your own stuff (PPM)
          </div>

          <div className="shipment_box_contents">
            <PPMStatusTimeline ppm={ppm} />
            <div className="step-contents">
              <div className="status_box usa-width-two-thirds">
                {paymentRequested ? (
                  isMissingWeightTicketDocuments ? (
                    <div className="step">
                      <div className="title">Next step: Contact the PPPO office</div>
                      <div>
                        You will need to go into the PPPO office in order to take care of your missing weight ticket.
                      </div>
                      <Link
                        data-cy="edit-payment-request"
                        to={ppmPaymentRequestReviewRoute}
                        className="usa-button usa-button-secondary"
                      >
                        Edit Payment Request
                      </Link>
                    </div>
                  ) : (
                    <div className="step">
                      <div className="title">What's next?</div>
                      <div>
                        We'll email you a link so you can see and download your final payment paperwork.
                        <br />
                        <br />
                        We've also sent your paperwork to Finance. They'll review it, determine a final amount, then
                        send your payment.
                      </div>
                      <Link
                        data-cy="edit-payment-request"
                        to={ppmPaymentRequestReviewRoute}
                        className="usa-button usa-button-secondary"
                      >
                        Edit Payment Request
                      </Link>
                    </div>
                  )
                ) : (
                  <div className="step">
                    {weightTicketSets.length ? (
                      <>
                        <div className="title">Next Step: Finish requesting payment</div>
                        <div>
                          Continue uploading your weight tickets and expense to get paid after your move is done.
                        </div>
                        <Link to={ppmPaymentRequestReviewRoute} className="usa-button usa-button-secondary">
                          Continue Requesting Payment
                        </Link>
                      </>
                    ) : (
                      <>
                        <div className="title">Next Step: Request payment</div>
                        <div>
                          Request a PPM payment, a storage payment, or an advance against your PPM payment before your
                          move is done.
                        </div>
                        <Link to={ppmPaymentRequestIntroRoute} className="usa-button usa-button-secondary">
                          Request Payment
                        </Link>
                      </>
                    )}
                  </div>
                )}
              </div>
              <div className="usa-width-one-third">
                <NewPPMMoveDetails ppm={ppm} />
              </div>
            </div>
            <div className="step-links" />
          </div>
        </div>
      </div>
    </Fragment>
  );
};

const mapStateToNewApprovedMoveSummaryProps = (state, { move }) => ({
  weightTicketSets: selectPPMCloseoutDocumentsForMove(state, move.id, ['WEIGHT_TICKET_SET']),
  incentiveEstimate: get(state, 'ppm.incentive_estimate_min'),
});

const NewApprovedMoveSummary = connect(mapStateToNewApprovedMoveSummaryProps)(NewApprovedMoveSummaryComponent);

export const ApprovedMoveSummary = props => {
  const { ppm, move, requestPaymentSuccess } = props;
  const paymentRequested = ppm.status === 'PAYMENT_REQUESTED';
  const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
  const ppmPaymentRequestIntroRoute = `moves/${move.id}/request-payment`;
  return (
    <Fragment>
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
                      Request a PPM payment, a storage payment, or an advance against your PPM payment before your move
                      is done.
                    </div>
                    <Link to={ppmPaymentRequestIntroRoute} className="usa-button usa-button-secondary">
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
};

//TODO remove redundant PPMMoveDetailsPanel component w/ ppmPaymentRequest flag
const NewPPMMoveDetailsPanel = ({ advance, ppm, isMissingWeightTicketDocuments }) => {
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString =
    ppm.has_requested_advance && advance && advance.requested_amount
      ? `Advance Requested: $${formatCents(advance.requested_amount)}`
      : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.currentPpm.weight_estimate} lbs</div>
      <div className="title" style={{ paddingTop: '0.5em' }}>
        Payment request
      </div>
      <div>Estimated payment: </div>
      {ppm.incentive_estimate_min ? (
        isMissingWeightTicketDocuments ? (
          <>
            <div className="missing-label">
              Unknown
              <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
            </div>
            <div style={{ fontSize: '0.90em', color: '#767676' }}>
              <em>Estimated payment will be given after resolving missing weight tickets.</em>
            </div>
          </>
        ) : (
          <>
            <div>${formatCents(ppm.incentive_estimate_min)}</div>
            <div style={{ fontSize: '0.90em', color: '#767676' }}>
              <em>Actual payment may vary, subject to Finance review.</em>
            </div>
          </>
        )
      ) : (
        <>
          Not ready yet{' '}
          <IconWithTooltip toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive." />
        </>
      )}

      {ppm.has_sit && <div>{hasSitString}</div>}
      {ppm.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const PPMMoveDetailsPanel = props => {
  const { advance, ppm } = props;
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString = ppm.has_requested_advance ? `Advance Requested: $${formatCents(advance.requested_amount)}` : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.currentPpm.weight_estimate} lbs</div>
      <div>
        Incentive (est.):{' '}
        {formatCentsRange(ppm.currentPpm.incentive_estimate_min, ppm.currentPpm.incentive_estimate_max)}
      </div>
      {ppm.has_sit && <div>{hasSitString}</div>}
      {ppm.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const mapStateToPPMMoveDetailsProps = (state, ownProps) => {
  const advance = selectReimbursement(state, ownProps.ppm.advance);
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, ownProps.ppm.move_id, [
    'WEIGHT_TICKET_SET',
  ]).some(doc => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);
  return { ppm: get(state, 'ppm', {}), advance, isMissingWeightTicketDocuments };
};

// TODO remove redundant function when remove ppmPaymentRequest flag
const PPMMoveDetails = connect(mapStateToPPMMoveDetailsProps)(PPMMoveDetailsPanel);
const NewPPMMoveDetails = connect(mapStateToPPMMoveDetailsProps)(NewPPMMoveDetailsPanel);

const FindWeightScales = () => (
  <span>
    <a href="https://www.move.mil/resources/locator-maps" target="_blank" rel="noopener noreferrer">
      Find Certified Weight Scales
    </a>
  </span>
);

const MoveInfoHeader = props => {
  const { orders, profile, move, entitlement, requestPaymentSuccess } = props;
  return (
    <Fragment>
      {requestPaymentSuccess && <Alert type="success" heading="Payment request submitted" />}

      <h2 className="move-summary-header">
        {get(orders, 'new_duty_station.name', 'New move')} (from {get(profile, 'current_station.name', '')})
      </h2>
      {get(move, 'locator') && <div>Move Locator: {get(move, 'locator')}</div>}
      {!isEmpty(entitlement) && (
        <div>
          Weight Entitlement: <span>{entitlement.sum.toLocaleString()} lbs</span>
        </div>
      )}
    </Fragment>
  );
};

// TODO revert this function to a constant when remove ppmPaymentRequest flag
const genPpmSummaryStatusComponents = context => {
  if (context && context.flags && context.flags.ppmPaymentRequest) {
    return {
      DRAFT: DraftMoveSummary,
      SUBMITTED: SubmittedPpmMoveSummary,
      APPROVED: NewApprovedMoveSummary,
      CANCELED: CanceledMoveSummary,
      PAYMENT_REQUESTED: ApprovedMoveSummary,
    };
  }
  return {
    DRAFT: DraftMoveSummary,
    SUBMITTED: SubmittedPpmMoveSummary,
    APPROVED: ApprovedMoveSummary,
    CANCELED: CanceledMoveSummary,
    PAYMENT_REQUESTED: ApprovedMoveSummary,
  };
};

const getPPMStatus = (moveStatus, ppm) => {
  // PPM status determination
  const ppmStatus = get(ppm, 'status', 'DRAFT');
  return moveStatus === 'APPROVED' && (ppmStatus === 'SUBMITTED' || ppmStatus === 'DRAFT') ? 'SUBMITTED' : moveStatus;
};

export class MoveSummaryComponent extends React.Component {
  componentDidMount() {
    this.props.getMoveDocumentsForMove(this.props.move.id).then(({ obj: documents }) => {
      const weightTicketNetWeight = calcNetWeight(documents);
      const netWeight =
        weightTicketNetWeight > this.props.entitlement.sum ? this.props.entitlement.sum : weightTicketNetWeight;
      this.props.getPpmWeightEstimate(
        this.props.ppm.actual_move_date || this.props.ppm.original_move_date,
        this.props.ppm.pickup_postal_code,
        this.props.originDutyStationZip,
        this.props.ppm.destination_postal_code,
        netWeight,
      );
    });
  }
  render() {
    const {
      context,
      profile,
      move,
      orders,
      ppm,
      editMove,
      entitlement,
      resumeMove,
      reviewProfile,
      requestPaymentSuccess,
      isMissingWeightTicketDocuments,
    } = this.props;
    const moveStatus = get(move, 'status', 'DRAFT');
    const PPMComponent = genPpmSummaryStatusComponents(context)[getPPMStatus(moveStatus, ppm)];
    return (
      <div className="move-summary">
        {move.status === 'CANCELED' && (
          <Alert type="info" heading="Your move was canceled">
            Your move from {get(profile, 'current_station.name')} to {get(orders, 'new_duty_station.name')} with the
            move locator ID {get(move, 'locator')} was canceled.
          </Alert>
        )}

        <div className="whole_box">
          {move.status !== 'CANCELED' && (
            <div>
              <MoveInfoHeader
                orders={orders}
                profile={profile}
                move={move}
                entitlement={entitlement}
                requestPaymentSuccess={requestPaymentSuccess}
              />
              <br />
            </div>
          )}
          {isMissingWeightTicketDocuments && ppm.status === 'PAYMENT_REQUESTED' && (
            <Alert type="warning" heading="Payment request is missing info">
              You will need to contact your local PPPO office to resolve your missing weight ticket.
            </Alert>
          )}
          <div className="usa-width-three-fourths">
            <PPMComponent
              context={context}
              className="status-component"
              ppm={ppm}
              orders={orders}
              profile={profile}
              move={move}
              entitlement={entitlement}
              resumeMove={resumeMove}
              reviewProfile={reviewProfile}
              requestPaymentSuccess={requestPaymentSuccess}
              isMissingWeightTicketDocuments={isMissingWeightTicketDocuments}
            />
          </div>

          <div className="sidebar usa-width-one-fourth">
            <div>
              <button
                className="usa-button-secondary"
                onClick={() => editMove(move)}
                disabled={includes(['DRAFT', 'CANCELED'], move.status)}
              >
                Edit Move
              </button>
            </div>
            <div className="contact_block">
              <div className="title">Contacts</div>
              <TransportationOfficeContactInfo dutyStation={profile.current_station} isOrigin={true} />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, ownProps.move.id, [
    'WEIGHT_TICKET_SET',
  ]).some(doc => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);

  return {
    isMissingWeightTicketDocuments,
    originDutyStationZip: get(state, 'serviceMember.currentServiceMember.current_station.address.postal_code'),
  };
}

const mapDispatchToProps = {
  getMoveDocumentsForMove,
  getPpmWeightEstimate,
};
export const MoveSummary = connect(
  mapStateToProps,
  mapDispatchToProps,
)(MoveSummaryComponent);
