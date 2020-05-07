import React from 'react';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents } from 'shared/formatters';
import { formatIncentiveRange } from 'shared/incentive';
import { selectReimbursement } from 'shared/Entities/modules/ppms';
import { selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';

const SubmittedPpmMoveDetails = (props) => {
  const { advance, ppm, currentPPM, tempCurrentPPM, hasEstimateError } = props;
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString = ppm.has_requested_advance ? `Advance Requested: $${formatCents(advance.requested_amount)}` : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;
  const incentiveRange = formatIncentiveRange(ppm);

  const weightEstimate = isEmpty(currentPPM) ? tempCurrentPPM.weight_estimate : currentPPM.weight_estimate;

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {weightEstimate} lbs</div>
      <div>
        Incentive (est.):{' '}
        {ppm.hasEstimateError || hasEstimateError ? (
          <>
            Not ready yet{' '}
            <IconWithTooltip
              // without this styling the tooltip is obstructed by the status timeline and z-index does not help because they don't share the same parent container
              toolTipStyles={{ position: 'absolute', top: '8.5em', left: '20.5em' }}
              toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive."
            />
          </>
        ) : (
          incentiveRange
        )}
      </div>
      {ppm.has_sit && <div>{hasSitString}</div>}
      {ppm.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const mapStateToProps = (state, ownProps) => {
  const advance = selectReimbursement(state, ownProps.ppm.advance);
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, ownProps.ppm.move_id, [
    'WEIGHT_TICKET_SET',
  ]).some((doc) => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);
  const moveID = state.moves.currentMove.id;

  const props = {
    currentPPM: selectActivePPMForMove(state, moveID),
    // TODO this is a work around till we refactor more SM data...
    tempCurrentPPM: get(state, 'ppm.currentPpm'),
    ppm: get(state, 'ppm', {}),
    advance,
    isMissingWeightTicketDocuments,
  };
  return props;
};

export default connect(mapStateToProps)(SubmittedPpmMoveDetails);
