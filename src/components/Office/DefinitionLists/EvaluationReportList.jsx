import React from 'react';
import classnames from 'classnames';

import { formatEvaluationReportLocation } from '../../../utils/formatters';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { EvaluationReportShape } from 'types/evaluationReport';

const capitalizeFirstLetterOnly = ([first, ...restOfString]) => {
  return first.toUpperCase() + restOfString.join('').toLowerCase();
};

const inspectionTypeFormatting = (inspectionType) => {
  if (inspectionType === 'DATA_REVIEW') {
    return 'Data review';
  }
  return capitalizeFirstLetterOnly(inspectionType);
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
        <div className={classnames(descriptionListStyles.row, descriptionListStyles.noBorder)}>
          <dt>Evaluation type</dt>
          <dd>{evaluationReport.inspectionType ? inspectionTypeFormatting(evaluationReport.inspectionType) : ''}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Evaluation location</dt>
          <dd>
            {formatEvaluationReportLocation(evaluationReport.location)}
            <br />
            {evaluationReport.locationDescription || ''}
          </dd>
        </div>
        {evaluationReport.travelTimeMinutes >= 0 && (
          <div className={descriptionListStyles.row}>
            <dt>Travel time to inspection</dt>
            <dd>{convertToHoursAndMinutes(evaluationReport.travelTimeMinutes)}</dd>
          </div>
        )}
        {evaluationReport.evaluationLengthMinutes >= 0 && (
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
