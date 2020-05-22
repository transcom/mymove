import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import '../../index.scss';
import '../../ghc_index.scss';
import { get } from 'lodash';
import CustomerInfoTable from 'components/Office/CustomerInfoTable';
import {
  getMoveOrder,
  getCustomer,
  selectMoveOrder,
  selectCustomer,
} from '../../shared/Entities/modules/moveTaskOrders';
import { loadOrders } from '../../shared/Entities/modules/orders';
import OrdersTable from '../../components/Office/OrdersTable';

class MoveDetails extends Component {
  componentDidMount() {
    /* eslint-disable */
    const { moveOrderId } = this.props.match.params;
    this.props.getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      this.props.getCustomer(moveOrder.customerID);
    });
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrder, customer } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
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
  const customerId = moveOrder.customerID;

  return {
    moveOrder,
    customer: selectCustomer(state, customerId),
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  loadOrders,
  getCustomer,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
