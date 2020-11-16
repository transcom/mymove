import React from 'react';
import { connect } from 'react-redux';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCents } from 'shared/formatters';
import { formatActualIncentiveRange, formatIncentiveRange } from 'shared/incentive';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faExclamationCircle } from '@fortawesome/free-solid-svg-icons/faExclamationCircle';
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
  const hasRangeReady = ppm.incentive_estimate_min || estimatedIncentiveRange;

  const incentiveNotReady = () => {
    return (
      <>
        Not ready yet{' '}
        <IconWithTooltip toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive." />
      </>
    );
  };

  return (
    <div className="titled_block">
      <div className={styles['detail-title']}>Estimated</div>
      <div>Weight: {ppm.weight_estimate} lbs</div>
      {hasRangeReady && isMissingWeightTicketDocuments ? (
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
          <div>
            <div>Payment: {hasRangeReady ? estimatedIncentiveRange : incentiveNotReady()}</div>
          </div>
          {ppm.status === 'PAYMENT_REQUESTED' && (
            <div className={styles['payment-details']}>
              <div className={styles['detail-title']}>Submitted</div>
              <div>Weight: {netWeight} lbs</div>
              <div>Payment request: {hasRangeReady ? actualIncentiveRange : incentiveNotReady()}</div>
            </div>
          )}
          <div className={styles.subText}>
            <em>Actual payment may vary, subject to Finance review.</em>
          </div>
        </>
      )}
      {ppm.has_sit && <div className={styles['payment-details']}>{hasSitString}</div>}
      {ppm.has_requested_advance && <div className={styles['payment-details']}>{advanceString}</div>}
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
