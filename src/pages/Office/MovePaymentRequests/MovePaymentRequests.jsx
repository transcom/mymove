import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import classnames from 'classnames';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';
import paymentRequestStatus from '../../../constants/paymentRequestStatus';

import LeftNav from 'components/LeftNav';
import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString, getShipmentModificationType } from 'utils/shipmentDisplay';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';

const sectionLabels = {
  'payment-requests': 'Payment requests',
};

const MovePaymentRequests = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setPendingPaymentRequestCount,
}) => {
  const { moveCode } = useParams();

  const { paymentRequests, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);
  const [activeSection, setActiveSection] = useState('');
  let sections = ['payment-requests'];

  useEffect(() => {
    const shipmentCount = mtoShipments
      ? mtoShipments.filter((shipment) => shipment.status === shipmentStatuses.SUBMITTED).length
      : 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    let serviceItemCount = 0;
    if (mtoShipments) {
      mtoShipments.forEach((shipment) => {
        if (shipment.status === shipmentStatuses.APPROVED) {
          serviceItemCount += shipment.mtoServiceItems?.filter(
            (serviceItem) => serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED,
          ).length;
        }
      });
    }
    setUnapprovedServiceItemCount(serviceItemCount);
  }, [mtoShipments, setUnapprovedServiceItemCount]);

  useEffect(() => {
    const pendingCount = paymentRequests?.filter((pr) => pr.status === paymentRequestStatus.PENDING).length;
    setPendingPaymentRequestCount(pendingCount);
  }, [paymentRequests, setPendingPaymentRequestCount]);

  const handleScroll = () => {
    const distanceFromTop = window.scrollY;
    let newActiveSection;

    sections.forEach((section) => {
      const sectionEl = document.querySelector(`#${section}`);
      if (sectionEl?.offsetTop <= distanceFromTop && sectionEl?.offsetTop + sectionEl?.offsetHeight > distanceFromTop) {
        newActiveSection = section;
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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentsInfo = [];

  if (paymentRequests.length) {
    mtoShipments.forEach((shipment) => {
      shipmentsInfo.push({
        mtoShipmentID: shipment.id,
        address: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
        departureDate: shipment.actualPickupDate,
        modificationType: getShipmentModificationType(shipment),
        mtoServiceItems: shipment.mtoServiceItems,
      });
    });
  }

  if (paymentRequests.length === 0) {
    sections = [];
  }

  return (
    <div className={txoStyles.tabContent}>
      <div className={txoStyles.container} data-testid="MovePaymentRequests">
        <LeftNav className={txoStyles.sidebar}>
          {sections.map((s) => {
            return (
              <a key={`sidenav_${s}`} href={`#${s}`} className={classnames({ active: s === activeSection })}>
                {sectionLabels[`${s}`]}
              </a>
            );
          })}
        </LeftNav>
        <GridContainer className={txoStyles.gridContainer} data-testid="tio-payment-request-details">
          <h1>Payment requests</h1>
          <div className={txoStyles.section} id="payment-requests">
            {paymentRequests.length ? (
              paymentRequests.map((paymentRequest) => (
                <PaymentRequestCard
                  paymentRequest={paymentRequest}
                  shipmentsInfo={shipmentsInfo}
                  key={paymentRequest.id}
                />
              ))
            ) : (
              <div className={txoStyles.emptyMessage}>
                <p>No payment requests have been submitted for this move yet.</p>
              </div>
            )}
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

MovePaymentRequests.propTypes = {
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
  setPendingPaymentRequestCount: func.isRequired,
};

export default MovePaymentRequests;
