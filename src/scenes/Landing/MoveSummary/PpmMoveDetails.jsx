import React from 'react';
import { connect } from 'react-redux';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents } from 'shared/formatters';
import { formatActualIncentiveRange, formatIncentiveRange } from 'shared/incentive';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { selectPPMEstimateRange, selectReimbursement } from 'shared/Entities/modules/ppms';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import styles from './PpmMoveDetails.module.scss';

const PpmMoveDetails = ({ advance, ppm, isMissingWeightTicketDocuments, estimateRange, netWeight }) => {
  const privateStorageString = ppm.estimated_storage_reimbursement
    ? `(up to ${ppm.estimated_storage_reimbursement})`
    : '';
  const advanceString =
    ppm.has_requested_advance && advance && advance.requested_amount
      ? `Advance Requested: $${formatCents(advance.requested_amount)}`
      : '';
  const hasSitString = `Temp. Storage: ${ppm.days_in_storage} days ${privateStorageString}`;
  const estimatedIncentiveRange = formatIncentiveRange(ppm, estimateRange);
  const actualIncentiveRange = formatActualIncentiveRange(estimateRange);

  return (
    <div className="titled_block">
      <div className="title">Details</div>
      <div>Weight (est.): {ppm.weight_estimate} lbs</div>
      {netWeight && <div>Weight (submitted): {netWeight} lbs</div>}
      <div className="title" style={{ paddingTop: '0.5em' }}>
        Payment request
      </div>
      <div>Estimated payment: </div>
      {ppm.incentive_estimate_min || estimatedIncentiveRange ? (
        isMissingWeightTicketDocuments ? (
          <>
            <div className="missing-label">
              Unknown
              <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
            </div>
            <div className={styles.subText}>
              <em>Estimated payment will be given after resolving missing weight tickets.</em>
            </div>
          </>
        ) : (
          <>
            <div className={styles['payment-details']}>
              <div>{estimatedIncentiveRange}</div>
            </div>

            <div className={styles['payment-details']}>
              <div>Submitted payment: </div>
              <div>{actualIncentiveRange}</div>
              <div className={styles.subText}>
                <em>Actual payment may vary, subject to Finance review.</em>
              </div>
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

const mapStateToProps = (state, ownProps) => {
  const advance = selectReimbursement(state, ownProps.ppm.advance);
  const isMissingWeightTicketDocuments = selectPPMCloseoutDocumentsForMove(state, ownProps.ppm.move_id, [
    'WEIGHT_TICKET_SET',
  ]).some((doc) => doc.empty_weight_ticket_missing || doc.full_weight_ticket_missing);
  return {
    advance,
    isMissingWeightTicketDocuments,
    estimateRange: selectPPMEstimateRange(state),
  };
};

export default connect(mapStateToProps)(PpmMoveDetails);
