import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';

import {
  getMoveOrder,
  getCustomer,
  getAllMoveTaskOrders,
  selectMoveOrder,
  selectCustomer,
} from '../../shared/Entities/modules/moveTaskOrders';
import { loadOrders } from '../../shared/Entities/modules/orders';
import OrdersTable from '../../components/Office/OrdersTable';

import 'pages/TOO/too.scss';

import CustomerInfoTable from 'components/Office/CustomerInfoTable';
import { getMTOShipments, selectMTOShipments } from 'shared/Entities/modules/mtoShipments';
import RequestedShipments from 'components/Office/RequestedShipments';
import AllowancesTable from 'components/Office/AllowancesTable';

class MoveDetails extends Component {
  componentDidMount() {
    /* eslint-disable */
    const { moveOrderId } = this.props.match.params;
    this.props.getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      this.props.getCustomer(moveOrder.customerID);
      this.props.getAllMoveTaskOrders(moveOrder.id).then(({ response: { body: moveTaskOrder } }) => {
        moveTaskOrder.forEach((item) => this.props.getMTOShipments(item.id));
      });
    });
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrder, allowances, customer, mtoShipments } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
          <RequestedShipments mtoShipments={mtoShipments} />
          <OrdersTable
            ordersInfo={{
              // eslint-disable-next-line react/prop-types
              newDutyStation: get(moveOrder.destinationDutyStation, 'name'),
              // eslint-disable-next-line react/prop-types
              currentDutyStation: get(moveOrder.originDutyStation, 'name'),
              // eslint-disable-next-line react/prop-types
              issuedDate: moveOrder.date_issued,
              // eslint-disable-next-line react/prop-types
              reportByDate: moveOrder.report_by_date,
              // eslint-disable-next-line react/prop-types
              departmentIndicator: moveOrder.department_indicator,
              // eslint-disable-next-line react/prop-types
              ordersNumber: moveOrder.order_number,
              // eslint-disable-next-line react/prop-types
              ordersType: moveOrder.order_type,
              // eslint-disable-next-line react/prop-types
              ordersTypeDetail: moveOrder.order_type_detail,
              // eslint-disable-next-line react/prop-types
              tacMDC: moveOrder.tac,
              // eslint-disable-next-line react/prop-types
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

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const allowances = get(moveOrder, 'entitlement', {});
  const customerId = moveOrder.customerID;

  return {
    moveOrder,
    allowances,
    customer: selectCustomer(state, customerId),
    mtoShipments: selectMTOShipments(state, moveOrderId),
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  loadOrders,
  getCustomer,
  getAllMoveTaskOrders,
  getMTOShipments,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
