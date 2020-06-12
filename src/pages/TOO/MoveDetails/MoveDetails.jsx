import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import 'pages/TOO/too.scss';

import {
  getMoveOrder as getMoveOrderAction,
  getCustomer as getCustomerAction,
  getAllMoveTaskOrders as getAllMoveTaskOrdersAction,
  selectMoveOrder,
  selectCustomer,
} from 'shared/Entities/modules/moveTaskOrders';
import { loadOrders } from 'shared/Entities/modules/orders';
import CustomerInfoTable from 'components/Office/CustomerInfoTable';
import { getMTOShipments as getMTOShipmentsAction, selectMTOShipments } from 'shared/Entities/modules/mtoShipments';
import RequestedShipments from 'components/Office/RequestedShipments';
import AllowancesTable from 'components/Office/AllowancesTable';
import OrdersTable from 'components/Office/OrdersTable';
import { MoveOrderShape, EntitlementShape, CustomerShape, MTOShipmentShape } from 'types/moveOrder';
import { MatchShape } from 'types/router';

export class MoveDetails extends Component {
  componentDidMount() {
    // TODO - API flow
    const { match, getMoveOrder, getCustomer, getAllMoveTaskOrders, getMTOShipments } = this.props;
    const { params } = match;
    const { moveOrderId } = params;

    getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      getCustomer(moveOrder.customerID);
      getAllMoveTaskOrders(moveOrder.id).then(({ response: { body: moveTaskOrder } }) => {
        moveTaskOrder.forEach((item) => getMTOShipments(item.id));
      });
    });
  }

  render() {
    const { moveOrder, allowances, customer, mtoShipments } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
          <RequestedShipments mtoShipments={mtoShipments} />
          <OrdersTable
            ordersInfo={{
              newDutyStation: moveOrder.destinationDutyStation?.name,
              currentDutyStation: moveOrder.originDutyStation?.name,
              issuedDate: moveOrder.date_issued,
              reportByDate: moveOrder.report_by_date,
              departmentIndicator: moveOrder.department_indicator,
              ordersNumber: moveOrder.order_number,
              ordersType: moveOrder.order_type,
              ordersTypeDetail: moveOrder.order_type_detail,
              tacMDC: moveOrder.tac,
              sacSDN: moveOrder.sacSDN,
            }}
          />
          <AllowancesTable
            info={{
              branch: customer.agency,
              rank: moveOrder.grade,
              weightAllowance: allowances.totalWeight,
              authorizedWeight: allowances.authorizedWeight,
              progear: allowances.proGearWeight,
              spouseProgear: allowances.proGearWeightSpouse,
              storageInTransit: allowances.storageInTransit,
              dependents: allowances.dependentsAuthorized,
            }}
          />
          <CustomerInfoTable
            customerInfo={{
              name: `${customer.last_name}, ${customer.first_name}`,
              dodId: customer.dodID,
              phone: `+1 ${customer.phone}`,
              email: customer.email,
              currentAddress: customer.current_address,
              destinationAddress: customer.destination_address,
              backupContactName: '',
              backupContactPhone: '',
              backupContactEmail: '',
            }}
          />
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
  getMTOShipments: PropTypes.func.isRequired,
  moveOrder: MoveOrderShape,
  allowances: EntitlementShape,
  customer: CustomerShape,
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape),
};

MoveDetails.defaultProps = {
  moveOrder: {},
  allowances: {},
  customer: {},
  mtoShipments: [],
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const allowances = moveOrder?.entitlement;
  const customerId = moveOrder.customerID;

  return {
    moveOrder,
    allowances,
    customer: selectCustomer(state, customerId),
    mtoShipments: selectMTOShipments(state, moveOrderId),
  };
};

const mapDispatchToProps = {
  getMoveOrder: getMoveOrderAction,
  loadOrders,
  getCustomer: getCustomerAction,
  getAllMoveTaskOrders: getAllMoveTaskOrdersAction,
  getMTOShipments: getMTOShipmentsAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
