import React, { useState } from 'react';
import 'styles/office.scss';
import { Button, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router-dom';

import EvaluationReportList from '../DefinitionLists/EvaluationReportList';
import PreviewRow from '../EvaluationReportPreview/PreviewRow/PreviewRow';

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

  // this is currently set to false to hide the "add appeal" button for future GSR work
  // TODO implement permissions check in order to render this button
  const [canLeaveAppeal] = useState(false);

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

  const hasViolations = reportViolations && reportViolations.length > 0;
  const showIncidentDescription = evaluationReport?.seriousIncident;

  return (
    <div className={styles.tabContent} data-testid="EvaluationReportPreview">
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
            <dl className={descriptionListStyles.descriptionList}>
              <div className={descriptionListStyles.row}>
                <dt data-testid="violationsObserved" className={styles.label}>
                  Violations observed
                </dt>
                {hasViolations ? (
                  <dd className={styles.violationsRemarks}>
                    {reportViolations.map((reportViolation) => (
                      <div className={styles.violation} key={`${reportViolation.id}-violation`}>
                        <div className={styles.violationHeader}>
                          <h5>{`${reportViolation?.violation?.paragraphNumber} ${reportViolation?.violation?.title}`}</h5>
                          {canLeaveAppeal && (
                            <Button unstyled className={styles.addAppealBtn}>
                              Add appeal
                            </Button>
                          )}
                        </div>
                        <p>
                          <small>{reportViolation?.violation?.requirementSummary}</small>
                        </p>
                      </div>
                    ))}
                  </dd>
                ) : (
                  <dd className={styles.violationsRemarks} data-testid="noViolationsObserved">
                    No
                  </dd>
                )}
              </div>
              <PreviewRow
                isShown={
                  'observedPickupSpreadStartDate' in evaluationReport &&
                  'observedPickupSpreadEndDate' in evaluationReport
                }
                label="Observed Pickup Spread Dates"
                data={`${formatDate(evaluationReport?.observedPickupSpreadStartDate, 'DD MMM YYYY')} - ${formatDate(
                  evaluationReport?.observedPickupSpreadEndDate,
                  'DD MMM YYYY',
                )}`}
              />
              <PreviewRow
                isShown={'observedClaimsResponseDate' in evaluationReport}
                label="Observed Claims Response Date"
                data={formatDate(evaluationReport?.observedClaimsResponseDate, 'DD MMM YYYY')}
              />
              <PreviewRow
                isShown={'observedPickupDate' in evaluationReport}
                label="Observed Pickup Date"
                data={formatDate(evaluationReport?.observedPickupDate, 'DD MMM YYYY')}
              />
              <PreviewRow
                isShown={'observedDeliveryDate' in evaluationReport}
                label="Observed Delivery Date"
                data={formatDate(evaluationReport?.observedDeliveryDate, 'DD MMM YYYY')}
              />
            </dl>
          </div>
          <div className={styles.section}>
            <div className={styles.seriousIncidentHeader}>
              <h3>Serious Incident</h3>
              {canLeaveAppeal && (
                <Button unstyled className={styles.addAppealBtn}>
                  Add appeal
                </Button>
              )}
            </div>
            <dl className={descriptionListStyles.descriptionList} data-testid="seriousIncidentSection">
              <div className={descriptionListStyles.row}>
                <dt className={styles.label}>Serious incident</dt>
                <dd className={styles.seriousIncidentRemarks} data-testid="seriousIncidentYesNo">
                  {showIncidentDescription ? 'Yes' : 'No'}
                </dd>
              </div>
              {showIncidentDescription && (
                <div className={descriptionListStyles.row} data-testid="seriousIncidentDescription">
                  <dt className={styles.label}>Description</dt>
                  <dd className={styles.seriousIncidentRemarks}>{evaluationReport?.seriousIncidentDesc}</dd>
                </div>
              )}
            </dl>
          </div>
          <div className={styles.section} data-testid="qaeRemarks">
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
