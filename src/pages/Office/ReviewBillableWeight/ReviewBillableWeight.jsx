import React from 'react';
import { Button, Alert } from '@trussworks/react-uswds';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import reviewBillableWeightStyles from './ReviewBillableWeight.module.scss';

import { MOVES, MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import { updateMTOShipment, updateMaxBillableWeightAsTIO } from 'services/ghcApi';
import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import EditBillableWeight from 'components/Office/BillableWeight/EditBillableWeight/EditBillableWeight';
import { useOrdersDocumentQueries, useMovePaymentRequestsQueries } from 'hooks/queries';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import {
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
  useCalculatedEstimatedWeight,
} from 'hooks/custom';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

export default function ReviewBillableWeight() {
  const [selectedShipmentIndex, setSelectedShipmentIndex] = React.useState(0);
  const [sidebarType, setSidebarType] = React.useState('MAX');

  const { moveCode } = useParams();
  const handleClickNextButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex + 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const handleClickBackButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex - 1;
    if (newSelectedShipmentIdx >= 0) {
      setSelectedShipmentIndex(newSelectedShipmentIdx);
    } else {
      setSidebarType('MAX');
    }
  };

  const history = useHistory();
  let documentsForViewer = [];
  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const { order, mtoShipments } = useMovePaymentRequestsQueries(moveCode);
  const isLastShipment = selectedShipmentIndex === mtoShipments?.length - 1;

  const totalBillableWeight = useCalculatedTotalBillableWeight(mtoShipments);
  const weightRequested = useCalculatedWeightRequested(mtoShipments);
  const totalEstimatedWeight = useCalculatedEstimatedWeight(mtoShipments);

  const maxBillableWeight = order.entitlement.authorizedWeight;
  const weightAllowance = order.entitlement.totalWeight;

  const handleClose = () => {
    history.push(generatePath(tioRoutes.PAYMENT_REQUESTS_PATH, { moveCode }));
  };

  const selectedShipment = mtoShipments[selectedShipmentIndex];

  const [mutateMTOShipment] = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const [mutateOrders] = useMutation(updateMaxBillableWeightAsTIO, {
    onSuccess: (data, variables) => {
      queryCache.invalidateQueries([MOVES, moveCode]);
      const updatedOrder = data.orders[variables.orderID];
      queryCache.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryCache.invalidateQueries([ORDERS, variables.orderID]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const editEntity = (formValues) => {
    if (sidebarType === 'MAX') {
      const payload = {
        orderID: order.id,
        ifMatchETag: order.eTag,
        body: {
          authorizedWeight: formValues.billableWeight,
          tioRemarks: formValues.billableWeightJustification,
        },
      };
      mutateOrders(payload);
    } else {
      const payload = {
        body: {
          ...formValues,
          billableWeightCap: formValues.billableWeight,
        },
        ifMatchETag: selectedShipment.eTag,
        moveTaskOrderID: selectedShipment.moveTaskOrderID,
        shipmentID: selectedShipment.id,
        normalize: false,
      };
      mutateMTOShipment(payload);
    }
  };

  if (upload) {
    documentsForViewer = Object.values(upload);
  }

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        <DocumentViewer files={documentsForViewer} />
      </div>
      <div className={styles.sidebar}>
        {sidebarType === 'MAX' ? (
          <DocumentViewerSidebar title="Review weights" subtitle="Edit max billable weight" onClose={handleClose}>
            <DocumentViewerSidebar.Content>
              {maxBillableWeight > weightAllowance && (
                <Alert slim type="error">
                  {`Max billable weight exceeded. \nPlease resolve.`}
                </Alert>
              )}
              <div className={reviewBillableWeightStyles.weightSummary}>
                <WeightSummary
                  maxBillableWeight={maxBillableWeight}
                  totalBillableWeight={totalBillableWeight}
                  weightRequested={weightRequested}
                  weightAllowance={weightAllowance}
                  shipments={mtoShipments}
                />
              </div>
              <EditBillableWeight
                title="Max billable weight"
                estimatedWeight={totalEstimatedWeight}
                maxBillableWeight={maxBillableWeight}
                weightAllowance={weightAllowance}
                editEntity={editEntity}
                billableWeightJustification={order.moveTaskOrder.tioRemarks}
              />
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer>
              <Button
                onClick={() => {
                  setSidebarType('SHIPMENT');
                }}
              >
                Review shipment weights
              </Button>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        ) : (
          <DocumentViewerSidebar
            title="Review weights"
            subtitle="Shipment weights"
            description={`Shipment ${selectedShipmentIndex + 1} of ${mtoShipments?.length}`}
            onClose={handleClose}
          >
            <DocumentViewerSidebar.Content>
              <div className={reviewBillableWeightStyles.contentContainer}>
                {((!selectedShipment.reweigh?.weight && selectedShipment.reweigh?.requestedAt) ||
                  !selectedShipment.primeEstimatedWeight) && (
                  <Alert slim type="warning">
                    Shipment missing information
                  </Alert>
                )}
                {shipmentIsOverweight(selectedShipment.primeEstimatedWeight, selectedShipment.primeActualWeight) && (
                  <Alert slim type="warning">
                    Shipment exceeds 110% of estimated weight.
                  </Alert>
                )}
                <div className={reviewBillableWeightStyles.weightSummary}>
                  <WeightSummary
                    maxBillableWeight={maxBillableWeight}
                    totalBillableWeight={totalBillableWeight}
                    weightRequested={weightRequested}
                    weightAllowance={weightAllowance}
                    totalBillableWeightFlag
                    shipments={mtoShipments}
                  />
                </div>
              </div>
              <div className={reviewBillableWeightStyles.contentContainer}>
                <ShipmentCard
                  editEntity={editEntity}
                  billableWeight={selectedShipment.billableWeightCap}
                  billableWeightJustification={selectedShipment.billableWeightJustification}
                  dateReweighRequested={selectedShipment.reweigh?.requestedAt}
                  departedDate={selectedShipment.actualPickupDate}
                  pickupAddress={selectedShipment.pickupAddress}
                  destinationAddress={selectedShipment.destinationAddress}
                  estimatedWeight={selectedShipment.primeEstimatedWeight}
                  originalWeight={selectedShipment.primeActualWeight}
                  reweighRemarks={selectedShipment.reweigh?.verificationReason}
                  reweighWeight={selectedShipment.reweigh?.weight}
                />
              </div>
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer className={reviewBillableWeightStyles.footer}>
              <div className={reviewBillableWeightStyles.flex}>
                <Button type="button" onClick={handleClickBackButton} secondary>
                  Back
                </Button>
                {!isLastShipment && (
                  <Button type="button" onClick={handleClickNextButton}>
                    Next Shipment
                  </Button>
                )}
              </div>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
