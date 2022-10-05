import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import { useEvaluationReportShipmentListQueries, usePWSViolationsQueries } from 'hooks/queries';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import EvaluationViolationsForm from 'components/Office/EvaluationViolationsForm/EvaluationViolationsForm';

const EvaluationViolations = ({ customerInfo }) => {
  const { reportId } = useParams();

  const { evaluationReport, reportViolations, mtoShipments } = useEvaluationReportShipmentListQueries(reportId);
  const { violations } = usePWSViolationsQueries();

  return (
    <div className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />

        <EvaluationViolationsForm
          violations={violations}
          evaluationReport={evaluationReport}
          reportViolations={reportViolations}
          customerInfo={customerInfo}
          mtoShipments={mtoShipments}
        />
      </GridContainer>
    </div>
  );
};

export default EvaluationViolations;
