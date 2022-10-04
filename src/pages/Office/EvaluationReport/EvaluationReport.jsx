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
import EvaluationReportShipmentInfo from 'components/Office/EvaluationReportShipmentInfo/EvaluationReportShipmentInfo';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const EvaluationReport = ({ customerInfo, grade }) => {
  const { reportId } = useParams();
  const { evaluationReport, reportViolations, mtoShipments, isLoading, isError } =
    useEvaluationReportShipmentListQueries(reportId);
  if (isLoading) {
    return <LoadingPlaceholder />;
  }
  if (isError) {
    return <SomethingWentWrong />;
  }
  let mtoShipmentsToShow;
  if (evaluationReport.shipmentID && mtoShipments) {
    mtoShipmentsToShow = [mtoShipments.find((shipment) => shipment.id === evaluationReport.shipmentID)];
  } else {
    mtoShipmentsToShow = mtoShipments;
  }

  return (
    <div className={classnames(styles.tabContent, evaluationReportStyles.tabContent)}>
      <GridContainer className={evaluationReportStyles.container}>
        <QaeReportHeader report={evaluationReport} />

        {mtoShipmentsToShow?.length > 0 && (
          <EvaluationReportShipmentInfo
            customerInfo={customerInfo}
            grade={grade}
            shipments={mtoShipmentsToShow}
            report={evaluationReport}
          />
        )}
        <EvaluationForm
          evaluationReport={evaluationReport}
          reportViolations={reportViolations}
          grade={grade}
          customerInfo={customerInfo}
          mtoShipments={mtoShipments}
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
