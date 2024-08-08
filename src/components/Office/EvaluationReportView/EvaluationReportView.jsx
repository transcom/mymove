import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router-dom';

import EvaluationReportList from '../DefinitionLists/EvaluationReportList';
import EvaluationReportViolationsList from '../DefinitionLists/EvaluationReportViolationsList';

import styles from './EvaluationReportView.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { useEvaluationReportShipmentListQueries } from 'hooks/queries';
import EvaluationReportShipmentInfo from 'components/Office/EvaluationReportShipmentInfo/EvaluationReportShipmentInfo';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DataTableWrapper from 'components/DataTableWrapper';
import DataTable from 'components/DataTable';
import { formatDate } from 'shared/dates';
import { formatDateFromIso } from 'utils/formatters';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';

const EvaluationReportView = ({ customerInfo, grade, destinationDutyLocationPostalCode }) => {
  const { reportId } = useParams();
  const { evaluationReport, reportViolations, mtoShipments, isLoading, isError } =
    useEvaluationReportShipmentListQueries(reportId);

  const isShipment = evaluationReport.type === EVALUATION_REPORT_TYPE.SHIPMENT;

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
    <div className={styles.tabContent}>
      <GridContainer className={styles.container}>
        <QaeReportHeader report={evaluationReport} />
        {mtoShipmentsToShow?.length > 0 && (
          <EvaluationReportShipmentInfo
            customerInfo={customerInfo}
            grade={grade}
            shipments={mtoShipmentsToShow}
            report={evaluationReport}
            destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
          />
        )}
        <div className={styles.section}>
          <h2>Evaluation report</h2>
          {isShipment && evaluationReport.location !== 'OTHER' ? (
            <div className={styles.section}>
              <h3>Information</h3>
              <div className={styles.sideBySideDetails}>
                <DataTableWrapper className={classnames(styles.detailsLeft, 'table--data-point-group')}>
                  {evaluationReport.location === 'ORIGIN' && (
                    <DataTable
                      columnHeaders={['Scheduled pickup', 'Observed pickup']}
                      dataRow={[
                        mtoShipments[0].scheduledPickupDate
                          ? formatDate(mtoShipments[0].scheduledPickupDate, 'DD MMM YYYY')
                          : '—',
                        evaluationReport.observedShipmentPhysicalPickupDate
                          ? formatDate(evaluationReport.observedShipmentPhysicalPickupDate, 'DD MMM YYYY')
                          : '—',
                      ]}
                    />
                  )}
                  {evaluationReport.location === 'DESTINATION' && (
                    <DataTable
                      columnHeaders={['Observed delivery']}
                      dataRow={[
                        evaluationReport.observedShipmentDeliveryDate
                          ? formatDate(evaluationReport.observedShipmentDeliveryDate, 'DD MMM YYYY')
                          : '—',
                      ]}
                    />
                  )}
                </DataTableWrapper>
                <DataTableWrapper className={classnames(styles.detailsRight, 'table--data-point-group')}>
                  <DataTable
                    columnHeaders={['Inspection date', 'Report submission']}
                    dataRow={[
                      evaluationReport.inspectionDate
                        ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY')
                        : '—',
                      evaluationReport.submittedAt ? formatDate(evaluationReport.submittedAt, 'DD MMM YYYY') : '—',
                    ]}
                  />
                </DataTableWrapper>
              </div>
              <EvaluationReportList evaluationReport={evaluationReport} />
            </div>
          ) : (
            <div className={styles.section}>
              <h3>Information</h3>
              <DataTableWrapper className={classnames(styles.detailsRight, 'table--data-point-group')}>
                <DataTable
                  columnHeaders={['Inspection date', 'Report submission']}
                  dataRow={[
                    evaluationReport.inspectionDate ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY') : '—',
                    evaluationReport.submittedAt
                      ? formatDateFromIso(evaluationReport.submittedAt, 'DD MMM YYYY')
                      : formatDate(new Date(), 'DD MMM YYYY'),
                  ]}
                />
              </DataTableWrapper>
              <EvaluationReportList evaluationReport={evaluationReport} />
            </div>
          )}
          <div className={styles.section}>
            <h3>Violations</h3>
            <EvaluationReportViolationsList evaluationReport={evaluationReport} reportViolations={reportViolations} />
          </div>
          <div className={styles.section}>
            <h3>Serious Incident</h3>
            <EvaluationReportViolationsList evaluationReport={evaluationReport} reportViolations={reportViolations} />
          </div>
          <div className={styles.section}>
            <h3>QAE remarks</h3>
            <dl className={descriptionListStyles.descriptionList}>
              <div className={descriptionListStyles.row}>
                <dt className={styles.label}>Evaluation remarks</dt>
                <dd className={styles.qaeRemarks}>{evaluationReport.remarks}</dd>
              </div>
            </dl>
          </div>
        </div>
      </GridContainer>
    </div>
  );
};

export default EvaluationReportView;
