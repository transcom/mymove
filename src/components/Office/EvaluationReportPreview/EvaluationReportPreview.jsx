import React from 'react';
import * as PropTypes from 'prop-types';
import { Grid } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './EvaluationReportPreview.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import EvaluationReportShipmentDisplay from 'components/Office/EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';
import DataTable from 'components/DataTable';
import DataTableWrapper from 'components/DataTableWrapper';
import EvaluationReportList from 'components/Office/DefinitionLists/EvaluationReportList';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';
import { formatDateFromIso, formatQAReportID } from 'utils/formatters';
import { formatDate } from 'shared/dates';
import { shipmentTypeLabels } from 'content/shipments';

const EvaluationReportPreview = ({ evaluationReport, mtoShipments, moveCode, customerInfo, grade }) => {
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
      {ORDERS_RANK_OPTIONS[grade]}
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
    <div className={styles.evaluationReportPreview} data-testid="EvaluationReportPreview">
      <div>
        <div className={styles.titleSection}>
          <div className={styles.pageHeader}>
            <h1>{`${isShipment ? 'Shipment' : 'Counseling'} report`}</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>REPORT ID {formatQAReportID(evaluationReport.id)}</h6>
              <h6>MOVE CODE #{moveCode}</h6>
              <h6>MTO REFERENCE ID #{evaluationReport.moveReferenceID}</h6>
            </div>
          </div>
        </div>
        <div className={styles.section}>
          <Grid row>
            <Grid col desktop={{ col: 8 }}>
              <h2>{`${isShipment ? 'Shipment' : 'Report'} information`}</h2>
              {mtoShipmentsToShow &&
                mtoShipmentsToShow.map((mtoShipment) => (
                  <div key={mtoShipment.id} className={styles.shipmentDisplayContainer}>
                    <EvaluationReportShipmentDisplay
                      isSubmitted
                      key={mtoShipment.id}
                      shipmentId={mtoShipment.id}
                      displayInfo={shipmentDisplayInfo(mtoShipment)}
                      shipmentType={mtoShipment.shipmentType}
                    />
                  </div>
                ))}
            </Grid>
            <Grid className={styles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
              <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
              <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
            </Grid>
          </Grid>
        </div>
      </div>
      <div className={styles.section}>
        <h2>Evaluation report</h2>
        {isShipment && evaluationReport.location !== 'OTHER' && (
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
                      evaluationReport.observedDate ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY') : '—',
                    ]}
                  />
                )}
                {evaluationReport.location === 'DESTINATION' && (
                  <DataTable
                    columnHeaders={['Observed delivery']}
                    dataRow={[
                      evaluationReport.observedDate ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY') : '—',
                    ]}
                  />
                )}
              </DataTableWrapper>
              <DataTableWrapper className={classnames(styles.detailsRight, 'table--data-point-group')}>
                <DataTable
                  columnHeaders={['Inspection date', 'Report submission']}
                  dataRow={[
                    evaluationReport.inspectionDate ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY') : '—',
                    evaluationReport.observedDate ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY') : '—',
                  ]}
                />
              </DataTableWrapper>
            </div>
            <EvaluationReportList evaluationReport={evaluationReport} />
          </div>
        )}
        {(!isShipment || evaluationReport.location === 'OTHER') && (
          <div className={styles.section}>
            <h3>Information</h3>
            <DataTableWrapper className={classnames(styles.counselingDetails, 'table--data-point-group')}>
              <DataTable
                columnHeaders={['Inspection date', 'Report submission']}
                dataRow={[
                  evaluationReport.inspectionDate ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY') : '—',
                  evaluationReport.submittedAt ? formatDateFromIso(evaluationReport.submittedAt, 'DD MMM YYYY') : '—',
                ]}
              />
            </DataTableWrapper>
            <EvaluationReportList evaluationReport={evaluationReport} />
          </div>
        )}
        <div className={styles.section}>
          <h3>QAE remarks</h3>
          <dl className={descriptionListStyles.descriptionList}>
            <div className={classnames(descriptionListStyles.row)}>
              <dt className={styles.qaeRemarksLabel}>Evaluation remarks</dt>
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
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
};

EvaluationReportPreview.defaultProps = {
  mtoShipments: null,
};

export default EvaluationReportPreview;
