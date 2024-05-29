import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import classnames from 'classnames';

import styles from './HeaderSection.module.scss';

import EditPPMHeaderSummaryModal from 'components/Office/PPM/PPMHeaderSummary/EditPPMHeaderSummaryModal';
import { formatDate, formatCents, formatWeight } from 'utils/formatters';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { updateMTOShipment } from 'services/ghcApi';
import { useEditShipmentQueries, usePPMShipmentDocsQueries } from 'hooks/queries';

export const sectionTypes = {
  incentives: 'incentives',
  shipmentInfo: 'shipmentInfo',
  incentiveFactors: 'incentiveFactors',
};

const HAUL_TYPES = {
  SHORTHAUL: 'Shorthaul',
  LINEHAUL: 'Linehaul',
};

const getSectionTitle = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.incentiveFactors:
      return `Incentive Factors`;
    default:
      return <Alert>Error getting section title!</Alert>;
  }
};

const OpenModalButton = ({ onClick }) => (
  <Button type="button" data-testid="editTextButton" className={styles['edit-btn']} onClick={onClick}>
    <span>
      <FontAwesomeIcon icon="pencil" style={{ marginRight: '5px' }} />
    </span>
  </Button>
);

// Returns the markup needed for a specific section
const getSectionMarkup = (sectionInfo, handleEditOnClick) => {
  const aoaRequestedValue = `$${formatCents(sectionInfo.advanceAmountRequested)}`;
  const aoaValue = `$${formatCents(sectionInfo.advanceAmountReceived)}`;

  const renderHaulType = (haulType) => {
    return haulType === HAUL_TYPES.LINEHAUL ? 'Linehaul' : 'Shorthaul';
  };

  switch (sectionInfo.type) {
    case sectionTypes.shipmentInfo:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Planned Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.plannedMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label>Actual Move Start Date</Label>
            <span className={styles.light}>
              {formatDate(sectionInfo.actualMoveDate, null, 'DD-MMM-YYYY')}
              <OpenModalButton onClick={() => handleEditOnClick(sectionInfo.type, 'actualMoveDate')} />
            </span>
          </div>
          <div>
            <Label>Starting ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualPickupPostalCode}</span>
          </div>
          <div>
            <Label>Ending ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualDestinationPostalCode}</span>
          </div>
          <div>
            <Label>Miles</Label>
            <span className={styles.light}>{sectionInfo.miles}</span>
          </div>
          <div>
            <Label>Estimated Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.estimatedWeight)}</span>
          </div>
          <div>
            <Label>Actual Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.actualWeight)}</span>
          </div>
        </div>
      );

    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Government Constructed Cost (GCC)</Label>
            <span data-testid="gcc" className={styles.light}>
              ${formatCents(sectionInfo.gcc)}
            </span>
          </div>
          <div>
            <Label>Gross Incentive</Label>
            <span data-testid="grossIncentive" className={styles.light}>
              ${formatCents(sectionInfo.grossIncentive)}
            </span>
          </div>
          <div>
            <Label>Advance Requested</Label>
            <span data-testid="advanceRequested" className={styles.light}>
              {aoaRequestedValue}
            </span>
          </div>
          <div>
            <Label>Advance Received</Label>
            <span data-testid="advanceReceived" className={styles.light}>
              {aoaValue}
              <OpenModalButton onClick={() => handleEditOnClick(sectionInfo.type, 'advanceAmountReceived')} />
            </span>
          </div>
          <div>
            <Label>Remaining Incentive</Label>
            <span data-testid="remainingIncentive" className={styles.light}>
              ${formatCents(sectionInfo.remainingIncentive)}
            </span>
          </div>
        </div>
      );

    case sectionTypes.incentiveFactors:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>{renderHaulType(sectionInfo.haulType)} Price</Label>
            <span data-testid="haulPrice" className={styles.light}>
              ${formatCents(sectionInfo.haulPrice)}
            </span>
          </div>
          <div>
            <Label>{renderHaulType(sectionInfo.haulType)} Fuel Rate Adjustment</Label>
            <span data-testid="haulFSC" className={styles.light}>
              {sectionInfo.haulFSC < 0 ? '-$' : '$'}
              {formatCents(Math.abs(sectionInfo.haulFSC))}
            </span>
          </div>
          <div>
            <Label>Packing Charge</Label>
            <span data-testid="packPrice" className={styles.light}>
              ${formatCents(sectionInfo.packPrice)}
            </span>
          </div>
          <div>
            <Label>Unpacking Charge</Label>
            <span data-testid="unpackPrice" className={styles.light}>
              ${formatCents(sectionInfo.unpackPrice)}
            </span>
          </div>
          <div>
            <Label>Origin Price</Label>
            <span data-testid="originPrice" className={styles.light}>
              ${formatCents(sectionInfo.dop)}
            </span>
          </div>
          <div>
            <Label>Destination Price</Label>
            <span data-testid="destinationPrice" className={styles.light}>
              ${formatCents(sectionInfo.ddp)}
            </span>
          </div>
        </div>
      );

    default:
      return <Alert>An error occured while getting section markup!</Alert>;
  }
};

export default function PPMHeaderSummary({ sectionInfo, dataTestId }) {
  const { shipmentId, moveCode } = useParams();
  const { mtoShipment, refetchMTOShipment } = usePPMShipmentDocsQueries(shipmentId);
  const queryClient = useQueryClient();
  const requestDetailsButtonTestId = `${sectionInfo.type}-showRequestDetailsButton`;
  const [showDetails, setShowDetails] = useState(false);
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [sectionName, setSectionName] = useState('');
  const [sectionType, setSectionType] = useState('');
  const [currentSectionInfo, setCurrentSectionInfo] = useState(sectionInfo);
  const [isUpdated, setIsUpdated] = useState(false);

  const showRequestDetailsButton = true;
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);
  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide details' : 'Show details';

  const handleEditOnClose = () => {
    setIsEditModalVisible(false);
    setSectionName('');
    setSectionType('');
  };
  // this is to avoid state issues
  useEffect(() => {
    if (!isUpdated) {
      setCurrentSectionInfo(sectionInfo);
    }
  }, [isUpdated, sectionInfo]);

  // fetch updated shipment data whenever edit modal is opened
  useEffect(() => {
    if (isEditModalVisible) {
      refetchMTOShipment();
    }
  }, [isEditModalVisible, refetchMTOShipment]);

  const { mtoShipments } = useEditShipmentQueries(moveCode);

  const { mutate: mutateMTOShipment } = useMutation(updateMTOShipment, {
    onSuccess: (updatedMTOShipments) => {
      const updatedMTOShipment = updatedMTOShipments.mtoShipments[shipmentId];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);

      setIsUpdated(true);

      setCurrentSectionInfo((prev) => ({
        ...prev,
        actualMoveDate: updatedMTOShipment.ppmShipment.actualMoveDate,
        advanceAmountReceived: updatedMTOShipment.ppmShipment.advanceAmountReceived,
      }));

      handleEditOnClose();
    },
  });

  const handleEditOnClick = (type, name) => {
    setIsEditModalVisible(true);
    setSectionName(name);
    setSectionType(type);
  };

  const handleEditSubmit = (values) => {
    let body = {};

    switch (sectionName) {
      case 'actualMoveDate':
        body = { actualMoveDate: formatDate(values.actualMoveDate, 'DD MMM YYYY', 'YYYY-MM-DD') };
        break;
      case 'advanceAmountReceived':
        if (values.advanceAmountReceived === '0') {
          body = {
            advanceAmountReceived: null,
            hasReceivedAdvance: false,
          };
        } else {
          body = {
            advanceAmountReceived: values.advanceAmountReceived * 100,
            hasReceivedAdvance: true,
          };
        }
        break;

      default:
        break;
    }

    mutateMTOShipment({
      moveTaskOrderID: mtoShipment.moveTaskOrderID,
      shipmentID: mtoShipment.id,
      ifMatchETag: mtoShipment.eTag,
      body: {
        ppmShipment: body,
      },
    });
  };

  return (
    <section className={classnames(styles.HeaderSection)} data-testid={dataTestId}>
      <header>
        <h4>{getSectionTitle(currentSectionInfo)}</h4>
      </header>
      <div className={styles.toggleDrawer}>
        {showRequestDetailsButton && (
          <Button
            aria-expanded={showDetails}
            data-testid={requestDetailsButtonTestId}
            type="button"
            unstyled
            onClick={handleToggleDetails}
          >
            <FontAwesomeIcon icon={showDetailsChevron} /> {showDetailsText}
          </Button>
        )}
      </div>
      {showDetails && getSectionMarkup(currentSectionInfo, handleEditOnClick)}
      {isEditModalVisible && (
        <EditPPMHeaderSummaryModal
          onClose={handleEditOnClose}
          onSubmit={handleEditSubmit}
          sectionInfo={currentSectionInfo}
          editSectionName={sectionName}
          sectionType={sectionType}
        />
      )}
    </section>
  );
}
