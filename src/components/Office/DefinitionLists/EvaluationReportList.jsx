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
        {evaluationReport.inspectionType === 'PHYSICAL' && (
          <PreviewRow label="Time departed for evaluation" data={evaluationReport.timeDepart} />
        )}
        <PreviewRow label="Time evaluation started" data={evaluationReport.evalStart} />
        <PreviewRow label="Time evaluation ended" data={evaluationReport.evalEnd} />
      </dl>
    </div>
  );
};
EvaluationReportList.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
};
export default EvaluationReportList;
