import React, { Component } from 'react';
import profileComplete from 'shared/images/profile-complete-gray-icon.png';
import ordersIcon from 'shared/images/orders-icon.png';
import 'scenes/ServiceMembers/ServiceMembers.css';

export class TransitionToOrders extends Component {
  render() {
    return (
      <div className="usa-grid">
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-3 desktop:margin-top-1" src={profileComplete} alt="profile-check" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">OK, your profile's complete!</h1>
          </div>
        </div>
        <div className="grid-row grid-gap">
          <div className="grid-col-3 desktop:grid-col-2 text-right">
            <img className="sm margin-top-6 desktop:margin-top-2" src={ordersIcon} alt="onto-move-orders" />
          </div>
          <div className="grid-col-9 desktop:grid-col-10">
            <h1 className="sm-heading">Now, we need to take a look at your move orders.</h1>
          </div>
        </div>
      </div>
    );
  }
}

export default TransitionToOrders;
