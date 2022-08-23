import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportStyles from './EvaluationReport.module.scss';

import EvaluationForm from 'components/Office/EvaluationForm/EvaluationForm';
import { useEvaluationReportQueries } from 'hooks/queries';
import { CustomerShape } from 'types';
import { OrdersShape } from 'types/customerShapes';
import EvaluationReportShipmentInfo from 'components/Office/EvaluationReportShipmentInfo/EvaluationReportShipmentInfo';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';

const EvaluationReport = ({ customerInfo, orders }) => {
  const { reportId } = useParams();
  const { evaluationReport, mtoShipments } = useEvaluationReportQueries(reportId);
  return (
    <div className={classnames(styles.tabContent, evaluationReportStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />

        {mtoShipments?.length > 0 && (
          <EvaluationReportShipmentInfo
            customerInfo={customerInfo}
            orders={orders}
            shipments={mtoShipments}
            report={evaluationReport}
          />
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
