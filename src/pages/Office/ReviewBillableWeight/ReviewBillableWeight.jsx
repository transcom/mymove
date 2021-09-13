import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { Alert, Button } from '@trussworks/react-uswds';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MatchShape } from 'types/router';
import shipmentIsOverweight from 'utils/shipmentIsOverweight';

export const ReviewBillableWeight = ({ match }) => {
  const [selectedShipmentIndex, setSelectedShipmentIndex] = useState(0);

  const { moveCode } = match.params;
  const { move, mtoShipments } = useMoveTaskOrderQueries(moveCode);
  const handleClickNextButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex + 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const handleClickBackButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex - 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const isLastShipment = selectedShipmentIndex === mtoShipments.length - 1;

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
        <DocumentViewerSidebar
          title="Review weights"
          subtitle="Shipment weights"
          description={`Shipment ${selectedShipmentIndex + 1} of ${mtoShipments.length}`}
          onClose={() => {}}
        >
          <DocumentViewerSidebar.Content>
            <div style={{ width: '350px', backgroundColor: 'white', marginBottom: '16px' }}>
              <Alert slim type="error">
                {`Max billable weight exceeded. \nPlease resolve.`}
              </Alert>
              {!mtoShipments[0].reweighWeight && (
                <Alert slim type="warning">
                  Shipment missing information
                </Alert>
              )}
              {shipmentIsOverweight(mtoShipments[0].estimatedWeight, mtoShipments[0].billableWeight) && (
                <Alert slim type="warning">
                  Shipment exceeds 110% of estimated weight.
                </Alert>
              )}
              <WeightSummary
                maxBillableWeight={move.maxBillableWeight}
                totalBillableWeight={move.totalBillableWeight}
                weightRequested={move.weightRequested}
                weightAllowance={move.weightAllowance}
                totalBillableWeightFlag
                shipments={mtoShipments}
              />
            </div>
            <div style={{ height: '100%', width: '350px' }}>
              <ShipmentCard
                billableWeight={mtoShipments[0].billableWeight}
                dateReweighRequested="mtoShipments[0].dateReweighRequested"
                departedDate={mtoShipments[0].departedDate}
                pickupAddress={mtoShipments[selectedShipmentIndex].pickupAddress}
                destinationAddress={mtoShipments[0].destinationAddress}
                estimatedWeight={mtoShipments[0].estimatedWeight}
                originalWeight={mtoShipments[0].originalWeight}
                reweighRemarks={mtoShipments[0].reweighRemarks}
                reweighWeight={mtoShipments[0].reweighWeight}
              />
            </div>
          </DocumentViewerSidebar.Content>
          <DocumentViewerSidebar.Footer
            style={{ position: 'fixed', bottom: '0', width: '100%', backgroundColor: 'white' }}
          >
            <div style={{ display: 'flex' }}>
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
      </div>
    </div>
  );
};

ReviewBillableWeight.propTypes = {
  match: MatchShape.isRequired,
};

export default withRouter(connect(() => ({}))(ReviewBillableWeight));
