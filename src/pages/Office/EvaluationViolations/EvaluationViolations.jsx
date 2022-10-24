import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useEvaluationReportShipmentListQueries, usePWSViolationsQueries } from 'hooks/queries';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import EvaluationViolationsForm from 'components/Office/EvaluationViolationsForm/EvaluationViolationsForm';

const EvaluationViolations = ({ customerInfo, grade, destinationDutyLocationPostalCode }) => {
  const { reportId } = useParams();

  const { evaluationReport, reportViolations, mtoShipments, isLoading, isError } =
    useEvaluationReportShipmentListQueries(reportId);
  const { violations, isLoading: isViolationsLoading, isError: isViolationsError } = usePWSViolationsQueries();

  if (isLoading || isViolationsLoading) {
    return <LoadingPlaceholder />;
  }
  if (isError || isViolationsError) {
    return <SomethingWentWrong />;
  }

  return (
    <div className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />

        <EvaluationViolationsForm
          violations={violations}
          evaluationReport={evaluationReport}
          reportViolations={reportViolations}
          customerInfo={customerInfo}
          grade={grade}
          mtoShipments={mtoShipments}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      </GridContainer>
    </div>
  );
};

export default EvaluationViolations;
