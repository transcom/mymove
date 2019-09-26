import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents, formatCentsRange } from 'shared/formatters';
import { selectReimbursement } from 'shared/Entities/modules/ppms';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';

const SubmittedPpmMoveDetails = props => {
  const { advance, ppm } = props;
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString = ppm.has_requested_advance ? `Advance Requested: $${formatCents(advance.requested_amount)}` : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;
  const incentiveRange = formatCentsRange(ppm.currentPpm.incentive_estimate_min, ppm.currentPpm.incentive_estimate_max);
  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.currentPpm.weight_estimate} lbs</div>
      <div>
        Incentive (est.):{' '}
        {ppm.hasEstimateError ? (
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
  ]).some(doc => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);
  return { ppm: get(state, 'ppm', {}), advance, isMissingWeightTicketDocuments };
};

export default connect(mapStateToProps)(SubmittedPpmMoveDetails);
