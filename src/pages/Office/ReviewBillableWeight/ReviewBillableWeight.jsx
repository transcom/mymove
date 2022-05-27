import React from 'react';
import { Button, Alert } from '@trussworks/react-uswds';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import reviewBillableWeightStyles from './ReviewBillableWeight.module.scss';

import { MOVES, MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import { updateMTOShipment, updateMaxBillableWeightAsTIO, updateTIORemarks } from 'services/ghcApi';
import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import EditBillableWeight from 'components/Office/BillableWeight/EditBillableWeight/EditBillableWeight';
import { useOrdersDocumentQueries, useMovePaymentRequestsQueries } from 'hooks/queries';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
  useCalculatedEstimatedWeight,
} from 'hooks/custom';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SHIPMENT_OPTIONS } from 'shared/constants';

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
  const { order, mtoShipments, move } = useMovePaymentRequestsQueries(moveCode);
  /* Only show shipments in statuses of approved, diversion requested, or cancellation requested */
  const filteredShipments = mtoShipments?.filter((shipment) => includedStatusesForCalculatingWeights(shipment.status));
  const isLastShipment = filteredShipments && selectedShipmentIndex === filteredShipments.length - 1;

  const totalBillableWeight = useCalculatedTotalBillableWeight(filteredShipments);
  const weightRequested = useCalculatedWeightRequested(filteredShipments);
  const totalEstimatedWeight = useCalculatedEstimatedWeight(filteredShipments);

  const maxBillableWeight = order?.entitlement?.authorizedWeight;
  const weightAllowance = order?.entitlement?.totalWeight;

  const shipmentsMissingInformation = filteredShipments?.filter((shipment) => {
    return !shipment.primeEstimatedWeight || (shipment.reweigh?.requestedAt && !shipment.reweigh?.weight);
  });

  const handleClose = () => {
    history.push(generatePath(tioRoutes.PAYMENT_REQUESTS_PATH, { moveCode }), {
      from: 'review-billable-weights',
    });
  };

  const selectedShipment = filteredShipments ? filteredShipments[selectedShipmentIndex] : {};

  const [mutateMTOShipment] = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      filteredShipments[filteredShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], filteredShipments);
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

  const [mutateMoves] = useMutation(updateTIORemarks, {
    onSuccess: (data, variables) => {
      const updatedMove = data.moves[variables.moveTaskOrderID];
      queryCache.setQueryData([MOVES, variables.moveTaskOrderID], {
        moves: {
          [`${variables.moveTaskOrderID}`]: updatedMove,
        },
      });
      queryCache.invalidateQueries([MOVES, move.locator]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const editEntity = (formValues) => {
    if (sidebarType === 'MAX') {
      const orderPayload = {
        orderID: order.id,
        ifMatchETag: order.eTag,
        body: {
          authorizedWeight: Number(formValues.billableWeight),
          tioRemarks: formValues.billableWeightJustification,
        },
      };
      const movePayload = {
        moveTaskOrderID: move.id,
        ifMatchETag: move.eTag,
        body: {
          tioRemarks: formValues.billableWeightJustification,
        },
      };
      mutateOrders(orderPayload);
      mutateMoves(movePayload);
    } else {
      const payload = {
        body: {
          ...formValues,
          billableWeightCap: Number(formValues.billableWeight),
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
              {totalBillableWeight > maxBillableWeight && (
                <Alert headingLevel="h4" slim type="error" data-testid="maxBillableWeightAlert">
                  {`Max billable weight exceeded. \nPlease resolve.`}
                </Alert>
              )}
              {shipmentsMissingInformation?.length > 0 && (
                <Alert headingLevel="h4" slim type="warning" data-testid="maxBillableWeightMissingShipmentWeightAlert">
                  Missing shipment weights may impact max billable weight.
                </Alert>
              )}
              <div className={reviewBillableWeightStyles.weightSummary}>
                <WeightSummary
                  maxBillableWeight={maxBillableWeight}
                  totalBillableWeight={totalBillableWeight}
                  weightRequested={weightRequested}
                  weightAllowance={weightAllowance}
                  totalBillableWeightFlag={totalBillableWeight > maxBillableWeight}
                  shipments={filteredShipments}
                />
              </div>
              <EditBillableWeight
                title="Max billable weight"
                estimatedWeight={totalEstimatedWeight}
                maxBillableWeight={maxBillableWeight}
                weightAllowance={weightAllowance}
                editEntity={editEntity}
                billableWeightJustification={move.tioRemarks}
                isNTSRShipment={selectedShipment.shipmentType === SHIPMENT_OPTIONS.NTSR}
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
            description={`Shipment ${selectedShipmentIndex + 1} of ${filteredShipments?.length}`}
            onClose={handleClose}
          >
            <DocumentViewerSidebar.Content>
              <div className={reviewBillableWeightStyles.contentContainer}>
                {totalBillableWeight > maxBillableWeight && (
                  <Alert headingLevel="h4" slim type="error" data-testid="maxBillableWeightAlert">
                    {`Max billable weight exceeded. \nPlease resolve.`}
                  </Alert>
                )}
                {((!selectedShipment?.reweigh?.weight && selectedShipment?.reweigh?.requestedAt) ||
                  !selectedShipment.primeEstimatedWeight) && (
                  <Alert headingLevel="h4" slim type="warning" data-testid="shipmentMissingInformation">
                    Shipment missing information
                  </Alert>
                )}
                {shipmentIsOverweight(
                  selectedShipment.primeEstimatedWeight,
                  selectedShipment.calculatedBillableWeight,
                ) &&
                  selectedShipment.shipmentType !== SHIPMENT_OPTIONS.NTSR && (
                    <Alert
                      headingLevel="h4"
                      slim
                      type="warning"
                      data-testid="shipmentBillableWeightExceeds110OfEstimated"
                    >
                      Shipment exceeds 110% of estimated weight.
                    </Alert>
                  )}
                <div className={reviewBillableWeightStyles.weightSummary}>
                  <WeightSummary
                    maxBillableWeight={maxBillableWeight}
                    totalBillableWeight={totalBillableWeight}
                    weightRequested={weightRequested}
                    weightAllowance={weightAllowance}
                    shipments={filteredShipments}
                  />
                </div>
              </div>
              <div className={reviewBillableWeightStyles.contentContainer}>
                <ShipmentCard
                  billableWeight={selectedShipment.calculatedBillableWeight}
                  editEntity={editEntity}
                  billableWeightJustification={selectedShipment.billableWeightJustification}
                  dateReweighRequested={selectedShipment?.reweigh?.requestedAt}
                  departedDate={selectedShipment.actualPickupDate}
                  pickupAddress={selectedShipment.pickupAddress}
                  destinationAddress={selectedShipment.destinationAddress}
                  estimatedWeight={selectedShipment.primeEstimatedWeight}
                  originalWeight={selectedShipment.primeActualWeight}
                  adjustedWeight={selectedShipment.billableWeightCap}
                  reweighRemarks={selectedShipment?.reweigh?.verificationReason}
                  reweighWeight={selectedShipment?.reweigh?.weight}
                  maxBillableWeight={maxBillableWeight}
                  totalBillableWeight={totalBillableWeight}
                  shipmentType={selectedShipment.shipmentType}
                  storageFacilityAddress={selectedShipment.storageFacility?.address}
                />
              </div>
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer className={reviewBillableWeightStyles.footer}>
              <div className={reviewBillableWeightStyles.flex}>
                <Button onClick={handleClickBackButton} secondary>
                  Back
                </Button>
                {!isLastShipment && <Button onClick={handleClickNextButton}>Next Shipment</Button>}
              </div>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
