import React from 'react';
import * as PropTypes from 'prop-types';
import { Button, Grid } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { ModalContainer, Overlay } from '../../MigratedModal/MigratedModal';
import Modal from '../../Modal/Modal';
import { formatDateFromIso, formatQAReportID } from '../../../utils/formatters';
import EvaluationReportShipmentDisplay from '../EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';
import DataTable from '../../DataTable';
import DataTableWrapper from '../../DataTableWrapper';
import { formatDate } from '../../../shared/dates';
import EvaluationReportList from '../DefinitionLists/EvaluationReportList';

import styles from './EvaluationReportContainer.module.scss';

import shipmentEvaluationReportStyles from 'pages/Office/ShipmentEvaluationReport/ShipmentEvaluationReport.module.scss';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { CustomerShape } from 'types';
import { shipmentTypeLabels } from 'content/shipments';
import { useViewEvaluationReportQueries } from 'hooks/queries';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const EvaluationReportContainer = ({
  evaluationReportId,
  shipmentId,
  moveCode,
  customerInfo,
  grade,
  setIsModalVisible,
}) => {
  const { evaluationReport, mtoShipments } = useViewEvaluationReportQueries(evaluationReportId);
  let mtoShipmentsToShow;

  if (shipmentId && mtoShipments) {
    mtoShipmentsToShow = [mtoShipments.find((shipment) => shipment.id === shipmentId)];
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
    <div className={styles.evaluationReportContainer} data-testid="EvaluationReportContainer">
      <Overlay />
      <ModalContainer>
        <Modal className={styles.evaluationReportModal}>
          <div>
            <div className={styles.titleSection}>
              <div className={styles.pageHeader}>
                <h1>{evaluationReport.type === 'SHIPMENT' ? 'Shipment' : 'Move'} report</h1>
                <div className={styles.pageHeaderDetails}>
                  <h6>REPORT ID {formatQAReportID(evaluationReportId)}</h6>
                  <h6>MOVE CODE #{moveCode}</h6>
                  <h6>MTO REFERENCE ID #{evaluationReport.moveReferenceID}</h6>
                </div>
              </div>
            </div>
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
                        />
                      </div>
                    ))}
                </Grid>
                <Grid className={shipmentEvaluationReportStyles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
                  <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
                  <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
                </Grid>
              </Grid>
            </div>
          </div>
          <div className={styles.section}>
            <h2>Evaluation report</h2>
            {evaluationReport.type === 'SHIPMENT' && evaluationReport.location !== 'OTHER' && (
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
                          evaluationReport.observedDate
                            ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY')
                            : '—',
                        ]}
                      />
                    )}
                    {evaluationReport.location === 'DESTINATION' && (
                      <DataTable
                        columnHeaders={['Observed delivery']}
                        dataRow={[
                          evaluationReport.observedDate
                            ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY')
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
                        evaluationReport.observedDate ? formatDate(evaluationReport.observedDate, 'DD MMM YYYY') : '—',
                      ]}
                    />
                  </DataTableWrapper>
                </div>
                <EvaluationReportList evaluationReport={evaluationReport} />
              </div>
            )}
            {(evaluationReport.type === 'COUNSELING' || evaluationReport.location === 'OTHER') && (
              <div className={styles.section}>
                <h3>Information</h3>
                <DataTableWrapper className={classnames(styles.detailsRight, 'table--data-point-group')}>
                  <DataTable
                    columnHeaders={['Inspection date', 'Report submission']}
                    dataRow={[
                      evaluationReport.inspectionDate
                        ? formatDate(evaluationReport.inspectionDate, 'DD MMM YYYY')
                        : '—',
                      evaluationReport.submittedAt
                        ? formatDateFromIso(evaluationReport.submittedAt, 'DD MMM YYYY')
                        : '—',
                    ]}
                  />
                </DataTableWrapper>
                <EvaluationReportList evaluationReport={evaluationReport} />
              </div>
            )}
            <div className={styles.section}>
              <h3>QAE remarks</h3>
              <dl className={descriptionListStyles.descriptionList}>
                <div className={descriptionListStyles.row}>
                  <dt>Evaluation remarks</dt>
                  <dd className={styles.qaeRemarks}>{evaluationReport.remarks}</dd>
                </div>
              </dl>
            </div>
          </div>
          <div className={styles.buttonsGroup}>
            <Button type="button" secondary onClick={() => setIsModalVisible(false)}>
              Close
            </Button>
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

EvaluationReportContainer.propTypes = {
  evaluationReportId: PropTypes.string.isRequired,
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  setIsModalVisible: PropTypes.func.isRequired,
  shipmentId: PropTypes.string,
};

EvaluationReportContainer.defaultProps = {
  shipmentId: '',
};

export default EvaluationReportContainer;
