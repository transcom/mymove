import React from 'react';
import { Button, Alert } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import reviewBillableWeightStyles from './ReviewBillableWeight.module.scss';

import { WEIGHT_ADJUSTMENT } from 'constants/shipments';
import { MOVES, MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import { updateMTOShipment, updateMaxBillableWeightAsTIO, updateTIORemarks } from 'services/ghcApi';
import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import EditBillableWeight from 'components/Office/BillableWeight/EditBillableWeight/EditBillableWeight';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { milmoveLogger } from 'utils/milmoveLog';
import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedEstimatedWeight,
  calculateWeightRequested,
  calculateEstimatedWeight,
} from 'hooks/custom';
import { shipmentIsOverweight, getDisplayWeight } from 'utils/shipmentWeights';
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

  const navigate = useNavigate();
  const { paymentRequests, order, mtoShipments, move, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);

  const getAllFiles = () => {
    const proofOfServiceDocs = paymentRequests.map((paymentRequest) => {
      return paymentRequest.proofOfServiceDocs ?? [];
    });

    const uploadedDocs = proofOfServiceDocs.flat();
    const uploadedFiles = uploadedDocs.flatMap((doc) => {
      return doc.uploads;
    });

    let uploads = [];
    uploads = uploads.concat(uploadedFiles);
    return uploads;
  };

  // filter out PPMs, as they're not including in TIO review
  const excludePPMShipments = mtoShipments?.filter((shipment) => shipment.shipmentType !== 'PPM');
  /* Only show shipments in statuses of approved, diversion requested, or cancellation requested */
  const filteredShipments = excludePPMShipments?.filter((shipment) =>
    includedStatusesForCalculatingWeights(shipment.status),
  );
  const isLastShipment = filteredShipments && selectedShipmentIndex === filteredShipments.length - 1;

  const totalBillableWeight = useCalculatedTotalBillableWeight(filteredShipments, WEIGHT_ADJUSTMENT);
  const weightRequested = calculateWeightRequested(filteredShipments);
  const totalEstimatedWeight = useCalculatedEstimatedWeight(filteredShipments);

  const maxBillableWeight = calculateEstimatedWeight(filteredShipments, undefined, WEIGHT_ADJUSTMENT);
  const weightAllowance = order?.entitlement?.totalWeight;

  const shipmentsMissingInformation = filteredShipments?.filter((shipment) => {
    return !shipment.primeEstimatedWeight || (shipment.reweigh?.requestedAt && !shipment.reweigh?.weight);
  });

  const handleClose = () => {
    navigate(`../${tioRoutes.PAYMENT_REQUESTS_PATH}`, {
      state: {
        from: 'review-billable-weights',
      },
    });
  };

  const selectedShipment = filteredShipments ? filteredShipments[selectedShipmentIndex] : {};

  const queryClient = useQueryClient();
  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      filteredShipments[filteredShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], filteredShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const { mutate: mutateOrders } = useMutation(updateMaxBillableWeightAsTIO, {
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [MOVES, moveCode] });
      queryClient.invalidateQueries({ queryKey: [ORDERS, variables.orderID] });
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const { mutate: mutateMoves } = useMutation(updateTIORemarks, {
    onSuccess: (data, variables) => {
      const updatedMove = data.moves[variables.moveTaskOrderID];
      queryClient.setQueryData([MOVES, variables.moveTaskOrderID], {
        moves: {
          [`${variables.moveTaskOrderID}`]: updatedMove,
        },
      });
      queryClient.invalidateQueries([MOVES, move.locator]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        <DocumentViewer files={getAllFiles()} />
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
                  billableWeight={getDisplayWeight(selectedShipment, WEIGHT_ADJUSTMENT)}
                  editEntity={editEntity}
                  billableWeightJustification={selectedShipment.billableWeightJustification}
                  dateReweighRequested={selectedShipment?.reweigh?.requestedAt}
                  departedDate={selectedShipment.actualPickupDate}
                  pickupAddress={selectedShipment.pickupAddress}
                  destinationAddress={selectedShipment.destinationAddress}
                  estimatedWeight={
                    selectedShipment.shipmentType !== SHIPMENT_OPTIONS.PPM
                      ? selectedShipment.primeEstimatedWeight
                      : selectedShipment.ppmShipment.estimatedWeight
                  }
                  primeActualWeight={
                    selectedShipment.shipmentType !== SHIPMENT_OPTIONS.PPM
                      ? selectedShipment.primeActualWeight
                      : weightRequested
                  }
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
                {isLastShipment && <Button onClick={handleClose}>Done</Button>}
              </div>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
