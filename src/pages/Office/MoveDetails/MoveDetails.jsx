import React, { useState, useEffect } from 'react';
import { useParams, useHistory, Link } from 'react-router-dom';
import { GridContainer, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { queryCache, useMutation } from 'react-query';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import 'styles/office.scss';
import { MOVES, MTO_SHIPMENTS, MTO_SERVICE_ITEMS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { shipmentStatuses } from 'constants/shipments';
import LeftNav from 'components/LeftNav';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import RequestedShipments from 'components/Office/RequestedShipments/RequestedShipments';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatus, updateMTOShipmentStatus } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const sectionLabels = {
  'requested-shipments': 'Requested shipments',
  'approved-shipments': 'Approved shipments',
  orders: 'Orders',
  allowances: 'Allowances',
  'customer-info': 'Customer info',
};

const MoveDetails = ({ setUnapprovedShipmentCount, setUnapprovedServiceItemCount }) => {
  const { moveCode } = useParams();
  const history = useHistory();

  const [activeSection, setActiveSection] = useState('');

  const { move, order, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveDetailsQueries(moveCode);

  let sections = ['orders', 'allowances', 'customer-info'];

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

  // use mutation calls
  const [mutateMoveStatus] = useMutation(updateMoveStatus, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);
      queryCache.invalidateQueries([MTO_SERVICE_ITEMS, data.id]);
    },
  });

  const [mutateMTOShipmentStatus] = useMutation(updateMTOShipmentStatus, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      queryCache.invalidateQueries([MTO_SERVICE_ITEMS, updatedMTOShipment.moveTaskOrderID]);
    },
  });

  const submittedShipments = mtoShipments?.filter(
    (shipment) => shipment.status === shipmentStatuses.SUBMITTED && !shipment.deletedAt,
  );
  const approvedOrCanceledShipments = mtoShipments?.filter(
    (shipment) => shipment.status === shipmentStatuses.APPROVED || shipment.status === shipmentStatuses.CANCELED,
  );

  useEffect(() => {
    const shipmentCount = submittedShipments?.length || 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, submittedShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    let serviceItemCount = 0;
    mtoServiceItems?.forEach((serviceItem) => {
      if (
        serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED &&
        serviceItem.mtoShipmentID &&
        approvedOrCanceledShipments?.find((shipment) => shipment.id === serviceItem.mtoShipmentID)
      ) {
        serviceItemCount += 1;
      }
    });
    setUnapprovedServiceItemCount(serviceItemCount);
  }, [approvedOrCanceledShipments, mtoServiceItems, setUnapprovedServiceItemCount]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;

  if (submittedShipments.length > 0 && approvedOrCanceledShipments.length > 0) {
    sections = ['requested-shipments', 'approved-shipments', ...sections];
  } else if (approvedOrCanceledShipments.length > 0) {
    sections = ['approved-shipments', ...sections];
  } else if (submittedShipments.length > 0) {
    sections = ['requested-shipments', ...sections];
  }

  const ordersInfo = {
    newDutyStation: order.destinationDutyStation,
    currentDutyStation: order.originDutyStation,
    issuedDate: order.date_issued,
    reportByDate: order.report_by_date,
    departmentIndicator: order.department_indicator,
    ordersNumber: order.order_number,
    ordersType: order.order_type,
    ordersTypeDetail: order.order_type_detail,
    tacMDC: order.tac,
    sacSDN: order.sac,
  };
  const allowancesInfo = {
    branch: customer.agency,
    rank: order.grade,
    weightAllowance: allowances.totalWeight,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
  };
  const customerInfo = {
    name: `${customer.last_name}, ${customer.first_name}`,
    dodId: customer.dodID,
    phone: `+1 ${customer.phone}`,
    email: customer.email,
    currentAddress: customer.current_address,
    backupContact: customer.backup_contact,
  };

  const requiredOrdersInfo = {
    ordersNumber: order.order_number,
    ordersType: order.order_type,
    ordersTypeDetail: order.order_type_detail,
    tacMDC: order.tac,
  };

  const hasMissingOrdersRequiredInfo = Object.values(requiredOrdersInfo).some((value) => !value || value === '');

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav className={styles.sidebar}>
          {sections.map((s) => {
            return (
              <a key={`sidenav_${s}`} href={`#${s}`} className={classnames({ active: s === activeSection })}>
                {sectionLabels[`${s}`]}
                {s === 'orders' && hasMissingOrdersRequiredInfo && (
                  <Tag className="usa-tag usa-tag--alert">
                    <FontAwesomeIcon icon="exclamation" />
                  </Tag>
                )}
                {s === 'requested-shipments' && (
                  <Tag className={styles.tag} data-testid="requestedShipmentsTag">
                    {submittedShipments?.length}
                  </Tag>
                )}
              </a>
            );
          })}
        </LeftNav>

        <GridContainer className={styles.gridContainer} data-testid="too-move-details">
          <h1>Move details</h1>
          {/* TODO - RequestedShipments could be simplified, if extra time we could tackle this or just write a story to track */}
          {submittedShipments.length > 0 && (
            <div className={styles.section} id="requested-shipments">
              <RequestedShipments
                mtoShipments={submittedShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                mtoServiceItems={mtoServiceItems}
                shipmentsStatus={shipmentStatuses.SUBMITTED}
                approveMTO={mutateMoveStatus}
                approveMTOShipment={mutateMTOShipmentStatus}
                moveTaskOrder={move}
                missingRequiredOrdersInfo={hasMissingOrdersRequiredInfo}
                handleAfterSuccess={history.push}
              />
            </div>
          )}
          {approvedOrCanceledShipments.length > 0 && (
            <div className={styles.section} id="approved-shipments">
              <RequestedShipments
                moveTaskOrder={move}
                mtoShipments={approvedOrCanceledShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                mtoServiceItems={mtoServiceItems}
                shipmentsStatus={shipmentStatuses.APPROVED}
              />
            </div>
          )}
          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              editButton={
                <Link className="usa-button usa-button--secondary" data-testid="edit-orders" to="orders">
                  Edit orders
                </Link>
              }
            >
              <OrdersList ordersInfo={ordersInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                <Link className="usa-button usa-button--secondary" data-testid="edit-allowances" to="allowances">
                  Edit allowances
                </Link>
              }
            >
              <AllowancesList info={allowancesInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel title="Customer info">
              <CustomerInfoList customerInfo={customerInfo} />
            </DetailsPanel>
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

MoveDetails.propTypes = {
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
};

export default MoveDetails;
