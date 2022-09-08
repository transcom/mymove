import React from 'react';

import { formatEvaluationReportLocation } from '../../../utils/formatters';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { EvaluationReportShape } from 'types/evaluationReport';

const capitalizeFirstLetterOnly = ([first, ...restOfString]) => {
  return first.toUpperCase() + restOfString.join('').toLowerCase();
};

const convertToHoursAndMinutes = (totalMinutes) => {
  // divide and round down to get hours
  const hours = Math.floor(totalMinutes / 60);
  // use modulus operator to get the remainder for minutes
  const minutes = totalMinutes % 60;
  return `${hours} hr ${minutes} min`;
};

const EvaluationReportList = ({ evaluationReport }) => {
  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Evaluation type</dt>
          <dd>{capitalizeFirstLetterOnly(evaluationReport.type)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Evaluation location</dt>
          <dd>
            {formatEvaluationReportLocation(evaluationReport.location)}
            <br />
            {evaluationReport.locationDescription || ''}
          </dd>
        </div>
        {evaluationReport.travelTimeMinutes && (
          <div className={descriptionListStyles.row}>
            <dt>Travel time to inspection</dt>
            <dd>{convertToHoursAndMinutes(evaluationReport.travelTimeMinutes)}</dd>
          </div>
        )}
        {evaluationReport.evaluationLengthMinutes && (
          <div className={descriptionListStyles.row}>
            <dt>Evaluation length</dt>
            <dd>{convertToHoursAndMinutes(evaluationReport.evaluationLengthMinutes)}</dd>
          </div>
        )}
      </dl>
    </div>
  );
};
EvaluationReportList.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
};
export default EvaluationReportList;
