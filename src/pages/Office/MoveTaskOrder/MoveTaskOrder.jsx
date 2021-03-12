import React, { useState, useEffect, useMemo } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { MTO_SERVICE_ITEMS } from 'constants/queryKeys';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MatchShape } from 'types/router';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';
import { MOVE_STATUSES } from 'shared/constants';
import { patchMTOServiceItemStatus } from 'services/ghcApi';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import dimensionTypes from 'constants/dimensionTypes';
import customerContactTypes from 'constants/customerContactTypes';
import { mtoShipmentTypes, shipmentStatuses } from 'constants/shipments';
import LeftNav from 'components/LeftNav';
import { shipmentSectionLabels } from 'content/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const weekday = new Intl.DateTimeFormat('en', { weekday: 'long' }).format(dateObj);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${weekday}, ${day} ${month} ${year}`;
}

function approvedFilter(shipment) {
  return shipment.status === shipmentStatuses.APPROVED || shipment.status === shipmentStatuses.CANCELLATION_REQUESTED;
}

export const MoveTaskOrder = ({ match, ...props }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);
  const [sections, setSections] = useState([]);
  const [activeSection, setActiveSection] = useState('');

  const { moveCode } = match.params;
  const { setUnapprovedShipmentCount } = props;

  const {
    moveOrders = {},
    moveTaskOrders,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
  } = useMoveTaskOrderQueries(moveCode);

  const mtoServiceItemsArr = Object.values(mtoServiceItems || {});
  const moveOrder = Object.values(moveOrders)?.[0];
  const moveTaskOrder = Object.values(moveTaskOrders || {})?.[0];

  const shipmentServiceItems = useMemo(() => {
    const serviceItemsForShipment = {};
    mtoServiceItemsArr?.forEach((item) => {
      // We're not interested in basic service items
      if (!item.mtoShipmentID) {
        return;
      }
      const newItem = { ...item };
      newItem.code = item.reServiceCode;
      newItem.serviceItem = item.reServiceName;
      newItem.details = {
        pickupPostalCode: item.pickupPostalCode,
        reason: item.reason,
        imgURL: '',
        description: item.description,
        itemDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.ITEM),
        crateDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.CRATE),
        firstCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.FIRST),
        secondCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.SECOND),
      };

      if (serviceItemsForShipment[`${newItem.mtoShipmentID}`]) {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`].push(newItem);
      } else {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`] = [newItem];
      }
    });
    return serviceItemsForShipment;
  }, [mtoServiceItemsArr]);

  const [mutateMTOServiceItemStatus] = useMutation(patchMTOServiceItemStatus, {
    onSuccess: (data, variables) => {
      const newMTOServiceItem = data.mtoServiceItems[variables.mtoServiceItemID];
      queryCache.setQueryData([MTO_SERVICE_ITEMS, variables.moveTaskOrderId, true], {
        mtoServiceItems: {
          ...mtoServiceItems,
          [`${variables.mtoServiceItemID}`]: newMTOServiceItem,
        },
      });
      queryCache.invalidateQueries(MTO_SERVICE_ITEMS, variables.moveTaskOrderId);
      setIsModalVisible(false);
      setSelectedServiceItem({});
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  const handleUpdateMTOServiceItemStatus = (mtoServiceItemID, mtoShipmentID, status, rejectionReason) => {
    const mtoServiceItemForRequest = shipmentServiceItems[`${mtoShipmentID}`]?.find((s) => s.id === mtoServiceItemID);

    mutateMTOServiceItemStatus({
      moveTaskOrderId: moveTaskOrder.id,
      mtoServiceItemID,
      status,
      rejectionReason,
      ifMatchEtag: mtoServiceItemForRequest.eTag,
    });
  };

  useEffect(() => {
    if (mtoShipments) {
      const shipmentCount = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.status === shipmentStatuses.SUBMITTED).length
        : 0;
      setUnapprovedShipmentCount(shipmentCount);
    }
  }, [mtoShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    const shipmentSections = [];
    mtoShipments?.forEach((shipment) => {
      if (shipment.status === shipmentStatuses.APPROVED) {
        shipmentSections.push({
          id: shipment.id,
          label: shipmentSectionLabels[`${shipment.shipmentType}`] || shipment.shipmentType,
        });
      }
    });
    setSections(shipmentSections);
  }, [mtoShipments]);

  const handleScroll = () => {
    const distanceFromTop = window.scrollY;
    let newActiveSection;

    sections.forEach((section) => {
      const sectionEl = document.querySelector(`#shipment-${section.id}`);
      if (sectionEl?.offsetTop <= distanceFromTop && sectionEl?.offsetTop + sectionEl?.offsetHeight > distanceFromTop) {
        newActiveSection = section.id;
      }
    });

    if (activeSection !== newActiveSection) {
      setActiveSection(newActiveSection);
    }
  };

  useEffect(() => {
    // attach scroll listener
    window.addEventListener('scroll', handleScroll);

    // remove scroll listener
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  });

  const handleShowRejectionDialog = (mtoServiceItemID, mtoShipmentID) => {
    const serviceItem = shipmentServiceItems[`${mtoShipmentID}`]?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsModalVisible(true);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  if (moveTaskOrder.status === MOVE_STATUSES.SUBMITTED || !mtoShipments.some(approvedFilter)) {
    return (
      <div className={styles.tabContent}>
        <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
          </div>
          <div className={styles.emptyMessage}>
            <p>This move does not have any approved shipments yet.</p>
          </div>
        </GridContainer>
      </div>
    );
  }

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav className={styles.sidebar}>
          {sections.map((s) => {
            const classes = classnames({ active: s.id === activeSection });
            return (
              <a key={`sidenav_${s.id}`} href={`#shipment-${s.id}`} className={classes}>
                {s.label}
              </a>
            );
          })}
        </LeftNav>
        <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
          {isModalVisible && (
            <RejectServiceItemModal
              serviceItem={selectedServiceItem}
              onSubmit={handleUpdateMTOServiceItemStatus}
              onClose={setIsModalVisible}
            />
          )}
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>MTO Reference ID #{moveTaskOrder?.referenceId}</h6>
              <h6>Contract #1234567890</h6> {/* TODO - need this value from the API */}
            </div>
          </div>
          {mtoShipments.map((mtoShipment) => {
            if (
              mtoShipment.status !== shipmentStatuses.APPROVED &&
              mtoShipment.status !== shipmentStatuses.CANCELLATION_REQUESTED
            ) {
              return false;
            }
            const serviceItemsForShipment = shipmentServiceItems[`${mtoShipment.id}`];
            const requestedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.SUBMITTED,
            );
            const approvedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.APPROVED,
            );
            const rejectedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.REJECTED,
            );
            // eslint-disable-next-line camelcase
            const dutyStationPostal = { postal_code: moveOrder.destinationDutyStation.address.postal_code };
            const { pickupAddress, destinationAddress } = mtoShipment;
            const formattedScheduledPickup = formatShipmentDate(mtoShipment.scheduledPickupDate);
            return (
              <div id={`shipment-${mtoShipment.id}`} key={mtoShipment.id}>
                <ShipmentContainer shipmentType={mtoShipment.shipmentType} className={styles.shipmentCard}>
                  <ShipmentHeading
                    shipmentInfo={{
                      shipmentType: mtoShipmentTypes[mtoShipment.shipmentType],
                      originCity: pickupAddress?.city,
                      originState: pickupAddress?.state,
                      originPostalCode: pickupAddress?.postal_code,
                      destinationAddress: destinationAddress || dutyStationPostal,
                      scheduledPickupDate: formattedScheduledPickup,
                      shipmentStatus: mtoShipment.status,
                    }}
                  />
                  <ImportantShipmentDates
                    requestedPickupDate={formatShipmentDate(mtoShipment.requestedPickupDate)}
                    scheduledPickupDate={formattedScheduledPickup}
                  />
                  <ShipmentAddresses
                    pickupAddress={pickupAddress}
                    destinationAddress={destinationAddress || dutyStationPostal}
                    originDutyStation={moveOrder.originDutyStation?.address}
                    destinationDutyStation={moveOrder.destinationDutyStation?.address}
                  />
                  <ShipmentWeightDetails
                    estimatedWeight={mtoShipment.primeEstimatedWeight}
                    actualWeight={mtoShipment.primeActualWeight}
                  />
                  {requestedServiceItems?.length > 0 && (
                    <RequestedServiceItemsTable
                      serviceItems={requestedServiceItems}
                      handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                      handleShowRejectionDialog={handleShowRejectionDialog}
                      statusForTableType={SERVICE_ITEM_STATUSES.SUBMITTED}
                    />
                  )}
                  {approvedServiceItems?.length > 0 && (
                    <RequestedServiceItemsTable
                      serviceItems={approvedServiceItems}
                      handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                      handleShowRejectionDialog={handleShowRejectionDialog}
                      statusForTableType={SERVICE_ITEM_STATUSES.APPROVED}
                    />
                  )}
                  {rejectedServiceItems?.length > 0 && (
                    <RequestedServiceItemsTable
                      serviceItems={rejectedServiceItems}
                      handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                      handleShowRejectionDialog={handleShowRejectionDialog}
                      statusForTableType={SERVICE_ITEM_STATUSES.REJECTED}
                    />
                  )}
                </ShipmentContainer>
              </div>
            );
          })}
        </GridContainer>
      </div>
    </div>
  );
};

MoveTaskOrder.propTypes = {
  match: MatchShape.isRequired,
  setUnapprovedShipmentCount: func.isRequired,
};

export default withRouter(MoveTaskOrder);
