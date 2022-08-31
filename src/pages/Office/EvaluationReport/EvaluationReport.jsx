import React from 'react';
import PropTypes from 'prop-types';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportStyles from './EvaluationReport.module.scss';

import EvaluationForm from 'components/Office/EvaluationForm/EvaluationForm';
import { useEvaluationReportShipmentListQueries } from 'hooks/queries';
import { CustomerShape } from 'types';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';
import EvaluationReportMoveInfo from 'components/Office/EvaluationReportMoveInfo/EvaluationReportMoveInfo';
import EvaluationReportShipmentInfo from 'components/Office/EvaluationReportShipmentInfo/EvaluationReportShipmentInfo';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';

const EvaluationReport = ({ customerInfo, grade }) => {
  const { reportId } = useParams();
  const { evaluationReport, mtoShipments } = useEvaluationReportShipmentListQueries(reportId);
  const { shipmentID } = evaluationReport;
  const isShipment = evaluationReport.type === EVALUATION_REPORT_TYPE.SHIPMENT;
  let singleShipment = null;

  if (isShipment) {
    singleShipment = [mtoShipments.find((shipment) => shipment.id === shipmentID)];
  }

  return (
    <div className={classnames(styles.tabContent, evaluationReportStyles.tabContent)}>
      <GridContainer className={evaluationReportStyles.container}>
        <QaeReportHeader report={evaluationReport} />
        {mtoShipments?.length > 0 && (
          <EvaluationReportShipmentInfo
            customerInfo={customerInfo}
            grade={grade}
            shipment={singleShipment}
            report={evaluationReport}
          />
        )}{' '}
        : (
        <EvaluationReportMoveInfo customerInfo={customerInfo} grade={grade} />
        )
        <EvaluationForm
          evaluationReport={evaluationReport}
          mtoShipments={mtoShipments}
          customerInfo={customerInfo}
          grade={grade}
        />
      </GridContainer>
    </div>
  );
};

EvaluationReport.propTypes = {
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
};

export default EvaluationReport;
