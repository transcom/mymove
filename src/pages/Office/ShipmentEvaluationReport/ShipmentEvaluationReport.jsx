import React, { useState } from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useHistory } from 'react-router';
import { queryCache, useMutation } from 'react-query';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import shipmentEvaluationReportStyles from './ShipmentEvaluationReport.module.scss';

import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { deleteEvaluationReport } from 'services/ghcApi';
import { SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';
import { useShipmentEvaluationReportQueries } from 'hooks/queries';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { shipmentTypeLabels } from 'content/shipments';
import DataTable from 'components/DataTable';
import { CustomerShape } from 'types';
import { OrdersShape } from 'types/customerShapes';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';

const ShipmentEvaluationReport = ({ customerInfo, orders }) => {
  const { moveCode, reportId } = useParams();
  const { evaluationReport, mtoShipment } = useShipmentEvaluationReportQueries(reportId);
  const mtoRefId = evaluationReport.moveReferenceID;

  const history = useHistory();

  const customerInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {ORDERS_RANK_OPTIONS[orders.grade]}
      <br />
      {ORDERS_BRANCH_OPTIONS[customerInfo.agency] ? ORDERS_BRANCH_OPTIONS[customerInfo.agency] : customerInfo.agency}
    </>
  );

  const officeUserInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {customerInfo.email}
    </>
  );

  const shipmentDisplayInfo = (shipment) => {
    return {
      ...shipment,
      heading: shipmentTypeLabels[shipment.shipmentType],
      isDiversion: shipment.diversion,
      shipmentStatus: shipment.status,
      destinationAddress: shipment.destinationAddress || '-',
    };
  };

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport, {
    onSuccess: async () => {
      await queryCache.invalidateQueries([SHIPMENT_EVALUATION_REPORTS, moveCode]);
    },
  });

  const [isDeleteModelOpen, setIsDeleteModelOpen] = useState(false);

  const toggleCancelModel = () => {
    setIsDeleteModelOpen(!isDeleteModelOpen);
  };

  const cancelReport = async () => {
    // Close the modal
    setIsDeleteModelOpen(!isDeleteModelOpen);

    // Mark as deleted in database
    await deleteEvaluationReportMutation(reportId);

    // Reroute back to eval report page, include flag to know to show alert
    history.push(`/moves/${moveCode}/evaluation-reports`, { showDeleteSuccess: true });
  };
  return (
    <>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModelOpen}
        closeModal={toggleCancelModel}
        submitModal={cancelReport}
      />
      <div className={classnames(styles.tabContent, shipmentEvaluationReportStyles.tabContent)}>
        <GridContainer>
          <div className={styles.pageHeader}>
            <h1>Evaluation report</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>REPORT ID #{reportId}</h6>
              <h6>MOVE CODE {moveCode}</h6>
              <h6>MTO REFERENCE ID #{mtoRefId}</h6>
            </div>
          </div>
        </GridContainer>
        <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8 }}>
              <h2>Shipment information</h2>
              {mtoShipment.id && (
                <ShipmentDisplay
                  isSubmitted
                  shipmentId={mtoShipment.id}
                  displayInfo={shipmentDisplayInfo(mtoShipment)}
                  shipmentType={mtoShipment.type}
                />
              )}
            </Grid>
            <Grid className={shipmentEvaluationReportStyles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
              <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
              <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
            </Grid>
          </Grid>
        </GridContainer>

        <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <h2>Evaluation form</h2>
            </Grid>
          </Grid>
        </GridContainer>
        <div style={{ display: 'flex', float: 'right' }}>
          <Button className="usa-button--unstyled" onClick={toggleCancelModel}>
            Cancel
          </Button>
          <Button className="usa-button--secondary">Save draft</Button>
          <Button type="submit">Submit</Button>
        </div>
      </div>
    </>
  );
};

ShipmentEvaluationReport.propTypes = {
  customerInfo: CustomerShape.isRequired,
  orders: OrdersShape.isRequired,
};

export default ShipmentEvaluationReport;
