import React, { Component } from 'react';
import profileComplete from 'shared/images/profile-complete-gray-icon.png';
import ordersIcon from 'shared/images/orders-icon.png';
import 'scenes/ServiceMembers/ServiceMembers.css';

export class TransitionToOrders extends Component {
  render() {
    return (
      <div className="usa-grid">
        <div className="grid-row">
          <div className="grid-col-12">
            <div className="lg center">
              <p> OK, your profile's complete!</p>
              <img className="sm" src={profileComplete} alt="profile-check" />
            </div>
          </div>
        </div>

        <div className="grid-row">
          <div className="grid-col-12">
            <div className="lg center">
              <p>Now, we need to take a look at your move orders.</p>
              <img className="sm" src={ordersIcon} alt="onto-move-orders" />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default TransitionToOrders;
