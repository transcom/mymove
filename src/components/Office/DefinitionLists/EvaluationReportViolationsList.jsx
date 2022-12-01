import React from 'react';

import styles from './EvaluationReportViolationsList.module.scss';

import PreviewRow from 'components/Office/EvaluationReportPreview/PreviewRow/PreviewRow';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { EvaluationReportShape } from 'types/evaluationReport';
import { formatDate } from 'shared/dates';

const EvaluationReportViolationsList = ({ evaluationReport, reportViolations }) => {
  const hasViolations = reportViolations && reportViolations.length > 0;
  const showIncidentDescription = evaluationReport?.seriousIncident;

  return (
    <dl className={descriptionListStyles.descriptionList}>
      <div className={descriptionListStyles.row}>
        <dt data-testid="violationsObserved" className={styles.label}>
          Violations observed
        </dt>
        {hasViolations ? (
          <dd className={styles.violationsRemarks}>
            {reportViolations.map((reportViolation) => (
              <div className={styles.violation} key={`${reportViolation.id}-violation`}>
                <h5>{`${reportViolation?.violation?.paragraphNumber} ${reportViolation?.violation?.title}`}</h5>
                <p>
                  <small>{reportViolation?.violation?.requirementSummary}</small>
                </p>
              </div>
            ))}
          </dd>
        ) : (
          <dd className={styles.violationsRemarks} data-testid="noViolationsObserved">
            No
          </dd>
        )}
      </div>
      <PreviewRow
        isShown={
          'observedPickupSpreadStartDate' in evaluationReport && 'observedPickupSpreadEndDate' in evaluationReport
        }
        label="Observed Pickup Spread Dates"
        data={`${formatDate(evaluationReport?.observedPickupSpreadStartDate, 'DD MMM YYYY')} - ${formatDate(
          evaluationReport?.observedPickupSpreadEndDate,
          'DD MMM YYYY',
        )}`}
      />
      <PreviewRow
        isShown={'observedClaimsResponseDate' in evaluationReport}
        label="Observed Claims Response Date"
        data={formatDate(evaluationReport?.observedClaimsResponseDate, 'DD MMM YYYY')}
      />
      <PreviewRow
        isShown={'observedPickupDate' in evaluationReport}
        label="Observed Pickup Date"
        data={formatDate(evaluationReport?.observedPickupDate, 'DD MMM YYYY')}
      />
      <PreviewRow
        isShown={'observedDeliveryDate' in evaluationReport}
        label="Observed Delivery Date"
        data={formatDate(evaluationReport?.observedDeliveryDate, 'DD MMM YYYY')}
      />
      <PreviewRow isShown={hasViolations} label="Serious incident" data={showIncidentDescription ? 'Yes' : 'No'} />
      <PreviewRow
        isShown={hasViolations && showIncidentDescription}
        label="Serious incident description"
        data={evaluationReport?.seriousIncidentDesc}
      />
    </dl>
  );
};
EvaluationReportViolationsList.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
};
export default EvaluationReportViolationsList;
