import React, { useEffect } from 'react';
import { Button, Alert } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import reviewBillableWeightStyles from './ReviewBillableWeight.module.scss';

import ReviewDocumentsSidePanel from 'components/Office/PPM/ReviewDocumentsSidePanel/ReviewDocumentsSidePanel';
import ShipmentModificationTag from 'components/ShipmentModificationTag/ShipmentModificationTag';
import { WEIGHT_ADJUSTMENT, shipmentModificationTypes } from 'constants/shipments';
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
} from 'hooks/custom';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default function ReviewBillableWeight() {
  const [selectedShipmentIndex, setSelectedShipmentIndex] = React.useState(0);
  const [selectedShipment, setSelectedShipment] = React.useState({});
  const [sidebarType, setSidebarType] = React.useState('MAX');
  const [ppmShipmentInfo, setPpmShipmentInfo] = React.useState({});

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

  const getAllPPMShipmentFiles = (ppmShipment) => {
    const weightTicketDocs = [];
    if (ppmShipment !== undefined) {
      if (ppmShipment.weightTickets !== undefined) {
        weightTicketDocs.push(
          ppmShipment.weightTickets.map((weightTicket) => {
            return weightTicket.emptyDocument ?? [];
          }),
          ppmShipment.weightTickets.map((weightTicket) => {
            return weightTicket.fullDocument ?? [];
          }),
        );
      }
      if (ppmShipment.proGearWeightTickets !== undefined) {
        weightTicketDocs.push(
          ppmShipment.proGearWeightTickets.map((proGearWeightTicket) => {
            return proGearWeightTicket.document ?? [];
          }),
        );
      }
      if (ppmShipment.movingExpenses !== undefined) {
        weightTicketDocs.push(
          ppmShipment.movingExpenses.map((movingExpense) => {
            return movingExpense.document ?? [];
          }),
        );
      }
    }
    const uploadedWeightDocs = weightTicketDocs.flat();
    const uploadedWeightFiles = uploadedWeightDocs.flatMap((doc) => {
      return doc.uploads;
    });

    let uploads = [];
    uploads = uploads.concat(uploadedWeightFiles);
    return uploads;
  };

  /* Only show shipments in statuses of approved, diversion requested, or cancellation requested */
  const filteredShipments = mtoShipments?.filter((shipment) => includedStatusesForCalculatingWeights(shipment.status));
  const readOnly = true;
  const isLastShipment = filteredShipments && selectedShipmentIndex === filteredShipments.length - 1;

  const totalBillableWeight = useCalculatedTotalBillableWeight(filteredShipments, WEIGHT_ADJUSTMENT);
  const weightRequested = calculateWeightRequested(filteredShipments);
  const totalEstimatedWeight = useCalculatedEstimatedWeight(filteredShipments);

  const maxBillableWeight = order?.entitlement?.authorizedWeight;
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

  useEffect(() => {
    setSelectedShipment(
      filteredShipments && filteredShipments.length > 0 ? filteredShipments[selectedShipmentIndex] : {},
    );
  }, [filteredShipments, selectedShipmentIndex]);

  useEffect(() => {
    if (!isLoading && selectedShipment.shipmentType === 'PPM') {
      let currentTotalWeight = 0;
      selectedShipment.ppmShipment.weightTickets.forEach((weight) => {
        currentTotalWeight += weight.fullWeight - weight.emptyWeight;
      });
      const updatedPpmShipmentInfo = {
        ...selectedShipment.ppmShipment,
        miles: selectedShipment.distance,
        actualWeight: currentTotalWeight,
      };
      setPpmShipmentInfo(updatedPpmShipmentInfo);
    }
  }, [isLoading, selectedShipment]);

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

  const getEstimatedWeight = () => {
    if (selectedShipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
      return selectedShipment.ppmShipment.estimatedWeight;
    }

    if (selectedShipment.shipmentType === SHIPMENT_OPTIONS.NTSR) {
      return selectedShipment.ntsRecordedWeight;
    }
    return selectedShipment.primeEstimatedWeight;
  };

  const getOriginalWeight = () => {
    if (selectedShipment.shipmentType === SHIPMENT_OPTIONS.NTSR) {
      return selectedShipment.ntsRecordedWeight;
    }
    return selectedShipment.primeActualWeight;
  };

  const selectedShipmentIsDiverted = selectedShipment.diversion;
  const moveContainsDivertedShipment =
    selectedShipmentIsDiverted || filteredShipments ? filteredShipments.filter((s) => s.diversion).length > 0 : false;

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const fileList =
    selectedShipment.shipmentType !== 'PPM' ? getAllFiles() : getAllPPMShipmentFiles(selectedShipment.ppmShipment);
  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        <DocumentViewer files={fileList} />
      </div>
      <div className={reviewBillableWeightStyles.reviewWeightSideBar}>
        {sidebarType === 'MAX' ? (
          <DocumentViewerSidebar
            title="Review weights"
            subtitle="Edit max billable weight"
            onClose={handleClose}
            titleTag={
              moveContainsDivertedShipment ? (
                <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />
              ) : null
            }
          >
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
              {selectedShipment.shipmentType !== 'PPM' && (
                <EditBillableWeight
                  title="Max billable weight"
                  estimatedWeight={totalEstimatedWeight}
                  maxBillableWeight={maxBillableWeight}
                  editEntity={editEntity}
                  billableWeightJustification={move.tioRemarks}
                  isNTSRShipment={selectedShipment.shipmentType === SHIPMENT_OPTIONS.NTSR}
                />
              )}
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
            subtitleTag={
              selectedShipmentIsDiverted ? (
                <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />
              ) : null
            }
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
              {selectedShipment.shipmentType !== 'PPM' ? (
                <div className={reviewBillableWeightStyles.contentContainer}>
                  <ShipmentCard
                    billableWeight={selectedShipment.calculatedBillableWeight}
                    editEntity={editEntity}
                    billableWeightJustification={selectedShipment.billableWeightJustification}
                    dateReweighRequested={selectedShipment?.reweigh?.requestedAt}
                    departedDate={selectedShipment.actualPickupDate}
                    pickupAddress={selectedShipment.pickupAddress}
                    destinationAddress={selectedShipment.destinationAddress}
                    estimatedWeight={getEstimatedWeight()}
                    primeActualWeight={
                      selectedShipment.shipmentType !== SHIPMENT_OPTIONS.PPM
                        ? selectedShipment.primeActualWeight
                        : weightRequested
                    }
                    originalWeight={getOriginalWeight()}
                    adjustedWeight={selectedShipment.billableWeightCap}
                    reweighRemarks={selectedShipment?.reweigh?.verificationReason}
                    reweighWeight={selectedShipment?.reweigh?.weight}
                    maxBillableWeight={maxBillableWeight}
                    totalBillableWeight={totalBillableWeight}
                    shipmentType={selectedShipment.shipmentType}
                    storageFacilityAddress={selectedShipment.storageFacility?.address}
                  />
                </div>
              ) : (
                <ReviewDocumentsSidePanel
                  ppmShipment={selectedShipment.ppmShipment}
                  ppmShipmentInfo={ppmShipmentInfo}
                  ppmNumber={selectedShipment.shipmentLocator}
                  weightTickets={selectedShipment.ppmShipment.weightTickets}
                  proGearTickets={selectedShipment.ppmShipment.proGearWeightTickets}
                  expenseTickets={selectedShipment.ppmShipment.movingExpenses}
                  readOnly={readOnly}
                  showAllFields={false}
                />
              )}
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
