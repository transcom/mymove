import React from 'react';
import { connect } from 'react-redux';
import { get, includes, isEmpty } from 'lodash';

import Alert from 'shared/Alert';
import TransportationOfficeContactInfo from 'shared/TransportationOffices/TransportationOfficeContactInfo';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { calcNetWeight } from 'scenes/Moves/Ppm/utility';
import ApprovedMoveSummary from 'scenes/PpmLanding/MoveSummary/ApprovedMoveSummary';
import CanceledMoveSummary from 'scenes/PpmLanding/MoveSummary/CanceledMoveSummary';
import DraftMoveSummary from 'scenes/PpmLanding/MoveSummary/DraftMoveSummary';
import PaymentRequestedSummary from 'scenes/PpmLanding/MoveSummary/PaymentRequestedSummary';
import SubmittedPpmMoveSummary from 'scenes/PpmLanding/MoveSummary/SubmittedPpmMoveSummary';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { calculatePPMEstimate } from 'services/internalApi';
import { updatePPMEstimate } from 'store/entities/actions';
import { setPPMEstimateError } from 'store/onboarding/actions';

import './PpmSummary.css';

const MoveInfoHeader = (props) => {
  const { orders, profile, move, entitlement } = props;
  return (
    <div>
      <h1>
        {get(orders, 'new_duty_location.name', 'New move')} (from {get(profile, 'current_station.name', '')})
      </h1>
      {get(move, 'locator') && <div>Move Locator: {get(move, 'locator')}</div>}
      {!isEmpty(entitlement) && (
        <div>
          Weight allowance:{' '}
          <span data-testid="move-header-weight-estimate">{entitlement.weight.toLocaleString()} lbs</span>
        </div>
      )}
    </div>
  );
};

const genPpmSummaryStatusComponents = {
  DRAFT: DraftMoveSummary,
  SUBMITTED: SubmittedPpmMoveSummary,
  APPROVED: ApprovedMoveSummary,
  CANCELED: CanceledMoveSummary,
  PAYMENT_REQUESTED: PaymentRequestedSummary,
};

const getPPMStatus = (moveStatus, ppm) => {
  // PPM status determination
  const ppmStatus = get(ppm, 'status', 'DRAFT');
  return moveStatus === 'APPROVED' && (ppmStatus === 'SUBMITTED' || ppmStatus === 'DRAFT') ? 'SUBMITTED' : moveStatus;
};

export class PpmSummaryComponent extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      hasEstimateError: false,
      netWeight: null,
    };
  }

  componentDidMount() {
    if (this.props.move.id) {
      this.props.getMoveDocumentsForMove(this.props.move.id).then(({ obj: documents }) => {
        const weightTicketNetWeight = calcNetWeight(documents);
        let netWeight =
          weightTicketNetWeight > this.props.entitlement.sum ? this.props.entitlement.sum : weightTicketNetWeight;

        if (netWeight === 0) {
          netWeight = this.props.ppm.weight_estimate;
        }
        if (!netWeight) {
          this.setState({ hasEstimateError: true });
        }
        if (!isEmpty(this.props.ppm) && netWeight) {
          calculatePPMEstimate(
            this.props.ppm.original_move_date,
            this.props.ppm.pickup_postal_code,
            this.props.originDutyLocationZip,
            this.props.orders.id,
            netWeight,
          )
            .then((response) => {
              this.props.updatePPMEstimate(response);
              this.props.setPPMEstimateError(null);
            })
            .catch((err) => {
              this.props.setPPMEstimateError(err);
              this.setState({ hasEstimateError: true });
            });

          this.setState({ netWeight: netWeight });
        }
      });
    }
  }
  render() {
    const {
      profile,
      move,
      orders,
      ppm,
      editMove,
      entitlement,
      resumeMove,
      reviewProfile,
      isMissingWeightTicketDocuments,
    } = this.props;
    const moveStatus = get(move, 'status', 'DRAFT');
    const ppmStatus = getPPMStatus(moveStatus, ppm);
    const PPMComponent = genPpmSummaryStatusComponents[ppmStatus];
    return (
      <div>
        {move.status === 'CANCELED' && (
          <Alert type="info" heading="Your move was canceled">
            Your move from {get(profile, 'current_station.name')} to {get(orders, 'new_duty_location.name')} with the
            move locator ID {get(move, 'locator')} was canceled.
          </Alert>
        )}

        <div className="grid-row grid-gap">
          <div className="grid-col-12">
            {move.status !== 'CANCELED' && (
              <div>
                <MoveInfoHeader orders={orders} profile={profile} move={move} entitlement={entitlement} />
                <br />
              </div>
            )}
            {isMissingWeightTicketDocuments && ppm.status === 'PAYMENT_REQUESTED' && (
              <Alert type="warning" heading="Payment request is missing info">
                You will need to contact your local PPPO office to resolve your missing weight ticket.
              </Alert>
            )}
          </div>
        </div>
        <div className="grid-row">
          <div className="st-wrapper tablet:grid-col-9 grid-col-12">
            <PPMComponent
              className="status-component"
              ppm={ppm}
              orders={orders}
              profile={profile}
              move={move}
              entitlement={entitlement}
              resumeMove={resumeMove}
              reviewProfile={reviewProfile}
              isMissingWeightTicketDocuments={isMissingWeightTicketDocuments}
              hasEstimateError={this.state.hasEstimateError}
              netWeight={this.state.netWeight}
            />
          </div>

          <div className="sidebar tablet:grid-col-3 grid-col-12 tablet:padding-left-3">
            <div>
              <button
                className="usa-button usa-button--secondary"
                onClick={() => editMove(move)}
                disabled={includes(['DRAFT', 'CANCELED'], move.status)}
                data-testid="edit-move"
              >
                Edit Move
              </button>
            </div>
            <div className="contact_block">
              <h2>Contacts</h2>
              <TransportationOfficeContactInfo dutyLocation={profile.current_station} isOrigin={true} />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, ownProps.move.id, [
    'WEIGHT_TICKET_SET',
  ]).some((doc) => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);

  return {
    isMissingWeightTicketDocuments,
    originDutyLocationZip: serviceMember?.current_station?.address?.postalCode,
  };
}

const mapDispatchToProps = {
  getMoveDocumentsForMove,
  updatePPMEstimate,
  setPPMEstimateError,
};

export const PpmSummary = connect(mapStateToProps, mapDispatchToProps)(PpmSummaryComponent);
