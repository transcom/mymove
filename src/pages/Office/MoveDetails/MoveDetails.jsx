import React, { useState } from 'react';
import classnames from 'classnames';
import { useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { patchMTOShipmentStatus } from 'shared/Entities/modules/mtoShipments';
import 'styles/office.scss';
import { updateMoveTaskOrderStatus } from 'shared/Entities/modules/moveTaskOrders';
import LeftNav from 'components/LeftNav';
import CustomerInfoTable from 'components/Office/CustomerInfoTable';
import RequestedShipments from 'components/Office/RequestedShipments/RequestedShipments';
import AllowancesTable from 'components/Office/AllowancesTable/AllowancesTable';
import OrdersTable from 'components/Office/OrdersTable/OrdersTable';
import { useMoveDetailsQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const sectionLabels = {
  'requested-shipments': 'Requested shipments',
  'approved-shipments': 'Approved shipments',
  orders: 'Orders',
  allowances: 'Allowances',
  'customer-info': 'Customer info',
};

// TODO - Convert to functional component
const MoveDetails = () => {
  const { moveCode } = useParams();

  // eslint-disable-next-line no-unused-vars
  const [activeSection, setActiveSection] = useState('');
  const [sections, setSections] = useState(['orders', 'allowances', 'customer-info']);

  const { move, moveOrder, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveDetailsQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = moveOrder;

  const approvedShipments = mtoShipments.filter((shipment) => shipment.status === 'APPROVED');
  const submittedShipments = mtoShipments.filter((shipment) => shipment.status === 'SUBMITTED');

  const hasSubmittedShipments = sections.includes('requested-shipments');
  const hasApprovedShipments = sections.includes('approved-shipments');

  if (submittedShipments.length > 0 && approvedShipments.length > 0) {
    if (!(hasApprovedShipments && hasSubmittedShipments)) {
      setSections(['approved-shipments', 'requested-shipments', 'orders', 'allowances', 'customer-info']);
    }
  } else if (approvedShipments.length > 0 && !hasApprovedShipments) {
    setSections(['approved-shipments', 'allowances', 'customer-info']);
  } else if (submittedShipments.length > 0 && !hasSubmittedShipments) {
    setSections(['requested-shipments', 'allowances', 'customer-info']);
  }

  /*
  componentDidMount() {
    // TODO - useEffects can be used for for this
    // attach scroll listener
    window.addEventListener('scroll', this.handleScroll);
  }
  */

  /*
  // TODO - useEffects can be used for for this
  componentWillUnmount() {
    // remove scroll listener
    window.removeEventListener('scroll', this.handleScroll);
  }
  */

  /*
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
  */

  const ordersInfo = {
    newDutyStation: moveOrder.destinationDutyStation,
    currentDutyStation: moveOrder.originDutyStation,
    issuedDate: moveOrder.date_issued,
    reportByDate: moveOrder.report_by_date,
    departmentIndicator: moveOrder.department_indicator,
    ordersNumber: moveOrder.order_number,
    ordersType: moveOrder.order_type,
    ordersTypeDetail: moveOrder.order_type_detail,
    tacMDC: moveOrder.tac,
    sacSDN: moveOrder.sac,
  };
  const allowancesInfo = {
    branch: customer.agency,
    rank: moveOrder.grade,
    weightAllowance: allowances.totalWeight,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
  };
  const customerInfo = {
    name: `${customer.last_name}, ${customer.first_name}`,
    dodId: customer.dodID,
    phone: `+1 ${customer.phone}`,
    email: customer.email,
    currentAddress: customer.current_address,
    backupContact: customer.backup_contact,
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav className={styles.sidebar}>
          {sections.map((s) => {
            const classes = classnames({ active: s === activeSection });
            return (
              <a key={`sidenav_${s}`} href={`#${s}`} className={classes}>
                {/* eslint-disable-next-line security/detect-object-injection */}
                {sectionLabels[s]}
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
                shipmentsStatus="SUBMITTED"
                approveMTO={updateMoveTaskOrderStatus}
                approveMTOShipment={patchMTOShipmentStatus}
                moveTaskOrder={move}
              />
            </div>
          )}
          {approvedShipments.length > 0 && (
            <div className={styles.section} id="approved-shipments">
              <RequestedShipments
                mtoShipments={approvedShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                mtoServiceItems={mtoServiceItems}
                shipmentsStatus="APPROVED"
                moveTaskOrder={move}
              />
            </div>
          )}
          <div className={styles.section} id="orders">
            <GridContainer>
              <Grid row gap>
                <Grid col>
                  <OrdersTable ordersInfo={ordersInfo} />
                </Grid>
              </Grid>
            </GridContainer>
          </div>
          <div className={styles.section} id="allowances">
            <GridContainer>
              <Grid row gap>
                <Grid col>
                  <AllowancesTable info={allowancesInfo} showEditBtn />
                </Grid>
              </Grid>
            </GridContainer>
          </div>
          <div className={styles.section} id="customer-info">
            <GridContainer>
              <Grid row gap>
                <Grid col>
                  <CustomerInfoTable customerInfo={customerInfo} />
                </Grid>
              </Grid>
            </GridContainer>
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

export default MoveDetails;
