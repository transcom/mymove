import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import '../../index.scss';
import '../../ghc_index.scss';
import { get } from 'lodash';
import RequestedShipments from 'components/Office/RequestedShipments';

import { getMTOShipments, selectMTOShipments } from 'shared/Entities/modules/mtoShipments';
import ShipmentDisplay from 'components/Office/ShipmentDisplay';

import { getMoveOrder, getAllMoveTaskOrders, selectMoveOrder } from '../../shared/Entities/modules/moveTaskOrders';
import { loadOrders } from '../../shared/Entities/modules/orders';
import OrdersTable from '../../components/Office/OrdersTable';

class MoveDetails extends Component {
  componentDidMount() {
    /* eslint-disable */
    const { moveOrderId } = this.props.match.params;
    this.props.getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      this.props.getAllMoveTaskOrders(moveOrder.id).then(({ response: { body: moveTaskOrder } }) => {
        moveTaskOrder.forEach((item) => this.props.getMTOShipments(item.id));
      });
    });
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrder, mtoShipments } = this.props;
    return (
      <div className="grid-container-desktop-lg" data-cy="too-move-details">
        <h1>Move details</h1>
        <div className="container">
          <RequestedShipments>
            {mtoShipments &&
              mtoShipments.map((shipment) => (
                <ShipmentDisplay
                  key={shipment.id}
                  shipmentType={shipment.shipmentType}
                  displayInfo={{
                    heading: shipment.shipmentType,
                    requestedMoveDate: shipment.requestedPickupDate,
                    currentAddress: shipment.pickupAddress,
                    destinationAddress: shipment.destinationAddress,
                  }}
                />
              ))}
          </RequestedShipments>
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
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;

  return {
    moveOrder: selectMoveOrder(state, moveOrderId),
    mtoShipments: selectMTOShipments(state, moveOrderId),
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  loadOrders,
  getAllMoveTaskOrders,
  getMTOShipments,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
