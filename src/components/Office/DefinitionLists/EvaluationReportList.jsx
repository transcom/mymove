import React from 'react';

import { formatEvaluationReportLocation } from '../../../utils/formatters';

import styles from './OfficeDefinitionLists.module.scss';

import PreviewRow from 'components/Office/EvaluationReportPreview/PreviewRow/PreviewRow';
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
        <PreviewRow
          label="Evaluation type"
          data={evaluationReport.inspectionType ? inspectionTypeFormatting(evaluationReport.inspectionType) : ''}
        />
        <PreviewRow
          label="Evaluation location"
          data={
            <>
              {formatEvaluationReportLocation(evaluationReport.location)}
              <br />
              {evaluationReport.locationDescription || ''}
            </>
          }
        />
        <PreviewRow
          isShown={evaluationReport.travelTimeMinutes >= 0}
          label="Travel time to inspection"
          data={convertToHoursAndMinutes(evaluationReport.travelTimeMinutes)}
        />
        <PreviewRow
          isShown={evaluationReport.evaluationLengthMinutes >= 0}
          label="Evaluation length"
          data={convertToHoursAndMinutes(evaluationReport.evaluationLengthMinutes)}
        />
      </dl>
    </div>
  );
};
EvaluationReportList.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
};
export default EvaluationReportList;
