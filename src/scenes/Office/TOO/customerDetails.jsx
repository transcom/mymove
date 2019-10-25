import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import { getEntitlements, updateMoveTaskOrderStatus } from 'shared/Entities/modules/moveTaskOrders';
import { selectServiceMember } from 'shared/Entities/modules/serviceMembers';
import { selectMoveTaskOrder } from '../../../shared/Entities/modules/moveTaskOrders';

class CustomerDetails extends React.Component {
  componentDidMount() {
    const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
    this.props.getEntitlements(fakeMoveTaskOrderID);
  }

  render() {
    const { entitlements, moveTaskOrder, customer } = this.props;
    console.log('customer:', customer);
    console.log('mto:', moveTaskOrder);
    const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
    const NTS = entitlements && entitlements.nonTemporaryStorage ? 'Y' : 'N';
    const POV = entitlements && entitlements.privatelyOwnedVehicle ? 'Y' : 'N';
    const moveTaskOrderNonTemporaryStorage = get(moveTaskOrder, 'entitlements.nonTemporaryStorage') ? 'Y' : 'N';
    const moveTaskOrderPrivatelyOwnedVehicle = get(moveTaskOrder, 'entitlements.privatelyOwnedVehicle') ? 'Y' : 'N';
    return (
      <>
        <h1>Customer Deets Page</h1>
        {entitlements && (
          <>
            <h2>Customer Entitlements</h2>
            <dl>
              <dt>Weight Entitlement</dt>
              <dd>{entitlements.totalWeightSelf}</dd>
              <dt>SIT Entitlement</dt>
              <dd>{entitlements.storageInTransit}</dd>
              <dt>NTS Entitlement</dt>
              <dd>{NTS}</dd>
              <dt>POV Entitlement</dt>
              <dd>{POV}</dd>
            </dl>
          </>
        )}
        {!isEmpty(moveTaskOrder) && (
          <>
            <h2>Move Task Order</h2>
            <dl>
              <dt>First Name</dt>
              <dd>{get(customer, 'first_name')}</dd>
              <dt>Last Name</dt>
              <dd>{get(customer, 'last_name')}</dd>
              <dt>Rank</dt>
              <dd>{get(customer, 'rank')}</dd>
              <dt>Email Address</dt>
              <dd>{get(customer, 'personal_email')}</dd>
              <dt>Phone</dt>
              <dd>{get(customer, 'telephone')}</dd>
              <dt>Origin Duty Station</dt>
              <dd>{get(moveTaskOrder, 'originDutyStation')}</dd>
              <dt>Destination Duty Station</dt>
              <dd>{get(moveTaskOrder, 'destinationDutyStation')}</dd>
              <dt>Pickup Address</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'pickupAddress'))}</dd>
              <dt>Destination Address</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'destinationAddress'))}</dd>
              <dt>Requested Pickup Date</dt>
              <dd>{get(moveTaskOrder, 'requestedPickupDate')}</dd>
              <dt>Customer Remarks</dt>
              <dd>{get(moveTaskOrder, 'remarks')}</dd>
              <dt>Service Items</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'serviceItems'))}</dd>
              <dt>Status</dt>
              <dd>{get(moveTaskOrder, 'status')}</dd>
              <dt>Weight Entitlement</dt>
              <dd>{get(moveTaskOrder, 'entitlements.totalWeightSelf')}</dd>
              <dt>SIT Entitlement</dt>
              <dd>{get(moveTaskOrder, 'entitlements.storageInTransit')}</dd>
              <dt>NTS Entitlement</dt>
              <dd>{moveTaskOrderNonTemporaryStorage}</dd>
              <dt>POV Entitlement</dt>
              <dd>{moveTaskOrderPrivatelyOwnedVehicle}</dd>
            </dl>
            {/*- Pickup Address*/}
            {/*- Destination Address (if known)*/}
          </>
        )}
        <div>
          <button onClick={() => this.props.updateMoveTaskOrderStatus(fakeMoveTaskOrderID, 'DRAFT')}>
            Generate MTO
          </button>
        </div>
      </>
    );
  }
}

const mapStateToProps = state => {
  //TODO hard coding mto for now
  const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
  const entitlements = get(state, 'entities.entitlements');
  const moveTaskOrder = selectMoveTaskOrder(state, fakeMoveTaskOrderID);
  const customerId = moveTaskOrder ? moveTaskOrder.customer : null;
  // TODO customer is service member for now
  const customer = selectServiceMember(state, customerId);
  return {
    entitlements: entitlements && Object.values(entitlements).length > 0 ? Object.values(entitlements)[0] : null,
    moveTaskOrder,
    customer,
  };
};
const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      getEntitlements,
      updateMoveTaskOrderStatus,
    },
    dispatch,
  );

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(CustomerDetails);
