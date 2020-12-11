import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { getMTOAgentList, selectMTOAgents } from 'shared/Entities/modules/mtoAgents';
import {
  patchMTOShipmentStatus as patchMTOShipmentStatusAction,
  getMTOShipments as getMTOShipmentsAction,
  selectMTOShipments,
} from 'shared/Entities/modules/mtoShipments';
import 'styles/office.scss';
import {
  getMoveOrder as getMoveOrderAction,
  getCustomer as getCustomerAction,
  getAllMoveTaskOrders as getAllMoveTaskOrdersAction,
  updateMoveTaskOrderStatus as updateMoveTaskOrderStatusAction,
  selectMoveTaskOrders,
  selectMoveOrder,
  selectCustomer,
} from 'shared/Entities/modules/moveTaskOrders';
import {
  getMTOServiceItems as getMTOServiceItemsAction,
  selectMTOServiceItems,
} from 'shared/Entities/modules/mtoServiceItems';
import { loadOrders } from 'shared/Entities/modules/orders';
import LeftNav from 'components/LeftNav';
import CustomerInfoTable from 'components/Office/CustomerInfoTable';
import RequestedShipments from 'components/Office/RequestedShipments/RequestedShipments';
import AllowancesTable from 'components/Office/AllowancesTable/AllowancesTable';
import OrdersTable from 'components/Office/OrdersTable/OrdersTable';
import {
  MoveOrderShape,
  EntitlementShape,
  CustomerShape,
  MTOShipmentShape,
  MTOAgentShape,
  MTOServiceItemShape,
  MoveTaskOrderShape,
} from 'types/moveOrder';
import { MatchShape } from 'types/router';

const sectionLabels = {
  'requested-shipments': 'Requested shipments',
  'approved-shipments': 'Approved shipments',
  orders: 'Orders',
  allowances: 'Allowances',
  'customer-info': 'Customer info',
};

export class MoveDetails extends Component {
  constructor(props) {
    super(props);

    this.state = {
      activeSection: '',
      sections: ['orders', 'allowances', 'customer-info'],
    };
  }

  componentDidMount() {
    // attach scroll listener
    window.addEventListener('scroll', this.handleScroll);

    // TODO - API flow
    const { match, getMoveOrder, getCustomer, getAllMoveTaskOrders, getMTOShipments, getMTOServiceItems } = this.props;
    const { params } = match;
    const { moveOrderId } = params;

    getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      getCustomer(moveOrder.customerID);
      getAllMoveTaskOrders(moveOrder.id).then(({ response: { body: moveTaskOrder } }) => {
        moveTaskOrder.forEach((item) =>
          getMTOShipments(item.id).then(({ response: { body: mtoShipments } }) => {
            mtoShipments.map((shipment) => getMTOAgentList(shipment.moveTaskOrderID, shipment.id));
            this.checkToAddShipmentsSections(mtoShipments);
            getMTOServiceItems(item.id);
          }),
        );
      });
    });
  }

  componentWillUnmount() {
    // remove scroll listener
    window.removeEventListener('scroll', this.handleScroll);
  }

  setActiveSection = (sectionId) => {
    this.setState({
      activeSection: sectionId,
    });
  };

  handleScroll = () => {
    const distanceFromTop = window.scrollY;
    const { sections, activeSection } = this.state;
    let newActiveSection;

    sections.forEach((section) => {
      const sectionEl = document.querySelector(`#${section}`);
      if (sectionEl?.offsetTop <= distanceFromTop && sectionEl?.offsetTop + sectionEl?.offsetHeight > distanceFromTop) {
        newActiveSection = section;
      }
    });

    if (activeSection !== newActiveSection) {
      this.setActiveSection(newActiveSection);
    }
  };

  checkToAddShipmentsSections = (shipments) => {
    const approvedShipments = shipments.filter((shipment) => shipment.status === 'APPROVED');
    const submittedShipments = shipments.filter((shipment) => shipment.status === 'SUBMITTED');

    if (submittedShipments.length > 0 && approvedShipments.length > 0) {
      this.setState((previousState) => ({
        sections: ['approved-shipments', 'requested-shipments', ...previousState.sections],
      }));
    } else if (approvedShipments.length > 0) {
      this.setState((previousState) => ({ sections: ['approved-shipments', ...previousState.sections] }));
    } else if (submittedShipments.length > 0) {
      this.setState((previousState) => ({ sections: ['requested-shipments', ...previousState.sections] }));
    }
  };

  render() {
    const {
      moveOrder,
      allowances,
      customer,
      mtoShipments,
      mtoAgents,
      mtoServiceItems,
      moveTaskOrder,
      updateMoveTaskOrderStatus,
      patchMTOShipmentStatus,
    } = this.props;

    const approvedShipments = mtoShipments.filter((shipment) => shipment.status === 'APPROVED');
    const submittedShipments = mtoShipments.filter((shipment) => shipment.status === 'SUBMITTED');

    const { activeSection, sections } = this.state;

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
            {submittedShipments.length > 0 && (
              <div className={styles.section} id="requested-shipments">
                <RequestedShipments
                  mtoShipments={submittedShipments}
                  ordersInfo={ordersInfo}
                  allowancesInfo={allowancesInfo}
                  customerInfo={customerInfo}
                  mtoAgents={mtoAgents}
                  mtoServiceItems={mtoServiceItems}
                  shipmentsStatus="SUBMITTED"
                  approveMTO={updateMoveTaskOrderStatus}
                  approveMTOShipment={patchMTOShipmentStatus}
                  moveTaskOrder={moveTaskOrder}
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
                  mtoAgents={mtoAgents}
                  mtoServiceItems={mtoServiceItems}
                  shipmentsStatus="APPROVED"
                  moveTaskOrder={moveTaskOrder}
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
  }
}

MoveDetails.propTypes = {
  match: MatchShape.isRequired,
  getMoveOrder: PropTypes.func.isRequired,
  getCustomer: PropTypes.func.isRequired,
  getAllMoveTaskOrders: PropTypes.func.isRequired,
  updateMoveTaskOrderStatus: PropTypes.func.isRequired,
  patchMTOShipmentStatus: PropTypes.func.isRequired,
  getMTOShipments: PropTypes.func.isRequired,
  getMTOServiceItems: PropTypes.func.isRequired,
  moveOrder: MoveOrderShape,
  allowances: EntitlementShape,
  customer: CustomerShape,
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape),
  mtoAgents: PropTypes.arrayOf(MTOAgentShape),
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
  moveTaskOrder: MoveTaskOrderShape,
};

MoveDetails.defaultProps = {
  moveOrder: {},
  allowances: {},
  customer: {},
  mtoShipments: [],
  mtoAgents: [],
  mtoServiceItems: [],
  moveTaskOrder: {},
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const allowances = moveOrder?.entitlement;
  const customerId = moveOrder.customerID;
  const moveTaskOrders = selectMoveTaskOrders(state, moveOrderId);
  const moveTaskOrder = moveTaskOrders[0];

  return {
    moveOrder,
    allowances,
    customer: selectCustomer(state, customerId),
    mtoShipments: selectMTOShipments(state, moveOrderId),
    mtoAgents: selectMTOAgents(state),
    mtoServiceItems: selectMTOServiceItems(state, moveOrderId),
    moveTaskOrder,
  };
};

const mapDispatchToProps = {
  getMoveOrder: getMoveOrderAction,
  loadOrders,
  getCustomer: getCustomerAction,
  getAllMoveTaskOrders: getAllMoveTaskOrdersAction,
  updateMoveTaskOrderStatus: updateMoveTaskOrderStatusAction,
  getMTOShipments: getMTOShipmentsAction,
  getMTOAgentList,
  getMTOServiceItems: getMTOServiceItemsAction,
  patchMTOShipmentStatus: patchMTOShipmentStatusAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
