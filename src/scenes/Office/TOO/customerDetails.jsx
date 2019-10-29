import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { getEntitlements, getCustomerInfo } from 'shared/Entities/modules/moveTaskOrders';
import { selectCustomer } from 'shared/Entities/modules/customer';

class CustomerDetails extends Component {
  componentDidMount() {
    this.props.getEntitlements('fake_move_task_order_id');
    this.props.getCustomerInfo(this.props.match.params.customerId);
  }

  render() {
    const { entitlements, customer } = this.props;
    const NTS = entitlements && entitlements.nonTemporaryStorage ? 'Y' : 'N';
    const POV = entitlements && entitlements.privatelyOwnedVehicle ? 'Y' : 'N';
    return (
      <>
        <h1>Customer Deets Page</h1>
        {customer && (
          <>
            <h2>Customer Info</h2>
            <dl>
              <dt>Full Name</dt>
              <dd>{customer.customer_name}</dd>
              <dt>Service Branch / Agency</dt>
              <dd>{customer.agency}</dd>
              <dt>Rank / Grade</dt>
              <dd>{customer.grade}</dd>
              <dt>Email</dt>
              <dd>{customer.email}</dd>
              <dt>Phone</dt>
              <dd>{customer.telephone}</dd>
              <dt>Origin Duty Station</dt>
              <dd>{customer.origin_duty_station}</dd>
              <dt>Destination Duty Station</dt>
              <dd>{customer.destination_duty_station}</dd>
              <dt>Pickup Address</dt>
              <dd>{customer.pickup_address}</dd>
              <dt>Dependents Authorized</dt>
              <dd>{customer.dependents_authorized ? 'Y' : 'N'}</dd>
            </dl>
          </>
        )}
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
      </>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const entitlements = get(state, 'entities.entitlements');
  return {
    entitlements: entitlements && Object.values(entitlements).length > 0 ? Object.values(entitlements)[0] : null,
    customer: selectCustomer(state, ownProps.match.params.customerId),
  };
};
const mapDispatchToProps = { getEntitlements, getCustomerInfo };

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(CustomerDetails);
