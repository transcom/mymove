import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportStyles from './EvaluationReport.module.scss';

import EvaluationReportMoveInfoDisplay from 'components/Office/EvaluationReportMoveInfo/EvaluationReportMoveInfo';
import EvaluationForm from 'components/Office/EvaluationForm/EvaluationForm';
import { useShipmentEvaluationReportQueries } from 'hooks/queries';
import { CustomerShape } from 'types';
import { OrdersShape } from 'types/customerShapes';
import EvaluationReportShipmentInfoDisplay from 'components/Office/EvaluationReportShipmentInfo/EvaluationReportShipmentInfo';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';

const EvaluationReport = ({ customerInfo, orders }) => {
  const { reportId } = useParams();
  const { evaluationReport, mtoShipment } = useShipmentEvaluationReportQueries(reportId);

  const isShipment = evaluationReport.type === EVALUATION_REPORT_TYPE.SHIPMENT;
  return (
    <div className={classnames(styles.tabContent, evaluationReportStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />

        {isShipment ? (
          <EvaluationReportShipmentInfoDisplay
            customerInfo={customerInfo}
            orders={orders}
            shipment={mtoShipment}
            report={evaluationReport}
          />
        ) : (
          <EvaluationReportMoveInfoDisplay customerInfo={customerInfo} orders={orders} />
        )}

        <EvaluationForm evaluationReport={evaluationReport} />
      </GridContainer>
    </div>
  );
};

EvaluationReport.propTypes = {
  customerInfo: CustomerShape.isRequired,
  orders: OrdersShape.isRequired,
};

export default EvaluationReport;
