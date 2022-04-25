import React from 'react';
import { connect } from 'react-redux';

import styles from './PpmMoveDetails.module.scss';

import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents } from 'utils/formatters';
import { getIncentiveRange } from 'utils/incentives';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { selectCurrentPPM, selectPPMEstimateRange, selectReimbursementById } from 'store/entities/selectors';
import { selectPPMEstimateError } from 'store/onboarding/selectors';

const SubmittedPpmMoveDetails = (props) => {
  const { advance, currentPPM, hasEstimateError, estimateRange } = props;
  const privateStorageString = currentPPM?.estimated_storage_reimbursement
    ? `(up to ${currentPPM.estimated_storage_reimbursement})`
    : '';
  const advanceString = currentPPM?.has_requested_advance
    ? `Advance Requested: $${formatCents(advance.requested_amount)}`
    : '';
  const hasSitString = `Temp. Storage: ${currentPPM?.days_in_storage} days ${privateStorageString}`;
  const incentiveRange = getIncentiveRange(currentPPM, estimateRange);

  const weightEstimate = currentPPM?.weight_estimate;
  return (
    <div className="titled_block">
      <div className={styles['detail-title']}>Estimated</div>
      <div>Weight: {weightEstimate} lbs</div>
      <div>
        Payment:{' '}
        {!incentiveRange || hasEstimateError ? (
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
      {currentPPM?.has_sit && <div>{hasSitString}</div>}
      {currentPPM?.has_requested_advance && <div>{advanceString}</div>}
    </div>
  );
};

const mapStateToProps = (state) => {
  const currentPPM = selectCurrentPPM(state) || {};
  const advance = selectReimbursementById(state, currentPPM?.advance) || {};
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, currentPPM?.move_id, [
    'WEIGHT_TICKET_SET',
  ]).some((doc) => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);

  const props = {
    currentPPM,
    advance,
    isMissingWeightTicketDocuments,
    estimateRange: selectPPMEstimateRange(state) || {},
    hasEstimateError: selectPPMEstimateError(state),
  };
  return props;
};

export default connect(mapStateToProps)(SubmittedPpmMoveDetails);
