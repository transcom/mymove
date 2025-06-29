import React from 'react';
import * as PropTypes from 'prop-types';
import { Grid } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './EvaluationReportPreview.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';
import EvaluationReportShipmentDisplay from 'components/Office/EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';
import DataTable from 'components/DataTable';
import DataTableWrapper from 'components/DataTableWrapper';
import EvaluationReportList from 'components/Office/DefinitionLists/EvaluationReportList';
import EvaluationReportViolationsList from 'components/Office/DefinitionLists/EvaluationReportViolationsList';
import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';
import { formatDateFromIso, formatQAReportID } from 'utils/formatters';
import { formatDate } from 'shared/dates';
import { shipmentTypeLabels } from 'content/shipments';

const EvaluationReportPreview = ({
  evaluationReport,
  reportViolations,
  mtoShipments,
  moveCode,
  customerInfo,
  grade,
  destinationDutyLocationPostalCode,
}) => {
  const isShipment = evaluationReport.type === EVALUATION_REPORT_TYPE.SHIPMENT;

  let mtoShipmentsToShow;
  if (evaluationReport.shipmentID) {
    mtoShipmentsToShow = [mtoShipments.find((shipment) => shipment.id === evaluationReport.shipmentID)];
  } else {
    mtoShipmentsToShow = mtoShipments;
  }

  const customerInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {grade}
      <br />
      {ORDERS_BRANCH_OPTIONS[customerInfo.agency] ? ORDERS_BRANCH_OPTIONS[customerInfo.agency] : customerInfo.agency}
    </>
  );

  const officeUserInfoTableBody = evaluationReport.officeUser ? (
    <>
      {evaluationReport.officeUser.lastName}, {evaluationReport.officeUser.firstName}
      <br />
      {evaluationReport.officeUser.phone}
      <br />
      {evaluationReport.officeUser.email}
    </>
  ) : (
    ''
  );

  const shipmentDisplayInfo = (shipment) => {
    return {
      ...shipment,
      heading: shipmentTypeLabels[shipment.shipmentType],
      isDiversion: shipment.diversion,
      shipmentStatus: shipment.status,
      destinationAddress: shipment.destinationAddress,
    };
  };

  return (
    <div data-testid="EvaluationReportPreview">
      <div>
        {/* Page Header */}
        <div className={styles.pageHeader}>
          <h1>{`${isShipment ? 'Shipment' : 'Counseling'} report`}</h1>
          <div className={styles.pageHeaderDetails}>
            <h6>REPORT ID {formatQAReportID(evaluationReport.id)}</h6>
            <h6>MOVE CODE #{moveCode}</h6>
            <h6>MTO REFERENCE ID #{evaluationReport.moveReferenceID}</h6>
          </div>
        </div>

        {/* Move and Customer/QAE info */}
        <div className={styles.section}>
          <Grid row>
            <Grid col desktop={{ col: 8 }}>
              <h2>Move information</h2>
              {mtoShipmentsToShow &&
                mtoShipmentsToShow.map((mtoShipment) => (
                  <div key={mtoShipment.id} className={styles.shipmentDisplayContainer}>
                    <EvaluationReportShipmentDisplay
                      isSubmitted
                      key={mtoShipment.id}
                      shipmentId={mtoShipment.id}
                      displayInfo={shipmentDisplayInfo(mtoShipment)}
                      shipmentType={mtoShipment.shipmentType}
                      destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
                    />
                  </div>
                ))}
            </Grid>
            <Grid className={styles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
              <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
              <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
            </Grid>
          </Grid>
        </div>
      </div>

      {/* Report content */}
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
                    evaluationReport.inspectionDate ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY') : '—',
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
          <h3>QAE remarks</h3>
          <dl className={descriptionListStyles.descriptionList}>
            <div className={descriptionListStyles.row}>
              <dt className={styles.label}>Evaluation remarks</dt>
              <dd className={styles.qaeRemarks}>{evaluationReport.remarks}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  );
};

EvaluationReportPreview.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  reportViolations: PropTypes.arrayOf(Object),
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  destinationDutyLocationPostalCode: PropTypes.string,
};

EvaluationReportPreview.defaultProps = {
  mtoShipments: null,
  reportViolations: null,
  destinationDutyLocationPostalCode: '',
};

export default EvaluationReportPreview;
